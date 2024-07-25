package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/zclconf/go-cty/cty"
)

// This rule checks for map literals with duplicate keys
type TerraformMapDuplicateKeysRule struct {
	tflint.DefaultRule
}

func NewTerraformMapDuplicateKeysRule() *TerraformMapDuplicateKeysRule {
	return &TerraformMapDuplicateKeysRule{}
}

func (r *TerraformMapDuplicateKeysRule) Name() string {
	return "terraform_map_duplicate_keys"
}

func (r *TerraformMapDuplicateKeysRule) Enabled() bool {
	return true
}

func (r *TerraformMapDuplicateKeysRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformMapDuplicateKeysRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *TerraformMapDuplicateKeysRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules
		return nil
	}

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(e hcl.Expression) hcl.Diagnostics {
		return r.checkObjectConsExpr(e, runner)
	}))
	if diags.HasErrors() {
		return diags
	}

	return nil
}

func (r *TerraformMapDuplicateKeysRule) checkObjectConsExpr(e hcl.Expression, runner tflint.Runner) hcl.Diagnostics {
	objExpr, ok := e.(*hclsyntax.ObjectConsExpr)
	if !ok {
		return nil
	}

	var diags hcl.Diagnostics
	keys := make(map[string]hcl.Range)

	for _, item := range objExpr.Items {
		expr := item.KeyExpr.(*hclsyntax.ObjectConsKeyExpr)
		var val cty.Value

		err := runner.EvaluateExpr(expr, &val, nil)
		if err != nil {
			logger.Debug("Failed to evaluate an expression, continuing", "error", err)

			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "failed to evaluate expression",
				Detail:   err.Error(),
			})
			continue
		}

		if !val.IsKnown() || val.IsNull() {
			// When trying to evaluate an expression
			// with a variable without a value,
			// runner.evaluateExpr() returns a null value.
			// Ignore this case since there's nothing we can do.
			logger.Debug("Unknown key, continuing", "range", expr.Range())
			continue
		}

		if declRange, exists := keys[val.AsString()]; exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Duplicate key: %q, first defined at %s", val.AsString(), declRange),
				expr.Range(),
			); err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "failed to call EmitIssue()",
					Detail:   err.Error(),
				})

				return diags
			}

			continue
		}

		keys[val.AsString()] = expr.Range()
	}

	return diags
}
