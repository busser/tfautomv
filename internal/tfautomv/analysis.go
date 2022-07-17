package tfautomv

import (
	"github.com/padok-team/tfautomv/internal/flatmap"
	"github.com/padok-team/tfautomv/internal/terraform"
)

// Analysis of resources planned for creation or destruction by Terraform and
// whether these resource's types and attributes match.
type Analysis struct {
	// For each resource, its comparison with resources of the same type and the
	// opposite planned operation.
	Comparisons map[*Resource][]Comparison

	// Index of resources based on whether they are planned for creation or
	// destruction.
	CreatedByType   map[string][]*Resource
	DestroyedByType map[string][]*Resource
}

func AnalysisFromPlan(plan *terraform.Plan) (*Analysis, error) {

	// We start with some preprocessing. We identify all ressources planned for
	// creation, or deletion, or both and ignore the rest. We flatten each of
	// those resources' attributes to simplify comparison. We also index each
	// resource by its type, because it only makes sense to compare resources of
	// the same type.

	createdByType := make(map[string][]*Resource)
	destroyedByType := make(map[string][]*Resource)

	for _, c := range plan.ResourceChanges {
		isCreated := sliceContains(c.Change.Actions, "create")
		isDestroyed := sliceContains(c.Change.Actions, "delete")

		if !isCreated && !isDestroyed {
			continue
		}

		if isCreated {
			flatAttributes, err := flatmap.Flatten(c.Change.After)
			if err != nil {
				return nil, err
			}

			r := Resource{
				Type:       c.Type,
				Address:    c.Address,
				Attributes: flatAttributes,
			}

			createdByType[r.Type] = append(createdByType[r.Type], &r)
		}

		if isDestroyed {
			flatAttributes, err := flatmap.Flatten(c.Change.Before)
			if err != nil {
				return nil, err
			}

			r := Resource{
				Type:       c.Type,
				Address:    c.Address,
				Attributes: flatAttributes,
			}

			destroyedByType[r.Type] = append(destroyedByType[r.Type], &r)
		}
	}

	// Then, we compare all resources planned for creation will all resources
	// planned for destruction of the same type.

	comparisons := make(map[*Resource][]Comparison)
	for typ := range createdByType {
		for _, created := range createdByType[typ] {
			for _, destroyed := range destroyedByType[typ] {
				// If both resources have the same address, then no move is
				// possible. This can happen when a resource requires changes or
				// has been tainted, for example.
				if created.Address == destroyed.Address {
					continue
				}

				comp := Compare(created, destroyed)
				comparisons[created] = append(comparisons[created], comp)
				comparisons[destroyed] = append(comparisons[destroyed], comp)
			}
		}
	}

	analysis := Analysis{
		Comparisons:     comparisons,
		CreatedByType:   createdByType,
		DestroyedByType: destroyedByType,
	}

	return &analysis, nil
}

func sliceContains[T comparable](s []T, v T) bool {
	for i := range s {
		if s[i] == v {
			return true
		}
	}
	return false
}
