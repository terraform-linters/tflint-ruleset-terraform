package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformMapDuplicateValues(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "No duplicates",
			Content: `
resource "null_resource" "test" {
	test = {
	  a = 1
	  b = 2
	  c = 3
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "duplicate values in map literal",
			Content: `
resource "null_resource" "test" {
    triggers = {
        a = "b"
        c = "b"
    }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMapDuplicateValuesRule(),
					Message: `Duplicate value: "b", first defined at module.tf:4,13-16`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 13},
						End:      hcl.Pos{Line: 5, Column: 16},
					},
				},
			},
		},
		{
			Name: "duplicate values with quoting",
			Content: `
resource "null_resource" "test" {
    triggers = {
        a = "b"
        c = "b"
    }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMapDuplicateValuesRule(),
					Message: `Duplicate value: "b", first defined at module.tf:4,13-16`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 13},
						End:      hcl.Pos{Line: 5, Column: 16},
					},
				},
			},
		},
		{
			Name: "Using variables as values",
			Content: `
variable "a" {
  type    = string
  default = "b"
}

resource "null_resource" "test" {
	map = {
	  key1 = var.a
	  key2 = "b"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMapDuplicateValuesRule(),
					Message: `Duplicate value: "b", first defined at module.tf:9,11-16`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 10, Column: 11},
						End:      hcl.Pos{Line: 10, Column: 14},
					},
				},
			},
		},
		{
			Name: "Using a variable as a value without a default",
			Content: `
variable "unknown" {
  type    = string
}

resource "null_resource" "test" {
	map = {
	  key1 = "x"
	  key2 = var.unknown
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Multiple duplicates in same map",
			Content: `
resource "null_resource" "test" {
	map = {
	  key1 = "a"
	  key2 = "a"
	  key3 = "a"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMapDuplicateValuesRule(),
					Message: `Duplicate value: "a", first defined at module.tf:4,11-14`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 11},
						End:      hcl.Pos{Line: 5, Column: 14},
					},
				},
				{
					Rule:    NewTerraformMapDuplicateValuesRule(),
					Message: `Duplicate value: "a", first defined at module.tf:4,11-14`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 6, Column: 11},
						End:      hcl.Pos{Line: 6, Column: 14},
					},
				},
			},
		},
		{
			Name: "Using same value in different maps is okay",
			Content: `
resource "null_resource" "test" {
	map1 = {
	  key1 = "x"
	}
	map2 = {
	  key2 = "x"
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Using sensitive variable values",
			Content: `
variable "sensitive" {
  default = "secret"
  sensitive = true
}

resource "null_resource" "test" {
  map = {
    key1 = var.sensitive
    key2 = "secret"
  }
}`,
			// Do not report sensitive duplicate values to prevent unintentional exposure of sensitive values
			Expected: helper.Issues{},
		},
		{
			Name: "Using non-string values",
			Content: `
resource "null_resource" "test" {
  map = {
    key1 = 1
    key2 = 1
    key3 = {}
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMapDuplicateValuesRule(),
					Message: `Duplicate value: "1", first defined at module.tf:4,12-13`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 13},
					},
				},
			},
		},
		{
			Name: "values in for expressions",
			Content: `
resource "null_resource" "test" {
  list = [for a in ["foo", "bar"] : {
    key1 = "${a}_baz"
	key2 = "foo_baz"
  }]
}`,
			// The current implementation cannot find duplicate values in for expressions.
			Expected: helper.Issues{},
		},
		{
			Name: "ignore boolean string values",
			Content: `
resource "null_resource" "test" {
  map = {
    key1 = true
    key2 = true
  }
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformMapDuplicateValuesRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := testRunner(t, map[string]string{"module.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
		})
	}
}
