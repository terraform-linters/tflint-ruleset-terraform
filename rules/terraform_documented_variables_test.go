package rules

import (
	"maps"
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDocumentedVariablesRule(t *testing.T) {
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
variable "no_description" {
  default = "default"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`no_description` variable has no description",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 26},
					},
				},
			},
		},
		{
			Name: "empty description",
			Content: `
variable "empty_description" {
  description = ""
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`empty_description` variable has no description",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 29},
					},
				},
			},
		},
		{
			Name: "with description",
			Content: `
variable "with_description" {
  description = "This is description"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate descriptions allowed by default",
			Content: `
variable "first" {
  description = "Shared"
}

variable "second" {
  description = "Shared"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate descriptions",
			Content: `
variable "first" {
  description = "Shared"
}

variable "second" {
  description = "Shared"
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`first` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`second` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 25},
					},
				},
			},
		},
		{
			Name: "three duplicate descriptions",
			Content: `
variable "first" {
  description = "Shared"
}

variable "second" {
  description = "Shared"
}

variable "third" {
  description = "Shared"
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`first` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`second` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`third` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 11, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 25},
					},
				},
			},
		},
		{
			Name: "duplicate descriptions in multiple files",
			Content: `
variable "first" {
  description = "Shared"
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			ExtraFiles: map[string]string{
				"other_variables.tf": `
variable "second" {
  description = "Shared"
}`,
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`first` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`second` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "other_variables.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
		},
		{
			Name: "unique descriptions",
			Content: `
variable "first" {
  description = "First description"
}

variable "second" {
  description = "Second description"
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "descriptions differing in case or whitespace",
			Content: `
variable "first" {
  description = "Shared"
}

variable "second" {
  description = "shared"
}

variable "third" {
  description = "Shared "
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "undocumented variable alongside duplicate descriptions",
			Content: `
variable "undocumented" {
}

variable "first" {
  description = "Shared"
}

variable "second" {
  description = "Shared"
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`undocumented` variable has no description",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 24},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`first` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 25},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`second` variable description is not unique: \"Shared\"",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 10, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 25},
					},
				},
			},
		},
		{
			Name: "empty descriptions are not duplicates",
			Content: `
variable "empty_one" {
  description = ""
}

variable "empty_two" {
  description = ""
}`,
			Config: `
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`empty_one` variable has no description",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 21},
					},
				},
				{
					Rule:    NewTerraformDocumentedVariablesRule(),
					Message: "`empty_two` variable has no description",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 21},
					},
				},
			},
		},
	}

	rule := NewTerraformDocumentedVariablesRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			files := map[string]string{"variables.tf": tc.Content, ".tflint.hcl": tc.Config}
			maps.Copy(files, tc.ExtraFiles)

			runner := helper.TestRunner(t, files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
