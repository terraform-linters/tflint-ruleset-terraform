package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformDocumentedVariablesRule checks whether variables have descriptions
type TerraformDocumentedVariablesRule struct {
	tflint.DefaultRule
}

// TerraformDocumentedVariablesRuleConfig is the config structure for the TerraformDocumentedVariablesRule
type TerraformDocumentedVariablesRuleConfig struct {
	Unique bool `hclext:"unique,optional"`
}

// NewTerraformDocumentedVariablesRule returns a new rule
func NewTerraformDocumentedVariablesRule() *TerraformDocumentedVariablesRule {
	return &TerraformDocumentedVariablesRule{}
}

// Name returns the rule name
func (r *TerraformDocumentedVariablesRule) Name() string {
	return "terraform_documented_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformDocumentedVariablesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformDocumentedVariablesRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformDocumentedVariablesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether variables have descriptions
func (r *TerraformDocumentedVariablesRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	config := TerraformDocumentedVariablesRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: "description"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	return checkDocumentedDescriptions(runner, r, body.Blocks, "variable", config.Unique)
}
