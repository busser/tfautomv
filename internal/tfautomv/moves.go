package tfautomv

import (
	"sort"

	"github.com/busser/tfautomv/internal/terraform"
)

// MovesFromAnalysis identifies which resources should be moved from one
// address to the other in Terraform's state, based on the provided analysis.
func MovesFromAnalysis(analysis *Analysis) []terraform.Move {

	// We choose to move a resource planned for destruction to a resource
	// planned for creation if and only if the resources match each other and
	// only each other.

	matchCountByResource := make(map[*Resource]int)
	for res, comps := range analysis.Comparisons {
		for _, c := range comps {
			if c.IsMatch() {
				matchCountByResource[res]++
			}
		}
	}

	var moves []terraform.Move

	for _, resources := range analysis.CreatedByType {
		for _, created := range resources {
			if matchCountByResource[created] != 1 {
				continue
			}

			var destroyed *Resource
			for _, comp := range analysis.Comparisons[created] {
				if comp.IsMatch() {
					destroyed = comp.Destroyed
				}
			}

			if matchCountByResource[destroyed] != 1 {
				continue
			}

			m := terraform.Move{
				From: destroyed.Address,
				To:   created.Address,
			}
			moves = append(moves, m)
		}
	}

	sort.Sort(terraform.InOrder(moves))

	return moves
}
