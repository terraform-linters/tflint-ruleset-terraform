# terraform_documented_outputs

Disallow `output` declarations without description.

## Configuration

Name | Description | Default | Type
--- | --- | --- | ---
unique | Require every output description to be unique | false | Boolean

```hcl
rule "terraform_documented_outputs" {
  enabled = true
  unique  = false # default
}
```

## Example

```hcl
output "no_description" {
  value = "value"
}

output "empty_description" {
  value = "value"
  description = ""
}

output "description" {
  value = "value"
  description = "This is description"
}
```

```
$ tflint
2 issue(s) found:

Notice: `no_description` output has no description (terraform_documented_outputs)

  on template.tf line 1:
   1: output "no_description" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_documented_outputs.md

Notice: `empty_description` output has no description (terraform_documented_outputs)

  on template.tf line 5:
   5: output "empty_description" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_documented_outputs.md
 
```

### Unique

When `unique = true`, output descriptions are compared using exact string equality. Every output sharing a duplicated description is reported.

```hcl
rule "terraform_documented_outputs" {
  enabled = true
  unique  = true
}
```

```hcl
output "first" {
  description = "Shared description"
}

output "second" {
  description = "Shared description"
}
```

```
$ tflint
2 issue(s) found:

Notice: `first` output description is not unique: "Shared description" (terraform_documented_outputs)

  on outputs.tf line 2:
   2:   description = "Shared description"

Notice: `second` output description is not unique: "Shared description" (terraform_documented_outputs)

  on outputs.tf line 6:
   6:   description = "Shared description"
```

## Why

Since `description` is optional value, it is not always necessary to write it. But this rule is useful if you want to force the writing of description. Especially it is useful when combined with [terraform-docs](https://github.com/terraform-docs/terraform-docs).

## How To Fix

Write a description other than an empty string. When `unique = true`, manually reword duplicated descriptions so each output has distinct, meaningful documentation.
