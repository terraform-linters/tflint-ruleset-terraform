package rules

import (
	"testing"

	"github.com/diofeher/tflint-ruleset-opentofu/opentofu"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testRunner(t *testing.T, files map[string]string) *opentofu.Runner {
	return opentofu.NewRunner(helper.TestRunner(t, files))
}
