# terraform_comment_syntax

Disallow `//` comments in favor of `#`.

## Configuration

Name | Default | Value
--- | --- | ---
enabled | true | Boolean
allow_multiline | false | Boolean

```hcl
rule "terraform_comment_syntax" {
  enabled = true

  allow_multiline = false
}
```

## Example

```hcl
# Good
// Bad
```

```
$ tflint
1 issue(s) found:

Warning: Single line comments should begin with # (terraform_comment_syntax)

  on main.tf line 2:
   2: // Bad

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.14.0/docs/rules/terraform_comment_syntax.md
```

### allow_multiline = false

Disallows usage of multi-line comments.

```hcl
# Good
/* Bad */
```

```
Warning: Multi-line comments are not allowed. Use single-line comments starting with # (terraform_comment_syntax)

  on main.tf line 1:
   2: /* Bad */

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.14.0/docs/rules/terraform_comment_syntax.md
```

## Why

The Terraform language supports two different syntaxes for single-line comments: `#` and `//`. However, `#` is the default comment style and should be used in most cases.

* [Configuration Syntax: Comments](https://developer.hashicorp.com/terraform/language/syntax/configuration#comments)

Terraform also supports multi-line comments using `/*` and `*/` as delimiters. However, there's rarely a use-case where
it makes sense to use multi-line comments over multiple single-line comments. Additionally, modern editors make it easy
to work with single-line comments, and therefore it makes sense to disallow multi-line comments so there's only one syntax 
for comments.

## How To Fix

Replace the leading double-slash (`//`) in your comment with the number sign (`#`).
Replace multi-line comments with multiple single-line comments if `allow_multiline = false` (default).
