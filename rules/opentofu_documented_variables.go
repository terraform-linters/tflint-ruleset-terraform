package rules

import (
	"fmt"

	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OpentofuDocumentedVariablesRule checks whether variables have descriptions
type OpentofuDocumentedVariablesRule struct {
	tflint.DefaultRule
}

// NewOpentofuDocumentedVariablesRule returns a new rule
func NewOpentofuDocumentedVariablesRule() *OpentofuDocumentedVariablesRule {
	return &OpentofuDocumentedVariablesRule{}
}

// Name returns the rule name
func (r *OpentofuDocumentedVariablesRule) Name() string {
	return "opentofu_documented_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *OpentofuDocumentedVariablesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *OpentofuDocumentedVariablesRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *OpentofuDocumentedVariablesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether variables have descriptions
func (r *OpentofuDocumentedVariablesRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: "description"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		attr, exists := variable.Body.Attributes["description"]
		if !exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("`%s` variable has no description", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		var description string
		diags := gohcl.DecodeExpression(attr.Expr, nil, &description)
		if diags.HasErrors() {
			return diags
		}

		if description == "" {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("`%s` variable has no description", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
