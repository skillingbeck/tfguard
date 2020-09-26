package tfguard

import (
	"encoding/json"
	"errors"
	"fmt"
)

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

	if !json.Valid(data) {
		return nil, errors.New("file is not valid json")
	}
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("plan format does not match schema: %w", err)
	}

	if len(plan.FormatVersion) == 0 {
		return nil, errors.New("format_version not present, check this is a terraform plan generated with terraform show")
	}

	return &plan, nil
}
