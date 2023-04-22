package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDeprecatedIndexRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deprecated dot index style",
			Content: `
locals {
  list = ["a"]
  value = list.0
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedIndexRule(),
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
		},
		{
			Name: "deprecated dot splat index style",
			Content: `
locals {
  maplist = [{a = "b"}]
  values = maplist.*.a
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 19,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 21,
						},
					},
				},
			},
		},
		{
			Name: "attribute access",
			Content: `
locals {
  map = {a = "b"}
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
%{ for ip in aws_instance.example[*].private_ip }
server ${ip}
%{ endfor }
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
		%{ for ip in aws_instance.example.*.private_ip }
		server ${ip}
		%{ endfor }
		EOF
		}
		`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 36,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 38,
						},
					},
				},
			},
		},
		{
			Name: "legacy splat and legacy index",
			Content: `
locals {
  nested_list = [["a"]]
  value = nested_list.*.0
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 22,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 24,
						},
					},
				},
				{
					Rule:    NewTerraformDeprecatedIndexRule(),
					Message: "List items should be accessed using square brackets",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 24,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 26,
						},
					},
				},
			},
		},
		{
			Name: "complex expression",
			Content: `
locals {
  create_namespace = true
  kubernetes_namespace = local.create_namespace ? join("", kubernetes_namespace.default.*.id) : var.kubernetes_namespace
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedIndexRule(),
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
		},
	}

	rule := NewTerraformDeprecatedIndexRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
