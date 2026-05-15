package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// conditionFile mirrors the YAML structure for loading conditions.
type conditionFile struct {
	Conditions []conditionEntry `yaml:"conditions"`
}

type conditionEntry struct {
	Path  string `yaml:"path"`
	Op    string `yaml:"op"`
	Value string `yaml:"value"`
}

// LoadConditions reads a YAML file and returns a slice of MergeCondition.
//
// Example YAML:
//
//	conditions:
//	  - path: env
//	    op: eq
//	    value: production
//	  - path: feature.flags
//	    op: exists
func LoadConditions(path string) ([]MergeCondition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load conditions: %w", err)
	}

	var cf conditionFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parse conditions: %w", err)
	}

	var conditions []MergeCondition
	for i, e := range cf.Conditions {
		if e.Path == "" {
			return nil, fmt.Errorf("condition[%d]: path is required", i)
		}
		op := ConditionOp(e.Op)
		switch op {
		case OpEquals, OpNotEquals, OpExists, OpNotExists:
			// valid
		default:
			return nil, fmt.Errorf("condition[%d]: unknown op %q", i, e.Op)
		}
		conditions = append(conditions, MergeCondition{
			Path:  e.Path,
			Op:    op,
			Value: e.Value,
		})
	}
	return conditions, nil
}
