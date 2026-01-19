package rules

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformPrivateModuleReferenceRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid private module reference",
			Content: `
module "foo" {
  source = "modules/foo"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "peer private module reference",
			Content: `
module "foo" {
  source = "../foo"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformPrivateModuleReferenceRule(),
					Message: "Private modules should not be referenced externally. Add a README.md to make the referenced module public or remove the reference.",
					Range: hcl.Range{
						Filename: "module.tf",
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
			Name: "valid private module reference without correct modules subdir",
			Content: `
module "foo" {
  source = "./foo"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "external private module reference",
			Content: `
module "bar" {
  source = "../another-root/modules/bar"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformPrivateModuleReferenceRule(),
					Message: "Private modules should not be referenced externally. Add a README.md to make the referenced module public or remove the reference.",
					Range: hcl.Range{
						Filename: "module.tf",
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
			Name: "valid public submodule reference",
			Content: `
module "baz" {
  source = "../another-root/modules/baz"
}
`,
			Expected: helper.Issues{},
		},
	}

	mockStat := func(name string) (os.FileInfo, error) {
		switch name {
		case "../another-root/modules/baz/README.md", "../foo", "./foo", "modules/foo", "modules/foo/bar", "../another-root/modules/bar":
			return nil, nil // File exists
		default:
			return nil, os.ErrNotExist // File doesn't exist
		}
	}

	rule := NewTerraformPrivateModuleReferenceRule()
	rule.statFunc = mockStat

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			files := map[string]string{}
			if tc.Content != "" {
				files = map[string]string{"module.tf": tc.Content}
			}
			runner := testRunner(t, files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
		})
	}
}
