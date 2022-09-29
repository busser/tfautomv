package format

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	tt := []struct {
		name string

		err     error
		noColor bool

		want string
	}{
		{
			name:    "simple",
			err:     errors.New("simple error"),
			noColor: false,
			want:    filepath.Join("testdata", "error", "simple.txt"),
		},
		{
			name:    "simple no color",
			err:     errors.New("simple error"),
			noColor: true,
			want:    filepath.Join("testdata", "error", "simple-no-color.txt"),
		},
		{
			name:    "multiline",
			err:     errors.New("multiple errors:\n  - first error\n  - second error"),
			noColor: false,
			want:    filepath.Join("testdata", "error", "multiline.txt"),
		},
		{
			name:    "multiline no color",
			err:     errors.New("multiple errors:\n  - first error\n  - second error"),
			noColor: true,
			want:    filepath.Join("testdata", "error", "multiline-no-color.txt"),
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

			actual := Error(tc.err)

			if *update {
				stringToFile(t, tc.want, actual)
			}

			want := stringFromFile(t, tc.want)

			const escapeSequence = "\x1b"
			if tc.noColor && strings.Contains(want, escapeSequence) {
				t.Errorf("Error() output contains espace sequence %q even though color is disabled:\n%q", escapeSequence, want)
			}

			if want != actual {
				t.Errorf("Error() mismatch\nWant:\n%s\nGot:\n%s", want, actual)
			}
		})
	}
}
