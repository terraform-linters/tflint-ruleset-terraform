package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type documentedDescription struct {
	name        string
	description string
	rng         hcl.Range
}

func checkDocumentedDescriptions(runner tflint.Runner, rule tflint.Rule, blocks []*hclext.Block, kind string, unique bool) error {
	descriptions := []documentedDescription{}
	counts := map[string]int{}

	for _, block := range blocks {
		attr, exists := block.Body.Attributes["description"]
		if !exists {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf("`%s` %s has no description", block.Labels[0], kind),
				block.DefRange,
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
				rule,
				fmt.Sprintf("`%s` %s has no description", block.Labels[0], kind),
				block.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		if unique {
			descriptions = append(descriptions, documentedDescription{
				name:        block.Labels[0],
				description: description,
				rng:         attr.Range,
			})
			counts[description]++
		}
	}

	for _, description := range descriptions {
		if counts[description.description] < 2 {
			continue
		}

		if err := runner.EmitIssue(
			rule,
			fmt.Sprintf("`%s` %s description is not unique: %q", description.name, kind, description.description),
			description.rng,
		); err != nil {
			return err
		}
	}

	return nil
}
