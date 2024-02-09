package pretty

import (
	"fmt"
	"testing"

	"github.com/busser/tfautomv/pkg/golden"
)

func TestBoxItems(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		color string
	}{
		{
			name:  "single item",
			items: []string{"lorem ipsum"},
			color: "red",
		},
		{
			name:  "multiple items",
			items: []string{"lorem ipsum", "lorem ipsum"},
			color: "green",
		},
		{
			name:  "multiline items",
			items: []string{"lorem ipsum\nlorem ipsum", "lorem ipsum\nlorem ipsum"},
			color: "yellow",
		},
		{
			name:  "multiline items with empty lines",
			items: []string{"lorem ipsum\n\nlorem ipsum", "lorem ipsum\n\nlorem ipsum"},
			color: "magenta",
		},
	}

	for _, colorsEnabled := range []bool{true, false} {
		t.Run(fmt.Sprintf("colors %v", colorsEnabled), func(t *testing.T) {
			setupColors(t, colorsEnabled)

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					golden.Equal(t, BoxItems(tt.items, tt.color))
				})
			}
		})
	}
}

func TestBoxSection(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		color   string
	}{
		{
			name:    "without title",
			title:   "",
			content: "lorem ipsum",
			color:   "red",
		},
		{
			name:    "with title",
			title:   "title",
			content: "lorem ipsum",
			color:   "blue",
		},
		{
			name:    "multiline content",
			title:   "title",
			content: "lorem ipsum\nlorem ipsum\nlorem ipsum",
			color:   "green",
		},
		{
			name:    "multiline content with empty lines",
			title:   "title",
			content: "lorem ipsum\n\nlorem ipsum\n\nlorem ipsum",
			color:   "yellow",
		},
	}

	for _, colorsEnabled := range []bool{true, false} {
		t.Run(fmt.Sprintf("colors %v", colorsEnabled), func(t *testing.T) {
			setupColors(t, colorsEnabled)

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					golden.Equal(t, BoxSection(tt.title, tt.content, tt.color))
				})
			}
		})
	}
}
