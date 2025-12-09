package opentofu

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	regaddr "github.com/opentofu/registry-address"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/zclconf/go-cty/cty"
)

// ModuleCall represents a "module" block.
type ModuleCall struct {
	Name        string
	DefRange    hcl.Range
	Source      string
	SourceAttr  *hclext.Attribute
	Version     version.Constraints
	VersionAttr *hclext.Attribute
}

// @see https://github.com/hashicorp/terraform/blob/v1.2.7/internal/configs/module_call.go#L36-L224
func decodeModuleCall(block *hclext.Block) (*ModuleCall, hcl.Diagnostics) {
	module := &ModuleCall{
		Name:     block.Labels[0],
		DefRange: block.DefRange,
	}
	diags := hcl.Diagnostics{}

	if source, exists := block.Body.Attributes["source"]; exists {
		module.SourceAttr = source
		sourceDiags := gohcl.DecodeExpression(source.Expr, nil, &module.Source)
		diags = diags.Extend(sourceDiags)
	}

	if versionAttr, exists := block.Body.Attributes["version"]; exists {
		module.VersionAttr = versionAttr

		var versionVal string
		versionDiags := gohcl.DecodeExpression(versionAttr.Expr, nil, &versionVal)
		diags = diags.Extend(versionDiags)
		if diags.HasErrors() {
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
	}

	return module, diags
}

// Local represents a single entry from a "locals" block.
type Local struct {
	Name      string
	Attribute *hcl.Attribute
	DefRange  hcl.Range
}

// ProviderRef represents a reference to a provider like `provider = google.europe` in a resource or module.
type ProviderRef struct {
	Name          string
	Alias         string
	AliasRange    *hcl.Range
	KeyExpression hcl.Expression
	DefRange      hcl.Range
}

func ExprIsNativeQuotedString(expr hcl.Expression) bool {
	_, ok := expr.(*hclsyntax.TemplateExpr)
	return ok
}

func IsProviderPartNormalized(str string) (bool, error) {
	normalized, err := regaddr.ParseProviderPart(str)
	if err != nil {
		return false, err
	}
	if str == normalized {
		return true, nil
	}
	return false, nil
}

func checkProviderNameNormalized(name string, declrange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics
	// verify that the provider local name is normalized
	normalized, err := IsProviderPartNormalized(name)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider local name",
			Detail:   fmt.Sprintf("%s is an invalid provider local name: %s", name, err),
			Subject:  &declrange,
		})
		return diags
	}
	if !normalized {
		// we would have returned this error already
		normalizedProvider, _ := regaddr.ParseProviderPart(name)
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider local name",
			Detail:   fmt.Sprintf("Provider names must be normalized. Replace %q with %q to fix this error.", name, normalizedProvider),
			Subject:  &declrange,
		})
	}
	return diags
}

func ConvertJSONExpressionToHCL(expr hcl.Expression) (hcl.Expression, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	// We can abuse the hcl json api and rely on the fact that calling
	// Value on a json expression with no EvalContext will return the
	// raw string. We can then parse that as normal hcl syntax, and
	// continue with the decoding.
	value, ds := expr.Value(nil)
	diags = append(diags, ds...)
	if diags.HasErrors() {
		return nil, diags
	}

	if value.Type() != cty.String || value.IsNull() {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected string expression",
			Detail:   fmt.Sprintf("This value must be a string, but got %s.", value.Type().FriendlyName()),
			Subject:  expr.Range().Ptr(),
		})
		return nil, diags
	}

	expr, ds = hclsyntax.ParseExpression([]byte(value.AsString()), expr.Range().Filename, expr.Range().Start)
	diags = append(diags, ds...)
	if diags.HasErrors() {
		return nil, diags
	}

	return expr, diags
}

// @see https://github.com/opentofu/opentofu/blob/3258c673194ecba26d856fb825d4eb4a7e36ab34/internal/configs/resource.go#L903
func decodeProviderRef(expr hcl.Expression, defRange hcl.Range) (*ProviderRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	var keyExpr hcl.Expression
	const (
		// name.alias[const_key]
		nameIndex  = 0
		aliasIndex = 1
		keyIndex   = 2
	)

	var maxTraversalLength = keyIndex + 1

	if ok := hcljson.IsJSONExpression(expr); ok {
		expr, diags = ConvertJSONExpressionToHCL(expr)
		if diags.HasErrors() {
			return nil, diags
		}
	}

	// name.alias[expr_key]
	if iex, ok := expr.(*hclsyntax.IndexExpr); ok {
		maxTraversalLength = aliasIndex + 1 // expr key found, no const key allowed

		keyExpr = iex.Key
		expr = iex.Collection
	}

	var shimDiags hcl.Diagnostics
	expr, shimDiags = shimTraversalInString(expr)
	diags = append(diags, shimDiags...)

	traversal, travDiags := hcl.AbsTraversalForExpr(expr)

	// AbsTraversalForExpr produces only generic errors, so we'll discard
	// the errors given and produce our own with extra context. If we didn't
	// get any errors then we might still have warnings, though.
	if !travDiags.HasErrors() {
		diags = append(diags, travDiags...)
	}

	if len(traversal) == 0 || len(traversal) > maxTraversalLength {
		// A provider reference was given as a string literal in the legacy
		// configuration language and there are lots of examples out there
		// showing that usage, so we'll sniff for that situation here and
		// produce a specialized error message for it to help users find
		// the new correct form.
		if ExprIsNativeQuotedString(expr) {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider configuration reference",
				Detail:   "A provider configuration reference must not be given in quotes.",
				Subject:  expr.Range().Ptr(),
			})
			return nil, diags
		}

		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider configuration reference",
			Detail:   fmt.Sprintf("The provider argument requires a provider type name, optionally followed by a period and then a configuration alias and optional instance key."),
			Subject:  expr.Range().Ptr(),
		})
		return nil, diags
	}

	// verify that the provider local name is normalized
	name := traversal.RootName()
	nameDiags := checkProviderNameNormalized(name, traversal[nameIndex].SourceRange())
	diags = append(diags, nameDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	ret := &ProviderRef{
		Name:          traversal.RootName(),
		DefRange:      defRange,
		KeyExpression: keyExpr,
		AliasRange:    nil,
	}

	if len(traversal) > aliasIndex {
		aliasStep, ok := traversal[aliasIndex].(hcl.TraverseAttr)
		if !ok {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider configuration reference",
				Detail:   "Provider name must either stand alone or be followed by a period and then a configuration alias.",
				Subject:  traversal[aliasIndex].SourceRange().Ptr(),
			})
			return ret, diags
		}

		ret.Alias = aliasStep.Name
		ret.AliasRange = aliasStep.SourceRange().Ptr()
	}

	if len(traversal) > keyIndex {
		indexStep, ok := traversal[keyIndex].(hcl.TraverseIndex)
		if !ok {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider configuration reference",
				Detail:   "Provider name must either stand alone or be followed by a period and then a configuration alias.",
				Subject:  traversal[keyIndex].SourceRange().Ptr(),
			})
			return ret, diags
		}

		ret.KeyExpression = hcl.StaticExpr(indexStep.Key, traversal.SourceRange())
	}

	if len(ret.Alias) == 0 && ret.KeyExpression != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider configuration reference",
			Detail:   "Provider assignment requires an alias when specifying an instance key, in the form of provider.name[instance_key]",
			Subject:  traversal.SourceRange().Ptr(),
		})
	}

	return ret, nil
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
