package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// TerraformLocalsOrderRule checks whether comments use the preferred syntax
type TerraformLocalsOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformLocalsOrderRule returns a new rule
func NewTerraformLocalsOrderRule() *TerraformLocalsOrderRule {
	return &TerraformLocalsOrderRule{}
}

// Name returns the rule name
func (r *TerraformLocalsOrderRule) Name() string {
	return "terraform_locals_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformLocalsOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformLocalsOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformLocalsOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether single line comments is used
func (r *TerraformLocalsOrderRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = r.checkFile(runner, file); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformLocalsOrderRule) checkFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	for _, block := range blocks {
		if block.Type != "locals" {
			continue
		}
		if err := r.checkLocalsOrder(runner, block); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformLocalsOrderRule) checkLocalsOrder(runner tflint.Runner, block *hclsyntax.Block) error {
	attributes := r.attributesInLines(block)
	if r.sorted(attributes) {
		return nil
	}
	file, err := runner.GetFile(block.Range().Filename)
	if err != nil {
		return err
	}
	return r.suggestOrder(runner, attributes, block, file)
}

func (r *TerraformLocalsOrderRule) sorted(attributes []*hclsyntax.Attribute) bool {
	var names []string
	for _, a := range attributes {
		names = append(names, a.Name)
	}
	return sort.StringsAreSorted(names)
}

func (r *TerraformLocalsOrderRule) suggestOrder(runner tflint.Runner, attributes []*hclsyntax.Attribute, block *hclsyntax.Block, file *hcl.File) error {
	sort.Slice(attributes, func(i, j int) bool {
		return attributes[i].Name < attributes[j].Name
	})
	var locals []string
	for _, a := range attributes {
		locals = append(locals, string(a.SrcRange.SliceBytes(file.Bytes)))
	}
	suggestedBlock := strings.Join(locals, "\n")
	suggestedBlock = fmt.Sprintf("%s {\n%s\n}", block.Type, suggestedBlock)
	formattedBlock := string(hclwrite.Format([]byte(suggestedBlock)))
	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended locals variable order:\n%s", formattedBlock),
		block.DefRange(),
	)
}

func (r *TerraformLocalsOrderRule) attributesInLines(block *hclsyntax.Block) []*hclsyntax.Attribute {
	var attributes []*hclsyntax.Attribute
	for _, a := range block.Body.Attributes {
		attributes = append(attributes, a)
	}
	sort.Slice(attributes, func(i, j int) bool {
		if attributes[i].Range().Start.Line == attributes[j].Range().Start.Line {
			return attributes[i].Range().Start.Column < attributes[j].Range().Start.Column
		}
		return attributes[i].Range().Start.Line < attributes[j].Range().Start.Line
	})
	return attributes
}
