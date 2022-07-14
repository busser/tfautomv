package tfautomv

import (
	"path/filepath"

	"github.com/padok-team/tfautomv/internal/flatmap"
	"github.com/padok-team/tfautomv/internal/terraform"
)

func GenerateMovedBlocks(workdir string) (int, error) {
	refactored := terraform.NewRunner(workdir)

	err := refactored.Init()
	if err != nil {
		return 0, err
	}

	plan, err := refactored.Plan()
	if err != nil {
		return 0, err
	}

	blocks, err := moveBlocksFromPlanBasedOnAttributes(plan)
	if err != nil {
		return 0, err
	}

	err = blocks.AppendTo(filepath.Join(workdir, "moves.tf"))
	if err != nil {
		return 0, err
	}

	return len(blocks.Moves), nil
}

func moveBlocksFromPlanBasedOnAttributes(plan *terraform.Plan) (terraform.MoveBlocks, error) {
	/*
		1. Catégoriser les ressources supprimées par types.
		2. Pour chaque ressource créée, comparer les valeurs connues de ses
		   attributs aux valeurs des mêmes attributs dans chaque ressource
		   supprimée du même type.
		3. Pour chaque ressource créée, construire une liste de matchs.
		4. Pour chaque ressource créée, si un seul match, générer un block.
	*/

	createdByType := make(map[string][]*resource)
	deletedByType := make(map[string][]*resource)

	for _, rc := range plan.ResourceChanges {

		if containsString(rc.Change.Actions, "create") {
			flatAttributes, err := flatmap.Flatten(rc.Change.After)
			if err != nil {
				return terraform.MoveBlocks{}, err
			}

			res := resource{
				typ:        rc.Type,
				address:    rc.Address,
				attributes: flatAttributes,
			}

			createdByType[res.typ] = append(createdByType[res.typ], &res)
		}

		if containsString(rc.Change.Actions, "delete") {
			flatAttributes, err := flatmap.Flatten(rc.Change.Before)
			if err != nil {
				return terraform.MoveBlocks{}, err
			}

			res := resource{
				typ:        rc.Type,
				address:    rc.Address,
				attributes: flatAttributes,
			}

			deletedByType[res.typ] = append(deletedByType[res.typ], &res)
		}
	}

	var blocks terraform.MoveBlocks

	for typ := range createdByType {
		var matches []resourceMatch

		for _, created := range createdByType[typ] {
			for _, deleted := range deletedByType[typ] {
				if createdMatchesDeleted(created, deleted) {
					matches = append(matches, resourceMatch{created, deleted})
				}
			}
		}

		matchesByResource := make(map[*resource][]resourceMatch)
		for _, m := range matches {
			matchesByResource[m.created] = append(matchesByResource[m.created], m)
			matchesByResource[m.deleted] = append(matchesByResource[m.deleted], m)
		}

		for _, created := range createdByType[typ] {
			if len(matchesByResource[created]) != 1 {
				continue
			}

			deleted := matchesByResource[created][0].deleted
			if len(matchesByResource[deleted]) != 1 {
				continue
			}

			move := terraform.Moved{From: deleted.address, To: created.address}
			blocks.Moves = append(blocks.Moves, move)
		}
	}

	return blocks, nil
}

type resource struct {
	typ        string
	address    string
	attributes map[string]interface{}
}

type resourceMatch struct {
	created, deleted *resource
}

func createdMatchesDeleted(created, deleted *resource) bool {
	if created.address == deleted.address {
		return false
	}

	for k, vc := range created.attributes {
		if vc == nil {
			continue
		}
		vd, ok := deleted.attributes[k]
		if !ok {
			return false
		}
		if vc != vd {
			return false
		}
	}
	return true
}

func containsString(vv []string, v string) bool {
	for _, s := range vv {
		if s == v {
			return true
		}
	}
	return false
}

func CountChanges(plan *terraform.Plan) int {
	count := 0

	for _, rc := range plan.ResourceChanges {
		if containsString(rc.Change.Actions, "create") || containsString(rc.Change.Actions, "delete") {
			count++
		}
	}

	return count
}
