# terraform_required_providers

Require that all providers have version constraints through `required_providers`.

> This rule is enabled by "recommended" preset.

## Configuration

```hcl
rule "terraform_required_providers" {
  enabled = true
}
```

## Examples

```hcl
provider "template" {}
```

```
$ tflint
1 issue(s) found:

Warning: Missing version constraint for provider "template" in "required_providers" (terraform_required_providers)

  on main.tf line 1:
   1: provider "template" {}

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md
```

<hr>

```hcl
provider "template" {
  version = "2"
}
```

```
$ tflint
2 issue(s) found:

Warning: provider.template: version constraint should be specified via "required_providers" (terraform_required_providers)

  on main.tf line 1:
   1: provider "template" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md

Warning: Missing version constraint for provider "template" in "required_providers" (terraform_required_providers)

  on main.tf line 1:
   1: provider "template" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md
```

## Why

Providers are plugins released on a separate rhythm from Terraform itself, and so they have their own version numbers. For production use, you should constrain the acceptable provider versions via configuration, to ensure that new versions with breaking changes will not be automatically installed by `terraform init` in future.

### Should we follow this rule in modules?

It depends on the module. Declaring the required versions is the recommended practice if the module is intended to be widely used. On the other hand, if the scope is limited, it may be sufficient to declare the required versions in the root module. You can ignore this rule if you find it redundant.

## How To Fix

Add the [`required_providers`](https://www.terraform.io/docs/configuration/terraform.html#specifying-required-provider-versions) block to the `terraform` configuration block and include current versions for all providers. For example:

```tf
terraform {
  required_providers {
    template = "~> 2.0"
  }
}
```

Provider version constraints can be specified using a [version argument within a provider block](https://www.terraform.io/docs/configuration/providers.html#provider-versions) for backwards compatability. This approach is now discouraged, particularly for child modules.
