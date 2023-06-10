package format

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/busser/tfautomv/internal/terraform"
)

func TestMoves(t *testing.T) {
	tt := []struct {
		name string

		moves   []terraform.Move
		noColor bool

		want string
	}{
		{
			name: "single",
			moves: []terraform.Move{
				{From: "random_pet.original", To: "random_pet.refactored"},
			},
			noColor: false,
			want:    filepath.Join("testdata", "moves", "single.txt"),
		},
		{
			name: "single no color",
			moves: []terraform.Move{
				{From: "random_pet.original", To: "random_pet.refactored"},
			},
			noColor: true,
			want:    filepath.Join("testdata", "moves", "single-no-color.txt"),
		},
		{
			name: "multiple",
			moves: []terraform.Move{
				{From: "random_id.original", To: "random_id.refactored"},
				{From: "random_pet.original", To: "random_pet.refactored"},
			},
			noColor: false,
			want:    filepath.Join("testdata", "moves", "multiple.txt"),
		},
		{
			name: "multiple no color",
			moves: []terraform.Move{
				{From: "random_id.original", To: "random_id.refactored"},
				{From: "random_pet.original", To: "random_pet.refactored"},
			},
			noColor: true,
			want:    filepath.Join("testdata", "moves", "multiple-no-color.txt"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Set NoColor for the duration of the test.
			originalNoColor := NoColor
			NoColor = tc.noColor
			defer func() {
				NoColor = originalNoColor
			}()

			actual := Moves(tc.moves)

			if *update {
				stringToFile(t, tc.want, actual)
			}

			want := stringFromFile(t, tc.want)

			const escapeSequence = "\x1b"
			if tc.noColor && strings.Contains(want, escapeSequence) {
				t.Errorf("Moves() output contains espace sequence %q even though color is disabled:\n%q", escapeSequence, want)
			}

			if want != actual {
				t.Errorf("Moves() mismatch\nWant:\n%s\nGot:\n%s", want, actual)
			}
		})
	}
}
