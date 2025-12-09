package opentofu

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/json"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// Runner is a custom runner that provides helper functions for this ruleset.
type Runner struct {
	tflint.Runner
}

// NewRunner returns a new custom runner.
func NewRunner(runner tflint.Runner) *Runner {
	return &Runner{Runner: runner}
}

// GetModuleCalls returns all "module" blocks, including uncreated module calls.
func (r *Runner) GetModuleCalls() ([]*ModuleCall, hcl.Diagnostics) {
	calls := []*ModuleCall{}
	diags := hcl.Diagnostics{}

	body, err := r.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "source"},
						{Name: "version"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return calls, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "failed to call GetModuleContent()",
				Detail:   err.Error(),
			},
		}
	}

	for _, block := range body.Blocks {
		call, decodeDiags := decodeModuleCall(block)
		diags = diags.Extend(decodeDiags)
		if decodeDiags.HasErrors() {
			continue
		}
		calls = append(calls, call)
	}

	return calls, diags
}

// GetLocals returns all entries in "locals" blocks.
func (r *Runner) GetLocals() (map[string]*Local, hcl.Diagnostics) {
	locals := map[string]*Local{}
	diags := hcl.Diagnostics{}

	files, err := r.GetFiles()
	if err != nil {
		return locals, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "failed to call GetFiles()",
				Detail:   err.Error(),
			},
		}
	}

	for _, file := range files {
		content, _, schemaDiags := file.Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{{Type: "locals"}},
		})
		diags = diags.Extend(schemaDiags)
		if schemaDiags.HasErrors() {
			continue
		}

		for _, block := range content.Blocks {
			attrs, localsDiags := block.Body.JustAttributes()
			diags = diags.Extend(localsDiags)
			if localsDiags.HasErrors() {
				continue
			}

			for name, attr := range attrs {
				locals[name] = &Local{
					Name:      attr.Name,
					Attribute: attr,
					DefRange:  attr.Range,
				}
			}
		}
	}

	return locals, diags
}

// GetProviderRefs returns all references to providers in resources, data, provider declarations, module calls, and provider-defined functinos.
func (r *Runner) GetProviderRefs() (map[string]*ProviderRef, hcl.Diagnostics) {
	providerRefs := map[string]*ProviderRef{}

	body, err := r.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "provider"},
					},
				},
			},
			{
				Type:       "ephemeral",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "provider"},
					},
				},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "provider"},
					},
				},
			},
			{
				Type:       "provider",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "providers"},
					},
				},
			},
			{
				Type:       "check",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type:       "data",
							LabelNames: []string{"type", "name"},
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "provider"},
								},
							},
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return providerRefs, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "failed to call `GetModuleContent()`",
				Detail:   err.Error(),
			},
		}
	}

	var diags hcl.Diagnostics
	for _, block := range body.Blocks {
		switch block.Type {
		case "resource", "ephemeral", "data":
			if attr, exists := block.Body.Attributes["provider"]; exists {
				ref, decodeDiags := decodeProviderRef(attr.Expr, block.DefRange)
				diags = diags.Extend(decodeDiags)
				if decodeDiags.HasErrors() {
					continue
				}
				providerRefs[ref.Name] = ref
			} else {
				providerName := block.Labels[0]
				if under := strings.Index(providerName, "_"); under != -1 {
					providerName = providerName[:under]
				}
				providerRefs[providerName] = &ProviderRef{
					Name:     providerName,
					DefRange: block.DefRange,
				}
			}
		case "provider":
			providerRefs[block.Labels[0]] = &ProviderRef{
				Name:     block.Labels[0],
				DefRange: block.DefRange,
			}
		case "module":
			if attr, exists := block.Body.Attributes["providers"]; exists {
				pairs, mapDiags := hcl.ExprMap(attr.Expr)
				diags = diags.Extend(mapDiags)
				if mapDiags.HasErrors() {
					continue
				}

				for _, pair := range pairs {
					ref, decodeDiags := decodeProviderRef(pair.Value, block.DefRange)
					diags = diags.Extend(decodeDiags)
					if decodeDiags.HasErrors() {
						continue
					}
					providerRefs[ref.Name] = ref
				}
			}
		case "check":
			for _, data := range block.Body.Blocks {
				if attr, exists := data.Body.Attributes["provider"]; exists {
					ref, decodeDiags := decodeProviderRef(attr.Expr, data.DefRange)
					diags = diags.Extend(decodeDiags)
					if decodeDiags.HasErrors() {
						continue
					}
					providerRefs[ref.Name] = ref
				} else {
					providerName := data.Labels[0]
					if under := strings.Index(providerName, "_"); under != -1 {
						providerName = providerName[:under]
					}
					providerRefs[providerName] = &ProviderRef{
						Name:     providerName,
						DefRange: data.DefRange,
					}
				}
			}
		default:
			panic("unreachable")
		}
	}

	walkDiags := r.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		// For JSON syntax, walker is not implemented,
		// so extract the hclsyntax.Node that we can walk on.
		// See https://github.com/hashicorp/hcl/issues/543
		nodes, diags := r.walkableNodesInExpr(expr)

		for _, node := range nodes {
			visitDiags := hclsyntax.VisitAll(node, func(n hclsyntax.Node) hcl.Diagnostics {
				if funcCallExpr, ok := n.(*hclsyntax.FunctionCallExpr); ok {
					parts := strings.Split(funcCallExpr.Name, "::")
					if len(parts) < 2 || parts[0] != "provider" || parts[1] == "" {
						return nil
					}
					providerRefs[parts[1]] = &ProviderRef{
						Name:     parts[1],
						DefRange: funcCallExpr.Range(),
					}
				}
				return nil
			})
			diags = diags.Extend(visitDiags)
		}
		return diags
	}))
	diags = diags.Extend(walkDiags)
	if walkDiags.HasErrors() {
		return providerRefs, diags
	}

	return providerRefs, diags
}

