# terraform_workspace_remote

`terraform.workspace` should not be used with a "remote" backend with remote execution in Terraform v1.0.x.

If remote operations are [disabled](https://www.terraform.io/docs/cloud/run/index.html#disabling-remote-operations) for your workspace, you can safely disable this rule:

```hcl
rule "terraform_workspace_remote" {
  enabled = false
}
```

This rule looks at `required_version` for Terraform version estimation. If the `required_version` is not declared, it is assumed that you are using a more recent version.

> This rule is enabled by "recommended" preset.

## Example

```hcl
terraform {
  required_version = ">= 1.0"
  backend "remote" {
    # ...
  }
}

resource "aws_instance" "a" {
  tags = {
    workspace = terraform.workspace
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: terraform.workspace should not be used with a 'remote' backend (terraform_workspace_remote)

  on example.tf line 8:
   9:   tags = {
  10:     workspace = terraform.workspace
  11:   }

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.5.0/docs/rules/terraform_workspace_remote.md
```

## Why

Terraform configuration may include the name of the [current workspace](https://developer.hashicorp.com/terraform/language/state/workspaces#current-workspace-interpolation) using the `${terraform.workspace}` interpolation sequence. However, when Terraform Cloud workspaces are executing Terraform runs remotely, the Terraform v1.0.x always uses the `default` workspace.

The [remote](https://developer.hashicorp.com/terraform/language/settings/backends/remote) backend is used with Terraform Cloud workspaces. Even if you set a `prefix` in the `workspaces` block, this value will be ignored during remote runs.

For more information, see the [`remote` backend workspaces documentation](https://developer.hashicorp.com/terraform/language/settings/backends/remote#workspace-names).

## How To Fix

If you still need support for Terarform v1.0.x, consider adding a variable to your configuration and setting it in each cloud workspace:

```tf
variable "workspace" {
  type        = string
  description = "The workspace name" 
}
```

You can also name the variable based on what the workspace suffix represents in your configuration (e.g. environment).

If you don't need support for Terraform v1.0.x, you can suppress the issue by updating the `required_version` to not contain 1.0.x.
