# terraform_module_shallow_clone

Require pinned Git-hosted Terraform modules to use shallow cloning.

## Example

```hcl
module "ssh" {
  source = "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0"
}

module "https" {
  source = "git::https://github.com/hashicorp/consul.git?ref=v1.0.0"
}

module "github" {
  source = "github.com/hashicorp/consul?ref=v1.0.0"
}

module "shallow" {
  source = "git::https://github.com/hashicorp/consul.git?depth=1&ref=v1.0.0"
}

module "commit" {
  source = "git::https://github.com/hashicorp/consul.git?ref=abc123def456789012345678901234567890abcd"
}

module "registry" {
  source = "hashicorp/consul"
  version = "~> 1.0"
}
```

```
$ tflint
3 issue(s) found:

Warning: Module source "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter (terraform_module_shallow_clone)

  on main.tf line 2:
   3:   source = "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.13.0/docs/rules/terraform_module_shallow_clone.md

Warning: Module source "git::https://github.com/hashicorp/consul.git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter (terraform_module_shallow_clone)

  on main.tf line 6:
   7:   source = "git::https://github.com/hashicorp/consul.git?ref=v1.0.0"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.13.0/docs/rules/terraform_module_shallow_clone.md

Warning: Module source "github.com/hashicorp/consul?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter (terraform_module_shallow_clone)

  on main.tf line 10:
  11:   source = "github.com/hashicorp/consul?ref=v1.0.0"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.13.0/docs/rules/terraform_module_shallow_clone.md

```

## Why

When using Git-hosted Terraform modules that are pinned to a specific version (tag or branch), enabling shallow cloning can significantly improve performance by reducing the amount of data that needs to be downloaded. This is especially beneficial in CI/CD pipelines where modules are downloaded frequently.

Shallow cloning downloads only the specific commit being referenced rather than the entire git history, which can save substantial time and bandwidth.

## How To Fix

Add the `depth=1` parameter to your pinned Git-hosted module sources.

Note: This rule only applies to Git-hosted modules that are pinned to a specific version. Unpinned modules and modules using raw Git commit hashes are not affected by this rule.
