## What's Changed

### Breaking Changes
* `terraform_comment_syntax`: Enforce `#` comment syntax for multiline by @AleksaC in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/298

### Enhancements
* terraform: Add support dynamic module sources and versions by @wata727 in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/356
  * This change is required for Terraform v1.15 compatibility.

### Bug Fixes
* unused_declarations: fix false positive on provider aliases in JSON files by @bendrucker in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/330

### Chores
* Bump actions/setup-go from 6.1.0 to 6.2.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/321
* Bump actions/checkout from 6.0.1 to 6.0.2 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/324
* Bump actions/attest-build-provenance from 3.1.0 to 3.2.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/325
* Bump goreleaser/goreleaser-action from 6.4.0 to 7.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/331
* Bump go.opentelemetry.io/otel/sdk from 1.39.0 to 1.40.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/332
* Bump actions/attest-build-provenance from 3.2.0 to 4.1.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/333
* Bump github.com/zclconf/go-cty from 1.17.0 to 1.18.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/335
* Bump actions/setup-go from 6.2.0 to 6.3.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/334
* deps: Bump Go version to 1.26 by @wata727 in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/338
* Bump github.com/hashicorp/go-getter from 1.8.4 to 1.8.5 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/339
* Bump google.golang.org/grpc from 1.78.0 to 1.79.3 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/340
* Bump github.com/terraform-linters/tflint-plugin-sdk from 0.23.1 to 0.24.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/341
* dependabot: Set cooldown period by @wata727 in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/342
* Bump github.com/go-jose/go-jose/v4 from 4.1.3 to 4.1.4 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/343
* Bump actions/setup-go from 6.3.0 to 6.4.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/344
* Bump github.com/hashicorp/go-version from 1.8.0 to 1.9.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/345
* Bump github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream from 1.7.4 to 1.7.8 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/346
* Bump github.com/aws/aws-sdk-go-v2/service/s3 from 1.96.0 to 1.97.3 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/347
* Bump go.opentelemetry.io/otel/sdk from 1.40.0 to 1.43.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/348
* Bump github.com/hashicorp/go-getter from 1.8.5 to 1.8.6 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/349
* Bump goreleaser/goreleaser-action from 7.0.0 to 7.1.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/350
* Bump github.com/zclconf/go-cty from 1.18.0 to 1.18.1 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/351
* build(deps): bump goreleaser/goreleaser-action from 7.1.0 to 7.2.1 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/352
* build(deps): bump github.com/Masterminds/semver/v3 from 3.4.0 to 3.5.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/353
* build(deps): bump goreleaser/goreleaser-action from 7.2.1 to 7.2.2 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/355

## New Contributors
* @AleksaC made their first contribution in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/298

**Full Changelog**: https://github.com/terraform-linters/tflint-ruleset-terraform/compare/v0.14.1...v0.15.0
