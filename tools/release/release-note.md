## What's Changed

Support for Cosign signatures has been removed from this release. The `checksums.txt.keyless.sig` and `checksums.txt.pem` will not be included in the release.
These files are not used in normal use cases, so in most cases this will not affect you, but if you are affected, you can use Artifact Attestations instead.

### Breaking Changes
* Bump github.com/terraform-linters/tflint-plugin-sdk from 0.22.0 to 0.23.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/289
  * Requires TFLint v0.46+

### Enhancements
* add terraform_json_syntax rule by @bendrucker in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/297
* `terraform_unused_declarations`: detect unused provider aliases by @bendrucker in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/304

### Bug Fixes
* `module_pinned_source`: handle directories in path by @bendrucker in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/296

### Chores
* Bump goreleaser/goreleaser-action from 6.3.0 to 6.4.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/275
* Bump github.com/hashicorp/go-getter from 1.7.8 to 1.7.9 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/277
* Bump actions/checkout from 4.2.2 to 5.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/276
* Bump github.com/zclconf/go-cty from 1.16.3 to 1.16.4 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/278
* Bump github.com/ulikunitz/xz from 0.5.10 to 0.5.14 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/279
* dependabot: allow actions writes by @wata727 in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/280
* Bump github.com/hashicorp/terraform-registry-address from 0.3.0 to 0.4.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/281
* Bump actions/attest-build-provenance from 2.4.0 to 3.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/282
* Bump actions/setup-go from 5.5.0 to 6.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/283
* Bump github.com/zclconf/go-cty from 1.16.4 to 1.17.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/284
* Bump github.com/hashicorp/go-getter from 1.7.9 to 1.8.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/285
* Bump sigstore/cosign-installer from 3.9.2 to 3.10.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/286
* Bump github.com/hashicorp/go-getter from 1.8.0 to 1.8.1 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/287
* Bump github.com/hashicorp/go-getter from 1.8.1 to 1.8.2 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/294
* Bump sigstore/cosign-installer from 3.10.0 to 4.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/299
* Bump github.com/terraform-linters/tflint-plugin-sdk from 0.23.0 to 0.23.1 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/300
* Bump github.com/hashicorp/go-getter from 1.8.2 to 1.8.3 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/301
* Bump Go version to v1.25 by @wata727 in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/302
* Bump golang.org/x/crypto from 0.42.0 to 0.45.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/306
* Bump actions/checkout from 5.0.0 to 6.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/307
* Bump actions/setup-go from 6.0.0 to 6.1.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/308
* Bump github.com/hashicorp/go-version from 1.7.0 to 1.8.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/309
* Bump actions/checkout from 6.0.0 to 6.0.1 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/311
* Bump actions/attest-build-provenance from 3.0.0 to 3.1.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/313
* Bump github.com/hashicorp/go-getter from 1.8.3 to 1.8.4 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/314
* Drop support for Cosign signatures by @wata727 in https://github.com/terraform-linters/tflint-ruleset-terraform/pull/315


**Full Changelog**: https://github.com/terraform-linters/tflint-ruleset-terraform/compare/v0.13.0...v0.14.0
