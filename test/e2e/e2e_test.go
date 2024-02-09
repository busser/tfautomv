package e2e

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestE2E_SimpleMatch(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_pet" "original_first" {
	length = 1
}
resource "random_pet" "original_second" {
	length = 2
}
resource "random_pet" "original_third" {
	length = 3
}`

	refactoredCode := `
resource "random_pet" "refactored_first" {
	length = 1
}
resource "random_pet" "refactored_second" {
	length = 2
}
resource "random_pet" "refactored_third" {
	length = 3
}`

	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)
	runTfautomv(t, workdir)
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, changeCount, 0)
}

func TestE2E_DependencyAnalysis(t *testing.T) {
	// This test is an example of a case where looking at resources one by one
	// isn't enough.
	//
	// In the cade below, each random_pet resource depends on a
	// random_integer resource. While refactoring, Terraform loses track of the
	// radnom_integer resources and tries to delete and recreate them. Because
	// of the dependency, Terraform also wants to delete and recreate the
	// random_pet resources.
	//
	// Writing a moved block for each random_integer resource would lead to an
	// empty plan. However, tfautomv does not find these moves because the
	// random_integer resources have identical attiributes.
	t.Skip("tfautomv cannot yet solve this case")

	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_integer" "first" {
	min = 1
	max = 5
}
resource "random_integer" "second" {
	min = 1
	max = 5
}
resource "random_pet" "first" {
	length    = random_integer.first.result
	separator = "-"
}
resource "random_pet" "second" {
	length    = random_integer.second.result
	separator = "+"
}`

	refactoredCode := `
resource "random_integer" "alpha" {
	min = 1
	max = 5
}
resource "random_integer" "beta" {
	min = 1
	max = 5
}
resource "random_pet" "first" {
	length    = random_integer.alpha.result
	separator = "-"
}
resource "random_pet" "second" {
	length    = random_integer.beta.result
	separator = "+"
}`

	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)
	runTfautomv(t, workdir)
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, changeCount, 0)
}

func TestE2E_MultipleTypes(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_pet" "original" {}
resource "random_uuid" "original" {}`

	refactoredCode := `
resource "random_pet" "refactored" {}
resource "random_uuid" "refactored" {}`

	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)
	runTfautomv(t, workdir)
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, changeCount, 0)
}

func TestE2E_DifferentAttributes(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_pet" "original" {
	length = 1
}`

	refactoredCode := `
resource "random_pet" "refactored" {
	length = 2
}`

	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)

	runTfautomv(t, workdir)
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, changeCount, 2)

	runTfautomv(t, workdir, "--ignore=everything:random_pet:length")
	changeCount = countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, changeCount, 1)
}

func TestE2E_Terragrunt(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")
	terragruntConfigPath := filepath.Join(workdir, "terragrunt.hcl")

	originalCode := `
variable "prefix" {
	type = string
}
resource "random_pet" "original_first" {
	prefix = var.prefix
	length = 1
}
resource "random_pet" "original_second" {
	prefix = var.prefix
	length = 2
}
resource "random_pet" "original_third" {
	prefix = var.prefix
	length = 3
}`

	refactoredCode := `
variable "prefix" {
	type = string
}
resource "random_pet" "refactored_first" {
	prefix = var.prefix
	length = 1
}
resource "random_pet" "refactored_second" {
	prefix = var.prefix
	length = 2
}
resource "random_pet" "refactored_third" {
	prefix = var.prefix
	length = 3
}`

	terragruntConfig := `
inputs = {
	prefix = "my-"
}`

	writeCode(t, codePath, originalCode)
	writeCode(t, terragruntConfigPath, terragruntConfig)
	terragruntInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)
	runTfautomv(t, workdir, "--terraform-bin=terragrunt")
	changeCount := countPlannedChanges(terragruntPlan(t, workdir))
	assert.Equal(t, changeCount, 0)
}

func TestE2E_TerraformCloud(t *testing.T) {
	tfVersion := terraformVersion(t)
	if tfVersion.LessThan(version.Must(version.NewVersion("1.6"))) {
		t.Skip("tfautomv requires Terraform 1.6 or later to run this test")
	}

	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
terraform {
	cloud {
		organization = "busser"
		workspaces {
			name = "tfautomv"
		}
	}
}
resource "random_pet" "original_first" {
	length = 1
}
resource "random_pet" "original_second" {
	length = 2
}
resource "random_pet" "original_third" {
	length = 3
}
`

	refactoredCode := `
terraform {
	cloud {
		organization = "busser"
		workspaces {
			name = "tfautomv"
		}
	}
}
resource "random_pet" "refactored_first" {
	length = 1
}
resource "random_pet" "refactored_second" {
	length = 2
}
resource "random_pet" "refactored_third" {
	length = 3
}
`

	// We don't apply the original code because we don't have a way of cleaning
	// up the workspace after the test. It's important that this test not change
	// the state of the Terraform Cloud workspace.
	writeCode(t, codePath, originalCode)
	terraformInit(t, workdir)
	writeCode(t, codePath, refactoredCode)
	runTfautomv(t, workdir)
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, changeCount, 0)
}
