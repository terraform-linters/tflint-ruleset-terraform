package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDeprecatedInterpolationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "deprecated single interpolation",
			Content: `
resource "null_resource" "a" {
  triggers = "${var.triggers}"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 31},
					},
				},
			},
			Fixed: `
resource "null_resource" "a" {
  triggers = var.triggers
}`,
		},
		{
			Name: "deprecated single interpolation in provider block",
			Content: `
provider "null" {
  foo = "${var.triggers["foo"]}"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 9},
						End:      hcl.Pos{Line: 3, Column: 33},
					},
				},
			},
			Fixed: `
provider "null" {
  foo = var.triggers["foo"]
}`,
		},
		{
			Name: "deprecated single interpolation in locals block",
			Content: `
locals {
  foo = "${var.triggers["foo"]}"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 9},
						End:      hcl.Pos{Line: 3, Column: 33},
					},
				},
			},
			Fixed: `
locals {
  foo = var.triggers["foo"]
}`,
		},
		{
			Name: "deprecated single interpolation in nested block",
			Content: `
resource "null_resource" "a" {
  provisioner "local-exec" {
    single = "${var.triggers["greeting"]}"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 4, Column: 14},
						End:      hcl.Pos{Line: 4, Column: 43},
					},
				},
			},
			Fixed: `
resource "null_resource" "a" {
  provisioner "local-exec" {
    single = var.triggers["greeting"]
  }
}`,
		},
		{
			Name: "interpolation as template",
			Content: `
resource "null_resource" "a" {
  triggers = "${var.triggers} "
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "interpolation in array",
			Content: `
resource "null_resource" "a" {
  triggers = ["${var.triggers}"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 15},
						End:      hcl.Pos{Line: 3, Column: 32},
					},
				},
			},
			Fixed: `
resource "null_resource" "a" {
  triggers = [var.triggers]
}`,
		},
		{
			Name: "new interpolation syntax",
			Content: `
resource "null_resource" "a" {
  triggers = var.triggers
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "nested wraps",
			Content: `
resource "null_resource" "a" {
  triggers = "${"${var.triggers}"}"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 36},
					},
				},
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 17},
						End:      hcl.Pos{Line: 3, Column: 34},
					},
				},
			},
			Fixed: `
resource "null_resource" "a" {
  triggers = var.triggers
}`,
		},
		{
			Name: "interpolation as an object key",
			Content: `
resource "null_resource" "a" {
  triggers = {
    "${var.triggers}" = "foo"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 22},
					},
				},
			},
			Fixed: `
resource "null_resource" "a" {
  triggers = {
    (var.triggers) = "foo"
  }
}`,
		},
		{
			Name: "interpolation in an object key",
			Content: `
resource "null_resource" "a" {
  triggers = {
    upper("${var.triggers}") = "foo"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDeprecatedInterpolationRule(),
					Message: "Interpolation-only expressions are deprecated in Terraform v0.12.14",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 4, Column: 11},
						End:      hcl.Pos{Line: 4, Column: 28},
					},
				},
			},
			Fixed: `
resource "null_resource" "a" {
  triggers = {
    upper(var.triggers) = "foo"
  }
}`,
		},
	}

	rule := NewTerraformDeprecatedInterpolationRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
			want := map[string]string{}
			if tc.Fixed != "" {
				want["config.tf"] = tc.Fixed
			}
			helper.AssertChanges(t, want, runner.Changes())
		})
	}
}
