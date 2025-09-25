package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	tfsdk "github.com/terraform-linters/tflint-plugin-sdk/terraform"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/terraform-linters/tflint-ruleset-terraform/terraform"
	"github.com/zclconf/go-cty/cty"
)

// TerraformRequiredProvidersRule checks whether Terraform sets version constraints for all configured providers
type TerraformRequiredProvidersRule struct {
	tflint.DefaultRule
}

// ProviderRequirement defines expected provider configuration for validation.
// Both Source and Version fields are optional - empty values skip validation for that field.
type ProviderRequirement struct {
	Source  *string `cty:"source"`  // Expected provider source (e.g., "hashicorp/aws")
	Version *string `cty:"version"` // Expected version constraint (e.g., "~> 4.0")
}

type terraformRequiredProvidersRuleConfig struct {
	// Source specifies whether the rule should assert the presence of a `source` attribute
	Source *bool `hclext:"source,optional"`
	// Version specifies whether the rule should assert the presence of a `version` attribute
	Version *bool `hclext:"version,optional"`
	// Providers defines the expected provider configurations mapped by provider name
	Providers map[string]ProviderRequirement `hclext:"providers,optional"`
	// ProviderWhitelist when true only allows providers defined in this rule,
	// when false validates only the configured providers (default: false)
	ProviderWhitelist *bool `hclext:"provider_whitelist,optional"`
}

// NewTerraformRequiredProvidersRule returns new rule with default attributes
func NewTerraformRequiredProvidersRule() *TerraformRequiredProvidersRule {
	return &TerraformRequiredProvidersRule{}
}

// Name returns the rule name
func (r *TerraformRequiredProvidersRule) Name() string {
	return "terraform_required_providers"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformRequiredProvidersRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformRequiredProvidersRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformRequiredProvidersRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// config returns the rule config, with defaults
func (r *TerraformRequiredProvidersRule) config(runner tflint.Runner) (*terraformRequiredProvidersRuleConfig, error) {
	config := &terraformRequiredProvidersRuleConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return nil, err
	}

	dv := true
	if config.Source == nil {
		config.Source = &dv
	}

	if config.Version == nil {
		config.Version = &dv
	}

	// Set default for ProviderWhitelist if not specified
	if config.ProviderWhitelist == nil {
		falseVal := false
		config.ProviderWhitelist = &falseVal
	}

	// Initialize empty map if not specified
	if config.Providers == nil {
		config.Providers = make(map[string]ProviderRequirement)
	}

	return config, nil
}

