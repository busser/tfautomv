package pretty

import (
	"fmt"
	"testing"

	"github.com/busser/tfautomv/pkg/engine"
	"github.com/busser/tfautomv/pkg/golden"
)

func TestSummary(t *testing.T) {
	tests := []struct {
		name      string
		verbosity int
	}{
		{
			name:      "verbosity 0",
			verbosity: 0,
		},
		{
			name:      "verbosity 1",
			verbosity: 1,
		},
		{
			name:      "verbosity 2",
			verbosity: 2,
		},
		{
			name:      "verbosity 3",
			verbosity: 3,
		},
	}

	moves, comparisons := testData()

	for _, colorsEnabled := range []bool{true, false} {
		t.Run(fmt.Sprintf("colors %v", colorsEnabled), func(t *testing.T) {
			setupColors(t, colorsEnabled)

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					summarizer := NewSummarizer(moves, comparisons, tt.verbosity)
					golden.Equal(t, summarizer.Summary())
				})
			}
		})
	}
}

func testData() ([]engine.Move, []engine.ResourceComparison) {
	return testDataMoves(), testDataComparisons()
}

func testDataResources() map[string]engine.Resource {
	return map[string]engine.Resource{

		// Resources Terraform plans to delete:

		"alpha": {
			ModuleID: "demo/module-a",
			Type:     "random_pet",
			Address:  "random_pet.alpha",
			Attributes: map[string]any{
				"id":        "alpha-unique-mayfly",
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "alpha",
				"separator": "-",
			},
		},

		"bravo": {
			ModuleID: "demo/module-a",
			Type:     "random_pet",
			Address:  "random_pet.bravo",
			Attributes: map[string]any{
				"id":        "bravo-smooth-weevil",
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "bravo",
				"separator": "-",
			},
		},

		"charlie": {
			ModuleID: "demo/module-a",
			Type:     "random_pet",
			Address:  "random_pet.charlie",
			Attributes: map[string]any{
				"id":        "charlie-stirring-guinea",
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "charlie",
				"separator": "-",
			},
		},

		"delta": {
			ModuleID: "demo/module-a",
			Type:     "random_pet",
			Address:  "random_pet.delta",
			Attributes: map[string]any{
				"id":        "delta-super-roughy",
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "delta",
				"separator": "-",
			},
		},

		"echo": {
			ModuleID: "demo/module-a",
			Type:     "random_pet",
			Address:  "random_pet.echo",
			Attributes: map[string]any{
				"id":        "echo-allowed-camel",
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "echo",
				"separator": "-",
			},
		},

		// Resources Terraform plans to create:

		"alice": {
			ModuleID: "demo/module-a",
			Type:     "random_pet",
			Address:  "random_pet.alice",
			Attributes: map[string]any{
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "alpha",
				"separator": "-",
			},
		},

		"bob": {
			ModuleID: "demo/module-b",
			Type:     "random_pet",
			Address:  "random_pet.bob",
			Attributes: map[string]any{
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "bravo",
				"separator": "-",
			},
		},

		"carol": {
			ModuleID: "demo/module-b",
			Type:     "random_pet",
			Address:  "random_pet.carol",
			Attributes: map[string]any{
				"keepers":   nil,
				"length":    float64(3),
				"prefix":    "charlie",
				"separator": "-",
			},
		},

		"david": {
			ModuleID: "demo/module-b",
			Type:     "random_pet",
			Address:  "random_pet.david",
			Attributes: map[string]any{
				"keepers":   nil,
				"length":    float64(3),
				"prefix":    "delta",
				"separator": "_",
			},
		},

		"daniel": {
			ModuleID: "demo/module-b",
			Type:     "random_pet",
			Address:  "random_pet.daniel",
			Attributes: map[string]any{
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "delta",
				"separator": "-",
			},
		},

		"felix": {
			ModuleID: "demo/module-b",
			Type:     "random_pet",
			Address:  "random_pet.felix",
			Attributes: map[string]any{
				"keepers":   nil,
				"length":    float64(2),
				"prefix":    "foxtrot",
				"separator": "-",
			},
		},
	}

}

