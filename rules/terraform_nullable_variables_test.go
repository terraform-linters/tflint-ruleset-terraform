package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformNullableVariablesRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "no nullable",
			Content: `
variable "no_nullable" {
  default = "default"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformNullableVariablesRule(),
					Message: "`no_nullable` variable has no nullable field",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 23},
					},
				},
			},
		},
		{
			Name: "nullable true",
			Content: `
variable "nullable" {
  nullable = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "nullable false",
			Content: `
variable "not_nullable" {
  nullable = false
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformNullableVariablesRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
