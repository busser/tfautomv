package golden

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func init() {
	// For compatibility with other packages that also define an -update parameter, only define the
	// flag if it's not already defined.
	if updateFlag := flag.Lookup("update"); updateFlag == nil {
		flag.Bool("update", false, "update .golden files, leaving unused)")
	}
}

func update() bool {
	return flag.Lookup("update").Value.(flag.Getter).Get().(bool)
}

// Compare the given value to the golden file, and update it if necessary and
// the user has specified the -update flag.
func Equal(t *testing.T, value string) {
	t.Helper()

	goldenFile := filepath.Join("testdata", filepath.FromSlash(t.Name()+".golden"))
	wantBytes, err := os.ReadFile(goldenFile)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to read golden file: %s", err)
	}
	want := string(wantBytes)

	if diff := cmp.Diff(want, value); diff != "" {
		if update() {
			if err := os.MkdirAll(filepath.Dir(goldenFile), 0755); err != nil {
				t.Fatalf("failed to create golden file directory: %s", err)
			}
			if err := os.WriteFile(goldenFile, []byte(value), 0644); err != nil {
				t.Fatalf("failed to update golden file: %s", err)
			}
		} else {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
