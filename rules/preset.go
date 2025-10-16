package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var PresetRules = map[string][]tflint.Rule{
	"all": {
		NewTerraformCommentSyntaxRule(),
		NewTerraformDeprecatedIndexRule(),
		NewTerraformDeprecatedInterpolationRule(),
		NewTerraformDeprecatedLookupRule(),
		NewTerraformDocumentedOutputsRule(),
		NewTerraformDocumentedVariablesRule(),
		NewTerraformEmptyListEqualityRule(),
		NewTerraformJSONSyntaxRule(),
		NewTerraformMapDuplicateKeysRule(),
		NewTerraformModulePinnedSourceRule(),
		NewTerraformModuleShallowCloneRule(),
		NewTerraformModuleVersionRule(),
		NewTerraformNamingConventionRule(),
		NewTerraformRequiredProvidersRule(),
		NewTerraformRequiredVersionRule(),
		NewTerraformStandardModuleStructureRule(),
		NewTerraformTypedVariablesRule(),
		NewTerraformUnusedDeclarationsRule(),
		NewTerraformUnusedRequiredProvidersRule(),
		NewTerraformWorkspaceRemoteRule(),
	},
	"recommended": {
		NewTerraformDeprecatedIndexRule(),
		NewTerraformDeprecatedInterpolationRule(),
		NewTerraformDeprecatedLookupRule(),
		NewTerraformEmptyListEqualityRule(),
		NewTerraformJSONSyntaxRule(),
		NewTerraformMapDuplicateKeysRule(),
		NewTerraformModulePinnedSourceRule(),
		NewTerraformModuleVersionRule(),
		NewTerraformRequiredProvidersRule(),
		NewTerraformRequiredVersionRule(),
		NewTerraformTypedVariablesRule(),
		NewTerraformUnusedDeclarationsRule(),
		NewTerraformWorkspaceRemoteRule(),
	},
}
