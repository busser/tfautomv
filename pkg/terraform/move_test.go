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

			err := WriteMoveCommands(buf, tt.moves, tt.options...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			golden.Equal(t, buf.String())
		})
	}
}
