package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/zclconf/go-cty/cty"
)

// TerraformStaticAttributeNotationRule checks if dot notation is used in static contexts
type TerraformStaticAttributeNotationRule struct {
	tflint.DefaultRule
}

// NewTerraformStaticAttributeNotationRule returns a new rule
func NewTerraformStaticAttributeNotationRule() *TerraformStaticAttributeNotationRule {
	return &TerraformStaticAttributeNotationRule{}
}

// Name returns the rule name
func (r *TerraformStaticAttributeNotationRule) Name() string {
	return "terraform_static_attribute_notation"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformStaticAttributeNotationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformStaticAttributeNotationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformStaticAttributeNotationRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check walks all expressions and emit issues if bracket notation is found in static contexts
func (r *TerraformStaticAttributeNotationRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules
		return nil
	}

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		if !isStaticContext(expr) {
			return nil
		}

		if attr, ok := expr.(*hclsyntax.ScopeTraversalExpr); ok {
			hasStaticIndex := false
			for _, traverser := range attr.Traversal {
				if t, ok := traverser.(hcl.TraverseIndex); ok {
					// Check if the index is a string literal
					if t.Key.Type().IsPrimitiveType() && t.Key.Type().Equals(cty.String) {
						hasStaticIndex = true
						break
					}
				}
			}

			if hasStaticIndex {
				if err := runner.EmitIssueWithFix(
					r,
					"Must use dot notation for static attributes",
					expr.Range(),
					func(f tflint.Fixer) error {
						var result string
						if root, ok := attr.Traversal[0].(hcl.TraverseRoot); ok {
							result = root.Name
						}
						for i := 1; i < len(attr.Traversal); i++ {
							if trav, ok := attr.Traversal[i].(hcl.TraverseAttr); ok {
								result += "." + trav.Name
							} else if trav, ok := attr.Traversal[i].(hcl.TraverseIndex); ok {
								if trav.Key.Type().IsPrimitiveType() && trav.Key.Type().Equals(cty.String) {
									result += "." + trav.Key.AsString()
								} else {
									result += "[" + trav.Key.AsString() + "]"
								}
							}
						}
						return f.ReplaceText(expr.Range(), result)
					},
				); err != nil {
					return hcl.Diagnostics{
						{
							Severity: hcl.DiagError,
							Summary:  "failed to call EmitIssueWithFix()",
							Detail:   err.Error(),
						},
					}
				}
			}
		}
		return nil
	}))

	if diags.HasErrors() {
		return diags
	}
	return nil
}

func isStaticContext(expr hcl.Expression) bool {
	// Check if we're inside a dynamic context
	vars := expr.Variables()
	for _, v := range vars {
		if len(v) > 0 {
			root, ok := v[0].(hcl.TraverseRoot)
			if ok && (root.Name == "each" || root.Name == "count") {
				return false
			}
		}
	}

	// Check parent context
	switch expr.(type) {
	case *hclsyntax.ForExpr, *hclsyntax.SplatExpr:
		return false
	}

	return true
}
