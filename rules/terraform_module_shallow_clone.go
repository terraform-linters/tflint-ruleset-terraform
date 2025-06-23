package rules

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/terraform-linters/tflint-ruleset-terraform/terraform"
)

var gitCommitRegex = regexp.MustCompile("^[a-f0-9]{40}$")

// TerraformModuleShallowCloneRule checks that Git-hosted Terraform modules use shallow cloning
type TerraformModuleShallowCloneRule struct {
	tflint.DefaultRule

	attributeName string
}

// NewTerraformModuleShallowCloneRule returns new rule with default attributes
func NewTerraformModuleShallowCloneRule() *TerraformModuleShallowCloneRule {
	return &TerraformModuleShallowCloneRule{
		attributeName: "source",
	}
}

// Name returns the rule name
func (r *TerraformModuleShallowCloneRule) Name() string {
	return "terraform_module_shallow_clone"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformModuleShallowCloneRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformModuleShallowCloneRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformModuleShallowCloneRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if Git-hosted Terraform modules use shallow cloning
func (r *TerraformModuleShallowCloneRule) Check(rr tflint.Runner) error {
	runner := rr.(*terraform.Runner)

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	calls, diags := runner.GetModuleCalls()
	if diags.HasErrors() {
		return diags
	}

	for _, call := range calls {
		if err := r.checkModule(runner, call); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformModuleShallowCloneRule) checkModule(runner tflint.Runner, module *terraform.ModuleCall) error {
	filename := module.DefRange.Filename
	source, err := getter.Detect(module.Source, filepath.Dir(filename), []getter.Detector{
		// https://github.com/hashicorp/terraform/blob/51b0aee36cc2145f45f5b04051a01eb6eb7be8bf/internal/getmodules/getter.go#L30-L52
		new(getter.GitHubDetector),
		new(getter.GitDetector),
		new(getter.BitBucketDetector),
		new(getter.GCSDetector),
		new(getter.S3Detector),
		new(getter.FileDetector),
	})
	if err != nil {
		return err
	}

	u, err := url.Parse(source)
	if err != nil {
		return err
	}

	// Only check Git-based sources
	if u.Scheme != "git" {
		return nil
	}

	if u.Opaque != "" {
		// for git:: pseudo-URLs, Opaque is :https, but query will still be parsed
		query := u.RawQuery
		u, err = url.Parse(strings.TrimPrefix(u.Opaque, ":"))
		if err != nil {
			return err
		}

		u.RawQuery = query
	}

	if u.Hostname() == "" {
		return nil
	}

	query := u.Query()

	// Check if module is pinned to a specific version
	ref := query.Get("ref")

	// Skip if not pinned at all
	if ref == "" {
		return nil
	}

	// Skip if it's a raw git commit ID (40 character hex string)
	if gitCommitRegex.MatchString(ref) {
		return nil
	}

	// Check if depth parameter is already set
	if query.Get("depth") == "1" {
		return nil
	}

	exprRange := module.SourceAttr.Expr.Range()

	if err := runner.EmitIssueWithFix(
		r,
		fmt.Sprintf(`Module source %q should enable shallow cloning by adding "depth=1" parameter`, module.Source),
		exprRange,
		func(f tflint.Fixer) error {
			// Find the position of "ref=" in the source string
			refPos := strings.Index(module.Source, "ref=")
			if refPos == -1 {
				return fmt.Errorf(`could not find "ref=" string in module source`)
			}

			// Create a range that includes the opening quote + source up to "ref="
			endPos := exprRange.Start
			// +1 for opening quote, +refPos for chars before "ref="
			endPos.Byte += 1 + refPos
			endPos.Column += 1 + refPos

			insertRange := hcl.Range{
				Filename: filename,
				Start:    exprRange.Start,
				End:      endPos,
			}

			return f.InsertTextAfter(insertRange, "depth=1&")
		},
	); err != nil {
		return hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "failed to call EmitIssueWithFix()",
				Detail:   err.Error(),
			},
		}
	}

	return nil
}
