package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformNoShortCircuitEvaluationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "short circuit with null check and &&",
			Content: `
resource "aws_instance" "example" {
  count = var.obj != null && var.obj.enabled ? 1 : 0
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformNoShortCircuitEvaluationRule(),
					Message: "Short-circuit evaluation is not supported in Terraform. Use a conditional expression (condition ? true : false) instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 45},
					},
				},
			},
		},
		{
			Name: "short circuit with null check and ||",
			Content: `
resource "aws_instance" "example" {
  count = var.obj == null || var.obj.enabled ? 1 : 0
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformNoShortCircuitEvaluationRule(),
					Message: "Short-circuit evaluation is not supported in Terraform. Use a conditional expression (condition ? true : false) instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 45},
					},
				},
			},
		},
		{
			Name: "correct conditional usage",
			Content: `
resource "aws_instance" "example" {
  count = var.obj == null ? 0 : (var.obj.enabled ? 1 : 0)
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid use of logical operators with independent values",
			Content: `
resource "aws_instance" "example" {
  count = var.value > 3 || var.other_value < 10 ? 1 : 0
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid use of logical operators with same object",
			Content: `
resource "aws_instance" "example" {
  count = var.obj > 3 || var.obj < 10 ? 1 : 0
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple null checks in one expression",
			Content: `
resource "aws_instance" "example" {
  count = var.obj != null && var.obj.enabled && var.obj.property > 0 ? 1 : 0
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformNoShortCircuitEvaluationRule(),
					Message: "Short-circuit evaluation is not supported in Terraform. Use a conditional expression (condition ? true : false) instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 45},
					},
				},
			},
		},
	}

	rule := NewTerraformNoShortCircuitEvaluationRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
