# terraform_private_module_reference

According to the [Standard Module Structure](https://developer.hashicorp.com/terraform/language/modules/develop/structure):

> Nested modules should exist under the modules/ subdirectory. Any nested module with a README.md is considered usable by an external user.

This rule only checks local path references and ignores remote module references.

## Example

```hcl
module "foo" {
  source = "../../another-root/foo"
}
```

```plain
$ tflint
1 issue(s) found:

Warning: Private modules should not be referenced externally. Add a README.md to make the referenced module public or remove the reference. (terraform_private_module_reference)

  on main.tf line 2:
   2: module "foo" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.9.2/docs/rules/terraform_private_module_reference.md
```

## Why

Terraform does not enforce the convention described by the [Standard Module Structure](https://developer.hashicorp.com/terraform/language/modules/develop/structure). This `tflint` rule can be used to enforce the described convention.

It is best not to have consumers of a module that was not intended to be used externally.

## How To Fix

Either add a README.md to the private module to make it public or remove the reference to the private module.
