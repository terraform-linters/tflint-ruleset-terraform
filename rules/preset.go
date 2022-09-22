package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var PresetRules = map[string][]tflint.Rule{
	"all": {
		NewTerraformCommentSyntaxRule(),
		NewTerraformDeprecatedIndexRule(),
		NewTerraformDeprecatedInterpolationRule(),
		NewTerraformDocumentedOutputsRule(),
		NewTerraformDocumentedVariablesRule(),
		NewTerraformEmptyListEqualityRule(),
		NewTerraformLocalsOrderRule(),
		NewTerraformModulePinnedSourceRule(),
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
		NewTerraformEmptyListEqualityRule(),
		NewTerraformLocalsOrderRule(),
		NewTerraformModulePinnedSourceRule(),
		NewTerraformModuleVersionRule(),
		NewTerraformRequiredProvidersRule(),
		NewTerraformRequiredVersionRule(),
		NewTerraformTypedVariablesRule(),
		NewTerraformUnusedDeclarationsRule(),
		NewTerraformWorkspaceRemoteRule(),
	},
}
