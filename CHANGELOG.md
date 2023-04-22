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
