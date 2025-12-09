package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_OpentofuJSONSyntaxRule(t *testing.T) {
	for _, tc := range []struct {
		name     string
		content  string
		filename string
		expected helper.Issues
		fixed    string
	}{
		{
			name:     "object syntax valid",
			content:  `{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}`,
			filename: "main.tf.json",
			expected: helper.Issues{},
		},
		{
			name:     "empty object valid",
			content:  `{}`,
			filename: "main.tf.json",
			expected: helper.Issues{},
		},
		{
			name:     "array syntax invalid",
			content:  `[{"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}}]`,
			filename: "main.tf.json",
			expected: helper.Issues{
				{
					Rule:    NewOpentofuJSONSyntaxRule(),
					Message: "JSON configuration uses array syntax at root, expected object",
					Range: hcl.Range{
						Filename: "main.tf.json",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 2},
					},
				},
			},
			fixed: `{
  "resource": {
    "aws_instance": {
      "example": {
        "ami": "ami-12345678"
      }
    }
  }
}
`,
		},
		{
			name:     "regular HCL file ignored",
			content:  `resource "aws_instance" "example" { ami = "ami-12345678" }`,
			filename: "main.tf",
			expected: helper.Issues{},
		},
		{
			name: "complex object syntax valid",
			content: `{
  "terraform": {
    "required_version": ">= 1.0"
  },
  "resource": {
    "aws_instance": {
      "example": {
        "ami": "ami-12345678"
      }
    }
  }
}`,
			filename: "main.tf.json",
			expected: helper.Issues{},
		},
		{
			name: "array with multiple objects invalid",
			content: `[
  {"resource": {"aws_instance": {"example": {"ami": "ami-12345678"}}}},
  {"variable": {"region": {"type": "string"}}}
]`,
			filename: "config.tf.json",
			expected: helper.Issues{
				{
					Rule:    NewOpentofuJSONSyntaxRule(),
					Message: "JSON configuration uses array syntax at root, expected object",
					Range: hcl.Range{
						Filename: "config.tf.json",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 2},
					},
				},
			},
			fixed: `{
  "resource": {
    "aws_instance": {
      "example": {
        "ami": "ami-12345678"
      }
    }
  },
  "variable": {
    "region": {
      "type": "string"
    }
  }
}
`,
		},
		{
			name: "array with multiple resources of same type",
			content: `[
  {"resource": {"aws_instance": {"foo": {"ami": "ami-11111111"}}}},
  {"resource": {"aws_instance": {"bar": {"ami": "ami-22222222"}}}}
]`,
			filename: "multi.tf.json",
			expected: helper.Issues{
				{
					Rule:    NewOpentofuJSONSyntaxRule(),
					Message: "JSON configuration uses array syntax at root, expected object",
					Range: hcl.Range{
						Filename: "multi.tf.json",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 2},
					},
				},
			},
			fixed: `{
  "resource": {
    "aws_instance": {
      "bar": {
        "ami": "ami-22222222"
      },
      "foo": {
        "ami": "ami-11111111"
      }
    }
  }
}
`,
		},
		{
			name: "array with multiple unlabeled blocks",
			content: `[
  {"import": {"id": "i-abc123", "to": "aws_instance.foo"}},
  {"import": {"id": "i-def456", "to": "aws_instance.bar"}}
]`,
			filename: "imports.tf.json",
			expected: helper.Issues{
				{
					Rule:    NewOpentofuJSONSyntaxRule(),
					Message: "JSON configuration uses array syntax at root, expected object",
					Range: hcl.Range{
						Filename: "imports.tf.json",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 2},
					},
				},
			},
			fixed: `{
  "import": [
    {
      "id": "i-abc123",
      "to": "aws_instance.foo"
    },
    {
      "id": "i-def456",
      "to": "aws_instance.bar"
    }
  ]
}
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewOpentofuJSONSyntaxRule()
			runner := helper.TestRunner(t, map[string]string{tc.filename: tc.content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.expected, runner.Issues)
			want := map[string]string{}
			if tc.fixed != "" {
				want[tc.filename] = tc.fixed
			}
			helper.AssertChanges(t, want, runner.Changes())
		})
	}
}
