# terraform_json_syntax

Enforce the official Terraform JSON syntax that uses a root object with keys for each block type.

## Example

```json
[{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}]
```

```
$ tflint
1 issue(s) found:

Warning: JSON configuration uses array syntax at root, expected object (terraform_json_syntax)

  on main.tf.json line 1:
   1: [{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}]

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_json_syntax.md
```

## Why

The [Terraform JSON syntax documentation](https://developer.hashicorp.com/terraform/language/syntax/json#json-file-structure) states:

> At the root of any JSON-based Terraform configuration is a JSON object. The properties of this object correspond to the top-level block types of the Terraform language.

While Terraform's underlying HCL parser supports flattening arrays, the documented and supported standard is to use a root object with top-level keys for Terraform's block types: `resource`, `variable`, `output`, etc. Using the official syntax ensures compatibility with third party tools that implement the documented standard.

## How To Fix

Convert your array-based JSON configuration to use a root object. Instead of wrapping configuration in an array, use an object with appropriate top-level keys.

### Before

```json
[
  {"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}},
  {"variable": {"region": {"type": "string"}}}
]
```

### After

```json
{
  "resource": {
    "aws_instance": {
      "example": {
        "ami": "ami-12345678"
      }
    }
  },
  "variable": {
    "region": {
      "type": "string"
    }
  }
}
```
