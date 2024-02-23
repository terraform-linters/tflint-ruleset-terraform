# terraform_deprecated_interpolation

Disallow deprecated (0.11-style) interpolation

> This rule is enabled by "recommended" preset.

## Example

```hcl
resource "aws_instance" "deprecated" {
    instance_type = "${var.type}"
}

resource "aws_instance" "new" {
    instance_type = var.type
}
```

```
$ tflint
1 issue(s) found:

Warning: Interpolation-only expressions are deprecated in Terraform v0.12.14 (terraform_deprecated_interpolation)

  on example.tf line 2:
   2:     instance_type = "${var.type}"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_deprecated_interpolation.md

```

## Why

Terraform v0.12 introduces a new interpolation syntax, but continues to support the old 0.11-style interpolation syntax for compatibility.

`terraform fmt` can replace this redundant interpolation, so although it is not deprecated in the latest Terraform version, this rule allows you to issue a warning similar to Terraform v0.12.14.

## How To Fix

Switch to the new interpolation syntax. See the release notes for Terraform 0.12.14 for details: https://github.com/hashicorp/terraform/releases/tag/v0.12.14
