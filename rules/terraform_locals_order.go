package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"sort"
)

// TerraformLocalsOrderRule checks whether all arguments inside a `locals` block are sortedByAlphabetOrder in alphabet order
type TerraformLocalsOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformLocalsOrderRule returns a new rule
func NewTerraformLocalsOrderRule() *TerraformLocalsOrderRule {
	return &TerraformLocalsOrderRule{}
}

// Name returns the rule name
func (r *TerraformLocalsOrderRule) Name() string {
	return "terraform_ordered_locals"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformLocalsOrderRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformLocalsOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformLocalsOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether all arguments inside a `locals` block are sortedByAlphabetOrder in alphabet order
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
	content, _, schemaDiags := file.Body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{{Type: "locals"}},
	})
	if schemaDiags.HasErrors() {
		return schemaDiags
	}

	for _, block := range content.Blocks {
		if err := r.checkLocalsOrder(runner, block); err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformLocalsOrderRule) checkLocalsOrder(runner tflint.Runner, block *hcl.Block) error {
	locals, err := r.attributesInLines(block)
	if err != nil {
		return err
	}
	if !r.sortedByAlphabetOrder(locals) {
		err = runner.EmitIssue(
			r,
			"local values must be in alphabetical order",
			block.DefRange,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformLocalsOrderRule) sortedByAlphabetOrder(attributes []*hcl.Attribute) bool {
	var names []string
	for _, a := range attributes {
		names = append(names, a.Name)
	}
	return sort.StringsAreSorted(names)
}

func (r *TerraformLocalsOrderRule) attributesInLines(block *hcl.Block) ([]*hcl.Attribute, error) {
	attributesMaps, diagnostics := block.Body.JustAttributes()
	if diagnostics.HasErrors() {
		return nil, diagnostics
	}
	var attributes []*hcl.Attribute
	for _, a := range attributesMaps {
		attributes = append(attributes, a)
	}
	sort.Slice(attributes, func(x, y int) bool {
		posX := attributes[x].Range.Start
		posY := attributes[y].Range.Start
		if posX.Line == posY.Line {
			return posX.Column < posY.Column
		}
		return posX.Line < posY.Line
	})
	return attributes, nil
}
