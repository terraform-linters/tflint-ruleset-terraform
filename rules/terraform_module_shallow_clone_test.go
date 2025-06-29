package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformModuleShallowClone(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "terraform registry module",
			Content: `
module "registry" {
  source = "hashicorp/consul"
  version = "~> 1.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "local module",
			Content: `
module "local" {
  source = "./modules/consul"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "unpinned git module",
			Content: `
module "unpinned" {
  source = "git://github.com/hashicorp/consul.git"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module with 4 char long commit ID",
			Content: `
module "short_commit_pinned" {
  source = "git://github.com/hashicorp/consul.git?ref=babe"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module with 40 char long (SHA-1) commit ID",
			Content: `
module "sha1_commit_pinned" {
  source = "git://github.com/hashicorp/consul.git?ref=abc123def456789012345678901234567890abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module with 64 char long (SHA-256) commit ID",
			Content: `
module "sha256_commit_pinned" {
  source = "git://github.com/hashicorp/consul.git?ref=abc123def456789012345678901234567890abc123def4567890123456789012"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module with shallow clone enabled (ref first)",
			Content: `
module "shallow_clone" {
  source = "git://github.com/hashicorp/consul.git?ref=v1.0.0&depth=1"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module with shallow clone enabled (depth first)",
			Content: `
module "shallow_clone" {
  source = "git://github.com/hashicorp/consul.git?depth=1&ref=v1.0.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "pinned git module with SSH protocol",
			Content: `
module "ssh_pinned" {
  source = "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleShallowCloneRule(),
					Message: `Module source "git::ssh://git@github.com/hashicorp/consul.git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 71},
					},
				},
			},
			Fixed: `
module "ssh_pinned" {
  source = "git::ssh://git@github.com/hashicorp/consul.git?depth=1&ref=v1.0.0"
}`,
		},
		{
			Name: "pinned git module with HTTPS protocol",
			Content: `
module "https_pinned" {
  source = "git::https://github.com/hashicorp/consul.git?ref=v1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleShallowCloneRule(),
					Message: `Module source "git::https://github.com/hashicorp/consul.git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 69},
					},
				},
			},
			Fixed: `
module "https_pinned" {
  source = "git::https://github.com/hashicorp/consul.git?depth=1&ref=v1.0.0"
}`,
		},
		{
			Name: "pinned github module with SSH protocol",
			Content: `
module "github_ssh" {
  source = "git@github.com:hashicorp/consul.git?ref=v1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleShallowCloneRule(),
					Message: `Module source "git@github.com:hashicorp/consul.git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 60},
					},
				},
			},
			Fixed: `
module "github_ssh" {
  source = "git@github.com:hashicorp/consul.git?depth=1&ref=v1.0.0"
}`,
		},
		{
			Name: "pinned github module with HTTPS protocol",
			Content: `
module "github_https" {
  source = "github.com/hashicorp/consul?ref=v1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleShallowCloneRule(),
					Message: `Module source "github.com/hashicorp/consul?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 52},
					},
				},
			},
			Fixed: `
module "github_https" {
  source = "github.com/hashicorp/consul?depth=1&ref=v1.0.0"
}`,
		},
		{
			Name: "pinned bitbucket module",
			Content: `
module "bitbucket" {
  source = "bitbucket.org/hashicorp/tf-test-git?ref=v1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleShallowCloneRule(),
					Message: `Module source "bitbucket.org/hashicorp/tf-test-git?ref=v1.0.0" should enable shallow cloning by adding "depth=1" parameter`,
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 60},
					},
				},
			},
			Fixed: `
module "bitbucket" {
  source = "bitbucket.org/hashicorp/tf-test-git?depth=1&ref=v1.0.0"
}`,
		},
	}

	rule := NewTerraformModuleShallowCloneRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "module.tf"
			runner := testRunner(t, map[string]string{filename: tc.Content})

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
