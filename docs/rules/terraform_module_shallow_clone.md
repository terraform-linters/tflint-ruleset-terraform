# terraform_module_shallow_clone

Require pinned Git-hosted Terraform modules to use shallow cloning.

## Example

```hcl
module "consul" {
  source = "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0"
}
```

```
$ tflint
1 issue(s) found:

Warning: Module source "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter (terraform_module_shallow_clone)

  on main.tf line 2:
   3:   source = "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0"

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.13.0/docs/rules/terraform_module_shallow_clone.md
```

## Why

https://developer.hashicorp.com/terraform/language/modules/sources#shallow-clone

When sourcing a Terraform module from a Git repository by tag or branch, enabling shallow cloning can significantly improve performance by reducing the amount of data that needs to be downloaded. This is especially beneficial in CI/CD pipelines where modules are downloaded frequently.

Shallow cloning only includes the most recent commit for a reference. Because it uses the `--branch` argument to `git clone`, it can only be used for named branches and tags, not raw commit IDs.

## How To Fix

Add the `depth=1` query parameter to enable shallow cloning.
