package engine

import (
	"fmt"
	"sort"
)

// Migration describes a proposed state transfer between two plan entries.
// The plan IDs are deliberately kept separate from resource addresses so a
// dry-run can explain exactly which state file is affected.
type Migration struct {
	SourcePlan      string
	SourceAddress   string
	DestinationPlan string
	DestinationAddr string
	ResourceType    string
	ImportID        string
}

// MigrationAnalysis contains all safe matches and the diagnostics that kept
// potentially unsafe matches from being emitted.
type MigrationAnalysis struct {
	Migrations       []Migration
	AmbiguousDeletes []Resource
	AmbiguousCreates []Resource
	UnmatchedDeletes []Resource
	UnmatchedCreates []Resource
}

// AnalyzeMigrations finds one-to-one resource matches across all supplied
// plans. A match is accepted only when both endpoints have exactly one valid
// candidate, which prevents accidental transfers when identical resources are
// present in more than one stack.
func AnalyzeMigrations(plans []Plan, rules []Rule) MigrationAnalysis {
	comparisons := CompareAll(MergePlans(plans), rules)
	createMatches := make(map[string][]Resource)
	deleteMatches := make(map[string][]Resource)
	for _, comparison := range comparisons {
		if !comparison.IsMatch() {
			continue
		}
		createMatches[comparison.ToCreate.ID()] = append(createMatches[comparison.ToCreate.ID()], comparison.ToDelete)
		deleteMatches[comparison.ToDelete.ID()] = append(deleteMatches[comparison.ToDelete.ID()], comparison.ToCreate)
	}

	analysis := MigrationAnalysis{}
	for _, plan := range plans {
		for _, resource := range plan.ToDelete {
			matches := deleteMatches[resource.ID()]
			switch len(matches) {
			case 0:
				analysis.UnmatchedDeletes = append(analysis.UnmatchedDeletes, resource)
			case 1:
				if len(createMatches[matches[0].ID()]) != 1 {
					analysis.AmbiguousDeletes = append(analysis.AmbiguousDeletes, resource)
				}
			default:
				analysis.AmbiguousDeletes = append(analysis.AmbiguousDeletes, resource)
			}
		}
		for _, resource := range plan.ToCreate {
			matches := createMatches[resource.ID()]
			switch len(matches) {
			case 0:
				analysis.UnmatchedCreates = append(analysis.UnmatchedCreates, resource)
			case 1:
				if len(deleteMatches[matches[0].ID()]) != 1 {
					analysis.AmbiguousCreates = append(analysis.AmbiguousCreates, resource)
				}
			default:
				analysis.AmbiguousCreates = append(analysis.AmbiguousCreates, resource)
			}
		}
	}

	for _, plan := range plans {
		for _, resource := range plan.ToDelete {
			matches := deleteMatches[resource.ID()]
			if len(matches) != 1 || len(createMatches[matches[0].ID()]) != 1 {
				continue
			}
			to := matches[0]
			analysis.Migrations = append(analysis.Migrations, Migration{
				SourcePlan: planID(resource), SourceAddress: resource.Address,
				DestinationPlan: planID(to), DestinationAddr: to.Address,
				ResourceType: resource.Type, ImportID: importID(resource),
			})
		}
	}

	sort.Slice(analysis.Migrations, func(i, j int) bool {
		a, b := analysis.Migrations[i], analysis.Migrations[j]
		if a.SourcePlan != b.SourcePlan {
			return a.SourcePlan < b.SourcePlan
		}
		if a.DestinationPlan != b.DestinationPlan {
			return a.DestinationPlan < b.DestinationPlan
		}
		if a.SourceAddress != b.SourceAddress {
			return a.SourceAddress < b.SourceAddress
		}
		return a.DestinationAddr < b.DestinationAddr
	})
	return analysis
}

// ValidateNetZero rejects plans containing an unequal number of planned
// creates and deletes. Callers can combine this with ambiguity checks to make
// migration mode fail closed before any generated configuration is written.
func (a MigrationAnalysis) ValidateNetZero(plans []Plan) error {
	var creates, deletes int
	for _, plan := range plans {
		creates += len(plan.ToCreate)
		deletes += len(plan.ToDelete)
	}
	if creates != deletes {
		return fmt.Errorf("net-zero validation failed: %d creates and %d deletes", creates, deletes)
	}
	if len(a.AmbiguousCreates) > 0 || len(a.AmbiguousDeletes) > 0 {
		return fmt.Errorf("migration validation failed: %d ambiguous creates and %d ambiguous deletes", len(a.AmbiguousCreates), len(a.AmbiguousDeletes))
	}
	if len(a.UnmatchedCreates) > 0 || len(a.UnmatchedDeletes) > 0 {
		return fmt.Errorf("migration validation failed: %d unmatched creates and %d unmatched deletes", len(a.UnmatchedCreates), len(a.UnmatchedDeletes))
	}
	return nil
}

func planID(resource Resource) string { return resource.ModuleID }

func importID(resource Resource) string {
	if id, ok := resource.Attributes["id"].(string); ok {
		return id
	}
	return ""
}
