# terraform_nullable_variables

Disallow `variable` declarations without `nullable` field.

## Example

```hcl
variable "no_nullable" {}

variable "enabled" {
  default     = false
  description = "This is description"
  nullable    = false
  type        = bool
}
```

```
$ tflint
1 issue(s) found:

Warning: `no_nullable` variable has no nullable field (terraform_nullable_variables)

  on template.tf line 1:
   1: variable "no_nullable" {}

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.11.0/docs/rules/terraform_nullable_variables.md
```

## Why

`nullable` field is optional and `true` by default. This rule forces explicit setting of the `nullable` field for variables.

## How To Fix

Add a `nullable` field to the variable. See https://developer.hashicorp.com/terraform/language/values/variables#disallowing-null-input-values for more details about `nullable`.
