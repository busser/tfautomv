package format

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDone(t *testing.T) {
	tt := []struct {
		name string

		msg     string
		noColor bool

		want string
	}{
		{
			name:    "simple",
			msg:     "simple message",
			noColor: false,
			want:    filepath.Join("testdata", "done", "simple.txt"),
		},
		{
			name:    "simple no color",
			msg:     "simple message",
			noColor: true,
			want:    filepath.Join("testdata", "done", "simple-no-color.txt"),
		},
		{
			name:    "multiline",
			msg:     "multiple messages:\n  - first message\n  - second message",
			noColor: false,
			want:    filepath.Join("testdata", "done", "multiline.txt"),
		},
		{
			name:    "multiline no color",
			msg:     "multiple messages:\n  - first message\n  - second message",
			noColor: true,
			want:    filepath.Join("testdata", "done", "multiline-no-color.txt"),
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

			actual := Done(tc.msg)

			if *update {
				stringToFile(t, tc.want, actual)
			}

			want := stringFromFile(t, tc.want)

			const escapeSequence = "\x1b"
			if tc.noColor && strings.Contains(want, escapeSequence) {
				t.Errorf("Done() output contains espace sequence %q even though color is disabled:\n%q", escapeSequence, want)
			}

			if want != actual {
				t.Errorf("Done() mismatch\nWant:\n%s\nGot:\n%s", want, actual)
			}
		})
	}
}
