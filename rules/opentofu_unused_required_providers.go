package rules

import (
	"fmt"

	"github.com/diofeher/tflint-ruleset-opentofu/opentofu"
	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformUnusedRequiredProvidersRule checks whether required providers are used in the module
type TerraformUnusedRequiredProvidersRule struct {
	tflint.DefaultRule
}

// NewTerraformUnusedRequiredProvidersRule returns new rule with default attributes
func NewTerraformUnusedRequiredProvidersRule() *TerraformUnusedRequiredProvidersRule {
	return &TerraformUnusedRequiredProvidersRule{}
}

// Name returns the rule name
func (r *TerraformUnusedRequiredProvidersRule) Name() string {
	return "terraform_unused_required_providers"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformUnusedRequiredProvidersRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformUnusedRequiredProvidersRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformUnusedRequiredProvidersRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether required providers are used
func (r *TerraformUnusedRequiredProvidersRule) Check(rr tflint.Runner) error {
	runner := rr.(*opentofu.Runner)

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	providerRefs, diags := runner.GetProviderRefs()
	if diags.HasErrors() {
		return diags
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	requiredProviders := hcl.Attributes{}
	for _, file := range files {
		content, _, schemaDiags := file.Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{{Type: "terraform"}},
		})
		diags = diags.Extend(schemaDiags)
		if diags.HasErrors() {
			continue
		}

		for _, block := range content.Blocks {
			content, _, schemaDiags = block.Body.PartialContent(&hcl.BodySchema{
				Blocks: []hcl.BlockHeaderSchema{{Type: "required_providers"}},
			})
			diags = diags.Extend(schemaDiags)
			if diags.HasErrors() {
				continue
			}

			for _, block := range content.Blocks {
				attributes, attrDiags := block.Body.JustAttributes()
				for k, v := range attributes {
					requiredProviders[k] = v
				}
				diags = diags.Extend(attrDiags)
				if diags.HasErrors() {
					continue
				}
			}
		}
	}
	if diags.HasErrors() {
		return diags
	}

	for _, required := range requiredProviders {
		if _, exists := providerRefs[required.Name]; !exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("provider '%s' is declared in required_providers but not used by the module", required.Name),
				required.Range,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
