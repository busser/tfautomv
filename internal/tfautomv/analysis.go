package tfautomv

import (
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/busser/tfautomv/internal/flatmap"
	"github.com/busser/tfautomv/internal/slices"
	"github.com/busser/tfautomv/internal/tfautomv/ignore"
)

// Analysis of resources planned for creation or destruction by Terraform and
// whether these resource's types and attributes match.
type Analysis struct {
	// Index of resources based on whether they are planned for creation or
	// destruction.
	CreatedByType   map[string][]*Resource
	DestroyedByType map[string][]*Resource

	// For each resource, its comparison with resources of the same type and the
	// opposite planned operation.
	Comparisons map[*Resource][]Comparison
}

// Resource contains information about a Terraform resource.
type Resource struct {
	// The resource's type.
	Type string

	// The resource's address in Terraform's state.
	Address string

	// The resource's attributes, flattened.
	Attributes map[string]interface{}
}

// AnalysisFromPlan reads the contents of plan and compares resources planned
// for creation with resources planned for destruction of the same type.
// Resources may match, depending on their attributes' values and the rules
// passed to AnalysisFromPlan.
func AnalysisFromPlan(plan *tfjson.Plan, rules []ignore.Rule) (*Analysis, error) {

	// We start with some preprocessing. We identify all ressources planned for
	// creation, or deletion, or both and ignore the rest. We flatten each of
	// those resources' attributes to simplify comparison. We also index each
	// resource by its type, because it only makes sense to compare resources of
	// the same type.

	createdByType := make(map[string][]*Resource)
	destroyedByType := make(map[string][]*Resource)

	for _, c := range plan.ResourceChanges {
		isCreated := slices.Contains(c.Change.Actions, tfjson.ActionCreate)
		isDestroyed := slices.Contains(c.Change.Actions, tfjson.ActionDelete)

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

				comp := Compare(created, destroyed, rules)
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
