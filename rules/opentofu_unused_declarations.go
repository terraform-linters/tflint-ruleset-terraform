package rules

import (
	"fmt"

	"github.com/diofeher/tflint-ruleset-opentofu/opentofu"
	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/addrs"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/lang"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformUnusedDeclarationsRule checks whether variables, data sources, or locals are declared but unused
type TerraformUnusedDeclarationsRule struct {
	tflint.DefaultRule
}

type declarations struct {
	Variables       map[string]*hclext.Block
	DataResources   map[string]*hclext.Block
	Locals          map[string]*opentofu.Local
	ProviderAliases map[string]*hclext.Block
}

// NewTerraformUnusedDeclarationsRule returns a new rule
func NewTerraformUnusedDeclarationsRule() *TerraformUnusedDeclarationsRule {
	return &TerraformUnusedDeclarationsRule{}
}

// Name returns the rule name
func (r *TerraformUnusedDeclarationsRule) Name() string {
	return "opentofu_unused_declarations"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformUnusedDeclarationsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformUnusedDeclarationsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformUnusedDeclarationsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check emits issues for any variables, locals, and data sources that are declared but not used
func (r *TerraformUnusedDeclarationsRule) Check(rr tflint.Runner) error {
	runner := rr.(*opentofu.Runner)

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	decl, err := r.declarations(runner)
	if err != nil {
		return err
	}
	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		r.checkForRefsInExpr(expr, decl)
		return nil
	}))
	if diags.HasErrors() {
		return diags
	}

	for _, variable := range decl.Variables {
		if err := runner.EmitIssueWithFix(
			r,
			fmt.Sprintf(`variable "%s" is declared but not used`, variable.Labels[0]),
			variable.DefRange,
			func(f tflint.Fixer) error { return f.RemoveExtBlock(variable) },
		); err != nil {
			return err
		}
	}
	for _, data := range decl.DataResources {
		if err := runner.EmitIssueWithFix(
			r,
			fmt.Sprintf(`data "%s" "%s" is declared but not used`, data.Labels[0], data.Labels[1]),
			data.DefRange,
			func(f tflint.Fixer) error { return f.RemoveExtBlock(data) },
		); err != nil {
			return err
		}
	}
	for _, local := range decl.Locals {
		if err := runner.EmitIssueWithFix(
			r,
			fmt.Sprintf(`local.%s is declared but not used`, local.Name),
			local.DefRange,
			func(f tflint.Fixer) error { return f.RemoveAttribute(local.Attribute) },
		); err != nil {
			return err
		}
	}
	for _, provider := range decl.ProviderAliases {
		aliasAttr := provider.Body.Attributes["alias"]
		var aliasName string
		if diags := gohcl.DecodeExpression(aliasAttr.Expr, nil, &aliasName); diags.HasErrors() {
			continue
		}
		if err := runner.EmitIssueWithFix(
			r,
			fmt.Sprintf(`provider "%s" with alias "%s" is declared but not used`, provider.Labels[0], aliasName),
			provider.DefRange,
			func(f tflint.Fixer) error { return f.RemoveExtBlock(provider) },
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformUnusedDeclarationsRule) declarations(runner *opentofu.Runner) (*declarations, error) {
	decl := &declarations{
		Variables:       map[string]*hclext.Block{},
		DataResources:   map[string]*hclext.Block{},
		Locals:          map[string]*opentofu.Local{},
		ProviderAliases: map[string]*hclext.Block{},
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: "validation",
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "condition"},
									{Name: "error_message"},
								},
							},
						},
					},
				},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "check",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type:       "data",
							LabelNames: []string{"type", "name"},
							Body:       &hclext.BodySchema{},
						},
					},
				},
			},
			{
				Type:       "provider",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "alias"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return decl, err
	}

	for _, block := range body.Blocks {
		switch block.Type {
		case "variable":
			decl.Variables[block.Labels[0]] = block
		case "data":
			decl.DataResources[fmt.Sprintf("data.%s.%s", block.Labels[0], block.Labels[1])] = block
		case "check":
			for _, data := range block.Body.Blocks {
				// Scoped data source addresses are unique in the module
				decl.DataResources[fmt.Sprintf("data.%s.%s", data.Labels[0], data.Labels[1])] = data
			}
		case "provider":
			// Only track providers with aliases
			if aliasAttr, exists := block.Body.Attributes["alias"]; exists {
				var aliasName string
				if diags := gohcl.DecodeExpression(aliasAttr.Expr, nil, &aliasName); !diags.HasErrors() {
					decl.ProviderAliases[fmt.Sprintf("%s.%s", block.Labels[0], aliasName)] = block
				}
			}
		default:
			panic("unreachable")
		}
	}

	locals, diags := runner.GetLocals()
	if diags.HasErrors() {
		return decl, diags
	}
	decl.Locals = locals

	return decl, nil
}

func (r *TerraformUnusedDeclarationsRule) checkForRefsInExpr(expr hcl.Expression, decl *declarations) {
	// Check for provider alias references (e.g., aws.west in provider = aws.west)
	if traversal, diags := hcl.AbsTraversalForExpr(expr); diags == nil && len(traversal) == 2 {
		// Provider aliases are referenced as <provider>.<alias> (2 parts)
		providerRef := fmt.Sprintf("%s.%s", traversal.RootName(), traversal[1].(hcl.TraverseAttr).Name)
		delete(decl.ProviderAliases, providerRef)
	}

ReferenceLoop:
	for _, ref := range lang.ReferencesInExpr(expr) {
		switch sub := ref.Subject.(type) {
		case addrs.InputVariable:
			// Input variables can refer to themselves as var.NAME inside validation blocks.
			// Do not mark such expressions as used, skip to next reference.
			if varBlock, exists := decl.Variables[sub.Name]; exists {
				for _, validationBlock := range varBlock.Body.Blocks {
					for _, attr := range validationBlock.Body.Attributes {
						if attr.Expr.Range().Overlaps(expr.Range()) {
							continue ReferenceLoop
						}
					}
				}
			}
			delete(decl.Variables, sub.Name)
		case addrs.LocalValue:
			delete(decl.Locals, sub.Name)
		case addrs.Resource:
			delete(decl.DataResources, sub.String())
		case addrs.ResourceInstance:
			delete(decl.DataResources, sub.Resource.String())
		}
	}
}
