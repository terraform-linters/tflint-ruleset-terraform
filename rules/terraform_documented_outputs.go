package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformDocumentedOutputsRule checks whether outputs have descriptions
type TerraformDocumentedOutputsRule struct {
	tflint.DefaultRule
}

// TerraformDocumentedOutputsRuleConfig is the config structure for the TerraformDocumentedOutputsRule
type TerraformDocumentedOutputsRuleConfig struct {
	Unique bool `hclext:"unique,optional"`
}

// NewTerraformDocumentedOutputsRule returns a new rule
func NewTerraformDocumentedOutputsRule() *TerraformDocumentedOutputsRule {
	return &TerraformDocumentedOutputsRule{}
}

// Name returns the rule name
func (r *TerraformDocumentedOutputsRule) Name() string {
	return "terraform_documented_outputs"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformDocumentedOutputsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformDocumentedOutputsRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformDocumentedOutputsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether outputs have descriptions
func (r *TerraformDocumentedOutputsRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	config := TerraformDocumentedOutputsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "output",
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

	return checkDocumentedDescriptions(runner, r, body.Blocks, "output", config.Unique)
}
