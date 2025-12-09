# terraform_required_version

Disallow `terraform` declarations without `required_version`.

> This rule is enabled by "recommended" preset.

## Configuration

```hcl
rule "terraform_required_version" {
  enabled = true
}
```

## Example

```hcl
terraform {
  required_version = ">= 1.0" 
}
```

```
$ tflint
1 issue(s) found:

Warning: terraform "required_version" attribute is required

Reference: https://github.com/diofeher/tflint-ruleset-opentofu/blob/v0.1.0/docs/rules/terraform_required_version.md 
```

## Why
The `required_version` setting can be used to constrain which versions of the OpenTofu CLI are compatible with your configuration. 
If the running version of OpenTofu doesn't match the constraints specified, OpenTofu will produce an error and exit without 
taking any further actions.

## How To Fix

Add the `required_version` attribute to the terraform configuration block.
