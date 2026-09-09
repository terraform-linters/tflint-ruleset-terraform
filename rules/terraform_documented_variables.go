package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformDocumentedVariablesRule checks whether variables have descriptions
type TerraformDocumentedVariablesRule struct {
	documentedRule
}

// NewTerraformDocumentedVariablesRule returns a new rule
func NewTerraformDocumentedVariablesRule() *TerraformDocumentedVariablesRule {
	return &TerraformDocumentedVariablesRule{
		documentedRule: documentedRule{
			name:      "terraform_documented_variables",
			blockType: "variable",
		},
	}
}

// Check checks whether variables have descriptions
func (r *TerraformDocumentedVariablesRule) Check(runner tflint.Runner) error {
	return r.check(runner, r)
}
