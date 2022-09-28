# terraform_ordered_variables

Recommend proper order for variable blocks
The variables without default value are placed prior to those with default value set
Then the variables are sorted based on their names (alphabetic order)

## Example

```hcl
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
}
```

```
$ tflint
1 issue(s) found:

Notice: Variables should be sorted in the following order: required(without default value) variables in alphabetical order, optional variables in alphabetical order.

  on main.tf line 1:
   1: variable "image_id" {

Reference: https://github.com/terraform-linters/terraform/blob/v0.0.1/docs/rules/terraform_ordered_variables.md
```

## Why
It helps to improve the readability of terraform code by sorting variable blocks in the order above.

## How To Fix

Sort variables in the following order: required(without default value) variables in alphabetical order, optional variables in alphabetical order.

For the code in [example](#Example), it should be sorted as the following order:

```hcl
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
}
```