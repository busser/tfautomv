//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/busser/tfautomv/internal/slices"
)

// ANSI escape sequence used for color output
const colorEscapeSequence = "\x1b"

func TestE2E(t *testing.T) {
	tt := []struct {
		name    string
		workdir string
		args    []string

		wantChanges       int
		wantOutputInclude []string
		wantOutputExclude []string

		skip       bool
		skipReason string
	}{
		{
			name:        "same attributes",
			workdir:     filepath.Join("testdata", "same-attributes"),
			wantChanges: 0,
			wantOutputInclude: []string{
				colorEscapeSequence,
			},
		},
		{
			name:        "requires dependency analysis",
			workdir:     filepath.Join("testdata", "requires-dependency-analysis"),
			wantChanges: 0,
			wantOutputInclude: []string{
				colorEscapeSequence,
			},
			skip:       true,
			skipReason: "tfautomv cannot yet solve this case",
		},
		{
			name:        "same type",
			workdir:     filepath.Join("testdata", "same-type"),
			wantChanges: 0,
			wantOutputInclude: []string{
				colorEscapeSequence,
			},
		},
		{
			name:        "different attributes",
			workdir:     filepath.Join("testdata", "different-attributes"),
			wantChanges: 2,
			wantOutputInclude: []string{
				colorEscapeSequence,
			},
		},
		{
			name:    "ignore different attributes",
			workdir: filepath.Join("testdata", "different-attributes"),
			args: []string{
				"-ignore=everything:random_pet:length",
			},
			wantChanges: 1,
			wantOutputInclude: []string{
				colorEscapeSequence,
			},
		},
		{
			name:    "no color",
			workdir: filepath.Join("testdata", "same-attributes"),
			args: []string{
				"-no-color",
			},
			wantChanges: 0,
			wantOutputExclude: []string{
				colorEscapeSequence,
			},
		},
		{
			name:        "terraform cloud",
			workdir:     filepath.Join("testdata", "terraform-cloud"),
			wantChanges: 0,
		},
		{
			name:    "terragrunt",
			workdir: filepath.Join("testdata", "terragrunt"),
			args: []string{
				"-terraform-bin=terragrunt",
			},
			wantChanges: 0,
			wantOutputInclude: []string{
				colorEscapeSequence,
			},
		},
	}

	binPath := buildBinary(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			for _, outputFormat := range []string{"blocks", "commands"} {
				t.Run(outputFormat, func(t *testing.T) {

					originalWorkdir := filepath.Join(tc.workdir, "original-code")
					refactoredWorkdir := filepath.Join(tc.workdir, "refactored-code")

					terraformBin := "terraform"
					for _, a := range tc.args {
						if strings.HasPrefix(a, "-terraform-bin=") {
							terraformBin = strings.TrimPrefix(a, "-terraform-bin=")
						}
					}

					/*
						Skip tests that serve as documentation of known limitations or
						use features incompatible with the Terraform CLI's version.
					*/

					if tc.skip {
						t.Skip(tc.skipReason)
					}

					if outputFormat == "blocks" {
						tf, err := tfexec.NewTerraform(originalWorkdir, terraformBin)
						if err != nil {
							t.Fatal(err)
						}
						tfVer, _, err := tf.Version(context.TODO(), false)
						if err != nil {
							t.Fatalf("failed to get terraform version: %v", err)
						}

						if tfVer.LessThan(version.Must(version.NewVersion("1.1"))) {
							t.Skip("terraform moves output format is only supported in terraform 1.1 and above")
						}
					}

					/*
						Create a fresh environment for each test.
					*/

					setupWorkdir(t, originalWorkdir, refactoredWorkdir, terraformBin)

					args := append(tc.args, fmt.Sprintf("-output=%s", outputFormat))

					/*
						Run tfautomv to generate `moved` blocks or `terraform state mv` commands.
					*/

					tfautomvCmd := exec.Command(binPath, args...)
					tfautomvCmd.Dir = refactoredWorkdir

					var tfautomvStdout bytes.Buffer
					var tfautomvCompleteOutput bytes.Buffer
					tfautomvCmd.Stdout = io.MultiWriter(&tfautomvStdout, &tfautomvCompleteOutput, os.Stderr)
					tfautomvCmd.Stderr = io.MultiWriter(&tfautomvCompleteOutput, os.Stderr)

					if err := tfautomvCmd.Run(); err != nil {
						t.Fatalf("running tfautomv: %v", err)
					}

					/*
						Validate that tfautomv produced the expected output.
					*/

					outputStr := tfautomvCompleteOutput.String()
					for _, s := range tc.wantOutputInclude {
						if !strings.Contains(outputStr, s) {
							t.Errorf("output should contain %q but does not", s)
						}
					}
					for _, s := range tc.wantOutputExclude {
						if strings.Contains(outputStr, s) {
							t.Errorf("output should not contain %q but does", s)
						}
					}

					/*
						If using `terraform state mv` commands, run them.
					*/

					if outputFormat == "commands" {
						cmd := exec.Command("/bin/sh")
						cmd.Dir = refactoredWorkdir

						cmd.Stdin = &tfautomvStdout
						cmd.Stdout = os.Stderr
						cmd.Stderr = os.Stderr

						if err := cmd.Run(); err != nil {
							t.Fatalf("running terraform state mv commands: %v", err)
						}
					}

					/*
						Count how many changes remain in Terraform's plan.
					*/

					tf, err := tfexec.NewTerraform(refactoredWorkdir, terraformBin)
					if err != nil {
						t.Fatal(err)
					}

					planFile, err := os.CreateTemp("", "tfautomv.*.plan")
					if err != nil {
						t.Fatal(err)
					}
					defer os.Remove(planFile.Name())
					if _, err := tf.Plan(context.TODO(), tfexec.Out(planFile.Name())); err != nil {
						t.Fatalf("terraform plan (after addings moves): %v", err)
					}
					plan, err := tf.ShowPlanFile(context.TODO(), planFile.Name())
					if err != nil {
						t.Fatalf("terraform show (after addings moves): %v", err)
					}

					changes := numChanges(plan)
					if changes != tc.wantChanges {
						t.Errorf("%d changes remaining, want %d", changes, tc.wantChanges)
					}
				})
			}
		})
	}
}

