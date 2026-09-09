package rules

import (
	"maps"
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
  value       = aws_alb.first.dns_name
  description = "Shared"
}

output "second" {
  value       = aws_alb.second.dns_name
  description = "Shared"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate descriptions",
			Content: `
output "first" {
  value       = aws_alb.first.dns_name
  description = "Shared"
}

output "second" {
  value       = aws_alb.second.dns_name
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
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 9, Column: 3},
						End:      hcl.Pos{Line: 9, Column: 25},
					},
				},
			},
		},
		{
			Name: "three duplicate descriptions",
			Content: `
output "first" {
  value       = aws_alb.first.dns_name
  description = "Shared"
}

output "second" {
  value       = aws_alb.second.dns_name
  description = "Shared"
}

output "third" {
  value       = aws_alb.third.dns_name
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
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 9, Column: 3},
						End:      hcl.Pos{Line: 9, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`third` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 14, Column: 3},
						End:      hcl.Pos{Line: 14, Column: 25},
					},
				},
			},
		},
		{
			Name: "duplicate descriptions in multiple files",
			Content: `
output "first" {
  value       = aws_alb.first.dns_name
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
  value       = aws_alb.second.dns_name
  description = "Shared"
}`,
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`first` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "other_outputs.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 25},
					},
				},
			},
		},
		{
			Name: "unique descriptions",
			Content: `
output "first" {
  value       = aws_alb.first.dns_name
  description = "First description"
}

output "second" {
  value       = aws_alb.second.dns_name
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
			Name: "descriptions differing in case or whitespace",
			Content: `
output "first" {
  value       = aws_alb.first.dns_name
  description = "Shared"
}

output "second" {
  value       = aws_alb.second.dns_name
  description = "shared"
}

output "third" {
  value       = aws_alb.third.dns_name
  description = "Shared "
}`,
			Config: `
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "undocumented output alongside duplicate descriptions",
			Content: `
output "undocumented" {
  value = aws_alb.undocumented.dns_name
}

output "first" {
  value       = aws_alb.first.dns_name
  description = "Shared"
}

output "second" {
  value       = aws_alb.second.dns_name
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
					Message: "`undocumented` output has no description",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 22},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`first` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 8, Column: 3},
						End:      hcl.Pos{Line: 8, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedOutputsRule(),
					Message: "`second` output description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 13, Column: 3},
						End:      hcl.Pos{Line: 13, Column: 25},
					},
				},
			},
		},
		{
			Name: "empty descriptions are not duplicates",
			Content: `
output "empty_one" {
  value       = aws_alb.first.dns_name
  description = ""
}

output "empty_two" {
  value       = aws_alb.second.dns_name
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
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 19},
					},
				},
			},
		},
	}

	rule := NewTerraformDocumentedOutputsRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			files := map[string]string{"outputs.tf": tc.Content, ".tflint.hcl": tc.Config}
			maps.Copy(files, tc.ExtraFiles)

			runner := helper.TestRunner(t, files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
