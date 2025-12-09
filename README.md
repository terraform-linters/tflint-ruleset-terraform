# TFLint Ruleset for OpenTofu Language
[![Build Status](https://github.com/terraform-linters/tflint-ruleset-terraform/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/terraform-linters/tflint-ruleset-terraform/actions)
[![GitHub release](https://img.shields.io/github/release/terraform-linters/tflint-ruleset-terraform.svg)](https://github.com/terraform-linters/tflint-ruleset-terraform/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

TFLint ruleset plugin for OpenTofu Language

This project is forked from https://github.com/terraform-linters/tflint-ruleset-terraform

This ruleset focus on possible errors and best practices about OpenTofu Language.

## Requirements

- TFLint v0.46+
- Go v1.25

## Installation

This ruleset is built into TFLint, so you usually don't need to worry about how to install it. You can check the built-in version with `tflint -v`:

```
$ tflint -v
TFLint version 0.52.0
+ ruleset.opentofu (0.8.0-bundled)
```

If you want to use a version different from the built-in version, you can declare `plugin` in `.tflint.hcl` as follows and install it with `tflint --init`:

```hcl
plugin "opentofu" {
    enabled = true
    version = "0.13.0"
    source  = "github.com/diofeher/tflint-ruleset-opentofu"
}
```

For more configuration about the plugin, see [Plugin Configuration](docs/configuration.md).

## Rules

See [Rules](docs/rules/README.md).

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

Note that if you install the plugin with `make install`, you must omit the `version` and `source` attributes in` .tflint.hcl`:

```hcl
plugin "opentofu" {
    enabled = true
}
```
