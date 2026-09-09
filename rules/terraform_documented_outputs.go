package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformDocumentedOutputsRule checks whether outputs have descriptions
type TerraformDocumentedOutputsRule struct {
	documentedRule
}

// NewTerraformDocumentedOutputsRule returns a new rule
func NewTerraformDocumentedOutputsRule() *TerraformDocumentedOutputsRule {
	return &TerraformDocumentedOutputsRule{
		documentedRule: documentedRule{
			name:      "terraform_documented_outputs",
			blockType: "output",
		},
	}
}

// Check checks whether outputs have descriptions
func (r *TerraformDocumentedOutputsRule) Check(runner tflint.Runner) error {
	return r.check(runner, r)
}
