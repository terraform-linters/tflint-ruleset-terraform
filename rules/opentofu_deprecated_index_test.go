package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_OpentofuDeprecatedIndexRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "deprecated dot index style",
			Content: `
locals {
  list  = ["a"]
  value = list.0
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 15,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 17,
						},
					},
				},
			},
			Fixed: `
locals {
  list  = ["a"]
  value = list[0]
}
`,
		},
		{
			Name: "deprecated dot splat index style",
			Content: `
locals {
  maplist = [{ a = "b" }]
  values  = maplist.*.a
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 20,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 22,
						},
					},
				},
			},
			Fixed: `
locals {
  maplist = [{ a = "b" }]
  values  = maplist[*].a
}
`,
		},
		{
			Name: "attribute access",
			Content: `
locals {
  map   = { a = "b" }
  value = map.a
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "fractional number",
			Content: `
locals {
  value = 1.5
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "directive: valid",
			Content: `
locals {
  servers = <<EOF
%{for ip in aws_instance.example[*].private_ip}
server ${ip}
%{endfor}
EOF
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "directive: invalid",
			Content: `
locals {
  servers = <<EOF
%{for ip in aws_instance.example.*.private_ip}
server ${ip}
%{endfor}
EOF
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 33,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 35,
						},
					},
				},
			},
			Fixed: `
locals {
  servers = <<EOF
%{for ip in aws_instance.example[*].private_ip}
server ${ip}
%{endfor}
EOF
}
`,
		},
		{
			Name: "legacy splat and legacy index",
			Content: `
locals {
  nested_list = [["a"]]
  value       = nested_list.*.0
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 28,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 30,
						},
					},
				},
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 30,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 32,
						},
					},
				},
			},
			Fixed: `
locals {
  nested_list = [["a"]]
  value       = nested_list[*][0]
}
`,
		},
		{
			Name: "complex expression",
			Content: `
locals {
  create_namespace     = true
  kubernetes_namespace = local.create_namespace ? join("", kubernetes_namespace.default.*.id) : var.kubernetes_namespace
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 88,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 90,
						},
					},
				},
			},
			Fixed: `
locals {
  create_namespace     = true
  kubernetes_namespace = local.create_namespace ? join("", kubernetes_namespace.default[*].id) : var.kubernetes_namespace
}
`,
		},
		{
			Name: "json invalid",
			JSON: true,
			Content: `
			{
				"locals": {
					"list": ["a"],
					"value": "${list.0}"
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf.json",
						Start: hcl.Pos{
							Line:   5,
							Column: 27,
						},
						End: hcl.Pos{
							Line:   5,
							Column: 29,
						},
					},
				},
			},
			Fixed: `
			{
				"locals": {
					"list": ["a"],
					"value": "${list[0]}"
				}
			}`,
		},
		{
			Name: "json valid",
			JSON: true,
			Content: `
			{
				"locals": {
					"list": ["a"],
					"value": "${list[0]}"
				}
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "json strings",
			JSON: true,
			Content: `
			{
				"locals": {
					"string": "foo",
					"bool": "${local.string == \"foo\"}"
				}
			}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewOpentofuDeprecatedIndexRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
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
