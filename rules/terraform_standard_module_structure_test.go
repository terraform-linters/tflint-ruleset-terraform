package rules

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformStandardModuleStructureRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  map[string]string
		Config   string
		Expected helper.Issues
	}{
		{
			Name:     "empty module",
			Content:  map[string]string{},
			Expected: helper.Issues{},
		},
		{
			Name: "non-standard module",
			Content: map[string]string{
				"foo.tf": "",
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformStandardModuleStructureRule(),
					Message: "Module should include a main.tf file as the primary entrypoint",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.InitialPos,
					},
				},
				{
					Rule:    NewTerraformStandardModuleStructureRule(),
					Message: "Module should include an empty variables.tf file",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.InitialPos,
					},
				},
				{
					Rule:    NewTerraformStandardModuleStructureRule(),
					Message: "Module should include an empty outputs.tf file",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.InitialPos,
					},
				},
			},
		},
		{
			Name: "directory in path",
			Content: map[string]string{
				"foo/main.tf": "",
				"foo/variables.tf": `
variable "v" {}
				`,
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformStandardModuleStructureRule(),
					Message: "Module should include an empty outputs.tf file",
					Range: hcl.Range{
						Filename: filepath.Join("foo", "outputs.tf"),
						Start:    hcl.InitialPos,
					},
				},
			},
		},
		{
			Name: "allow no outputs.tf",
			Content: map[string]string{
				"foo/main.tf": "",
				"foo/variables.tf": `
variable "v" {}
				`,
				".tflint.hcl": `
rule "terraform_standard_module_structure" {
  enabled = true

  outputs = false
}
`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "move variable",
			Content: map[string]string{
				"main.tf": `
variable "v" {}
`,
				"variables.tf": "",
				"outputs.tf":   "",
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformStandardModuleStructureRule(),
					Message: `variable "v" should be moved from main.tf to variables.tf`,
					Range: hcl.Range{
						Filename: "main.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 13,
						},
					},
				},
			},
		},
		{
			Name: "move output",
			Content: map[string]string{
				"main.tf": `
output "o" { value = null }
`,
				"variables.tf": "",
				"outputs.tf":   "",
			},
			Expected: helper.Issues{
				{
					Rule:    NewTerraformStandardModuleStructureRule(),
					Message: `output "o" should be moved from main.tf to outputs.tf`,
					Range: hcl.Range{
						Filename: "main.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 11,
						},
					},
				},
			},
		},
		{
			Name: "json only",
			Content: map[string]string{
				"main.tf.json": "{}",
			},
			Expected: helper.Issues{},
		},
		{
			Name: "json variable",
			Content: map[string]string{
				"main.tf.json": `{"variable": {"v": {}}}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "json output",
			Content: map[string]string{
				"main.tf.json": `{"output": {"o": {"value": null}}}`,
			},
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformStandardModuleStructureRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.Content)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
