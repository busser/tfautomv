//go:build e2e

package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/padok-team/tfautomv/internal/terraform"
)

func TestE2E(t *testing.T) {
	tt := []struct {
		name    string
		workdir string
		args    []string

		wantChanges int

		skip       bool
		skipReason string
	}{
		{
			name:        "same attributes",
			workdir:     filepath.Join("testdata", "same-attributes"),
			wantChanges: 0,
		},
		{
			name:        "requires dependency analysis",
			workdir:     filepath.Join("testdata", "requires-dependency-analysis"),
			wantChanges: 0,
			skip:        true,
			skipReason:  "tfautomv cannot yet solve this case",
		},
		{
			name:        "same type",
			workdir:     filepath.Join("testdata", "same-type"),
			wantChanges: 0,
		},
		{
			name:        "different attributes",
			workdir:     filepath.Join("testdata", "different-attributes"),
			wantChanges: 2,
		},
		{
			name:    "ignore different attributes",
			workdir: filepath.Join("testdata", "different-attributes"),
			args: []string{
				"--ignore=everything:random_pet:length",
			},
			wantChanges: 1,
		},
		{
			name:        "terraform cloud",
			workdir:     filepath.Join("testdata", "terraform-cloud"),
			wantChanges: 0,
			skip:        true,
			skipReason:  "tfautomv is currently incompatible with Terraform Cloud workspaces with the \"Remote\" execution mode.\nFor more details, see https://github.com/padok-team/tfautomv/issues/17",
		},
	}

	binPath := buildBinary(t)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip(tc.skipReason)
			}

			setupWorkdir(t, tc.workdir)

			workdir := filepath.Join(tc.workdir, "refactored-code")

			tfautomvCmd := exec.Command(binPath, tc.args...)
			tfautomvCmd.Dir = workdir
			tfautomvCmd.Stdout = os.Stderr
			tfautomvCmd.Stderr = os.Stderr

			if err := tfautomvCmd.Run(); err != nil {
				t.Fatalf("running tfautomv: %v", err)
			}

			plan, err := terraform.NewRunner(workdir).Plan()
			if err != nil {
				t.Fatalf("terraform plan (after addings moves): %v", err)
			}

			changes := plan.NumChanges()
			if changes != tc.wantChanges {
				t.Errorf("%d changes remaining, want %d", changes, tc.wantChanges)
			}
		})
	}
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

func setupWorkdir(t *testing.T, workdir string) {
	t.Helper()

	originalWorkdir := filepath.Join(workdir, "original-code")
	refactoredWorkdir := filepath.Join(workdir, "refactored-code")

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

	original := terraform.NewRunner(originalWorkdir)

	if err := original.Init(); err != nil {
		t.Fatal(err)
	}
	if err := original.Apply(); err != nil {
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
