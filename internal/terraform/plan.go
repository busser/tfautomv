package terraform

import "github.com/padok-team/tfautomv/internal/slices"

type Plan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions []string               `json:"actions"`
	Before  map[string]interface{} `json:"before"`
	After   map[string]interface{} `json:"after"`
}

const (
	CreateAction = "create"
	DeleteAction = "delete"
)

func (p *Plan) NumChanges() int {
	count := 0

	for _, rc := range p.ResourceChanges {
		if slices.Contains(rc.Change.Actions, "create") || slices.Contains(rc.Change.Actions, "delete") {
			count++
		}
	}

	return count
}
