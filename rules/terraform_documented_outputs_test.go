package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDocumentedOutputsRule(t *testing.T) {
	cases := []struct {
		Name       string
		Content    string
		Config     string
		ExtraFiles map[string]string
		Expected   helper.Issues
	}{
		{
			Name: "no description",
			Content: `
output "endpoint" {
  value = aws_alb.main.dns_name
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`endpoint` output has no description",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
			},
		},
		{
			Name: "empty description",
			Content: `
output "endpoint" {
  value = aws_alb.main.dns_name
  description = ""
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`endpoint` output has no description",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
			},
		},
		{
			Name: "with description",
			Content: `
output "endpoint" {
  value = aws_alb.main.dns_name
  description = "DNS Endpoint"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate descriptions allowed by default",
			Content: `
output "first" {
  description = "Shared description"
}

output "second" {
  description = "Shared description"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate descriptions explicitly allowed",
			Content: `
output "first" {
  description = "Shared description"
}

output "second" {
  description = "Shared description"
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = false
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate descriptions",
			Content: `
output "first" {
  description = "Shared description"
}

output "second" {
  description = "Shared description"
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`first` output description is not unique: \"Shared description\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 37},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared description\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 37},
					},
				},
			},
		},
		{
			Name: "three duplicate descriptions",
			Content: `
output "first" {
  description = "Shared"
}

output "second" {
  description = "Shared"
}

output "third" {
  description = "Shared"
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`first` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`third` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 11, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 25},
					},
				},
			},
		},
		{
			Name: "duplicate descriptions in multiple files",
			Content: `
output "first" {
  description = "Shared"
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			ExtraFiles: map[string]string{
				"other_outputs.tf": `
output "second" {
  description = "Shared"
}`,
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`first` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "other_outputs.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
		},
		{
			Name: "unique descriptions",
			Content: `
output "first" {
  description = "First description"
}

output "second" {
  description = "Second description"
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "empty and populated descriptions",
			Content: `
output "empty" {
  description = ""
}

output "populated" {
  description = "Description"
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`empty` output has no description",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 15},
					},
				},
			},
		},
		{
			Name: "two empty descriptions",
			Content: `
output "empty_one" {
  description = ""
}

output "empty_two" {
  description = ""
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`empty_one` output has no description",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`empty_two` output has no description",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 19},
					},
				},
			},
		},
	}

	rule := NewTerraformDocumentedOutputsRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			files := map[string]string{"outputs.tf": tc.Content}
			if tc.Config != "" {
				files[".tflint.hcl"] = tc.Config
			}
			for name, content := range tc.ExtraFiles {
				files[name] = content
			}

			runner := helper.TestRunner(t, files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
