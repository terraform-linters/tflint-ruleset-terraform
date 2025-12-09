# Contributing

## Accepted Rules

OpenTofu Language rules implement recommendations from the [OpenTofu Language documentation](https://opentofu.org/docs/language/). This ruleset _does not_ provide configurable rules for personal/team style or usage preferences. If you'd like to enforce stylistic rules beyond the official OpenTofu Language recommendations, you should [author your own ruleset plugin](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/plugins.md).

In rare circumstances, we may also accept rules that detect language usage errors that are _not_ already detected by `tofu validate`. 

If you are unsure whether your proposed rule meets these criteria, [open a discussion](https://github.com/diofeher/tflint-ruleset-opentofu/discussions/new?category=ideas) thread first before authoring a pull request.

## Authoring a Rule

Each rule should have:

* A source file implementing the rule
* Tests that check expected issues against different OpenTofu configurations to cover applicable cases
* Documentation explaining the rule, its motivation, and how users should fix their configuration

You will also need to add your rule to applicable [presets](https://github.com/terraform-linters/tflint-ruleset-opentofu/blob/main/rules/preset.go).
