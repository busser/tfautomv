package terraform

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
