package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
)

// documentedRule is the shared implementation of the rules that require a
// description on every block of a given type.
type documentedRule struct {
	tflint.DefaultRule

	name      string
	blockType string
}

// documentedRuleConfig is the config structure for rules embedding documentedRule
type documentedRuleConfig struct {
	Unique bool `hclext:"unique,optional"`
}

// documentedBlock is a block that declares a non-empty description
type documentedBlock struct {
	name             string
	description      string
	descriptionRange hcl.Range
}

// Name returns the rule name
func (r *documentedRule) Name() string {
	return r.name
}

// Enabled returns whether the rule is enabled by default
func (r *documentedRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *documentedRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *documentedRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// check emits an issue for every block without a description, and when unique
// is configured, for every description shared by more than one block. Issues
// are attributed to rule, which is the rule embedding documentedRule.
func (r *documentedRule) check(runner tflint.Runner, rule tflint.Rule) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	config := documentedRuleConfig{}
	if err := runner.DecodeRuleConfig(r.name, &config); err != nil {
		return err
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       r.blockType,
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

	documented, err := r.checkDescriptions(runner, rule, body.Blocks)
	if err != nil {
		return err
	}

	if !config.Unique {
		return nil
	}

	return r.checkUniqueDescriptions(runner, rule, documented)
}

// checkDescriptions emits an issue for each block that declares no description
// or an empty one, and returns the blocks that declare a description.
func (r *documentedRule) checkDescriptions(runner tflint.Runner, rule tflint.Rule, blocks []*hclext.Block) ([]documentedBlock, error) {
	documented := make([]documentedBlock, 0, len(blocks))

	for _, block := range blocks {
		name := block.Labels[0]

		if attr, exists := block.Body.Attributes["description"]; exists {
			var description string
			if diags := gohcl.DecodeExpression(attr.Expr, nil, &description); diags.HasErrors() {
				return nil, diags
			}

			if description != "" {
				documented = append(documented, documentedBlock{
					name:             name,
					description:      description,
					descriptionRange: attr.Range,
				})
				continue
			}
		}

		if err := runner.EmitIssue(
			rule,
			fmt.Sprintf("`%s` %s has no description", name, r.blockType),
			block.DefRange,
		); err != nil {
			return nil, err
		}
	}

	return documented, nil
}

// checkUniqueDescriptions emits an issue for every block whose description is
// shared with another block. Every participant is reported so that the output
// does not depend on the order blocks were parsed in.
func (r *documentedRule) checkUniqueDescriptions(runner tflint.Runner, rule tflint.Rule, blocks []documentedBlock) error {
	counts := make(map[string]int, len(blocks))
	for _, block := range blocks {
		counts[block.description]++
	}

	for _, block := range blocks {
		if counts[block.description] < 2 {
			continue
		}

		if err := runner.EmitIssue(
			rule,
			fmt.Sprintf("`%s` %s description is not unique: %q", block.name, r.blockType, block.description),
			block.descriptionRange,
		); err != nil {
			return err
		}
	}

	return nil
}
