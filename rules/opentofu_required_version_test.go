package rules

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformRequiredVersionRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "empty module",
			Expected: helper.Issues{},
		},
		{
			Name: "unset",
			Content: `
terraform {}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionRule(),
					Message: "terraform \"required_version\" attribute is required",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 10,
						},
					},
				},
			},
		},
		{
			Name: "set",
			Content: `
terraform {
  required_version = "~> 0.12"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple blocks",
			Content: `
terraform {
	cloud {
		workspaces {
			name = "foo"
		}
	}
}
terraform {
  required_version = "~> 0.12"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no terraform block",
			Content: `
locals {
	foo = "bar"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionRule(),
					Message: "terraform \"required_version\" attribute is required",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 1,
						},
					},
				},
			},
		},
	}

	rule := NewTerraformRequiredVersionRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			files := map[string]string{}
			if tc.Content != "" {
				files = map[string]string{"module.tf": tc.Content}
			}
			runner := helper.TestRunner(t, files)

			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

func Test_TerraformRequiredVersionRuleMultipleFiles(t *testing.T) {
	cases := []struct {
		Name     string
		Files    []string
		Expected helper.Issues
	}{
		{
			Name:  "has terraform.tf and main.tf",
			Files: []string{"modules/foo/main.tf", "modules/foo/terraform.tf"},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionRule(),
					Message: "terraform \"required_version\" attribute is required",
					Range: hcl.Range{
						Filename: filepath.FromSlash("modules/foo/terraform.tf"),
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 1,
						},
					},
				},
			},
		},
		{
			Name:  "has main.tf",
			Files: []string{"modules/foo/outputs.tf", "modules/foo/main.tf", "modules/foo/variables.tf"},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionRule(),
					Message: "terraform \"required_version\" attribute is required",
					Range: hcl.Range{
						Filename: filepath.FromSlash("modules/foo/main.tf"),
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 1,
						},
					},
				},
			},
		},
		{
			Name:  "has neither terraform.tf or main.tf",
			Files: []string{"modules/foo/variables.tf", "modules/foo/outputs.tf"},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionRule(),
					Message: "terraform \"required_version\" attribute is required",
					Range: hcl.Range{
						Filename: filepath.FromSlash("modules/foo/terraform.tf"),
					},
				},
			},
		},
	}

	rule := NewTerraformRequiredVersionRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var files = map[string]string{}
			for _, filename := range tc.Files {
				files[filepath.FromSlash(filename)] = ""
			}
			runner := helper.TestRunner(t, files)
			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
