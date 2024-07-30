# terraform_map_duplicate_keys

Disallow duplicate keys in a map object.

> This rule is enabled by "recommended" preset.

## Example

```hcl
locals {
  map = {
    foo = 1
    bar = 2
    bar = 3 // duplicate key
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: Duplicate key: "bar", first defined at main.tf:4,5-8 (terraform_map_duplicate_keys)

  on main.tf line 5:
   5:     bar = 3 // duplicated key

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.9.0/docs/rules/terraform_map_duplicate_keys.md
```

## Why

In the Terraform language, duplicate map keys are overwritten rather than throwing an error. However, in most cases this behavior is not what you want and is often caused by a mistake. This rule will catch such mistakes early.

See also https://github.com/hashicorp/terraform/issues/28727

## How To Fix

Remove the duplicate keys and leave the correct value.
