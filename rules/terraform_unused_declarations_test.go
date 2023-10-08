package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformUnusedDeclarationsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "unused variable",
			Content: `
variable "not_used" {}
variable "used" {}
output "u" { value = var.used }
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformUnusedDeclarationsRule(),
					Message: `variable "not_used" is declared but not used`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 20},
					},
				},
			},
			Fixed: `
variable "used" {}
output "u" { value = var.used }
`,
		},
		{
			Name: "unused data source",
			Content: `
data "null_data_source" "not_used" {}
data "null_data_source" "used" {}
output "u" { value = data.null_data_source.used }
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformUnusedDeclarationsRule(),
					Message: `data "null_data_source" "not_used" is declared but not used`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 35},
					},
				},
			},
			Fixed: `
data "null_data_source" "used" {}
output "u" { value = data.null_data_source.used }
`,
		},
		{
			Name: "unused local source",
			Content: `
locals {
  not_used = ""
  used = ""
}
output "u" { value = local.used }
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformUnusedDeclarationsRule(),
					Message: `local.not_used is declared but not used`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 16},
					},
				},
			},
			Fixed: `
locals {
  used = ""
}
output "u" { value = local.used }
`,
		},
		{
			Name: "variable used in resource",
			Content: `
variable "used" {}
resource "null_resource" "n" {
  triggers = {
    u = var.used
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable used in module",
			Content: `
variable "used" {}
module "m" {
  source = "./module"
  u = var.used
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable used in module",
			Content: `
variable "used" {}
module "m" {
  source = "./module"
  u = var.used
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "local used in module",
			Content: `
locals { used = "used" }
module "m" {
  source = "./module"
  u = local.used
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable used in provider",
			Content: `
variable "aws_region" {}
provider "aws" {
  region = var.aws_region
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "meta-arguments",
			Content: `
variable "used" {}
resource "null_resource" "n" {
  triggers = {
    u = var.used
  }

  lifecycle {
    ignore_changes = [triggers]
  }
  providers = {
    null = null
  }
  depends_on = [aws_instance.foo]
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "additional traversal",
			Content: `
variable "v" {
  type = object({ foo = string })
}
output "v" {
  value = var.v.foo
}
data "terraform_remote_state" "d" {}
output "d" {
  value = data.terraform_remote_state.d.outputs.foo
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable used in validation block",
			Content: `
variable "unused" {
  validation {
    condition     = var.unused != ""
    error_message = "variable should be empty string. got: ${var.unused}"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformUnusedDeclarationsRule(),
					Message: `variable "unused" is declared but not used`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
			},
			Fixed: `
`,
		},
		{
			Name: "unused scoped data source",
			Content: `
check "unused" {
  data "null_data_source" "unused" {}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformUnusedDeclarationsRule(),
					Message: `data "null_data_source" "unused" is declared but not used`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 35},
					},
				},
			},
			Fixed: `
check "unused" {
}
`,
		},
		{
			Name: "json",
			JSON: true,
			Content: `
{
  "resource": {
    "foo": {
      "bar": {
        "nested": [{
          "${var.again}": []
        }]
      }
    }
  },
  "variable": {
    "again": {}
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "json with unused variable",
			JSON: true,
			Content: `
{
  "variable": {
    "again": {}
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformUnusedDeclarationsRule(),
					Message: `variable "again" is declared but not used`,
					Range: hcl.Range{
						Filename: "config.tf.json",
						Start:    hcl.Pos{Line: 4, Column: 14},
						End:      hcl.Pos{Line: 4, Column: 15},
					},
				},
			},
		},
	}

	rule := NewTerraformUnusedDeclarationsRule()

	for _, tc := range cases {
		filename := "config.tf"
		if tc.JSON {
			filename += ".json"
		}

		t.Run(tc.Name, func(t *testing.T) {
			runner := testRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helperRunner := runner.Runner.(*helper.Runner)

			helper.AssertIssues(t, tc.Expected, helperRunner.Issues)
			want := map[string]string{}
			if tc.Fixed != "" {
				want[filename] = tc.Fixed
			}
			helper.AssertChanges(t, want, helperRunner.Changes())
		})
	}
}
