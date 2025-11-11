# terraform_unused_declarations

Disallow variables, data sources, locals, and provider aliases that are declared but never used.

> This rule is enabled by "recommended" preset.

## Example

```hcl
variable "not_used" {}

variable "used" {}
output "out" {
  value = var.used
}
```

```
$ tflint
1 issue(s) found:

Warning: variable "not_used" is declared but not used (terraform_unused_declarations)

  on config.tf line 1:
   1: variable "not_used" {

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_unused_declarations.md
 
```

Provider aliases example:

```hcl
provider "azurerm" {
  features {}
  alias           = "test_123"
  subscription_id = ""
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
  provider = azurerm.test_123
}
```

```
$ tflint
0 issue(s) found
```

Without the resource using the aliased provider:

```hcl
provider "azurerm" {
  features {}
  alias           = "test_123"
  subscription_id = ""
}
```

```
$ tflint
1 issue(s) found:

Warning: provider "azurerm" with alias "test_123" is declared but not used (terraform_unused_declarations)

  on config.tf line 1:
   1: provider "azurerm" {
```

## Why

Terraform will ignore variables and locals that are not used. It will refresh declared data sources regardless of usage. However, unreferenced variables and provider aliases likely indicate either a bug (and should be referenced) or removed code (and should be removed).

## How To Fix

Remove the declaration. For `variable`, `data`, and `provider` (with alias), remove the entire block. For a `local` value, remove the attribute from the `locals` block.

While data sources should generally not have side effects, take greater care when removing them. For example, removing `data "http"` will cause Terraform to no longer perform an HTTP `GET` request during each plan. If a data source is being used for side effects, add an annotation to ignore it:

```tf
# tflint-ignore: terraform_unused_declarations
data "http" "example" {
  url = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
}
```
