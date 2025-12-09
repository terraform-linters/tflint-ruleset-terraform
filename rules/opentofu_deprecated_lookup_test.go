package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_OpentofuDeprecatedLookupRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "deprecated lookup",
			Content: `
locals {
  map   = { a = 0 }
  value = lookup(local.map, "a")
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedLookupRule(),
					Message: "Lookup with 2 arguments is deprecated",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 11,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 33,
						},
					},
				},
			},
			Fixed: `
locals {
  map   = { a = 0 }
  value = local.map["a"]
}
`,
		},
		{
			Name: "deprecated lookup nested",
			Content: `
locals {
  map   = { a = { b = 0 } }
  value = lookup(lookup(local.map, "a"), "b")
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedLookupRule(),
					Message: "Lookup with 2 arguments is deprecated",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 11,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 46,
						},
					},
				},
				{
					Rule:    NewOpentofuDeprecatedLookupRule(),
					Message: "Lookup with 2 arguments is deprecated",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 18,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 40,
						},
					},
				},
			},
			Fixed: `
locals {
  map   = { a = { b = 0 } }
  value = local.map["a"]["b"]
}
`,
		},
	}

	rule := NewOpentofuDeprecatedLookupRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"

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
