package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableOrderRule(t *testing.T) {
	expectedIssue := &helper.Issue{
		Rule:    NewTerraformOrderedVariablesRule(),
		Message: `Variables should be sorted in the following order: required(without default value) variables in alphabetical order, optional variables in alphabetical order.`,
	}
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. no variable",
			Content: `
terraform{}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. correct variable order",
			Content: `
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "3. sorting based on default value",
			Content: `
variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "image_id" {
  type = string
}`,
			Expected: helper.Issues{expectedIssue},
		},
		{
			Name: "4. sorting in alphabetic order",
			Content: `
variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}`,
			Expected: helper.Issues{expectedIssue},
		},
		{
			Name: "5. mixed",
			Content: `
variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "image_id" {
  type = string
}`,
			Expected: helper.Issues{expectedIssue},
		},
		{
			Name: "6. required only",
			Content: `
variable "availability_zone_names" {
  type = list(string)
}

variable "image_id" {
  type = string
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "7. optional only",
			Content: `
variable "availability_zone_names" {
  type    = list(string)
  default = ["ap-northeast-1"]
}

variable "image_id" {
  type    = string
  default = "ami-063a9ea2ff5685f7f"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "8. incorrect required only",
			Content: `
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type = list(string)
}
`,
			Expected: helper.Issues{expectedIssue},
		},
		{
			Name: "9. incorrect optional only",
			Content: `
variable "image_id" {
  type    = string
  default = "ami-063a9ea2ff5685f7f"
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["ap-northeast-1"]
}
`,
			Expected: helper.Issues{expectedIssue},
		},
	}
	rule := NewTerraformOrderedVariablesRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
			if tc.JSON {
				filename = "config.tf.json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}
