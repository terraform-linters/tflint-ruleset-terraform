# terraform_documented_variables

Disallow `variable` declarations without description.

## Configuration

Name | Description | Default | Type
--- | --- | --- | ---
unique | Require every variable description to be unique | false | Boolean

```hcl
rule "terraform_documented_variables" {
  enabled = true
  unique  = false # default
}
```

## Example

```hcl
variable "no_description" {
  default = "value"
}

variable "empty_description" {
  default = "value"
  description = ""
}

variable "description" {
  default = "value"
  description = "This is description"
}
```

```
$ tflint
2 issue(s) found:

Notice: `no_description` variable has no description (terraform_documented_variables)

  on template.tf line 1:
   1: variable "no_description" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_documented_variables.md

Notice: `empty_description` variable has no description (terraform_documented_variables)

  on template.tf line 5:
   5: variable "empty_description" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_documented_variables.md

```

### Unique

When `unique = true`, variable descriptions are compared using exact string equality. Every variable sharing a duplicated description is reported.

```hcl
rule "terraform_documented_variables" {
  enabled = true
  unique  = true
}
```

```hcl
variable "first" {
  description = "Shared description"
}

variable "second" {
  description = "Shared description"
}
```

```
$ tflint
2 issue(s) found:

Notice: `first` variable description is not unique: "Shared description" (terraform_documented_variables)

  on variables.tf line 2:
   2:   description = "Shared description"

Notice: `second` variable description is not unique: "Shared description" (terraform_documented_variables)

  on variables.tf line 6:
   6:   description = "Shared description"
```

## Why

Since `description` is optional value, it is not always necessary to write it. But this rule is useful if you want to force the writing of description. Especially it is useful when combined with [terraform-docs](https://github.com/terraform-docs/terraform-docs).

## How To Fix

Write a description other than an empty string. When `unique = true`, manually reword duplicated descriptions so each variable has distinct, meaningful documentation.
