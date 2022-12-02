package rules

import (
	"bytes"
	"fmt"
	"html/template"
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
							Column: 11,
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
							Column: 12,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 23,
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
							Column: 16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 49,
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

func BenchmarkTerraformDeprecatedInterpolation(b *testing.B) {
	cases := []struct {
		Size    int
		Content string
	}{
		{
			Size: 10,
		},
		{
			Size: 100,
		},
		{
			Size: 1000,
		},
		// {
		// 	Size: 10000,
		// },
	}

	for _, tc := range cases {
		tc := tc

		// Generate a list of n objects as locals, where each has keys a-z
		// and values 0-n
		ct, err := template.New("test").Parse(`
				locals {
					alphas = [
						{{- range .Size }}
							{
								{{- range $i, $c := $.Alpha }}
								{{ $c }} = {{ $i }},
								{{- end }}
							},
						{{- end }}
					]
				}`)

		if err != nil {
			b.Fatalf("Error rendering content template: %v", err)
		}

		buf := bytes.Buffer{}

		err = ct.Execute(&buf, struct {
			Size  []struct{}
			Alpha []string
		}{
			Size:  make([]struct{}, tc.Size),
			Alpha: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		})

		if err != nil {
			b.Fatalf("Error executing content template: %v", err)
		}

		tc.Content = buf.String()

		b.Run(fmt.Sprintf("size=%d", tc.Size), func(b *testing.B) {
			rule := NewTerraformDeprecatedIndexRule()
			runner := helper.TestRunner(b, map[string]string{"config.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				b.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(b, helper.Issues{}, runner.Issues)
		})
	}
}