func numChanges(p *tfjson.Plan) int {
	count := 0

	for _, rc := range p.ResourceChanges {
		if slices.Contains(rc.Change.Actions, tfjson.ActionCreate) || slices.Contains(rc.Change.Actions, tfjson.ActionDelete) {
			count++
		}
	}

	return count
}

func buildBinary(t *testing.T) string {
	t.Helper()

	rootDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("could not get root directory: %v", err)
	}

	buildCmd := exec.Command("make", "build")
	buildCmd.Dir = rootDir
	buildCmd.Stdout = os.Stderr
	buildCmd.Stderr = os.Stderr

	t.Log("Building tfautomv binary...")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("make build: %v", err)
	}

	binPath := filepath.Join(rootDir, "bin", "tfautomv")
	return binPath
}

func setupWorkdir(t *testing.T, originalWorkdir, refactoredWorkdir, terraformBin string) {
	t.Helper()

	filesToRemove := []string{
		filepath.Join(originalWorkdir, "terraform.tfstate"),
		filepath.Join(originalWorkdir, ".terraform.lock.hcl"),
		filepath.Join(refactoredWorkdir, "terraform.tfstate"),
		filepath.Join(refactoredWorkdir, ".terraform.lock.hcl"),
		filepath.Join(refactoredWorkdir, "moves.tf"),
	}
	for _, f := range filesToRemove {
		ensureFileRemoved(t, f)
	}

	directoriesToRemove := []string{
		filepath.Join(originalWorkdir, ".terraform"),
		filepath.Join(refactoredWorkdir, ".terraform"),
	}
	for _, d := range directoriesToRemove {
		ensureDirectoryRemoved(t, d)
	}

	original, err := tfexec.NewTerraform(originalWorkdir, terraformBin)
	if err != nil {
		t.Fatal(err)
	}

	if err := original.Init(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if err := original.Apply(context.TODO()); err != nil {
		t.Fatal(err)
	}

	os.Rename(
		filepath.Join(originalWorkdir, "terraform.tfstate"),
		filepath.Join(refactoredWorkdir, "terraform.tfstate"),
	)
}

func ensureFileRemoved(t *testing.T, path string) {
	t.Helper()

	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("could not remove file %q: %v", path, err)
	}
}

func ensureDirectoryRemoved(t *testing.T, path string) {
	t.Helper()

	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("could not remove directory %q: %v", path, err)
	}
}
