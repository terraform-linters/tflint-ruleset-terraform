# terraform_dynamic_attribute_notation

Enforce bracket notation for dynamic attribute access in Terraform configurations.

> This rule is enabled by "recommended" preset.

## Example

```hcl
resource "aws_instance" "web" {
  for_each = local.instances
  subnet_id = each.value.subnet_id  # dot notation in dynamic context
}
```

```
$ tflint
1 issue(s) found:

Error: Must use bracket notation [] for dynamic attributes. Use each.value["subnet_id"] instead (terraform_dynamic_attribute_notation)

  on main.tf line 3:
   3:   subnet_id = each.value.subnet_id

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.9.0/docs/rules/terraform_dynamic_attribute_notation.md
```

## Why

In dynamic contexts (e.g., inside a `for_each` block, `for` expression, or when using `count`), using dot notation may lead to ambiguous results. Using bracket notation makes the attribute access explicit and prevents potential confusion.

## How To Fix

Replace dot notation with bracket notation in dynamic contexts:

```hcl
resource "aws_instance" "web" {
  for_each = local.instances
  subnet_id = each.value["subnet_id"]
}
``` 