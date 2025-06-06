package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

const (
	tfautomvBin   = "../../bin/tfautomv"
	terraformBin  = "terraform"
	terragruntBin = "terragrunt"
)

func writeCode(t *testing.T, path string, code string) {
	t.Helper()

	t.Logf("Writing code to %q", path)
	err := os.WriteFile(path, []byte(code), 0644)
	if err != nil {
		t.Fatalf("could not write file %q: %v", path, err)
	}
}

func runVersion(t *testing.T, executable string) *version.Version {
	t.Helper()

	runner, err := tfexec.NewTerraform(".", executable)
	if err != nil {
		t.Fatalf("could not create terraform runner: %v", err)
	}

	t.Logf("Running %s version", executable)
	v, _, err := runner.Version(context.Background(), false)
	if err != nil {
		t.Fatalf("%s version failed: %v", executable, err)
	}

	return v
}

func terraformVersion(t *testing.T) *version.Version {
	t.Helper()
	return runVersion(t, terraformBin)
}

func runInit(t *testing.T, workdir, executable string) {
	t.Helper()

	runner, err := tfexec.NewTerraform(workdir, executable)
	if err != nil {
		t.Fatalf("could not create terraform runner: %v", err)
	}

	t.Logf("Running %s init", executable)
	err = runner.Init(context.Background())
	if err != nil {
		t.Fatalf("%s init failed: %v", executable, err)
	}
}

func terraformInit(t *testing.T, workdir string) {
	t.Helper()
	runInit(t, workdir, terraformBin)
}

func runApply(t *testing.T, workdir, executable string) {
	t.Helper()

	runner, err := tfexec.NewTerraform(workdir, executable)
	if err != nil {
		t.Fatalf("could not create terraform runner: %v", err)
	}

	t.Logf("Running %s apply", executable)
	err = runner.Apply(context.Background())
	if err != nil {
		t.Fatalf("%s apply failed: %v", executable, err)
	}
}

func runInitAndApply(t *testing.T, workdir, executable string) {
	t.Helper()
	runInit(t, workdir, executable)
	runApply(t, workdir, executable)
}

func terraformInitAndApply(t *testing.T, workdir string) {
	t.Helper()
	runInitAndApply(t, workdir, terraformBin)
}

func terragruntInitAndApply(t *testing.T, workdir string) {
	t.Helper()
	runInitAndApply(t, workdir, terragruntBin)
}

func runPlan(t *testing.T, workdir, executable string) *tfjson.Plan {
	t.Helper()

	runner, err := tfexec.NewTerraform(workdir, executable)
	if err != nil {
		t.Fatalf("could not create terraform runner: %v", err)
	}

	planFile, err := os.CreateTemp(workdir, "tfplan-*.bin")
	if err != nil {
		t.Fatalf("could not create plan file: %v", err)
	}

	t.Logf("Running %s plan", executable)
	_, err = runner.Plan(context.Background(), tfexec.Out(planFile.Name()))
	if err != nil {
		t.Fatalf("%s plan failed: %v", executable, err)
	}

	t.Logf("Parsing %s plan", executable)
	plan, err := runner.ShowPlanFile(context.Background(), planFile.Name())
	if err != nil {
		t.Fatalf("could not parse %s plan: %v", executable, err)
	}

	return plan
}

func terraformPlan(t *testing.T, workdir string) *tfjson.Plan {
	t.Helper()
	return runPlan(t, workdir, terraformBin)
}

func terragruntPlan(t *testing.T, workdir string) *tfjson.Plan {
	t.Helper()
	return runPlan(t, workdir, terragruntBin)
}

func runTfautomv(t *testing.T, workdir string, args []string) string {
	t.Helper()

	tfautomvBinAbsPath, err := filepath.Abs(tfautomvBin)
	if err != nil {
		t.Fatalf("could not get absolute path to tfautomv binary: %v", err)
	}

	cmd := exec.Command(tfautomvBinAbsPath, args...)
	cmd.Dir = workdir

	var stdout bytes.Buffer
	var allOut bytes.Buffer

	cmd.Stdout = io.MultiWriter(&stdout, &allOut, os.Stderr)
	cmd.Stderr = io.MultiWriter(&allOut, os.Stderr)

	t.Log("Running tfautomv")
	if err := cmd.Run(); err != nil {
		t.Fatalf("tfautomv failed: %v", err)
	}

	return stdout.String()
}

func runShellCommands(t *testing.T, workdir string, commands string) {
	t.Helper()

	cmd := exec.Command("sh")
	cmd.Dir = workdir
	cmd.Stdin = strings.NewReader(commands)

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	t.Log("Running shell commands")
	if err := cmd.Run(); err != nil {
		t.Fatalf("shell commands failed: %v", err)
	}
}

func runTfautomvPipeSh(t *testing.T, workdir string, args []string) {
	t.Helper()

	commands := runTfautomv(t, workdir, args)

	if commands != "" {
		runShellCommands(t, workdir, commands)
	}
}

func countPlannedChanges(plan *tfjson.Plan) int {
	count := 0

	for _, change := range plan.ResourceChanges {
		actions := change.Change.Actions
		if slices.Contains(actions, tfjson.ActionCreate) || slices.Contains(actions, tfjson.ActionDelete) {
			count++
		}
	}

	return count
}

// createPlanFile creates a terraform plan file in the given directory
func createPlanFile(t *testing.T, workdir, filename string) {
	t.Helper()

	runner, err := tfexec.NewTerraform(workdir, terraformBin)
	if err != nil {
		t.Fatalf("could not create terraform runner: %v", err)
	}

	planPath := filepath.Join(workdir, filename)
	t.Logf("Creating plan file at %s", planPath)

	_, err = runner.Plan(context.Background(), tfexec.Out(planPath))
	if err != nil {
		t.Fatalf("terraform plan failed: %v", err)
	}
}

// createJSONPlanFile creates a JSON terraform plan file in the given directory
func createJSONPlanFile(t *testing.T, workdir, filename string) {
	t.Helper()

	runner, err := tfexec.NewTerraform(workdir, terraformBin)
	if err != nil {
		t.Fatalf("could not create terraform runner: %v", err)
	}

	// First create a binary plan file
	tempPlanPath := filepath.Join(workdir, "temp.plan")
	_, err = runner.Plan(context.Background(), tfexec.Out(tempPlanPath))
	if err != nil {
		t.Fatalf("terraform plan failed: %v", err)
	}

	// Convert to JSON using terraform show
	plan, err := runner.ShowPlanFile(context.Background(), tempPlanPath)
	if err != nil {
		t.Fatalf("terraform show failed: %v", err)
	}

	// Write JSON to the target file
	jsonPath := filepath.Join(workdir, filename)
	jsonData, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("failed to marshal plan to JSON: %v", err)
	}

	t.Logf("Creating JSON plan file at %s", jsonPath)
	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		t.Fatalf("failed to write JSON plan file: %v", err)
	}

	// Clean up temp file
	os.Remove(tempPlanPath)
}
