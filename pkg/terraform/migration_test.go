package terraform

import (
	"bytes"
	"testing"

	"github.com/busser/tfautomv/pkg/engine"
)

func TestWriteMigrationBlocksForPlan(t *testing.T) {
	migrations := []engine.Migration{
		{SourcePlan: "source", SourceAddress: "aws_instance.old", DestinationPlan: "destination", DestinationAddr: "aws_instance.new", ImportID: "i-123"},
		{SourcePlan: "same", SourceAddress: "aws_instance.before", DestinationPlan: "same", DestinationAddr: "aws_instance.after"},
	}
	var output bytes.Buffer
	if err := WriteMigrationBlocksForPlan(&output, migrations, "destination"); err != nil {
		t.Fatal(err)
	}
	want := "import {\n  to = aws_instance.new\n  id = \"i-123\"\n}\n"
	if output.String() != want {
		t.Fatalf("got %q, want %q", output.String(), want)
	}
}
