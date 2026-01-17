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
			Name:    "double-slash comment",
			Content: `// foo`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
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
			// Single-line /* */ comments can appear mid-expression (C-style),
			// e.g. `x = 1 /* comment */ + 2` evaluates to 3. Autofixing to #
			// would comment out the rest of the line, changing behavior.
			Name:    "single-line block comment",
			Content: `x = 1 /* comment */ + 2`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   1,
							Column: 7,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 20,
						},
					},
				},
			},
		},
		{
			Name: "multi-line comment",
			Content: `
/*
	This comment spans multiple lines
*/
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
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
			Fixed: `
#
#	This comment spans multiple lines
#
`,
		},
		{
			Name:    "double-slash comment without space",
			Content: `//foo`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 6,
						},
					},
				},
			},
			Fixed: `#foo`,
		},
		{
			Name: "end-of-line double-slash comment",
			Content: `
variable "foo" {
  type = string // a string
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   3,
							Column: 17,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 1,
						},
					},
				},
			},
			Fixed: `
variable "foo" {
  type = string # a string
}
`,
		},
		{
			Name: "multi-line comment with asterisk prefix",
			Content: `
/*
 * This is a comment
 * with asterisks
*/
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   5,
							Column: 3,
						},
					},
				},
			},
			Fixed: `
#
# * This is a comment
# * with asterisks
#
`,
		},
		{
			Name: "multiple comments",
			Content: `// first
// second`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 1,
						},
					},
				},
				{
					Rule:    NewTerraformCommentSyntaxRule(),
					Message: "Comments should begin with #",
					Range: hcl.Range{
						Filename: "variables.tf",
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
			Fixed: `# first
# second`,
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

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

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
