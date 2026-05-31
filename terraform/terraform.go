package terraform

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/zclconf/go-cty/cty"
)

// ModuleCall represents a "module" block.
type ModuleCall struct {
	Name         string
	DefRange     hcl.Range
	Source       string
	SourceKnown  bool
	SourceAttr   *hclext.Attribute
	Version      version.Constraints
	VersionKnown bool
	VersionAttr  *hclext.Attribute
}

func decodeModuleCall(runner *Runner, block *hclext.Block) (*ModuleCall, hcl.Diagnostics) {
	module := &ModuleCall{
		Name:     block.Labels[0],
		DefRange: block.DefRange,
	}
	diags := hcl.Diagnostics{}

	if source, exists := block.Body.Attributes["source"]; exists {
		module.SourceAttr = source

		sourceVal, sourceKnown, sourceNull, sourceDiags := evalModuleAttribute(runner, source.Expr)
		module.Source = sourceVal
		module.SourceKnown = sourceKnown
		if sourceNull {
			module.SourceAttr = nil
		}
		diags = diags.Extend(sourceDiags)
	} else {
		module.SourceKnown = true
	}

	if versionAttr, exists := block.Body.Attributes["version"]; exists {
		module.VersionAttr = versionAttr

		versionVal, versionKnown, versionNull, versionDiags := evalModuleAttribute(runner, versionAttr.Expr)
		diags = diags.Extend(versionDiags)
		if diags.HasErrors() {
			return module, diags
		}

		if !versionKnown {
			return module, diags
		}
		module.VersionKnown = true

		if versionNull {
			module.VersionAttr = nil
			return module, diags
		}

		constraints, err := version.NewConstraint(versionVal)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid version constraint",
				Detail:   "This string does not use correct version constraint syntax.",
				Subject:  versionAttr.Expr.Range().Ptr(),
			})
		}
		module.Version = constraints
	} else {
		module.VersionKnown = true
	}

	return module, diags
}

func evalModuleAttribute(runner *Runner, expr hcl.Expression) (val string, known bool, null bool, diags hcl.Diagnostics) {
	var ret cty.Value
	err := runner.EvaluateExpr(expr, &ret, nil)
	if err != nil {
		return "", false, false, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to evaluate expression",
			Detail:   err.Error(),
			Subject:  expr.Range().Ptr(),
		}}
	}

	// sensitive values are treated as unknown
	if !ret.IsKnown() || ret.IsMarked() {
		return "", false, false, nil
	}
	if ret.IsNull() {
		return "", true, true, nil
	}
	if ret.Type() != cty.String {
		return "", true, false, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Expected a string",
			Detail:   "The expression should evaluate to a string.",
			Subject:  expr.Range().Ptr(),
		}}
	}

	return ret.AsString(), true, false, nil
}

// Local represents a single entry from a "locals" block.
type Local struct {
	Name      string
	Attribute *hcl.Attribute
	DefRange  hcl.Range
}

// ProviderRef represents a reference to a provider like `provider = google.europe` in a resource or module.
type ProviderRef struct {
	Name     string
	DefRange hcl.Range
}

// @see https://github.com/hashicorp/terraform/blob/v1.2.7/internal/configs/resource.go#L624-L695
func decodeProviderRef(expr hcl.Expression, defRange hcl.Range) (*ProviderRef, hcl.Diagnostics) {
	expr, diags := shimTraversalInString(expr)
	if diags.HasErrors() {
		return nil, diags
	}

	traversal, diags := hcl.AbsTraversalForExpr(expr)
	if diags.HasErrors() {
		return nil, diags
	}

	return &ProviderRef{
		Name:     traversal.RootName(),
		DefRange: defRange,
	}, nil
}

// @see https://github.com/hashicorp/terraform/blob/v1.2.5/internal/configs/compat_shim.go#L34
func shimTraversalInString(expr hcl.Expression) (hcl.Expression, hcl.Diagnostics) {
	// ObjectConsKeyExpr is a special wrapper type used for keys on object
	// constructors to deal with the fact that naked identifiers are normally
	// handled as "bareword" strings rather than as variable references. Since
	// we know we're interpreting as a traversal anyway (and thus it won't
	// matter whether it's a string or an identifier) we can safely just unwrap
	// here and then process whatever we find inside as normal.
	if ocke, ok := expr.(*hclsyntax.ObjectConsKeyExpr); ok {
		expr = ocke.Wrapped
	}

	if _, ok := expr.(*hclsyntax.TemplateExpr); !ok {
		return expr, nil
	}

	strVal, diags := expr.Value(nil)
	if diags.HasErrors() || strVal.IsNull() || !strVal.IsKnown() {
		// Since we're not even able to attempt a shim here, we'll discard
		// the diagnostics we saw so far and let the caller's own error
		// handling take care of reporting the invalid expression.
		return expr, nil
	}

	// The position handling here isn't _quite_ right because it won't
	// take into account any escape sequences in the literal string, but
	// it should be close enough for any error reporting to make sense.
	srcRange := expr.Range()
	startPos := srcRange.Start // copy
	startPos.Column++          // skip initial quote
	startPos.Byte++            // skip initial quote

	traversal, tDiags := hclsyntax.ParseTraversalAbs(
		[]byte(strVal.AsString()),
		srcRange.Filename,
		startPos,
	)
	diags = append(diags, tDiags...)

	return &hclsyntax.ScopeTraversalExpr{
		Traversal: traversal,
		SrcRange:  srcRange,
	}, diags
}
