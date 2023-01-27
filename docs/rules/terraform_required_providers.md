# terraform_required_providers

Require that all providers specify a `source` and `version` constraint through `required_providers`.

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

Warning: Missing version constraint for provider "template" in `required_providers` (terraform_required_providers)

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

Warning: Missing version constraint for provider "template" in `required_providers` (terraform_required_providers)

  on main.tf line 1:
   1: provider "template" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md
```

<hr>

```hcl
provider "template" {}

terraform {
  required_providers {
    template = "~> 2"
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: Legacy version constraint for provider "template" in `required_providers` (terraform_required_providers)

  on main.tf line 5:
   5:     template = "~> 2"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md
```

<hr>

```hcl
provider "template" {}

terraform {
  required_providers {
    template = {
      version = "~> 2"
    }
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: Legacy version constraint for provider "template" in `required_providers` (terraform_required_providers)

  on main.tf line 5:
   5:     template = "~> 2"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md
```

<hr>

```hcl
provider "template" {}

terraform {
  required_providers {
    template = {
      version = "~> 2"
    }
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: Missing `source` for provider "template" in `required_providers` (terraform_required_providers)

  on main.tf line 5:
   5:     template = {
   6:       version = "~> 2"
   7:     }

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_required_providers.md
```

## Why

Providers are plugins released on a separate rhythm from Terraform itself, and so they have their own version numbers. For production use, you should constrain the acceptable provider versions via configuration, to ensure that new versions with breaking changes will not be automatically installed by `terraform init` in future.

Terraform supports multiple provider registries/namespaces through the [`source` address](https://developer.hashicorp.com/terraform/language/providers/requirements#source-addresses) attribute. While this is optional for providers in `registry.terraform.io` under the `hashicorp` namespace (the defaults), it is required for all other providers. Omitting `source` is a common error when using third-party providers and using explicit source addresses for all providers is recommended.

## How To Fix

Add the [`required_providers`](https://developer.hashicorp.com/terraform/language/providers/requirements#requiring-providers) block to the `terraform` configuration block and include current versions for all providers. For example:

```tf
terraform {
  required_providers {
    template = {
      source  = "hashicorp/template"
      version = "~> 2"
    }
  }
}
```

Provider version constraints can be specified using a [version argument within a provider block](https://www.terraform.io/docs/configuration/providers.html#provider-versions) for backwards compatibility. This approach is now discouraged, particularly for child modules.
