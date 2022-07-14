package terraform

type Plan struct {
	ResourceChanges []struct {
		Address string `json:"address"`
		Type    string `json:"type"`
		Change  struct {
			Actions []string               `json:"actions"`
			Before  map[string]interface{} `json:"before"`
			After   map[string]interface{} `json:"after"`
		} `json:"change"`
	} `json:"resource_changes"`
}
