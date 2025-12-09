package rules

import (
	"fmt"

	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OpentofuTypedVariablesRule checks whether variables have a type declared
type OpentofuTypedVariablesRule struct {
	tflint.DefaultRule
}

// NewOpentofuTypedVariablesRule returns a new rule
func NewOpentofuTypedVariablesRule() *OpentofuTypedVariablesRule {
	return &OpentofuTypedVariablesRule{}
}

// Name returns the rule name
func (r *OpentofuTypedVariablesRule) Name() string {
	return "opentofu_typed_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *OpentofuTypedVariablesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *OpentofuTypedVariablesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *OpentofuTypedVariablesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether variables have type
func (r *OpentofuTypedVariablesRule) Check(runner tflint.Runner) error {
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
					Attributes: []hclext.AttributeSchema{{Name: "type"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		if _, exists := variable.Body.Attributes["type"]; !exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("`%v` variable has no type", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
