package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDynamicAttributeNotationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "dot notation in for_each",
			Content: `
resource "aws_instance" "web" {
  for_each = local.instances
  subnet_id = each.value.subnet_id
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformDynamicAttributeNotationRule(),
					Message: "Must use bracket notation [] for dynamic attributes",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 4, Column: 15},
						End:      hcl.Pos{Line: 4, Column: 35},
					},
				},
			},
			Fixed: `
resource "aws_instance" "web" {
  for_each  = local.instances
  subnet_id = each.value["subnet_id"]
}`,
		},
		{
			Name: "valid bracket notation in for_each",
			Content: `
resource "aws_instance" "web" {
  for_each = local.instances
  subnet_id = each.value["subnet_id"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "dot notation in for expression",
			Content: `
locals {
  ips = [for instance in aws_instance.web : instance.private_ip]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDynamicAttributeNotationRule(),
					Message: "Must use bracket notation [] for dynamic attributes",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 3, Column: 50},
						End:      hcl.Pos{Line: 3, Column: 70},
					},
				},
			},
			Fixed: `
locals {
  ips = [for instance in aws_instance.web : instance["private_ip"]]
}`,
		},
		{
			Name: "dot notation in count",
			Content: `
resource "aws_instance" "web" {
  count = 2
  subnet_id = aws_subnet.main[count.index].id
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformDynamicAttributeNotationRule(),
					Message: "Must use bracket notation [] for dynamic attributes",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 4, Column: 31},
						End:      hcl.Pos{Line: 4, Column: 42},
					},
				},
			},
			Fixed: `
resource "aws_instance" "web" {
  count = 2
  subnet_id = aws_subnet.main[count.index]["id"]
}`,
		},
	}

	rule := NewTerraformDynamicAttributeNotationRule()

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
