package rules

import (
	"reflect"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformVariableOrderRule checks whether the variables are sorted in expected order
type TerraformVariableOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformOrderedVariablesRule returns a new rule
func NewTerraformOrderedVariablesRule() *TerraformVariableOrderRule {
	return &TerraformVariableOrderRule{}
}

// Name returns the rule name
func (r *TerraformVariableOrderRule) Name() string {
	return "terraform_ordered_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVariableOrderRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformVariableOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVariableOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are sorted in expected order
func (r *TerraformVariableOrderRule) Check(runner tflint.Runner) error {
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

func (r *TerraformVariableOrderRule) variablesGroupByFile(body *hclext.BodyContent) map[string]hclext.Blocks {
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

func (r *TerraformVariableOrderRule) checkVariableOrder(runner tflint.Runner, blocks []*hclext.Block) error {
	requiredVars := r.getSortedVariableNames(blocks, false)
	optionalVars := r.getSortedVariableNames(blocks, true)
	sortedVariableNames := append(requiredVars, optionalVars...)

	variableNames := r.getVariableNames(blocks)
	if reflect.DeepEqual(variableNames, sortedVariableNames) {
		return nil
	}

	return runner.EmitIssue(
		r,
		"Variables should be sorted in the following order: required(without default value) variables in alphabetical order, optional variables in alphabetical order.",
		*r.issueRange(blocks),
	)
}

func (r *TerraformVariableOrderRule) issueRange(blocks hclext.Blocks) *hcl.Range {
	requiredVariables := r.getVariables(blocks, false)
	optionalVariables := r.getVariables(blocks, true)

	for i, b := range requiredVariables {
		if i > 0 && (b.Labels[0] < requiredVariables[i-1].Labels[0]) {
			return &b.DefRange
		}
	}
	for i, b := range optionalVariables {
		if i > 0 && (b.Labels[0] < optionalVariables[i-1].Labels[0]) {
			return &b.DefRange
		}
	}

	firstOptional := optionalVariables[0].DefRange
	lastRequired := requiredVariables[len(requiredVariables)-1].DefRange

	if firstOptional.Start.Line < lastRequired.Start.Line {
		return &lastRequired
	}

	return nil
}

func (r *TerraformVariableOrderRule) getVariableNames(blocks hclext.Blocks) []string {
	var variableNames []string
	for _, b := range blocks {
		variableNames = append(variableNames, b.Labels[0])
	}
	return variableNames
}

func (r *TerraformVariableOrderRule) getSortedVariableNames(blocks hclext.Blocks, defaultWanted bool) []string {
	var names []string
	filteredBlocks := r.getVariables(blocks, defaultWanted)
	for _, b := range filteredBlocks {
		names = append(names, b.Labels[0])
	}
	sort.Strings(names)
	return names
}

func (r *TerraformVariableOrderRule) getVariables(blocks hclext.Blocks, defaultWanted bool) hclext.Blocks {
	var c hclext.Blocks
	for _, b := range blocks {
		if _, hasDefault := b.Body.Attributes["default"]; hasDefault == defaultWanted {
			c = append(c, b)
		}
	}
	return c
}
