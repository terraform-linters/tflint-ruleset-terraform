package main

import (
	"github.com/diofeher/tflint-ruleset-opentofu/opentofu"
	"github.com/diofeher/tflint-ruleset-opentofu/project"
	"github.com/diofeher/tflint-ruleset-opentofu/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &opentofu.RuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "opentofu",
				Version: project.Version,
			},
			PresetRules: rules.PresetRules,
		},
	})
}
