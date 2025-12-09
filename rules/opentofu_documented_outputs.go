package rules

import (
	"fmt"

	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OpentofuDocumentedOutputsRule checks whether outputs have descriptions
type OpentofuDocumentedOutputsRule struct {
	tflint.DefaultRule
}

// NewOpentofuDocumentedOutputsRule returns a new rule
func NewOpentofuDocumentedOutputsRule() *OpentofuDocumentedOutputsRule {
	return &OpentofuDocumentedOutputsRule{}
}

// Name returns the rule name
func (r *OpentofuDocumentedOutputsRule) Name() string {
	return "opentofu_documented_outputs"
}

// Enabled returns whether the rule is enabled by default
func (r *OpentofuDocumentedOutputsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *OpentofuDocumentedOutputsRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *OpentofuDocumentedOutputsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether outputs have descriptions
func (r *OpentofuDocumentedOutputsRule) Check(runner tflint.Runner) error {
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
				Type:       "output",
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

	for _, output := range body.Blocks {
		attr, exists := output.Body.Attributes["description"]
		if !exists {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("`%s` output has no description", output.Labels[0]),
				output.DefRange,
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
				fmt.Sprintf("`%s` output has no description", output.Labels[0]),
				output.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
