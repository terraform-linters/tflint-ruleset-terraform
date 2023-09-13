# terraform_deprecated_lookup

Disallow deprecated `lookup()` function with only 2 arguments.

## Example

```hcl
locals {
  map   = { a = 0 }
  value = lookup(local.map, "a")
}
```

```
$ tflint
1 issue(s) found:

Warning: [Fixable] Lookup with 2 arguments is deprecated (terraform_deprecated_lookup)

  on main.tf line 3:
   3:   value = lookup(local.map, "a")

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.5.0/docs/rules/terraform_deprecated_lookup.md
```

## Why

Calling `lookup()` with 2 arguments is deprecated since Terraform v0.7. `lookup(map, key)` is equivalent to the native index syntax `map[key]`

* [lookup() documentation](https://developer.hashicorp.com/terraform/language/functions/lookup)

## How To Fix

Use the natice index syntax:

Example:

```hcl
locals {
  map   = { a = 0 }
  value = lookup(local.map, "a")
}
```

Change this to: 

```hcl
locals {
  map   = { a = 0 }
  value = local.map["a"]
}
```