// Check Checks whether provider required version is set
func (r *TerraformRequiredProvidersRule) Check(rr tflint.Runner) error {
	runner := rr.(*terraform.Runner)

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	config, err := r.config(runner)
	if err != nil {
		return fmt.Errorf("failed to parse rule config: %w", err)
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "provider",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "version"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, provider := range body.Blocks {
		if _, exists := provider.Body.Attributes["version"]; exists {
			if err := runner.EmitIssue(
				r,
				"provider version constraint should be specified via `required_providers`",
				provider.DefRange,
			); err != nil {
				return err
			}
		}
	}

	providerRefs, diags := runner.GetProviderRefs()
	if diags.HasErrors() {
		return diags
	}

	requiredProvidersSchema := []hclext.AttributeSchema{}
	for name := range providerRefs {
		requiredProvidersSchema = append(requiredProvidersSchema, hclext.AttributeSchema{Name: name})
	}

	body, err = runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "terraform",
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: "required_providers",
							Body: &hclext.BodySchema{
								Attributes: requiredProvidersSchema,
							},
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	requiredProviders := hclext.Attributes{}
	for _, terraform := range body.Blocks {
		for _, requiredProvidersBlock := range terraform.Body.Blocks {
			for name, attr := range requiredProvidersBlock.Body.Attributes {
				requiredProviders[name] = attr
			}
		}
	}

	// Check provider whitelist if enabled
	if *config.ProviderWhitelist {
		for name, ref := range providerRefs {
			if name == "terraform" {
				continue // Skip builtin provider
			}
			if _, exists := config.Providers[name]; !exists {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Provider %q is not in the allowed provider list", name),
					ref.DefRange,
				); err != nil {
					return err
				}
			}
		}
	}

	for name, ref := range providerRefs {
		if name == "terraform" {
			// "terraform" provider is a builtin provider
			// @see https://github.com/hashicorp/terraform/blob/v1.2.5/internal/addrs/provider.go#L106-L112
			continue
		}

		requiredProvider, exists := requiredProviders[name]
		if !exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Missing version constraint for provider %q in `required_providers`", name),
				ref.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		val, diags := requiredProvider.Expr.Value(&hcl.EvalContext{
			Variables: map[string]cty.Value{
				// configuration_aliases can declare additional provider instances
				// required provider "foo" could have: configuration_aliases = [foo.a, foo.b]
				// @see https://www.terraform.io/language/modules/develop/providers#provider-aliases-within-modules
				name: cty.DynamicVal,
			},
		})
		if diags.HasErrors() {
			return diags
		}

		if val.Type() == cty.String {
			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf("Legacy version constraint for provider %q in `required_providers`", name),
				requiredProvider.Expr.Range(),
				func(f tflint.Fixer) error {
					if tfsdk.IsJSONFilename(requiredProvider.Expr.Range().Filename) {
						return tflint.ErrFixNotSupported
					}

					return f.ReplaceText(requiredProvider.Expr.Range(), fmt.Sprintf(`{
						source  = "hashicorp/%s"
						version = %s
					}`, name, f.TextAt(requiredProvider.Expr.Range()).Bytes))
				},
			); err != nil {
				return err
			}

			// Check if we need to validate provider constraints for this legacy format
			if expectedProvider, hasConstraints := config.Providers[name]; hasConstraints {
				if err := r.validateProvider(runner, name, expectedProvider, requiredProvider, val); err != nil {
					return err
				}
			}

			continue
		}

		vm := val.AsValueMap()

		if source, exists := vm["source"]; exists {
			p, err := tfaddr.ParseProviderSource(source.AsString())
			if err != nil {
				return err
			}

			if p.IsBuiltIn() {
				continue
			}
		} else if *config.Source {
			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf("Missing `source` for provider %q in `required_providers`", name),
				requiredProvider.Expr.Range(),
				func(f tflint.Fixer) error {
					if tfsdk.IsJSONFilename(requiredProvider.Expr.Range().Filename) {
						return tflint.ErrFixNotSupported
					}

					kvs, diags := hcl.ExprMap(requiredProvider.Expr)
					if diags.HasErrors() {
						return diags
					}
					if len(kvs) == 0 {
						return f.ReplaceText(requiredProvider.Expr.Range(), fmt.Sprintf(`{
							source = "hashicorp/%s"
						}`, name))
					}
					return f.InsertTextBefore(kvs[0].Key.StartRange(), fmt.Sprintf(`source = "hashicorp/%s"`+"\n", name))
				},
			); err != nil {
				return err
			}
		}

		if _, exists := vm["version"]; !exists && *config.Version {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Missing version constraint for provider %q in `required_providers`", name),
				requiredProvider.Expr.Range(),
			); err != nil {
				return err
			}
		}

		// Validate provider constraints if configured
		if expectedProvider, hasConstraints := config.Providers[name]; hasConstraints {
			if err := r.validateProvider(runner, name, expectedProvider, requiredProvider, val); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateProvider checks a single provider against expected configuration
func (r *TerraformRequiredProvidersRule) validateProvider(runner tflint.Runner, name string, expected ProviderRequirement, actual *hclext.Attribute, val cty.Value) error {
	// Handle legacy string format
	if val.Type() == cty.String {
		actualVersion := val.AsString()
		// For legacy format, we can only validate version constraints
		return r.validateVersionConstraint(runner, name, expected, actualVersion, "", actual.Expr.Range())
	}

	// Handle object format
	if !val.Type().IsObjectType() {
		return nil
	}

	// Extract source and version from the value map
	vm := val.AsValueMap()
	actualSource, actualVersion := "", ""
	if source, exists := vm["source"]; exists && !source.IsNull() {
		actualSource = source.AsString()
	}
	if version, exists := vm["version"]; exists && !version.IsNull() {
		actualVersion = version.AsString()
	}

	// Validate source if expected
	if expected.Source != nil && *expected.Source != "" && actualSource != *expected.Source {
		message := fmt.Sprintf("Provider %q missing required source %q", name, *expected.Source)
		if actualSource != "" {
			message = fmt.Sprintf("Provider %q has incorrect source (expected: %q, found: %q)", name, *expected.Source, actualSource)
		}
		expectedVersion := ""
		if expected.Version != nil {
			expectedVersion = *expected.Version
		}
		expectedSource := ""
		if expected.Source != nil {
			expectedSource = *expected.Source
		}
		return runner.EmitIssueWithFix(
			r,
			message,
			actual.Expr.Range(),
			func(f tflint.Fixer) error {
				return r.fixProviderConstraint(f, name, expectedVersion, expectedSource, actual.Expr.Range())
			},
		)
	}

	// Validate version constraint
	return r.validateVersionConstraint(runner, name, expected, actualVersion, actualSource, actual.Expr.Range())
}

// validateVersionConstraint validates version constraints using structural equality
func (r *TerraformRequiredProvidersRule) validateVersionConstraint(runner tflint.Runner, name string, expected ProviderRequirement, actualVersion, actualSource string, rng hcl.Range) error {
	// Early returns for cases where validation isn't needed
	if expected.Version == nil || *expected.Version == "" || actualVersion == "" {
		return nil
	}

	// Check structural equality using go-version
	expectedConstraint, err := version.NewConstraint(*expected.Version)
	if err != nil {
		return fmt.Errorf("invalid expected constraint %q: %w", *expected.Version, err)
	}

	actualConstraint, err := version.NewConstraint(actualVersion)
	if err != nil {
		// Invalid actual constraint - report as mismatch
		return runner.EmitIssueWithFix(
			r,
			fmt.Sprintf("Provider %q has invalid version constraint %q", name, actualVersion),
			rng,
			func(f tflint.Fixer) error {
				return r.fixProviderConstraint(f, name, *expected.Version, actualSource, rng)
			},
		)
	}

	// Check if constraints are structurally equal (exact match, not just equivalent ranges)
	if !expectedConstraint.Equals(actualConstraint) {
		return runner.EmitIssueWithFix(
			r,
			fmt.Sprintf("Provider %q version constraint does not match expected (expected: %q, found: %q)", name, *expected.Version, actualVersion),
			rng,
			func(f tflint.Fixer) error {
				return r.fixProviderConstraint(f, name, *expected.Version, actualSource, rng)
			},
		)
	}

	return nil
}

// fixProviderConstraint provides autofix for provider constraints
func (r *TerraformRequiredProvidersRule) fixProviderConstraint(f tflint.Fixer, name, expectedVersion, source string, rng hcl.Range) error {
	if tfsdk.IsJSONFilename(rng.Filename) {
		return tflint.ErrFixNotSupported
	}

	// Handle legacy string format (just version)
	existingText := f.TextAt(rng)
	isLegacyFormat := len(existingText.Bytes) > 0 && !strings.Contains(string(existingText.Bytes), "{")

	if isLegacyFormat {
		if expectedVersion != "" {
			return f.ReplaceText(rng, fmt.Sprintf("%q", expectedVersion))
		}
		if source != "" {
			return f.ReplaceText(rng, fmt.Sprintf(`{
      source  = %q
      version = %s
    }`, source, existingText.Bytes))
		}
	}

	// For object format, we need to parse the existing value to preserve unchanged fields
	// When only fixing source, preserve existing version
	// When only fixing version, preserve existing source

	// Build replacement for object format
	if source == "" && expectedVersion == "" {
		return fmt.Errorf("no source or version to fix")
	}

	// Try to extract existing values from the current text
	currentText := string(existingText.Bytes)
	hasVersion := strings.Contains(currentText, "version")

	// Extract current version if it exists
	var currentVersion string
	if hasVersion && expectedVersion == "" {
		// Try to extract version from current text
		versionMatch := strings.Index(currentText, `version = "`)
		if versionMatch != -1 {
			versionStart := versionMatch + len(`version = "`)
			versionEnd := strings.Index(currentText[versionStart:], `"`)
			if versionEnd != -1 {
				currentVersion = currentText[versionStart:versionStart+versionEnd]
			}
		}
	}

	replacement := ""
	if source != "" && expectedVersion != "" {
		replacement = fmt.Sprintf(`{
      source  = %q
      version = %q
    }`, source, expectedVersion)
	} else if source != "" && currentVersion != "" {
		// Fixing only source, preserve current version
		replacement = fmt.Sprintf(`{
      source  = %q
      version = %q
    }`, source, currentVersion)
	} else if source != "" {
		replacement = fmt.Sprintf(`{
      source = %q
    }`, source)
	} else if expectedVersion != "" {
		replacement = fmt.Sprintf(`{
      version = %q
    }`, expectedVersion)
	}

	return f.ReplaceText(rng, replacement)
}
