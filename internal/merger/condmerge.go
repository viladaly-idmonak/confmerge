package merger

import (
	"fmt"
	"strings"
)

// ConditionOp represents a comparison operator for a merge condition.
type ConditionOp string

const (
	OpEquals    ConditionOp = "eq"
	OpNotEquals ConditionOp = "neq"
	OpExists    ConditionOp = "exists"
	OpNotExists ConditionOp = "not_exists"
)

// MergeCondition defines a condition that must be satisfied before a layer is merged.
type MergeCondition struct {
	Path  string      // dot-separated path into the base map
	Op    ConditionOp // comparison operator
	Value string      // expected value (unused for exists/not_exists)
}

// String returns a human-readable representation of the condition.
func (c MergeCondition) String() string {
	if c.Op == OpExists || c.Op == OpNotExists {
		return fmt.Sprintf("%s %s", c.Op, c.Path)
	}
	return fmt.Sprintf("%s %s %q", c.Path, c.Op, c.Value)
}

// ConditionalMerger merges layers only when conditions on the base are satisfied.
type ConditionalMerger struct {
	conditions []MergeCondition
}

// NewConditionalMerger creates a ConditionalMerger with the given conditions.
func NewConditionalMerger(conditions []MergeCondition) *ConditionalMerger {
	return &ConditionalMerger{conditions: conditions}
}

// MergeIfSatisfied merges overlay into base only when all conditions pass.
// Returns (merged, true, nil) on success, (base, false, nil) when conditions
// are not met, or (nil, false, err) on evaluation error.
func (cm *ConditionalMerger) MergeIfSatisfied(base, overlay map[string]interface{}) (map[string]interface{}, bool, error) {
	for _, cond := range cm.conditions {
		ok, err := evalCondition(base, cond)
		if err != nil {
			return nil, false, fmt.Errorf("condition %s: %w", cond, err)
		}
		if !ok {
			return base, false, nil
		}
	}
	merged := Merge(base, overlay)
	return merged, true, nil
}

// evalCondition evaluates a single MergeCondition against data.
func evalCondition(data map[string]interface{}, cond MergeCondition) (bool, error) {
	parts := strings.Split(cond.Path, ".")
	val, found := walkPath(data, parts)

	switch cond.Op {
	case OpExists:
		return found, nil
	case OpNotExists:
		return !found, nil
	case OpEquals:
		if !found {
			return false, nil
		}
		return fmt.Sprintf("%v", val) == cond.Value, nil
	case OpNotEquals:
		if !found {
			return true, nil
		}
		return fmt.Sprintf("%v", val) != cond.Value, nil
	default:
		return false, fmt.Errorf("unknown operator %q", cond.Op)
	}
}

// walkPath traverses nested maps following parts, returning the value and whether it was found.
func walkPath(data map[string]interface{}, parts []string) (interface{}, bool) {
	var current interface{} = data
	for _, p := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		v, exists := m[p]
		if !exists {
			return nil, false
		}
		current = v
	}
	return current, true
}
