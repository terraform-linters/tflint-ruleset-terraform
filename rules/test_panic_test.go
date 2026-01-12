package rules

import (
"testing"
"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformMapDuplicateKeys_Setproduct(t *testing.T) {
cases := []struct {
Name     string
Content  string
Expected helper.Issues
}{
{
Name: "setproduct with for expression",
Content: `
locals {
  data = { for item in setproduct(
    [true],
    ["foo", "bar"]
  ) : item[1] => { foo = item[0], bar = item[1] }
  }
}`,
Expected: helper.Issues{},
},
{
Name: "map comprehension with indexing",
Content: `
locals {
  data = {
    for k, v in { foo = { momo = "bar" } } : k => v["momo"]
  }
}`,
Expected: helper.Issues{},
},
}

rule := NewTerraformMapDuplicateKeysRule()

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
