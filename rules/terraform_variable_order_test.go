package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableOrderRule(t *testing.T) {
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
			Expected: helper.Issues{
				{
					Rule: NewTerraformVariableOrderRule(),
					Message: `Recommended variable order:
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}`,
				},
			},
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
			Expected: helper.Issues{
				{
					Rule: NewTerraformVariableOrderRule(),
					Message: `Recommended variable order:
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
				},
			},
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
			Expected: helper.Issues{
				{
					Rule: NewTerraformVariableOrderRule(),
					Message: `Recommended variable order:
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
				},
			},
		},
	}
	rule := NewTerraformVariableOrderRule()

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
