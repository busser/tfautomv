//go:build e2e

package e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/padok-team/tfautomv/internal/format"
	"github.com/padok-team/tfautomv/internal/terraform"
	"github.com/padok-team/tfautomv/internal/tfautomv"
)

func TestE2E(t *testing.T) {
	tt := []struct {
		name        string
		workdir     string
		want        []terraform.Move
		wantChanges int

		skip       bool
		skipReason string
	}{
		{
			name:    "same attributes",
			workdir: filepath.Join("testdata", "same-attributes"),
			want: []terraform.Move{
				{From: "random_pet.original_first", To: "random_pet.refactored_first"},
				{From: "random_pet.original_second", To: "random_pet.refactored_second"},
				{From: "random_pet.original_third", To: "random_pet.refactored_third"},
			},
		},
		{
			name:    "requires dependency analysis",
			workdir: filepath.Join("testdata", "requires-dependency-analysis"),
			want: []terraform.Move{
				{From: "random_integer.first", To: "random_integer.alpha"},
				{From: "random_integer.second", To: "random_integer.beta"},
			},
			skip:       true,
			skipReason: "tfautomv cannot yet solve this case",
		},
		{
			name:    "same type",
			workdir: filepath.Join("testdata", "same-type"),
			want: []terraform.Move{
				{From: "random_id.original", To: "random_id.refactored"},
				{From: "random_pet.original", To: "random_pet.refactored"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip(tc.skipReason)
			}

			setupWorkdir(t, tc.workdir)

			workdir := filepath.Join(tc.workdir, "refactored-code")
			tf := terraform.NewRunner(workdir)

			err := tf.Init()
			if err != nil {
				t.Fatalf("terraform init: %v", err)
			}
			plan, err := tf.Plan()
			if err != nil {
				t.Fatalf("terraform plan: %v", err)
			}

			analysis, err := tfautomv.AnalysisFromPlan(plan)
			if err != nil {
				t.Fatalf("AnalysisFromPlan(): %v", err)
			}

			t.Logf("\n%s", format.Analysis(analysis))

			moves := tfautomv.MovesFromAnalysis(analysis)

			if diff := cmp.Diff(tc.want, moves); diff != "" {
				t.Errorf("Moves mismatch (-want +got):\n%s", diff)
			}

			err = terraform.AppendMovesToFile(moves, filepath.Join(workdir, "moves.tf"))
			if err != nil {
				t.Fatalf("AppendMovesToFile(): %v", err)
			}

			plan, err = tf.Plan()
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
