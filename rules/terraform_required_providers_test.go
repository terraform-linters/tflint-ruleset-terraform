package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformRequiredProvidersRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Config   string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "no version",
			Content: `
provider "template" {}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 20,
						},
					},
				},
			},
		},
		{
			Name: "implicit provider - resource",
			Content: `
resource "random_string" "foo" {
  length = 16
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"random\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 31,
						},
					},
				},
			},
		},
		{
			Name: "implicit provider - resource",
			Content: `
ephemeral "random_string" "foo" {
  length = 16
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"random\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 32,
						},
					},
				},
			},
		},
		{
			Name: "implicit provider - data source",
			Content: `
data "template_file" "foo" {
  template = ""
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 27,
						},
					},
				},
			},
		},
		{
			Name: "required_providers object",
			Content: `
terraform {
  required_providers {
    template = {
      source  = "hashicorp/template"
      version = "~> 2"
    }
  }
}
provider "template" {}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "legacy required_providers string",
			Content: `
terraform {
  required_providers {
    template = "~> 2"
  }
}
provider "template" {}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Legacy version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 22,
						},
					},
				},
			},
			Fixed: `
terraform {
  required_providers {
    template = {
      source  = "hashicorp/template"
      version = "~> 2"
    }
  }
}
provider "template" {}
`,
		},
		{
			Name: "required_providers object missing version",
			Content: `
terraform {
  required_providers {
    template = {
      source = "hashicorp/template"
    }
  }
}

provider "template" {}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 16,
						},
						End: hcl.Pos{
							Line:   6,
							Column: 6,
						},
					},
				},
			},
		},
		{
			Name: "required_providers object missing version ignored",
			Content: `
terraform {
  required_providers {
    template = {
      source = "hashicorp/template"
    }
  }
}

provider "template" {}
`,
			Config: `
rule "terraform_required_providers" {
  enabled = true

  version = false
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "required_providers object missing source",
			Content: `
terraform {
  required_providers {
    template = {
      version = "~> 2"
    }
  }
}

provider "template" {}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing `source` for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 16,
						},
						End: hcl.Pos{
							Line:   6,
							Column: 6,
						},
					},
				},
			},
			Fixed: `
terraform {
  required_providers {
    template = {
      source  = "hashicorp/template"
      version = "~> 2"
    }
  }
}

provider "template" {}
`,
		},
		{
			Name: "required_providers object missing source ignored",
			Content: `
terraform {
  required_providers {
    template = {
      version = "~> 2"
    }
  }
}

provider "template" {}
`,
			Config: `
rule "terraform_required_providers" {
  enabled = true

  source = false
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "required_providers empty object",
			Content: `
terraform {
  required_providers {
    template = {}
  }
}

provider "template" {}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing `source` for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 18,
						},
					},
				},
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 18,
						},
					},
				},
			},
			Fixed: `
terraform {
  required_providers {
    template = {
      source = "hashicorp/template"
    }
  }
}

provider "template" {}
`,
		},
		{
			Name: "single provider with alias",
			Content: `
provider "template" {
  alias = "b"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 20,
						},
					},
				},
			},
		},
		{
			Name: "version set",
			Content: `
terraform {
  required_providers {
    template = {
      source = "hashicorp/template"
      version = "~> 2"
    }
  }
}

provider "template" {
  version = "~> 2"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "provider version constraint should be specified via `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   11,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   11,
							Column: 20,
						},
					},
				},
			},
		},
		{
			Name: "version set with configuration_aliases",
			Content: `
terraform {
  required_providers {
    template = {
      source = "hashicorp/template"
      version = "~> 2"
      configuration_aliases = [template.alias]
    }
  }
}

data "template_file" "foo" {
  provider = template.alias
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "version set with alias",
			Content: `
terraform {
  required_providers {
    template = {
      source = "hashicorp/template"
      version = "~> 2"
    }
  }
}

provider "template" {
  alias   = "foo"
  version = "~> 2"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "provider version constraint should be specified via `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   11,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   11,
							Column: 20,
						},
					},
				},
			},
		},
		{
			Name: "terraform provider",
			Content: `
data "terraform_remote_state" "foo" {}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "builtin provider",
			Content: `
terraform {
  required_providers {
    test = {
      source = "terraform.io/builtin/test"
    }
  }
}
resource "test_assertions" "foo" {}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource provider ref",
			Content: `
terraform {
  required_providers {
    google = {
      version = "~> 4.27.0"
	}
  }
}

resource "google_compute_instance" "foo" {
  provider = google-beta
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"google-beta\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   10,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   10,
							Column: 41,
						},
					},
				},
			},
		},
		{
			Name: "resource provider ref as string",
			Content: `
terraform {
  required_providers {
    google = {
      version = "~> 4.27.0"
    }
  }
}

resource "google_compute_instance" "foo" {
  provider = "google-beta"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"google-beta\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   10,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   10,
							Column: 41,
						},
					},
				},
			},
		},
		{
			Name: "JSON syntax",
			Content: `
{
  "terraform": {
    "required_providers": {
      "template": "~> 2"
	}
  },
  "provider": {
    "template": {}
  }
}`,
			JSON: true,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Legacy version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf.json",
						Start: hcl.Pos{
							Line:   5,
							Column: 19,
						},
						End: hcl.Pos{
							Line:   5,
							Column: 25,
						},
					},
				},
			},
		},
		{
			Name: "provider-defined function",
			Content: `
output "foo" {
	value = provider::time::rfc3339_parse("2023-07-25T23:43:16Z")
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Missing version constraint for provider \"time\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   3,
							Column: 10,
						},
						End: hcl.Pos{
							Line:   3,
							Column: 63,
						},
					},
				},
			},
		},
		{
			Name: "multiple required providers",
			Content: `
terraform {
  required_providers {
    template = "~> 2"
  }

  required_providers {
    aws = "~> 5.0"
  }
}

provider "template" {}
provider "aws" {}
provider "google" {}

terraform {
  required_providers {
    google = "~> 6.0"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Legacy version constraint for provider \"template\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 22,
						},
					},
				},
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Legacy version constraint for provider \"aws\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   8,
							Column: 11,
						},
						End: hcl.Pos{
							Line:   8,
							Column: 19,
						},
					},
				},
				{
					Rule:    NewTerraformRequiredProvidersRule(),
					Message: "Legacy version constraint for provider \"google\" in `required_providers`",
					Range: hcl.Range{
						Filename: "module.tf",
						Start: hcl.Pos{
							Line:   18,
							Column: 14,
						},
						End: hcl.Pos{
							Line:   18,
							Column: 22,
						},
					},
				},
			},
			Fixed: `
terraform {
  required_providers {
    template = {
      source  = "hashicorp/template"
      version = "~> 2"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "template" {}
provider "aws" {}
provider "google" {}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}
`,
		},
	}

	rule := NewTerraformRequiredProvidersRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "module.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := testRunner(t, map[string]string{
				filename:      tc.Content,
				".tflint.hcl": tc.Config,
			})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
			want := map[string]string{}
			if tc.Fixed != "" {
				want[filename] = tc.Fixed
			}
			helper.AssertChanges(t, want, runner.Runner.(*helper.Runner).Changes())
		})
	}
}
