# terraform_static_attribute_notation

Enforce dot notation for static attribute access in Terraform configurations.

> This rule is enabled by "recommended" preset.

## Example

```hcl
resource "aws_instance" "web" {
  instance_type = var.instance["type"]  # bracket notation in static context
}
```

```
$ tflint
1 issue(s) found:

Error: Must use dot notation for static attributes (terraform_static_attribute_notation)

  on main.tf line 2:
   2:   instance_type = var.instance["type"]

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.9.0/docs/rules/terraform_static_attribute_notation.md
```

## Why

In static contexts (when accessing literal attribute keys), using bracket notation is unnecessary and reduces readability. Using dot notation makes the code cleaner and easier to understand.

## How To Fix

Replace bracket notation with dot notation in static contexts:

```hcl
resource "aws_instance" "web" {
  instance_type = var.instance.type
}
``` 