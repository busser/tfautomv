package tfautomv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/padok-team/tfautomv/internal/terraform"
	"github.com/padok-team/tfautomv/internal/tfautomv"
)

func TestDetermineMovedBlocks(t *testing.T) {
	tt := []struct {
		name    string
		workdir string
	}{
		{
			"attributes",
			filepath.Join("./testdata", "/based-on-attributes"),
		},
		// This test fails.
		// tfautomv cannot yet solve this case.
		// {
		// 	"dependencies",
		// 	filepath.Join("./testdata", "/based-on-dependencies"),
		// },
		{
			"type",
			filepath.Join("./testdata", "/based-on-type"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			setupWorkdir(t, tc.workdir)

			workdir := filepath.Join(tc.workdir, "refactored-code")
			report, err := tfautomv.GenerateReport(workdir)
			if err != nil {
				t.Fatalf("GenerateReport(%q): %v", workdir, err)
			}

			err = terraform.AppendMovesToFile(report.Moves, filepath.Join(workdir, "moves.tf"))
			if err != nil {
				t.Fatalf("AppendMovesToFile(): %v", err)
			}

			tf := terraform.NewRunner(workdir)
			plan, err := tf.Plan()
			if err != nil {
				t.Fatal(err)
			}

			changes := tfautomv.CountChanges(plan)
			if changes > 0 {
				t.Errorf("%d changes remaining", changes)
			}
		})
	}
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
