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

	runTfautomvPipeSh(t, workdir, nil)

	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

func TestE2E_OutputBlocks(t *testing.T) {
	tfVersion := terraformVersion(t)
	if tfVersion.LessThan(version.Must(version.NewVersion("1.1"))) {
		t.Skip("tfautomv requires Terraform 1.6 or later to run this test")
	}

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

	commands := runTfautomv(t, workdir, []string{"--output=blocks"})
	assert.Empty(t, commands)
	assert.FileExists(t, filepath.Join(workdir, "moves.tf"))

	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

func TestE2E_OutputCommands(t *testing.T) {
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

	commands := runTfautomv(t, workdir, []string{"--output=commands"})
	assert.NotEmpty(t, commands)
	assert.NoFileExists(t, filepath.Join(workdir, "moves.tf"))
	runShellCommands(t, workdir, commands)

	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
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

	runTfautomvPipeSh(t, workdir, nil)

	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
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

	runTfautomvPipeSh(t, workdir, nil)

	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
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

	runTfautomvPipeSh(t, workdir, nil)
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 2, changeCount)

	runTfautomvPipeSh(t, workdir, []string{"--ignore=everything:random_pet:length"})
	changeCount = countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 1, changeCount)
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

	runTfautomvPipeSh(t, workdir, []string{"--terraform-bin=terragrunt"})

	changeCount := countPlannedChanges(terragruntPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

func TestE2E_TerragruntOutputCommands(t *testing.T) {
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

	commands := runTfautomv(t, workdir, []string{
		"--terraform-bin=terragrunt",
		"--output=commands"})
	assert.NotEmpty(t, commands)
	assert.Contains(t, commands, "terragrunt state mv")
	assert.NotContains(t, commands, "terraform state mv")
	assert.NoFileExists(t, filepath.Join(workdir, "moves.tf"))
	runShellCommands(t, workdir, commands)

	changeCount := countPlannedChanges(terragruntPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

func TestE2E_OpenTofu(t *testing.T) {
	checkOpentofuAvailable(t)
	
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
	opentofuInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)

	runTfautomvPipeSh(t, workdir, []string{"--terraform-bin=tofu"})

	changeCount := countPlannedChanges(opentofuPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

func TestE2E_OpenTofuOutputCommands(t *testing.T) {
	checkOpentofuAvailable(t)
	
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
	opentofuInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)

	commands := runTfautomv(t, workdir, []string{
		"--terraform-bin=tofu",
		"--output=commands"})
	assert.NotEmpty(t, commands)
	assert.Contains(t, commands, "tofu state mv")
	assert.NotContains(t, commands, "terraform state mv")
	assert.NoFileExists(t, filepath.Join(workdir, "moves.tf"))
	runShellCommands(t, workdir, commands)

	changeCount := countPlannedChanges(opentofuPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
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

	runTfautomvPipeSh(t, workdir, nil)

	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

func TestE2E_MultipleModules(t *testing.T) {
	tfVersion := terraformVersion(t)
	if tfVersion.LessThan(version.Must(version.NewVersion("0.14"))) {
		t.Skip("tfautomv requires Terraform 0.14 or later to run this test")
	}

	workdirA := t.TempDir()
	codePathA := filepath.Join(workdirA, "main.tf")

	workdirB := t.TempDir()
	codePathB := filepath.Join(workdirB, "main.tf")

	workdirC := t.TempDir()
	codePathC := filepath.Join(workdirC, "main.tf")

	originalCodeA := `
resource "random_pet" "original_first" {
	length = 1
}
resource "random_pet" "original_second" {
	length = 2
}
`
	originalCodeB := `
resource "random_pet" "original_third" {
	length = 3
}
resource "random_pet" "original_fourth" {
	length = 4
}
`
	originalCodeC := `
resource "random_pet" "original_fifth" {
	length = 5
}
resource "random_pet" "original_sixth" {
	length = 6
}
`

	refactoredCodeA := `
resource "random_pet" "refactored_first" {
	length = 1
}
resource "random_pet" "refactored_fourth" {
	length = 4
}
`
	refactoredCodeB := `
resource "random_pet" "refactored_third" {
	length = 3
}
resource "random_pet" "refactored_sixth" {
	length = 6
}
`
	refactoredCodeC := `
resource "random_pet" "refactored_fifth" {
	length = 5
}
resource "random_pet" "refactored_second" {
	length = 2
}
`

	writeCode(t, codePathA, originalCodeA)
	writeCode(t, codePathB, originalCodeB)
	writeCode(t, codePathC, originalCodeC)

	terraformInitAndApply(t, workdirA)
	terraformInitAndApply(t, workdirB)
	terraformInitAndApply(t, workdirC)

	writeCode(t, codePathA, refactoredCodeA)
	writeCode(t, codePathB, refactoredCodeB)
	writeCode(t, codePathC, refactoredCodeC)

	// We don't care where we run tfautomv from, since we explicitly pass the
	// workdirs as arguments.
	runTfautomvPipeSh(t, t.TempDir(), []string{workdirA, workdirB, workdirC})

	changeCount := 0
	changeCount += countPlannedChanges(terraformPlan(t, workdirA))
	changeCount += countPlannedChanges(terraformPlan(t, workdirB))
	changeCount += countPlannedChanges(terraformPlan(t, workdirC))
	assert.Equal(t, 0, changeCount)
}

// Test --preplanned flag with single directory and default plan file name
func TestE2E_Preplanned_SingleDirectory_DefaultFile(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_pet" "original_first" {
	length = 1
}
resource "random_pet" "original_second" {
	length = 2
}`

	refactoredCode := `
resource "random_pet" "refactored_first" {
	length = 1
}
resource "random_pet" "refactored_second" {
	length = 2
}`

	// Set up infrastructure
	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)

	// Create plan file
	createPlanFile(t, workdir, "tfplan.bin")

	// Run tfautomv with --preplanned (default filename)
	runTfautomvPipeSh(t, workdir, []string{"--preplanned"})

	// Verify no changes needed
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

// Test --preplanned flag with custom plan file name
func TestE2E_Preplanned_SingleDirectory_CustomFile(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_pet" "original_first" {
	length = 1
}
resource "random_pet" "original_second" {
	length = 2
}`

	refactoredCode := `
resource "random_pet" "refactored_first" {
	length = 1
}
resource "random_pet" "refactored_second" {
	length = 2
}`

	// Set up infrastructure
	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)

	// Create plan file with custom name
	createPlanFile(t, workdir, "my-plan.bin")

	// Run tfautomv with custom plan file
	runTfautomvPipeSh(t, workdir, []string{"--preplanned", "--preplanned-file=my-plan.bin"})

	// Verify no changes needed
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

// Test --preplanned flag with JSON plan file
func TestE2E_Preplanned_JSONPlanFile(t *testing.T) {
	workdir := t.TempDir()
	codePath := filepath.Join(workdir, "main.tf")

	originalCode := `
resource "random_pet" "original_first" {
	length = 1
}
resource "random_pet" "original_second" {
	length = 2
}`

	refactoredCode := `
resource "random_pet" "refactored_first" {
	length = 1
}
resource "random_pet" "refactored_second" {
	length = 2
}`

	// Set up infrastructure
	writeCode(t, codePath, originalCode)
	terraformInitAndApply(t, workdir)
	writeCode(t, codePath, refactoredCode)

	// Create JSON plan file
	createJSONPlanFile(t, workdir, "tfplan.json")

	// Run tfautomv with JSON plan file
	runTfautomvPipeSh(t, workdir, []string{"--preplanned", "--preplanned-file=tfplan.json"})

	// Verify no changes needed
	changeCount := countPlannedChanges(terraformPlan(t, workdir))
	assert.Equal(t, 0, changeCount)
}

// Test --preplanned flag with multiple directories
func TestE2E_Preplanned_MultipleDirectories(t *testing.T) {
	tfVersion := terraformVersion(t)
	if tfVersion.LessThan(version.Must(version.NewVersion("0.14"))) {
		t.Skip("tfautomv requires Terraform 0.14 or later to run this test")
	}

	workdirA := t.TempDir()
	workdirB := t.TempDir()
	codePathA := filepath.Join(workdirA, "main.tf")
	codePathB := filepath.Join(workdirB, "main.tf")

	originalCodeA := `
resource "random_pet" "original_first" {
	length = 1
}`
	originalCodeB := `
resource "random_pet" "original_second" {
	length = 2
}`

	refactoredCodeA := `
resource "random_pet" "refactored_first" {
	length = 1
}`
	refactoredCodeB := `
resource "random_pet" "refactored_second" {
	length = 2
}`

	// Set up infrastructure in both directories
	writeCode(t, codePathA, originalCodeA)
	writeCode(t, codePathB, originalCodeB)
	terraformInitAndApply(t, workdirA)
	terraformInitAndApply(t, workdirB)
	writeCode(t, codePathA, refactoredCodeA)
	writeCode(t, codePathB, refactoredCodeB)

	// Create plan files in both directories
	createPlanFile(t, workdirA, "tfplan.bin")
	createPlanFile(t, workdirB, "tfplan.bin")

	// Run tfautomv with --preplanned on multiple directories
	runTfautomvPipeSh(t, t.TempDir(), []string{"--preplanned", workdirA, workdirB})

	// Verify no changes needed in either directory
	changeCountA := countPlannedChanges(terraformPlan(t, workdirA))
	changeCountB := countPlannedChanges(terraformPlan(t, workdirB))
	assert.Equal(t, 0, changeCountA+changeCountB)
}
