package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformCommentSyntaxRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		JSON     bool
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name:     "hash comment",
			Content:  `# foo`,
			Expected: helper.Issues{},
		},
		{
			Name: "multi-line comment allowed",
			Content: `
/*
	This comment spans multiple lines
*/			
`,
			Config: `
rule "terraform_comment_syntax" {
  enabled = true

  allow_multiline = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "multi-line comment disallowed",
			Content: `
/*
	This comment spans multiple lines
*/			
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Multi-line comments are not allowed. Use single-line comments starting with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 3,
						},
					},
				},
			},
		},
		{
			Name:    "double-slash comment",
			Content: `// foo`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Single line comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 7,
						},
					},
				},
			},
			Fixed: `# foo`,
		},
		{
			Name: "end-of-line hash comment",
			Content: `
variable "foo" {
	type = string # a string
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:     "JSON",
			Content:  `{"variable": {"foo": {"type": "string"}}}`,
			JSON:     true,
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformCommentSyntaxRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content, ".tflint.hcl": tc.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
			want := map[string]string{}
			if tc.Fixed != "" {
				want[filename] = tc.Fixed
			}
			helper.AssertChanges(t, want, runner.Changes())
		})
	}
}
