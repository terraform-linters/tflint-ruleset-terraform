package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_OpentofuMapDuplicateKeys(t *testing.T) {
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
			Name: "duplicate keys in map literal",
			Content: `
resource "null_resource" "test" {
    triggers = {
        a = "b"
        a = "c"
    }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuMapDuplicateKeysRule(),
					Message: `Duplicate key: "a", first defined at module.tf:4,9-10`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 9},
						End:      hcl.Pos{Line: 5, Column: 10},
					},
				},
			},
		},
		{
			Name: "duplicate keys with quoting",
			Content: `
resource "null_resource" "test" {
    triggers = {
        a = "b"
        "a" = "c"
    }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuMapDuplicateKeysRule(),
					Message: `Duplicate key: "a", first defined at module.tf:4,9-10`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 9},
						End:      hcl.Pos{Line: 5, Column: 12},
					},
				},
			},
		},
		{
			Name: "Using variables as keys",
			Content: `
variable "a" {
  type    = string
  default = "b"
}

resource "null_resource" "test" {
	map = {
	  (var.a) = 5
	  b       = 8
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuMapDuplicateKeysRule(),
					Message: `Duplicate key: "b", first defined at module.tf:9,4-11`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 10, Column: 4},
						End:      hcl.Pos{Line: 10, Column: 5},
					},
				},
			},
		},
		{
			Name: "Using a variable as a key without a default",
			Content: `
variable "unknown" {
  type    = string
}

resource "null_resource" "test" {
	map = {
	  x             = 8
	  (var.unknown) = 5
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Multiple duplicates in same map",
			Content: `
resource "null_resource" "test" {
	map = {
	  a = 7
	  a = 8
	  a = 9
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuMapDuplicateKeysRule(),
					Message: `Duplicate key: "a", first defined at module.tf:4,4-5`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 4},
						End:      hcl.Pos{Line: 5, Column: 5},
					},
				},
				{
					Rule:    NewOpentofuMapDuplicateKeysRule(),
					Message: `Duplicate key: "a", first defined at module.tf:4,4-5`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 6, Column: 4},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
			},
		},
		{
			Name: "Using same key in different maps is okay",
			Content: `

resource "null_resource" "test" {
	map = {
	  x = 1
	}
	map2 = {
	  x = 2
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Using sensitive variable keys",
			Content: `
variable "sensitive" {
  default = "secret"
  sensitive = true
}

resource "null_resource" "test" {
  map = {
    (var.sensitive) = 1
    "secret" = 2
  }
}`,
			// Do not report duplicate keys to prevent unintentional exposure of sensitive values
			Expected: helper.Issues{},
		},
		{
			Name: "Using non-string keys",
			Content: `
resource "null_resource" "test" {
  map = {
    1 = 1
    1 = 2
    {} = 3
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewOpentofuMapDuplicateKeysRule(),
					Message: `Duplicate key: "1", first defined at module.tf:4,5-6`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 5},
						End:      hcl.Pos{Line: 5, Column: 6},
					},
				},
			},
		},
		{
			Name: "keys in for expressions",
			Content: `
resource "null_resource" "test" {
  list = [for a in ["foo", "bar"] : {
    "${a}_baz" = 1
	"foo_baz" = 2
  }]
}`,
			// The current implementation cannot find duplicate keys in for expressions.
			Expected: helper.Issues{},
		},
	}

	rule := NewOpentofuMapDuplicateKeysRule()

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
