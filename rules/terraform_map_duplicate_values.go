package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// This rule checks for map literals with duplicate values
type TerraformMapDuplicateValuesRule struct {
	tflint.DefaultRule
}

func NewTerraformMapDuplicateValuesRule() *TerraformMapDuplicateValuesRule {
	return &TerraformMapDuplicateValuesRule{}
}

func (r *TerraformMapDuplicateValuesRule) Name() string {
	return "terraform_map_duplicate_values"
}

func (r *TerraformMapDuplicateValuesRule) Enabled() bool {
	return true
}

func (r *TerraformMapDuplicateValuesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformMapDuplicateValuesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *TerraformMapDuplicateValuesRule) Check(runner tflint.Runner) error {
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

func (r *TerraformMapDuplicateValuesRule) checkObjectConsExpr(e hcl.Expression, runner tflint.Runner) hcl.Diagnostics {
	objExpr, ok := e.(*hclsyntax.ObjectConsExpr)
	if !ok {
		return nil
	}

	var diags hcl.Diagnostics
	values := make(map[string]hcl.Range)

	for _, item := range objExpr.Items {
		valExpr := item.ValueExpr
		var val cty.Value

		err := runner.EvaluateExpr(valExpr, &val, nil)
		if err != nil {
			logger.Debug("Failed to evaluate value. The value will be ignored", "range", valExpr.Range(), "error", err.Error())
			continue
		}

		if !val.IsKnown() || val.IsNull() || val.IsMarked() {
			logger.Debug("Unprocessable value, continuing", "range", valExpr.Range())
			continue
		}
		// Map values must be strings, but some values ​​can be converted to strings and become valid values,
		// so try to convert them here.
		if converted, err := convert.Convert(val, cty.String); err == nil {
			val = converted
		}

		// ignore unprocessable values and boolean values
		if val.Type() != cty.String || val.AsString() == "true" || val.AsString() == "false" {
			logger.Debug("Unprocessable value, continuing", "range", valExpr.Range())
			continue
		}

		if declRange, exists := values[val.AsString()]; exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Duplicate value: %q, first defined at %s", val.AsString(), declRange),
				valExpr.Range(),
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

		values[val.AsString()] = valExpr.Range()
	}

	return diags
}
