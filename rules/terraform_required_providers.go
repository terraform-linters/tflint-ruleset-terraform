package rules

import (
	"fmt"

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

type terraformRequiredProvidersRuleConfig struct {
	// Source specifies whether the rule should assert the presence of a `source` attribute
	Source *bool `hclext:"source,optional"`
	// Version specifies whether the rule should assert the presence of a `version` attribute
	Version *bool `hclext:"version,optional"`
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

		// required_providers entries are objects like `{ source = "...", version = "..." }`.
		// Each attribute is evaluated individually instead of evaluating the whole object,
		// because configuration_aliases contains provider references (e.g. [foo.a, foo.b])
		// that cannot be evaluated as expressions.
		// @see https://www.terraform.io/language/modules/develop/providers#provider-aliases-within-modules
		kvs, diags := hcl.ExprMap(requiredProvider.Expr)
		if diags.HasErrors() {
			// Not an object. Check whether this is the legacy version constraint string.
			var val cty.Value
			if err := runner.EvaluateExpr(requiredProvider.Expr, &val, nil); err != nil {
				return err
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
						// Autofix that replaces the whole expression is only supported for literal expressions
						// to avoid rewriting expressions that reference variables.
						if len(requiredProvider.Expr.Variables()) > 0 {
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
			}

			continue
		}

		attrs := map[string]hcl.Expression{}
		for _, kv := range kvs {
			// Keys are static (naked identifiers or string literals)
			key, diags := kv.Key.Value(nil)
			if diags.HasErrors() || !key.IsKnown() || key.IsNull() || key.Type() != cty.String {
				continue
			}
			attrs[key.AsString()] = kv.Value
		}

		if source, exists := attrs["source"]; exists {
			// Since Terraform v1.17, the source can contain variables (e.g. "app.terraform.io/${var.org}/foo").
			// Unknown, null, and sensitive values are ignored because the callback is not called for them.
			var builtin bool
			if err := runner.EvaluateExpr(source, func(source string) error {
				p, err := tfaddr.ParseProviderSource(source)
				if err != nil {
					return err
				}
				builtin = p.IsBuiltIn()
				return nil
			}, nil); err != nil {
				return err
			}

			if builtin {
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

		if _, exists := attrs["version"]; !exists && *config.Version {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Missing version constraint for provider %q in `required_providers`", name),
				requiredProvider.Expr.Range(),
			); err != nil {
				return err
			}
		}
	}

	return nil
}
