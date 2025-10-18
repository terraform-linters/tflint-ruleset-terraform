# terraform_required_providers

Require that all providers specify a `source` and `version` constraint through `required_providers`. This rule can also enforce specific provider constraints and restrict which providers can be used.

> This rule is enabled by "recommended" preset.

## Configuration

```hcl
rule "terraform_required_providers" {
  enabled = true

  # defaults
  source = true
  version = true

  # Optional: Set to true to only allow listed providers
  provider_whitelist = false

  # Optional: Define constraints for specific providers
  providers = {
    aws = {
      source  = "hashicorp/aws"     # Require official AWS provider
      version = "~> 5.0"             # Require version 5.x
    }
    azurerm = {
      source  = "hashicorp/azurerm"  # Require official Azure provider
      version = "~> 3.0"             # Require version 3.x
    }
  }
}
```

## Configuration Options

### Basic Options

- **`source`** (boolean, default: true): Require all providers to specify a source attribute
- **`version`** (boolean, default: true): Require all providers to specify a version constraint

### Advanced Options

- **`provider_whitelist`** (boolean, default: false): When set to `true`, modules can only use providers that are explicitly defined in your `providers` configuration. This creates an "allowlist" of approved providers.

- **`providers`** (map, default: {}): Defines specific constraints for each provider. For each provider, you can specify:
  - **`source`** (optional): The required source address (e.g., `"hashicorp/aws"`)
  - **`version`** (optional): The required version constraint pattern (e.g., `"~> 5.0"`)

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

Provider version constraints can be specified using a [version argument within a provider block](https://developer.hashicorp.com/terraform/language/providers/configuration#provider-versions) for backwards compatibility. This approach is now discouraged, particularly for child modules.

Optionally, you can disable enforcement of either `source` or `version` by setting the corresponding attribute in the rule configuration to `false`.

### Provider Constraints Examples

#### Enforcing Version Constraints

**Rule Configuration:**
```hcl
rule "terraform_required_providers" {
  enabled = true
  providers = {
    aws = {
      version = "~> 5.0"  # All modules must use this exact pattern
    }
  }
}
```

**Problem Code:**
```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0, < 6.0"  # Different pattern, but similar effect
    }
  }
}
```

**Error Message:**
```
Error: Provider "aws" version constraint does not match expected (expected: "~> 5.0", found: ">= 5.0, < 6.0")
```

The rule requires the exact pattern `"~> 5.0"`, not just any constraint that accepts version 5.x.

#### Enforcing Provider Sources

**Rule Configuration:**
```hcl
rule "terraform_required_providers" {
  enabled = true
  providers = {
    aws = {
      source  = "hashicorp/aws"  # Must use official HashiCorp source
    }
  }
}
```

**Problem Code:**
```hcl
terraform {
  required_providers {
    aws = {
      source  = "custom/aws"  # Using a potentially unsafe custom source
      version = "~> 5.0"
    }
  }
}
```

**Error Message:**
```
Error: Provider "aws" has incorrect source (expected: "hashicorp/aws", found: "custom/aws")
```

#### Restricting to Approved Providers

**Rule Configuration:**
```hcl
rule "terraform_required_providers" {
  enabled = true
  provider_whitelist = true  # Only allow listed providers

  providers = {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

**Terraform Code:**
```hcl
resource "random_string" "example" {  # Using unapproved provider
  length = 16
}
```

**Result:**
```
Error: Provider "random" is not in the allowed provider list
```
