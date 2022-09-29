package rules

import (
	"github.com/google/go-cmp/cmp"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformOrderedVariablesRule checks whether the variables are sorted in expected order
type TerraformOrderedVariablesRule struct {
	tflint.DefaultRule
}

// NewTerraformOrderedVariablesRule returns a new rule
func NewTerraformOrderedVariablesRule() *TerraformOrderedVariablesRule {
	return &TerraformOrderedVariablesRule{}
}

// Name returns the rule name
func (r *TerraformOrderedVariablesRule) Name() string {
	return "terraform_ordered_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformOrderedVariablesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformOrderedVariablesRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformOrderedVariablesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are sorted in expected order
func (r *TerraformOrderedVariablesRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: "default"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{IncludeNotCreated: true})
	if err != nil {
		return err
	}

	variables := r.variablesGroupByFile(body)
	for _, blocks := range variables {
		if err = r.checkVariableOrder(runner, blocks); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformOrderedVariablesRule) variablesGroupByFile(body *hclext.BodyContent) map[string]hclext.Blocks {
	variables := make(map[string]hclext.Blocks)
	for _, b := range body.Blocks {
		variables[b.DefRange.Filename] = append(variables[b.DefRange.Filename], b)
	}
	for _, blocks := range variables {
		sort.Slice(blocks, func(i, j int) bool {
			return blocks[i].DefRange.Start.Line < blocks[j].DefRange.Start.Line
		})
	}
	return variables
}

func (r *TerraformOrderedVariablesRule) checkVariableOrder(runner tflint.Runner, blocks []*hclext.Block) error {
	requiredVars := r.getSortedVariableNames(blocks, false)
	optionalVars := r.getSortedVariableNames(blocks, true)
	sortedVariableNames := append(requiredVars, optionalVars...)

	variableNames := r.getVariableNames(blocks)
	if cmp.Equal(variableNames, sortedVariableNames) {
		return nil
	}

	return runner.EmitIssue(
		r,
		"Variables should be sorted in the following order: required(without default value) variables in alphabetical order, optional variables in alphabetical order.",
		r.issueRange(blocks),
	)
}

func (r *TerraformOrderedVariablesRule) issueRange(blocks hclext.Blocks) hcl.Range {
	requiredVariables := r.getVariables(blocks, false)
	optionalVariables := r.getVariables(blocks, true)

	if r.overlapped(requiredVariables, optionalVariables) {
		return optionalVariables[0].DefRange
	}

	firstRange := r.firstNonSortedBlockRange(requiredVariables)
	if firstRange != nil {
		return *firstRange
	}
	firstRange = r.firstNonSortedBlockRange(optionalVariables)
	if firstRange != nil {
		return *firstRange
	}
	panic("expected issue not found")
}

func (r *TerraformOrderedVariablesRule) overlapped(requiredVariables, optionalVariables hclext.Blocks) bool {
	if len(requiredVariables) == 0 || len(optionalVariables) == 0 {
		return false
	}

	firstOptional := optionalVariables[0].DefRange
	lastRequired := requiredVariables[len(requiredVariables)-1].DefRange

	return firstOptional.Start.Line < lastRequired.Start.Line
}

func (r *TerraformOrderedVariablesRule) firstNonSortedBlockRange(blocks hclext.Blocks) *hcl.Range {
	for i, b := range blocks {
		if i == 0 {
			continue
		}
		previousVariableName := blocks[i-1].Labels[0]
		currentVariableName := b.Labels[0]
		if currentVariableName < previousVariableName {
			return &b.DefRange
		}
	}
	return nil
}

func (r *TerraformOrderedVariablesRule) getVariableNames(blocks hclext.Blocks) []string {
	var variableNames []string
	for _, b := range blocks {
		variableNames = append(variableNames, b.Labels[0])
	}
	return variableNames
}

func (r *TerraformOrderedVariablesRule) getSortedVariableNames(blocks hclext.Blocks, defaultWanted bool) []string {
	var names []string
	filteredBlocks := r.getVariables(blocks, defaultWanted)
	for _, b := range filteredBlocks {
		names = append(names, b.Labels[0])
	}
	sort.Strings(names)
	return names
}

func (r *TerraformOrderedVariablesRule) getVariables(blocks hclext.Blocks, defaultWanted bool) hclext.Blocks {
	var c hclext.Blocks
	for _, b := range blocks {
		if _, hasDefault := b.Body.Attributes["default"]; hasDefault == defaultWanted {
			c = append(c, b)
		}
	}
	return c
}
