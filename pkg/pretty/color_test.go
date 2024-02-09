package pretty

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const colorEscape = "\x1b"

func setupColors(t *testing.T, enabled bool) {
	t.Helper()

	previous := colorsDisabled
	colorsDisabled = !enabled
	t.Cleanup(func() {
		colorsDisabled = previous
	})
}

func TestColor(t *testing.T) {
	setupColors(t, true)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "red",
			in:   "[red]foo",
			want: colorEscape + "[31mfoo" + colorEscape + "[0m",
		},
		{
			name: "bold",
			in:   "[bold]foo",
			want: colorEscape + "[1mfoo" + colorEscape + "[0m",
		},
		{
			name: "red bold",
			in:   "[red][bold]foo",
			want: colorEscape + "[31m" + colorEscape + "[1mfoo" + colorEscape + "[0m",
		},
		{
			name: "reset",
			in:   "[red]foo[reset]bar",
			want: colorEscape + "[31mfoo" + colorEscape + "[0mbar" + colorEscape + "[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Color(tt.in)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNoColor(t *testing.T) {
	setupColors(t, false)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "red",
			in:   "[red]foo",
			want: "foo",
		},
		{
			name: "bold",
			in:   "[bold]foo",
			want: "foo",
		},
		{
			name: "red bold",
			in:   "[red][bold]foo",
			want: "foo",
		},
		{
			name: "reset",
			in:   "[red]foo[reset]bar",
			want: "foobar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Color(tt.in)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
