# terraform_json_syntax

Enforce the official Terraform JSON syntax that uses a root object instead of a root array.

## Example

```json
# Good
{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}

# Bad
[{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}]
```

```
$ tflint
1 issue(s) found:

Warning: JSON configuration uses array syntax at root. The official Terraform JSON syntax uses a root object. See https://developer.hashicorp.com/terraform/language/syntax/json (terraform_json_syntax)

  on main.tf.json line 1:
   1: [{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}]

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.1.0/docs/rules/terraform_json_syntax.md
```

## Why

According to the official Terraform documentation, "At the root of any JSON-based Terraform configuration is a JSON object." While Terraform may parse array-based syntax in some cases, the documented and supported standard is to use a root object with top-level keys like `resource`, `variable`, `output`, etc.

Using the official syntax ensures compatibility and consistency with Terraform's expected JSON structure.

* [JSON Configuration Syntax](https://developer.hashicorp.com/terraform/language/syntax/json)

## How To Fix

Convert your array-based JSON configuration to use a root object. Instead of wrapping configuration in an array, use an object with appropriate top-level keys:

```json
# Before
[
  {"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}},
  {"variable": {"region": {"type": "string"}}}
]

# After
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