// walkableNodesInExpr returns hclsyntax.Node from the given expression.
// If the expression is an hclsyntax expression, it is returned as is.
// If the expression is a JSON expression, it is parsed and
// hclsyntax.Node it contains is returned.
func (r *Runner) walkableNodesInExpr(expr hcl.Expression) ([]hclsyntax.Node, hcl.Diagnostics) {
	nodes := []hclsyntax.Node{}

	expr = hcl.UnwrapExpressionUntil(expr, func(expr hcl.Expression) bool {
		_, native := expr.(hclsyntax.Expression)
		return native || json.IsJSONExpression(expr)
	})
	if expr == nil {
		return nil, nil
	}

	if json.IsJSONExpression(expr) {
		// HACK: For JSON expressions, we can get the JSON value as a literal
		//       without any prior HCL parsing by evaluating it in a nil context.
		//       We can take advantage of this property to walk through cty.Value
		//       that may contain HCL expressions instead of walking through
		//       expression nodes directly.
		//       See https://github.com/hashicorp/hcl/issues/642
		val, diags := expr.Value(nil)
		if diags.HasErrors() {
			return nodes, diags
		}

		err := cty.Walk(val, func(path cty.Path, v cty.Value) (bool, error) {
			if v.Type() != cty.String || v.IsNull() || !v.IsKnown() {
				return true, nil
			}

			node, parseDiags := hclsyntax.ParseTemplate([]byte(v.AsString()), expr.Range().Filename, expr.Range().Start)
			if diags.HasErrors() {
				diags = diags.Extend(parseDiags)
				return true, nil
			}

			nodes = append(nodes, node)
			return true, nil
		})
		if err != nil {
			return nodes, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to walk the expression value",
				Detail:   err.Error(),
				Subject:  expr.Range().Ptr(),
			}}
		}

		return nodes, diags
	}

	// The JSON syntax is already processed, so it's guaranteed to be native syntax.
	nodes = append(nodes, expr.(hclsyntax.Expression))

	return nodes, nil
}
