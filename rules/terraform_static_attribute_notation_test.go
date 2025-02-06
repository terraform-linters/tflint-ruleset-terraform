package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformStaticAttributeNotationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "bracket notation in static context",
			Content: `resource "aws_instance" "web" {
  instance_type = var.instance["type"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformStaticAttributeNotationRule(),
					Message: "Must use dot notation for static attributes",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 2, Column: 19},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
			Fixed: `resource "aws_instance" "web" {
  instance_type = var.instance.type
}`,
		},
		{
			Name: "dot notation in static context (valid)",
			Content: `resource "aws_instance" "web" {
  instance_type = var.instance.type
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "bracket notation with variable index (valid)",
			Content: `resource "aws_instance" "web" {
  instance_type = var.instance[var.env]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "bracket notation with number index (valid)",
			Content: `resource "aws_instance" "web" {
  instance_type = var.instance[0]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple static attributes with bracket notation",
			Content: `resource "aws_instance" "web" {
  instance_type = var.instance["type"]["size"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformStaticAttributeNotationRule(),
					Message: "Must use dot notation for static attributes",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 2, Column: 19},
						End:      hcl.Pos{Line: 2, Column: 47},
					},
				},
			},
			Fixed: `resource "aws_instance" "web" {
  instance_type = var.instance.type.size
}`,
		},
	}

	rule := NewTerraformStaticAttributeNotationRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
			if tc.Fixed != "" {
				helper.AssertChanges(t, map[string]string{"config.tf": tc.Fixed}, runner.Changes())
			}
		})
	}
}
