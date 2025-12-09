package rules

import (
	"strings"

	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OpentofuCommentSyntaxRule checks whether comments use the preferred syntax
type OpentofuCommentSyntaxRule struct {
	tflint.DefaultRule
}

// NewOpentofuCommentSyntaxRule returns a new rule
func NewOpentofuCommentSyntaxRule() *OpentofuCommentSyntaxRule {
	return &OpentofuCommentSyntaxRule{}
}

// Name returns the rule name
func (r *OpentofuCommentSyntaxRule) Name() string {
	return "opentofu_comment_syntax"
}

// Enabled returns whether the rule is enabled by default
func (r *OpentofuCommentSyntaxRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *OpentofuCommentSyntaxRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *OpentofuCommentSyntaxRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether single line comments is used
func (r *OpentofuCommentSyntaxRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for name, file := range files {
		if err := r.checkComments(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *OpentofuCommentSyntaxRule) checkComments(runner tflint.Runner, filename string, file *hcl.File) error {
	if strings.HasSuffix(filename, ".json") {
		return nil
	}

	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenComment {
			continue
		}

		if strings.HasPrefix(string(token.Bytes), "//") {
			if err := runner.EmitIssueWithFix(
				r,
				"Single line comments should begin with #",
				token.Range,
				func(f tflint.Fixer) error {
					return f.ReplaceText(f.RangeTo("//", filename, token.Range.Start), "#")
				},
			); err != nil {
				return err
			}
		}
	}

	return nil
}
