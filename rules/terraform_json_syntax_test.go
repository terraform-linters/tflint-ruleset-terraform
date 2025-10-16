package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformJSONSyntaxRule(t *testing.T) {
	for _, tc := range []struct {
		name     string
		content  string
		filename string
		expected helper.Issues
	}{
		{
			name:     "object syntax valid",
			content:  `{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}`,
			filename: "main.tf.json",
			expected: helper.Issues{},
		},
		{
			name:     "empty object valid",
			content:  `{}`,
			filename: "main.tf.json",
			expected: helper.Issues{},
		},
		{
			name:     "array syntax invalid",
			content:  `[{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}]`,
			filename: "main.tf.json",
			expected: helper.Issues{
				{
					Rule:    NewTerraformJSONSyntaxRule(),
					Message: "JSON configuration uses array syntax at root. The official Terraform JSON syntax uses a root object. See https://developer.hashicorp.com/terraform/language/syntax/json",
					Range: hcl.Range{
						Filename: "main.tf.json",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 2},
					},
				},
			},
		},
		{
			name:     "regular HCL file ignored",
			content:  `resource "aws_instance" "example" { ami = "ami-12345678" }`,
			filename: "main.tf",
			expected: helper.Issues{},
		},
		{
			name: "complex object syntax valid",
			content: `{
  "terraform": {
    "required_version": ">= 1.0"
  },
  "resource": {
    "aws_instance": {
      "example": {
        "ami": "ami-12345678"
      }
    }
  }
}`,
			filename: "main.tf.json",
			expected: helper.Issues{},
		},
		{
			name: "array with multiple objects invalid",
			content: `[
  {"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}},
  {"variable": {"region": {"type": "string"}}}
]`,
			filename: "config.tf.json",
			expected: helper.Issues{
				{
					Rule:    NewTerraformJSONSyntaxRule(),
					Message: "JSON configuration uses array syntax at root. The official Terraform JSON syntax uses a root object. See https://developer.hashicorp.com/terraform/language/syntax/json",
					Range: hcl.Range{
						Filename: "config.tf.json",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 2},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewTerraformJSONSyntaxRule()
			runner := helper.TestRunner(t, map[string]string{tc.filename: tc.content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.expected, runner.Issues)
		})
	}
}