func testDataMoves() []engine.Move {
	resources := testDataResources()

	return []engine.Move{
		{
			SourceModule:       resources["alpha"].ModuleID,
			DestinationModule:  resources["alice"].ModuleID,
			SourceAddress:      resources["alpha"].Address,
			DestinationAddress: resources["alice"].Address,
		},
		{
			SourceModule:       resources["bravo"].ModuleID,
			DestinationModule:  resources["bob"].ModuleID,
			SourceAddress:      resources["bravo"].Address,
			DestinationAddress: resources["bob"].Address,
		},
		{
			SourceModule:       resources["charlie"].ModuleID,
			DestinationModule:  resources["carol"].ModuleID,
			SourceAddress:      resources["charlie"].Address,
			DestinationAddress: resources["carol"].Address,
		},
	}
}

func testDataComparisons() []engine.ResourceComparison {
	resources := testDataResources()

	return []engine.ResourceComparison{
		{
			ToCreate:              resources["alice"],
			ToDelete:              resources["alpha"],
			MatchingAttributes:    []string{"length", "prefix", "separator"},
			MismatchingAttributes: nil,
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["alice"],
			ToDelete:              resources["bravo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["alice"],
			ToDelete:              resources["charlie"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["alice"],
			ToDelete:              resources["delta"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["alice"],
			ToDelete:              resources["echo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["bob"],
			ToDelete:              resources["alpha"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["bob"],
			ToDelete:              resources["bravo"],
			MatchingAttributes:    []string{"length", "prefix", "separator"},
			MismatchingAttributes: nil,
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["bob"],
			ToDelete:              resources["charlie"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["bob"],
			ToDelete:              resources["delta"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["bob"],
			ToDelete:              resources["echo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["carol"],
			ToDelete:              resources["alpha"],
			MatchingAttributes:    []string{"separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length"},
		},
		{
			ToCreate:              resources["carol"],
			ToDelete:              resources["bravo"],
			MatchingAttributes:    []string{"separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length"},
		},
		{
			ToCreate:              resources["carol"],
			ToDelete:              resources["charlie"],
			MatchingAttributes:    []string{"prefix", "separator"},
			MismatchingAttributes: nil,
			IgnoredAttributes:     []string{"length"},
		},
		{
			ToCreate:              resources["carol"],
			ToDelete:              resources["delta"],
			MatchingAttributes:    []string{"separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length"},
		},
		{
			ToCreate:              resources["carol"],
			ToDelete:              resources["echo"],
			MatchingAttributes:    []string{"separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length"},
		},
		{
			ToCreate:              resources["daniel"],
			ToDelete:              resources["alpha"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["daniel"],
			ToDelete:              resources["bravo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["daniel"],
			ToDelete:              resources["charlie"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["daniel"],
			ToDelete:              resources["delta"],
			MatchingAttributes:    []string{"length", "prefix", "separator"},
			MismatchingAttributes: nil,
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["daniel"],
			ToDelete:              resources["echo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["david"],
			ToDelete:              resources["alpha"],
			MatchingAttributes:    nil,
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length", "separator"},
		},
		{
			ToCreate:              resources["david"],
			ToDelete:              resources["bravo"],
			MatchingAttributes:    nil,
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length", "separator"},
		},
		{
			ToCreate:              resources["david"],
			ToDelete:              resources["charlie"],
			MatchingAttributes:    nil,
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length", "separator"},
		},
		{
			ToCreate:              resources["david"],
			ToDelete:              resources["delta"],
			MatchingAttributes:    []string{"prefix"},
			MismatchingAttributes: nil,
			IgnoredAttributes:     []string{"length", "separator"},
		},
		{
			ToCreate:              resources["david"],
			ToDelete:              resources["echo"],
			MatchingAttributes:    nil,
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     []string{"length", "separator"},
		},
		{
			ToCreate:              resources["felix"],
			ToDelete:              resources["alpha"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["felix"],
			ToDelete:              resources["bravo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["felix"],
			ToDelete:              resources["charlie"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["felix"],
			ToDelete:              resources["delta"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
		{
			ToCreate:              resources["felix"],
			ToDelete:              resources["echo"],
			MatchingAttributes:    []string{"length", "separator"},
			MismatchingAttributes: []string{"prefix"},
			IgnoredAttributes:     nil,
		},
	}
}
