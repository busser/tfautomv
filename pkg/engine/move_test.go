package engine

import (
	"slices"
	"testing"
)

func TestDetermineMoves(t *testing.T) {
	tests := []struct {
		name        string
		comparisons []ResourceComparison

		wantMoves []Move
	}{
		{
			// This is the base case: there is a single logical move.
			name: "single resource with single match",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
			},
			wantMoves: []Move{
				{
					SourceModule:       "that_module",
					SourceAddress:      "that_address",
					DestinationModule:  "this_module",
					DestinationAddress: "this_address",
				},
			},
		},

		{
			// When there are no valid moves to make, we don't make any moves.
			name: "single resource with no match",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: []string{"foo", "bar"},
				},
			},
			wantMoves: nil,
		},

		{
			// When there is a single logical move, we make it.
			name: "multiple resources with single match each",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("other_module", "", "other_address"),
					MismatchingAttributes: []string{"foo", "bar"},
				},
			},
			wantMoves: []Move{
				{
					SourceModule:       "that_module",
					SourceAddress:      "that_address",
					DestinationModule:  "this_module",
					DestinationAddress: "this_address",
				},
			},
		},

		{
			// When there are multiple logical moves for a single resource, we
			// don't make any moves.
			name: "single resource with multiple matches",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("other_module", "", "other_address"),
					MismatchingAttributes: nil,
				},
			},
			wantMoves: nil,
		},

		{
			// When there are no valid moves to make, we don't make any moves.
			name: "single resource with multiple comparisons and no match",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: []string{"foo", "bar"},
				},
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("other_module", "", "other_address"),
					MismatchingAttributes: []string{"baz", "tuq"},
				},
			},
			wantMoves: nil,
		},

		{
			// When working with multiple resources, we only move those that
			// have a single valid move.
			name: "multiple resources one with single match other no match",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
				{
					ToCreate:              dummyResource("other_module", "", "other_address"),
					ToDelete:              dummyResource("another_module", "", "another_address"),
					MismatchingAttributes: []string{"foo", "bar"},
				},
			},
			wantMoves: []Move{
				{
					SourceModule:       "that_module",
					SourceAddress:      "that_address",
					DestinationModule:  "this_module",
					DestinationAddress: "this_address",
				},
			},
		},

		{
			// When working with multiple resources, we only move those that
			// have a single valid move.
			name: "one resource with single match another with multiple matches",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
				{
					ToCreate:              dummyResource("other_module", "", "other_module"),
					ToDelete:              dummyResource("another_module", "", "another_module"),
					MismatchingAttributes: nil,
				},
				{
					ToCreate:              dummyResource("other_module", "", "other_module"),
					ToDelete:              dummyResource("yet_another_module", "", "yet_another_address"),
					MismatchingAttributes: nil,
				},
			},
			wantMoves: []Move{
				{
					SourceModule:       "that_module",
					SourceAddress:      "that_address",
					DestinationModule:  "this_module",
					DestinationAddress: "this_address",
				},
			},
		},

		{
			// When working with multiple resources, we don't move any of them
			// if any of the invalid moves are for the same resource.
			name: "two resources have a single match with the same resource",
			comparisons: []ResourceComparison{
				{
					ToCreate:              dummyResource("this_module", "", "this_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
				{
					ToCreate:              dummyResource("other_module", "", "other_address"),
					ToDelete:              dummyResource("that_module", "", "that_address"),
					MismatchingAttributes: nil,
				},
			},
			wantMoves: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := DetermineMoves(tt.comparisons)
			if !slices.Equal(actual, tt.wantMoves) {
				t.Errorf("got %v, want %v", actual, tt.wantMoves)
			}
		})
	}
}
