package tfguard

import "encoding/json"

// PlanRepresentation is an unmarshalled terraform plan
type PlanRepresentation struct {
	FormatVersion   string           `json:"format_version"`
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address       string               `json:"address"`
	ModuleAddress string               `json:"module_address"`
	Mode          string               `json:"mode"`
	Type          string               `json:"type"`
	Name          string               `json:"name"`
	Index         interface{}          `json:"index"`
	Deposed       string               `json:"deposed"`
	Change        ChangeRepresentation `json:"change"`
}

type ChangeRepresentation struct {
	Actions []string    `json:"actions"`
	Before  interface{} `json:"before"`
	After   interface{} `json:"after"`
}

// ReadPlan unmarshalls terraform plan JSON to an object
func ReadPlan(data []byte) (*PlanRepresentation, error) {
	var plan PlanRepresentation
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}
