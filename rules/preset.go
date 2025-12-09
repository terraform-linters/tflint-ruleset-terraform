package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var PresetRules = map[string][]tflint.Rule{
	"all": {
		NewOpentofuCommentSyntaxRule(),
		NewOpentofuDeprecatedIndexRule(),
		NewOpentofuDeprecatedInterpolationRule(),
		NewOpentofuDeprecatedLookupRule(),
		NewOpentofuDocumentedOutputsRule(),
		NewOpentofuDocumentedVariablesRule(),
		NewOpentofuEmptyListEqualityRule(),
		NewOpentofuJSONSyntaxRule(),
		NewOpentofuMapDuplicateKeysRule(),
		NewOpentofuModulePinnedSourceRule(),
		NewOpentofuModuleShallowCloneRule(),
		NewOpentofuModuleVersionRule(),
		NewOpentofuNamingConventionRule(),
		NewOpentofuRequiredProvidersRule(),
		NewOpentofuRequiredVersionRule(),
		NewOpentofuStandardModuleStructureRule(),
		NewOpentofuTypedVariablesRule(),
		NewOpentofuUnusedDeclarationsRule(),
		NewOpentofuUnusedRequiredProvidersRule(),
		NewOpentofuWorkspaceRemoteRule(),
	},
	"recommended": {
		NewOpentofuDeprecatedIndexRule(),
		NewOpentofuDeprecatedInterpolationRule(),
		NewOpentofuDeprecatedLookupRule(),
		NewOpentofuEmptyListEqualityRule(),
		NewOpentofuJSONSyntaxRule(),
		NewOpentofuMapDuplicateKeysRule(),
		NewOpentofuModulePinnedSourceRule(),
		NewOpentofuModuleVersionRule(),
		NewOpentofuRequiredProvidersRule(),
		NewOpentofuRequiredVersionRule(),
		NewOpentofuTypedVariablesRule(),
		NewOpentofuUnusedDeclarationsRule(),
		NewOpentofuWorkspaceRemoteRule(),
	},
}
