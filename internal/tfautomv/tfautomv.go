package tfautomv

import (
	"github.com/padok-team/tfautomv/internal/terraform"
)

func GenerateReport(workdir string) (*Report, error) {
	refactored := terraform.NewRunner(workdir)

	err := refactored.Init()
	if err != nil {
		return nil, err
	}

	plan, err := refactored.Plan()
	if err != nil {
		return nil, err
	}

	analysis, err := AnalysisFromPlan(plan)
	if err != nil {
		return nil, err
	}

	moves := MovesFromAnalysis(analysis)

	report := Report{
		InitialPlan: plan,
		Analysis:    analysis,
		Moves:       moves,
	}

	return &report, nil
}

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

			m := terraform.Move{
				From: destroyed.Address,
				To:   created.Address,
			}
			moves = append(moves, m)
		}
	}

	return moves
}

type Resource struct {
	Type       string
	Address    string
	Attributes map[string]interface{}
}

func CountChanges(plan *terraform.Plan) int {
	count := 0

	for _, rc := range plan.ResourceChanges {
		if sliceContains(rc.Change.Actions, "create") || sliceContains(rc.Change.Actions, "delete") {
			count++
		}
	}

	return count
}
