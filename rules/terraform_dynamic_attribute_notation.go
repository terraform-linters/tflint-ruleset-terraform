package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformDynamicAttributeNotationRule checks if bracket notation is used in dynamic contexts
type TerraformDynamicAttributeNotationRule struct {
	tflint.DefaultRule
}

// NewTerraformDynamicAttributeNotationRule returns a new rule
func NewTerraformDynamicAttributeNotationRule() *TerraformDynamicAttributeNotationRule {
	return &TerraformDynamicAttributeNotationRule{}
}

// Name returns the rule name
func (r *TerraformDynamicAttributeNotationRule) Name() string {
	return "terraform_dynamic_attribute_notation"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformDynamicAttributeNotationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformDynamicAttributeNotationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformDynamicAttributeNotationRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check walks all expressions and emit issues if dot notation is found in dynamic contexts
func (r *TerraformDynamicAttributeNotationRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules
		return nil
	}

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		if !isDynamicContext(expr) {
			return nil
		}

		if attr, ok := expr.(*hclsyntax.ScopeTraversalExpr); ok {
			// Skip if this is not a traversal with at least one attribute
			if len(attr.Traversal) < 2 {
				return nil
			}

			// Skip if already using valid bracket notation for dynamic attributes in 'each' context
			if len(attr.Traversal) >= 2 {
				if root, ok := attr.Traversal[0].(hcl.TraverseRoot); ok && root.Name == "each" {
					if token, ok := attr.Traversal[1].(hcl.TraverseAttr); ok && token.Name == "value" {
						valid := true
						for i := 2; i < len(attr.Traversal); i++ {
							if _, ok := attr.Traversal[i].(hcl.TraverseIndex); !ok {
								valid = false
								break
							}
						}
						if valid {
							return nil
						}
					}
				}
			}

			// Find the first dot notation attribute
			hasDotNotation := false
			for i := 1; i < len(attr.Traversal); i++ {
				if _, ok := attr.Traversal[i].(hcl.TraverseAttr); ok {
					// Skip if this is preceded by an index traversal
					if i > 1 {
						if _, ok := attr.Traversal[i-1].(hcl.TraverseIndex); ok {
							continue
						}
					}
					hasDotNotation = true
					break
				}
			}

			if !hasDotNotation {
				return nil
			}

			// Build the fixed expression
			result := ""
			if len(attr.Traversal) > 0 {
				if root, ok := attr.Traversal[0].(hcl.TraverseRoot); ok && root.Name == "each" {
					result = "each.value"
					for i := 1; i < len(attr.Traversal); i++ {
						if i == 1 {
							if attr, ok := attr.Traversal[i].(hcl.TraverseAttr); ok && attr.Name == "value" {
								continue
							}
						}
						switch token := attr.Traversal[i].(type) {
						case hcl.TraverseAttr:
							result += fmt.Sprintf("[\"%s\"]", token.Name)
						case hcl.TraverseIndex:
							result += fmt.Sprintf("[%s]", token.Key.AsString())
						}
					}
				} else {
					var useDot = true
					switch token := attr.Traversal[0].(type) {
					case hcl.TraverseRoot:
						result = token.Name
					case hcl.TraverseAttr:
						result = token.Name
					}
					for i := 1; i < len(attr.Traversal); i++ {
						switch token := attr.Traversal[i].(type) {
						case hcl.TraverseAttr:
							if useDot {
								result += "." + token.Name
							} else {
								result += fmt.Sprintf("[%c%s%c]", '"', token.Name, '"')
							}
						case hcl.TraverseIndex:
							useDot = false
							result += fmt.Sprintf("[%s]", token.Key.AsString())
						}
					}
				}
			}

			fmt.Printf("Transforming traversal: %#v\n", attr.Traversal)
			fmt.Printf("Initial result: %s\n", result)
			for i, token := range attr.Traversal {
				fmt.Printf("Token %d: %T = %+v\n", i, token, token)
			}

			if err := runner.EmitIssueWithFix(
				r,
				"Must use bracket notation [] for dynamic attributes",
				attr.Range(),
				func(f tflint.Fixer) error {
					return f.ReplaceText(attr.Range(), result)
				},
			); err != nil {
				return hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to emit issue",
						Detail:   err.Error(),
					},
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

func isDynamicContext(expr hcl.Expression) bool {
	// Check if we're inside a dynamic context
	vars := expr.Variables()
	for _, v := range vars {
		if len(v) > 0 {
			root, ok := v[0].(hcl.TraverseRoot)
			if ok && (root.Name == "each" || root.Name == "count") {
				return true
			}
		}
	}

	// Check parent context
	switch expr.(type) {
	case *hclsyntax.ForExpr, *hclsyntax.SplatExpr:
		return true
	}

	return false
}
