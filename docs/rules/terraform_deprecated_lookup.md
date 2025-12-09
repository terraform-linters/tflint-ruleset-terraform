# terraform_deprecated_lookup

Disallow deprecated [`lookup` function](https://developer.hashicorp.com/terraform/language/functions/lookup) usage without a default.

> This rule is enabled by "recommended" preset.

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

Reference: https://github.com/diofeher/tflint-ruleset-opentofu/blob/v0.5.0/docs/rules/terraform_deprecated_lookup.md
```

## Why

Calling [`lookup`](https://developer.hashicorp.com/terraform/language/functions/lookup) with 2 arguments has been deprecated since Terraform v0.7. `lookup(map, key)` is equivalent to the native index syntax `map[key]`. `lookup` should only be used with the third `default` argument, even though it is optional for backward compatibility. 

## How To Fix

Use the native index syntax:

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
