package format

import (
	"flag"
	"os"
	"testing"
)

// The package's functions produce human-readable output meant to be printed to
// a terminal. In these tests, the output is compared with the contents of what
// we call golden files. These files may include encoding to print text with
// color and other formatting. To view the result of these encodings, use the
// `cat` command to have your terminal render the file's contents in color.

var update = flag.Bool("update", false, "update golden files")

func stringFromFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func stringToFile(t *testing.T, path string, data string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
}
