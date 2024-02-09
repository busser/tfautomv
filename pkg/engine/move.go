package engine

import "sort"

// A Move represents a Terraform resource that we should move from one address
// to another. A resource can be moved within the same module or to a different
// module.
type Move struct {
	// The module the resource is being moved from.
	SourceModule string
	// The module the resource is being moved to. This is equal to SourceModule
	// when the resource is being moved within the same module.
	DestinationModule string

	// The resource's address before the move.
	SourceAddress string
	// The resource's address after the move.
	DestinationAddress string
}

func DetermineMoves(comparisons []ResourceComparison) []Move {

	// We choose to move a resource planned for deletion to a resource planned
	// for creation if and only if the resources match each other and
	// only each other.

	planToCreateMatchCount := make(map[string]int)
	planToDeleteMatchCount := make(map[string]int)
	for _, comparison := range comparisons {
		if comparison.IsMatch() {
			planToCreateMatchCount[comparison.ToCreate.ID()]++
			planToDeleteMatchCount[comparison.ToDelete.ID()]++
		}
	}

	var moves []Move

	for _, comparison := range comparisons {
		if !comparison.IsMatch() {
			continue
		}

		if planToCreateMatchCount[comparison.ToCreate.ID()] != 1 {
			continue
		}

		if planToDeleteMatchCount[comparison.ToDelete.ID()] != 1 {
			continue
		}

		m := Move{
			SourceModule:       comparison.ToDelete.ModuleID,
			SourceAddress:      comparison.ToDelete.Address,
			DestinationModule:  comparison.ToCreate.ModuleID,
			DestinationAddress: comparison.ToCreate.Address,
		}
		moves = append(moves, m)
	}

	// Sort the moves so that the result is deterministic.
	sortMoves(moves)

	return moves
}

func sortMoves(moves []Move) {
	sort.Slice(moves, func(i, j int) bool {
		a, b := moves[i], moves[j]

		switch {
		case a.SourceModule != b.SourceModule:
			return a.SourceModule < b.SourceModule

		case a.DestinationModule != b.DestinationModule:
			return a.DestinationModule < b.DestinationModule

		case a.SourceAddress != b.SourceAddress:
			return a.SourceAddress < b.SourceAddress

		case a.DestinationAddress != b.DestinationAddress:
			return a.DestinationAddress < b.DestinationAddress

		default:
			return false // a == b so it doesn't matter what we return here
		}
	})
}
