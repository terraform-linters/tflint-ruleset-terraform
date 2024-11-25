# terraform_map_duplicate_values

Disallow duplicate values in a map object.

## Example

```hcl
locals {
  map = {
    foo = 1
    bar = 1 // duplicate value
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: Duplicate key: "bar", first defined at main.tf:4,5-8 (terraform_map_duplicate_values)

  on main.tf line 5:
   5:     bar = 3 // duplicated value

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.11.0/docs/rules/terraform_map_duplicate_values.md
```

## Why

Sometimes, you want to maintain a map that contains only unique values (e.g., do not want to get duplicated SSM parameters values). This rule will catch such mistakes early.
The map structure is not a set, so it is possible to have duplicate values in a map, so make sure you run this rule only on files where you want to enforce unique values.

## How To Fix

Remove the duplicate values and leave the correct value.
