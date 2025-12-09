package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformEmptyListEqualityRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "comparing with [] is not recommended",
			Content: `
resource "aws_db_instance" "mysql" {
  count          = [] == [] ? 0 : 1
  instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 20},
						End:      hcl.Pos{Line: 3, Column: 28},
					},
				},
			},
			Fixed: `
resource "aws_db_instance" "mysql" {
  count          = length([]) == 0 ? 0 : 1
  instance_class = "m4.2xlarge"
}`,
		},
		{
			Name: "multiple comparisons with [] are not recommended",
			Content: `
resource "aws_db_instance" "mysql" {
  count          = [] == [] || [] == [] ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 20},
						End:      hcl.Pos{Line: 3, Column: 28},
					},
				},
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 32},
						End:      hcl.Pos{Line: 3, Column: 40},
					},
				},
			},
			Fixed: `
resource "aws_db_instance" "mysql" {
  count          = length([]) == 0 || length([]) == 0 ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
		},
		{
			Name: "comparing with [] inside parenthesis is not recommended",
			Content: `
resource "aws_db_instance" "mysql" {
  count          = ([] == []) ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 21},
						End:      hcl.Pos{Line: 3, Column: 29},
					},
				},
			},
			Fixed: `
resource "aws_db_instance" "mysql" {
  count          = (length([]) == 0) ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
		},
		{
			Name: "negatively comparing with [] is not recommended",
			Content: `
variable "my_list" {
  type = list(string)
}
resource "aws_db_instance" "mysql" {
  count          = var.my_list != [] ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 20},
						End:      hcl.Pos{Line: 6, Column: 37},
					},
				},
			},
			Fixed: `
variable "my_list" {
  type = list(string)
}
resource "aws_db_instance" "mysql" {
  count          = length(var.my_list) != 0 ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
		},
		{
			Name: "checking if length is 0 is recommended",
			Content: `
variable "my_list" {
  type = list(string)
}
resource "aws_db_instance" "mysql" {
  count          = length(var.my_list) == 0 ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple comparisons with [] are not recommended",
			Content: `
resource "aws_db_instance" "mysql" {
  count          = [] == ([] == [] ? [] : []) ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 20},
						End:      hcl.Pos{Line: 3, Column: 46},
					},
				},
				{
					Rule:    NewTerraformEmptyListEqualityRule(),
					Message: "Comparing a collection with an empty list is invalid. To detect an empty collection, check its length.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 27},
						End:      hcl.Pos{Line: 3, Column: 35},
					},
				},
			},
			Fixed: `
resource "aws_db_instance" "mysql" {
  count          = length((length([]) == 0 ? [] : [])) == 0 ? 1 : 0
  instance_class = "m4.2xlarge"
}`,
		},
	}

	rule := NewTerraformEmptyListEqualityRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
			want := map[string]string{}
			if tc.Fixed != "" {
				want["resource.tf"] = tc.Fixed
			}
			helper.AssertChanges(t, want, runner.Changes())
		})
	}
}
