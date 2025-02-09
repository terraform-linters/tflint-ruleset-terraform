# terraform_no_short_circuit_evaluation

Disallow using logical operators (`&&`, `||`) with null checks that could lead to errors due to lack of short-circuit evaluation.

> This rule is enabled by "recommended" preset.

## Example

```hcl
# This will error if var.obj is null
resource "aws_instance" "example" {
  count = var.obj != null && var.obj.enabled ? 1 : 0
}

# This is the safe way to write it
resource "aws_instance" "example" {
  count = var.obj != null ? var.obj.enabled ? 1 : 0 : 0
}
```

```
$ tflint
1 issue(s) found:

Warning: Short-circuit evaluation is not supported in Terraform. Use a conditional expression (condition ? true : false) instead. (terraform_no_short_circuit_evaluation)

  on main.tf line 3:
   3:   count = var.obj != null && var.obj.enabled ? 1 : 0

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_no_short_circuit_evaluation.md
```

## Why

Unlike many programming languages, Terraform's logical operators (`&&` and `||`) do not short-circuit. This means that in an expression like `var.obj != null && var.obj.enabled`, both sides will be evaluated even if `var.obj` is null, which will result in an error.

This is a common source of confusion for users coming from other programming languages where short-circuit evaluation is standard behavior. The issue is particularly problematic when checking for null before accessing object attributes.

## How To Fix

Use nested conditional expressions instead of logical operators when you need short-circuit behavior. For example:

```hcl
# Instead of this:
var.obj != null && var.obj.enabled

# Use this:
var.obj != null ? var.obj.enabled : false

# For more complex conditions:
var.obj != null ? (var.obj.enabled ? var.obj.value > 0 : false) : false
```

You can also use the `try()` function in some cases, though this may mask errors you want to catch:

```hcl
try(var.obj.enabled, false)
```

For more information, see [hashicorp/terraform#24128](https://github.com/hashicorp/terraform/issues/24128). 