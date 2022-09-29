package format

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/padok-team/tfautomv/internal/tfautomv"
)

func TestAnalysis(t *testing.T) {
	tt := []struct {
		name string

		analysis *tfautomv.Analysis
		noColor  bool

		want string
	}{
		{
			name:     "empty",
			analysis: &tfautomv.Analysis{},
			noColor:  false,
			want:     filepath.Join("testdata", "analysis", "empty.txt"),
		},
		{
			name:     "empty no color",
			analysis: &tfautomv.Analysis{},
			noColor:  true,
			want:     filepath.Join("testdata", "analysis", "empty-no-color.txt"),
		},
		// TODO(busser): add test cases for a complete analysis.
		// {
		// 	name:     "complete",
		// 	analysis: &tfautomv.Analysis{},
		// 	noColor:  false,
		// 	want:     filepath.Join("testdata", "analysis", "complete.txt"),
		// },
		// {
		// 	name:     "complete no color",
		// 	analysis: &tfautomv.Analysis{},
		// 	noColor:  true,
		// 	want:     filepath.Join("testdata", "analysis", "complete-no-color.txt"),
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Set NoColor for the duration of the test.
			originalNoColor := NoColor
			NoColor = tc.noColor
			defer func() {
				NoColor = originalNoColor
			}()

			actual := Analysis(tc.analysis)

			if *update {
				stringToFile(t, tc.want, actual)
			}

			want := stringFromFile(t, tc.want)

			const escapeSequence = "\x1b"
			if tc.noColor && strings.Contains(want, escapeSequence) {
				t.Errorf("Analysis() output contains espace sequence %q even though color is disabled:\n%q", escapeSequence, want)
			}

			if want != actual {
				t.Errorf("Analysis() mismatch\nWant:\n%s\nGot:\n%s", want, actual)
			}
		})
	}
}
