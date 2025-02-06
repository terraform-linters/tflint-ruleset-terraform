# Rules

Terraform Language rules implement recommendations from the [Terraform Language documentation](https://www.terraform.io/language). This ruleset _does not_ provide configurable rules for personal/team style or usage preferences. If you'd like to enforce stylistic rules beyond the official Terraform Language recommendations, you should [author your own ruleset plugin](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/plugins.md).

All rules are enabled by default, but by setting `preset = "recommended"`, you can enable only the rules marked "Recommended" among the following rules. See [Configuration](../configuration.md) for details.

|Rule|Description|Recommended|
| --- | --- | --- |
|[terraform_comment_syntax](terraform_comment_syntax.md)|Disallow `//` comments in favor of `#`||
|[terraform_deprecated_index](terraform_deprecated_index.md)|Disallow legacy dot index syntax|✔|
|[terraform_deprecated_interpolation](terraform_deprecated_interpolation.md)|Disallow deprecated (0.11-style) interpolation|✔|
|[terraform_deprecated_lookup](terraform_deprecated_lookup.md)|Disallow deprecated `lookup()` function with only 2 arguments.|✔|
|[terraform_documented_outputs](terraform_documented_outputs.md)|Disallow `output` declarations without description||
|[terraform_documented_variables](terraform_documented_variables.md)|Disallow `variable` declarations without description||
|[terraform_dynamic_attribute_notation](terraform_dynamic_attribute_notation.md)|Enforce bracket notation for dynamic attribute access|✔|
|[terraform_empty_list_equality](terraform_empty_list_equality.md)|Disallow comparisons with `[]` when checking if a collection is empty|✔|
|[terraform_map_duplicate_keys](terraform_map_duplicate_keys.md)|Disallow duplicate keys in a map object|✔|
|[terraform_module_pinned_source](terraform_module_pinned_source.md)|Disallow specifying a git or mercurial repository as a module source without pinning to a version|✔|
|[terraform_module_version](terraform_module_version.md)|Checks that Terraform modules sourced from a registry specify a version|✔|
|[terraform_naming_convention](terraform_naming_convention.md)|Enforces naming conventions for resources, data sources, etc||
|[terraform_required_providers](terraform_required_providers.md)|Require that all providers have version constraints through required_providers|✔|
|[terraform_required_version](terraform_required_version.md)|Disallow `terraform` declarations without require_version|✔|
|[terraform_standard_module_structure](terraform_standard_module_structure.md)|Ensure that a module complies with the Terraform Standard Module Structure||
|[terraform_static_attribute_notation](terraform_static_attribute_notation.md)|Enforce dot notation for static attribute access|✔|
|[terraform_typed_variables](terraform_typed_variables.md)|Disallow `variable` declarations without type|✔|
|[terraform_unused_declarations](terraform_unused_declarations.md)|Disallow variables, data sources, and locals that are declared but never used|✔|
|[terraform_unused_required_providers](terraform_unused_required_providers.md)|Check that all `required_providers` are used in the module||
|[terraform_workspace_remote](terraform_workspace_remote.md)|`terraform.workspace` should not be used with a "remote" backend with remote execution in Terraform v1.0.x|✔|
