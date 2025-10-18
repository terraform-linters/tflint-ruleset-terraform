# terraform_comment_syntax

Enforce usage of `#` for comments.

## Example

```hcl
# Good
// Bad
/* Bad */
```

```
$ tflint
2 issue(s) found:

Warning: [Fixable] Comments should begin with # (terraform_comment_syntax)

  on t.tf line 2:
   2: // Bad
   3: /* Bad */

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.14.1/docs/rules/terraform_comment_syntax.md

Warning: Comments should begin with # (terraform_comment_syntax)

  on t.tf line 3:
   3: /* Bad */

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.14.1/docs/rules/terraform_comment_syntax.md
```

## Why

The Terraform language supports two different syntaxes for single-line comments: `#` and `//` as well as `/*` `*/` for
multiline comments. However `#` are considered to be idiomatic for both single and multi-line comments.

* [Configuration Syntax: Comments](https://developer.hashicorp.com/terraform/language/syntax/configuration#comments)
* [Code Style](https://developer.hashicorp.com/terraform/language/style#code-style)

## How To Fix

Replace the leading double-slash (`//`) in your comment with the number sign (`#`) for single-line comments which can
also be fixed by running tflint with `--fix` flag. For multiline comments remove the `/*` and `*/` and put `#` at the
start of each line of the comment. This is not fixed by tflint since multiline comments can be put inside expressions,
so you need to move them to the end of the line or line above the expression they're in.
