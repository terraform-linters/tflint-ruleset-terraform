## 0.6.0 (2024-02-24)

### Enhancements

- [#156](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/156): workspace_remote: Suppress issues in Terraform v1.1+

### Chores

- [#140](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/140): Bump golang.org/x/net from 0.13.0 to 0.17.0
- [#141](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/141): Bump github.com/hashicorp/go-getter from 1.7.2 to 1.7.3
- [#142](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/142): Bump github.com/google/go-cmp from 0.5.9 to 0.6.0
- [#143](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/143): Bump github.com/hashicorp/hcl/v2 from 2.18.1 to 2.19.1
- [#144](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/144): Bump google.golang.org/grpc from 1.57.0 to 1.57.1
- [#147](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/147): Bump github.com/hashicorp/terraform-registry-address from 0.2.2 to 0.2.3
- [#152](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/152): Bump actions/setup-go from 4 to 5
- [#153](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/153): Bump github.com/zclconf/go-cty from 1.14.1 to 1.14.2
- [#155](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/155): deps: Go 1.22
- [#158](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/158): Revise rules documentation

## 0.5.0 (2023-10-09)

### Enhancements

- [#128](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/128) [#132](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/132): new rule: `terraform_deprecated_lookup`
- [#131](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/131): `terraform_naming_convention`: Add support for checks and scoped data sources
- [#135](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/135): `terraform_unused_declarations`: Add support for scoped data sources
- [#136](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/136): Add support for provider refs in scoped data sources

## BugFixes

- [#133](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/133): `terraform_unused_declarations`: Make unused variable checks aware of validation blocks

### Chores

- [#106](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/106) [#117](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/117): Bump github.com/hashicorp/terraform-registry-address from 0.2.0 to 0.2.2
- [#108](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/108): Bump github.com/hashicorp/go-getter from 1.7.1 to 1.7.2
- [#109](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/109): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.17.0 to 0.18.0
- [#114](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/114): Add raw binary entries to checksums.txt
- [#115](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/115): Fix typo in rule documentation
- [#122](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/122) [#123](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/123) [#137](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/137): Bump github.com/zclconf/go-cty from 1.13.2 to 1.14.1
- [#124](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/124) [#138](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/138): Bump github.com/hashicorp/hcl/v2 from 2.17.0 to 2.18.1
- [#126](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/126): deps: Go 1.21
- [#127](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/127): Bump actions/checkout from 3 to 4
- [#129](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/129): Bump goreleaser/goreleaser-action from 4 to 5

## 0.4.0 (2023-06-18)

### Breaking Changes

- [#104](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/104): Bump tflint-plugin-sdk to v0.17.0
  - This change drops support for TFLint v0.40/v0.41

### Enhancements

- [#93](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/93): Add autofix support
  - `terraform_comment_syntax`
  - `terraform_deprecated_index`
  - `terraform_deprecated_interpolation`
  - `terraform_empty_list_equality`
  - `terraform_required_provider`
    - However, only issues with missing `source` can be fixed
  - `terraform_unused_declarations`
    - HCL native syntax only

### BugFixes

- [#101](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/101): deprecated_index: restore evaluation of JSON expressions

### Chores

- [#96](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/96): terraform_deprecated_index: add example of fix
- [#99](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/99): Bump github.com/zclconf/go-cty from 1.13.1 to 1.13.2
- [#102](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/102): Bump github.com/hashicorp/hcl/v2 from 2.16.2 to 2.17.0

## 0.3.0 (2023-04-22)

### Breaking Changes

- [#64](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/64): required_providers: warn on legacy version syntax, missing source

### BugFixes

- [#63](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/63): required_providers: use required provider entry as range if present
- [#90](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/90): terraform_deprecated_index: Emit issues based on expression types

### Chores

- [#57](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/57): Fix typo in rule documentation
- [#65](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/65) [#70](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/70) [#79](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/79): Bump github.com/hashicorp/hcl/v2 from 2.15.0 to 2.16.2
- [#67](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/67): Bump golang.org/x/net from 0.3.0 to 0.7.0
- [#68](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/68): Bump github.com/aws/aws-sdk-go from 1.15.78 to 1.34.0
- [#69](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/69) [#82](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/82): Bump github.com/hashicorp/go-getter from 1.6.2 to 1.7.1
- [#76](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/76) [#81](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/81): Bump github.com/zclconf/go-cty from 1.12.1 to 1.13.1
- [#78](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/78): Bump sigstore/cosign-installer from 2 to 3
- [#80](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/80): Bump actions/setup-go from 3 to 4
- [#83](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/83): Bump github.com/hashicorp/terraform-registry-address from 0.1.0 to 0.2.0
- [#85](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/85) [#88](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/88): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.15.0 to 0.16.1
- [#87](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/87): Bump github.com/Masterminds/semver/v3 from 3.2.0 to 3.2.1

## 0.2.2 (2022-12-26)

### BugFixes

- [#49](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/49): terraform_deprecated_index: improve perf for files with many expressions

### Chores

- [#45](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/45): Add signatures for keyless signing
- [#46](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/46): Bump github.com/hashicorp/hcl/v2 from 2.14.1 to 2.15.0
- [#47](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/47): Bump github.com/zclconf/go-cty from 1.11.1 to 1.12.1
- [#50](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/50): Bump github.com/Masterminds/semver/v3 from 3.1.1 to 3.2.0
- [#51](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/51): Use # comment syntax in configuration example
- [#53](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/53): Bump goreleaser/goreleaser-action from 3 to 4
- [#54](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/54): Bump tflint-plugin-sdk to v0.15.0
- [#55](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/55): Bump terraform-registry-address to v0.1.0

## 0.2.1 (2022-10-26)

### BugFixes

- [#43](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/43): terraform_deprecated_index: handle Terraform directives (HCL template)

### Chores

- [#39](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/39): Bump github.com/zclconf/go-cty from 1.11.0 to 1.11.1

## 0.2.0 (2022-10-23)

### Enhancements

- [#34](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/34): deprecated_index: reject legacy splat expressions

### Chores

- [#30](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/30): Bump github.com/hashicorp/hcl/v2 from 2.14.0 to 2.14.1
- [#31](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/31): Fix typo in configuration.md
- [#38](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/38): Bump tflint-plugin-sdk to v0.14.0

## 0.1.1 (2022-09-17)

### BugFixes

- [#26](https://github.com/terraform-linters/tflint-ruleset-terraform/pull/26): Prefer --only option over other rules config
  - TFLint v0.40.1+ is required to apply this bug fix.

## 0.1.0 (2022-09-08)

Initial release 🎉

The rules have been migrated from TFLint v0.39.3. See the TFLint's CHANGELOG for a history of previous changes to these rules.
