package terraform

import (
	"bytes"
	"testing"

	"github.com/busser/tfautomv/pkg/golden"
)

func TestWriteMovedBlocks(t *testing.T) {
	tests := []struct {
		name    string
		moves   []Move
		wantErr bool
	}{
		{
			name: "moves within same workdir",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir",
					ToWorkdir:   "/path/to/workdir",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
			},
		},
		{
			name: "moves between different workdirs",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir1",
					ToWorkdir:   "/path/to/workdir2",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
			},
			wantErr: true,
		},
		{
			name: "multiple moves within same workdir",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir",
					ToWorkdir:   "/path/to/workdir",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
				{
					FromWorkdir: "/path/to/workdir",
					ToWorkdir:   "/path/to/workdir",
					FromAddress: "aws_instance.baz",
					ToAddress:   "aws_instance.qux",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)

			err := WriteMovedBlocks(buf, tt.moves)

			// Check if the error is as expected
			if err != nil && !tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && tt.wantErr {
				t.Fatalf("expected error but got none")
			}

			golden.Equal(t, buf.String())
		})
	}
}

func TestWriteMoveCommands(t *testing.T) {
	tests := []struct {
		name    string
		moves   []Move
		options []Option
	}{
		{
			name: "moves within current workdir",
			moves: []Move{
				{
					FromWorkdir: ".",
					ToWorkdir:   ".",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
			},
		},
		{
			name: "moves within same workdir",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir",
					ToWorkdir:   "/path/to/workdir",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
			},
		},
		{
			name: "moves between different workdirs",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir1",
					ToWorkdir:   "/path/to/workdir2",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
			},
		},
		{
			name: "multiple moves within same workdir",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir",
					ToWorkdir:   "/path/to/workdir",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
				{
					FromWorkdir: "/path/to/workdir",
					ToWorkdir:   "/path/to/workdir",
					FromAddress: "aws_instance.baz",
					ToAddress:   "aws_instance.qux",
				},
			},
		},
		{
			name: "multiple moves between different workdirs",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir1",
					ToWorkdir:   "/path/to/workdir2",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
				{
					FromWorkdir: "/path/to/workdir3",
					ToWorkdir:   "/path/to/workdir4",
					FromAddress: "aws_instance.baz",
					ToAddress:   "aws_instance.qux",
				},
			},
		},
		{
			name: "non-default terraform binary",
			moves: []Move{
				{
					FromWorkdir: "/path/to/workdir1",
					ToWorkdir:   "/path/to/workdir2",
					FromAddress: "aws_instance.foo",
					ToAddress:   "aws_instance.bar",
				},
			},
			options: []Option{
				WithTerraformBin("terragrunt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)

			var settings settings
			settings.apply(append(defaultOptions(), tt.options...))

			// We purposefully don't validate settings here. Whether the
			// settings are valid depends on the testing environment.

			err := writeMoveCommands(buf, tt.moves, settings.terraformBin)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			golden.Equal(t, buf.String())
		})
	}
}

func TestWriteMigrationBlocks(t *testing.T) {
	var buf bytes.Buffer
	moves := []Move{{
		FromWorkdir: "/source",
		ToWorkdir:   "/destination",
		FromAddress: "aws_instance.old",
		ToAddress:   "aws_instance.new",
		ImportID:    "i-123",
	}}
	if err := WriteMigrationBlocks(&buf, moves); err != nil {
		t.Fatal(err)
	}
	want := "removed {\n  from = aws_instance.old\n\n  lifecycle {\n    destroy = false\n  }\n}\n\nimport {\n  to = aws_instance.new\n  id = \"i-123\"\n}\n"
	if got := buf.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWriteImportCommandsImportsBeforeRemoving(t *testing.T) {
	var buf bytes.Buffer
	moves := []Move{{
		FromWorkdir: "/source",
		ToWorkdir:   "/destination",
		FromAddress: "aws_instance.old",
		ToAddress:   "aws_instance.new",
		ImportID:    "i-123",
	}}
	if err := WriteImportCommands(&buf, moves, "terraform"); err != nil {
		t.Fatal(err)
	}
	want := "terraform -chdir=\"/destination\" import \"aws_instance.new\" \"i-123\"\nterraform -chdir=\"/source\" state rm \"aws_instance.old\"\n"
	if got := buf.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
