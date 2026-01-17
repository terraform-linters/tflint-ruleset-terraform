# terraform_comment_syntax

Enforce usage of `#` for comments.

## Example

```hcl
# Good
// Bad
/*
  Bad
*/
```

```
$ tflint
2 issue(s) found:

Warning: [Fixable] Comments should begin with # (terraform_comment_syntax)

  on t.tf line 2:
   2: // Bad

Warning: [Fixable] Comments should begin with # (terraform_comment_syntax)

  on t.tf line 3:
   3: /*
```

## Why

The Terraform language supports two different syntaxes for single-line comments: `#` and `//` as well as `/*` `*/` for
multiline comments. However `#` is considered idiomatic for both single and multi-line comments.

* [Configuration Syntax: Comments](https://developer.hashicorp.com/terraform/language/syntax/configuration#comments)
* [Code Style](https://developer.hashicorp.com/terraform/language/style#code-style)

## How To Fix

Run `tflint --fix` to automatically replace `//` comments and multi-line `/* */` comments with `#` comments.

Single-line `/* */` comments are ignored because they can appear mid-expression (e.g., `x = 1 /* comment */ + 2`),
where converting to `#` would comment out the rest of the line.
