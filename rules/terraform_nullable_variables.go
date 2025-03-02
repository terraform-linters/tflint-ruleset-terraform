package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformNullableVariablesRule checks whether variables have a nullable field declared
type TerraformNullableVariablesRule struct {
	tflint.DefaultRule
}

// NewTerraformNullableVariablesRule returns a new rule
func NewTerraformNullableVariablesRule() *TerraformNullableVariablesRule {
	return &TerraformNullableVariablesRule{}
}

// Name returns the rule name
func (r *TerraformNullableVariablesRule) Name() string {
	return "terraform_nullable_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformNullableVariablesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformNullableVariablesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformNullableVariablesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether variables have nullable field
func (r *TerraformNullableVariablesRule) Check(runner tflint.Runner) error {
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
					Attributes: []hclext.AttributeSchema{{Name: "nullable"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		if _, exists := variable.Body.Attributes["nullable"]; !exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("`%v` variable has no nullable field", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
