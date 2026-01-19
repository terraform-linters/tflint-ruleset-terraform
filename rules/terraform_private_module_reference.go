package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/terraform-linters/tflint-ruleset-terraform/terraform"
)

type StatFunc func(name string) (os.FileInfo, error)

// TerraformPrivateModuleReferenceRule checks whether private are referenced externally
type TerraformPrivateModuleReferenceRule struct {
	tflint.DefaultRule
	statFunc StatFunc
}

func NewTerraformPrivateModuleReferenceRule() *TerraformPrivateModuleReferenceRule {
	return &TerraformPrivateModuleReferenceRule{
		statFunc: os.Stat,
	}
}

func (r *TerraformPrivateModuleReferenceRule) Name() string {
	return "terraform_private_module_reference"
}

func (r *TerraformPrivateModuleReferenceRule) Enabled() bool {
	return true
}

func (r *TerraformPrivateModuleReferenceRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformPrivateModuleReferenceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *TerraformPrivateModuleReferenceRule) Check(rr tflint.Runner) error {
	runner := rr.(*terraform.Runner)

	moduleCalls, diags := runner.GetModuleCalls()
	if diags.HasErrors() {
		return diags
	}

	for _, call := range moduleCalls {
		// Get the current file path
		currentFile := call.DefRange.Filename

		// Get the module source path
		modulePath := call.Source

		// If modulePath is not a local path its a remote reference and we should not continue checking.
		if _, err := r.statFunc(modulePath); os.IsNotExist(err) {
			return nil
		}

		// Check if the module is referenced from outside the root
		isSubDir, err := isSubdirectory(currentFile, modulePath)
		if err != nil {
			return err
		}

		if !isSubDir {
			// Check for README.md
			readmePath := filepath.Join(modulePath, "README.md")
			if _, err := r.statFunc(readmePath); os.IsNotExist(err) {
				runner.EmitIssue(
					r,
					"Private modules should not be referenced externally. Add a README.md to make the referenced module public or remove the reference.",
					call.DefRange,
				)
			}
		}
	}

	return nil
}

func isSubdirectory(currentFile, modulePath string) (bool, error) {
	absCurrentFile, err := filepath.Abs(currentFile)
	if err != nil {
		return false, err
	}

	absCurrentFilePath := filepath.Dir(absCurrentFile)

	absModulePath, err := filepath.Abs(modulePath)
	if err != nil {
		return false, err
	}

	relPath, err := filepath.Rel(absCurrentFilePath, absModulePath)
	if err != nil {
		return false, err
	}

	return !strings.HasPrefix(relPath, ".."), nil
}
