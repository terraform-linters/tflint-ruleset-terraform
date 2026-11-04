package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformNoShortCircuitEvaluationRule checks for attempts to use short-circuit evaluation
type TerraformNoShortCircuitEvaluationRule struct {
	tflint.DefaultRule
}

// NewTerraformNoShortCircuitEvaluationRule returns a new rule
func NewTerraformNoShortCircuitEvaluationRule() *TerraformNoShortCircuitEvaluationRule {
	return &TerraformNoShortCircuitEvaluationRule{}
}

// Name returns the rule name
func (r *TerraformNoShortCircuitEvaluationRule) Name() string {
	return "terraform_no_short_circuit_evaluation"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformNoShortCircuitEvaluationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformNoShortCircuitEvaluationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformNoShortCircuitEvaluationRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks for attempts to use short-circuit evaluation
func (r *TerraformNoShortCircuitEvaluationRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules
		return nil
	}

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		if binaryOpExpr, ok := expr.(*hclsyntax.BinaryOpExpr); ok {
			if binaryOpExpr.Op == hclsyntax.OpLogicalAnd || binaryOpExpr.Op == hclsyntax.OpLogicalOr {
				// Check if left side is a null check
				if isNullCheck(binaryOpExpr.LHS) {
					// Check if right side references the same variable as the left side
					if referencesNullCheckedVar(binaryOpExpr.LHS, binaryOpExpr.RHS) {
						if err := runner.EmitIssue(
							r,
							"Short-circuit evaluation is not supported in Terraform. Use a conditional expression (condition ? true : false) instead.",
							binaryOpExpr.Range(),
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

// isNullCheck determines if an expression is checking for null
func isNullCheck(expr hcl.Expression) bool {
	if binaryOpExpr, ok := expr.(*hclsyntax.BinaryOpExpr); ok {
		if binaryOpExpr.Op == hclsyntax.OpEqual || binaryOpExpr.Op == hclsyntax.OpNotEqual {
			// Check if either side is a null literal
			if isNullLiteral(binaryOpExpr.RHS) || isNullLiteral(binaryOpExpr.LHS) {
				return true
			}
		}
	}
	return false
}

// isNullLiteral checks if the expression is a null literal
func isNullLiteral(expr hcl.Expression) bool {
	if literalExpr, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		return literalExpr.Val.IsNull()
	}
	return false
}

// referencesNullCheckedVar checks if the right side expression references the same variable that was null checked
func referencesNullCheckedVar(nullCheck, expr hcl.Expression) bool {
	// Get the variable name from the null check
	var varName string
	if binaryOpExpr, ok := nullCheck.(*hclsyntax.BinaryOpExpr); ok {
		// Try to get variable name from LHS
		if scopeTraversalExpr, ok := binaryOpExpr.LHS.(*hclsyntax.ScopeTraversalExpr); ok {
			if len(scopeTraversalExpr.Traversal) > 0 {
				varName = scopeTraversalExpr.Traversal.RootName()
			}
		}
		// If not found in LHS, try RHS
		if varName == "" {
			if scopeTraversalExpr, ok := binaryOpExpr.RHS.(*hclsyntax.ScopeTraversalExpr); ok {
				if len(scopeTraversalExpr.Traversal) > 0 {
					varName = scopeTraversalExpr.Traversal.RootName()
				}
			}
		}
	}

	// If we couldn't find a variable name, return false
	if varName == "" {
		return false
	}

	// Check if the expression references the same variable
	vars := expr.Variables()
	for _, v := range vars {
		if len(v) > 0 && v.RootName() == varName {
			return true
		}
	}
	return false
} 