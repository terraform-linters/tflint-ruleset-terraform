package rules

import (
	stdjson "encoding/json"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/json"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformJSONSyntaxRule checks whether JSON configuration uses the official syntax
type TerraformJSONSyntaxRule struct {
	tflint.DefaultRule
}

// NewTerraformJSONSyntaxRule returns a new rule
func NewTerraformJSONSyntaxRule() *TerraformJSONSyntaxRule {
	return &TerraformJSONSyntaxRule{}
}

// Name returns the rule name
func (r *TerraformJSONSyntaxRule) Name() string {
	return "terraform_json_syntax"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformJSONSyntaxRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformJSONSyntaxRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformJSONSyntaxRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether JSON configurations use object syntax at root
func (r *TerraformJSONSyntaxRule) Check(runner tflint.Runner) error {
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
		if err := r.checkJSONSyntax(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformJSONSyntaxRule) checkJSONSyntax(runner tflint.Runner, filename string, file *hcl.File) error {
	if !strings.HasSuffix(filename, ".tf.json") {
		return nil
	}

	// Check if this is a JSON body
	if !json.IsJSONBody(file.Body) {
		return nil
	}

	// Unmarshal the file bytes to detect the root type
	var root any
	if err := stdjson.Unmarshal(file.Bytes, &root); err != nil {
		// If we can't parse it, skip (HCL will report the error)
		return nil
	}

	// Check if root is an array
	if _, isArray := root.([]any); isArray {
		if err := runner.EmitIssue(
			r,
			"JSON configuration uses array syntax at root, expected object",
			file.Body.MissingItemRange(),
		); err != nil {
			return err
		}
	}

	return nil
}
