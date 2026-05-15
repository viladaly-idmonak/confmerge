package merger

import (
	"testing"
)

func TestConditionalMerger_ConditionMet(t *testing.T) {
	base := map[string]interface{}{"env": "production", "port": 8080}
	overlay := map[string]interface{}{"port": 9090}

	cm := NewConditionalMerger([]MergeCondition{
		{Path: "env", Op: OpEquals, Value: "production"},
	})
	result, applied, err := cm.MergeIfSatisfied(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !applied {
		t.Fatal("expected merge to be applied")
	}
	if result["port"] != 9090 {
		t.Errorf("expected port 9090, got %v", result["port"])
	}
}

func TestConditionalMerger_ConditionNotMet(t *testing.T) {
	base := map[string]interface{}{"env": "staging", "port": 8080}
	overlay := map[string]interface{}{"port": 9090}

	cm := NewConditionalMerger([]MergeCondition{
		{Path: "env", Op: OpEquals, Value: "production"},
	})
	result, applied, err := cm.MergeIfSatisfied(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if applied {
		t.Fatal("expected merge to be skipped")
	}
	if result["port"] != 8080 {
		t.Errorf("expected port 8080 (unchanged), got %v", result["port"])
	}
}

func TestConditionalMerger_ExistsOp(t *testing.T) {
	base := map[string]interface{}{"feature": map[string]interface{}{"flags": true}}
	overlay := map[string]interface{}{"debug": true}

	cm := NewConditionalMerger([]MergeCondition{
		{Path: "feature.flags", Op: OpExists},
	})
	_, applied, err := cm.MergeIfSatisfied(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !applied {
		t.Fatal("expected merge to be applied when path exists")
	}
}

func TestConditionalMerger_NotExistsOp(t *testing.T) {
	base := map[string]interface{}{"port": 8080}
	overlay := map[string]interface{}{"debug": true}

	cm := NewConditionalMerger([]MergeCondition{
		{Path: "missing.key", Op: OpNotExists},
	})
	_, applied, err := cm.MergeIfSatisfied(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !applied {
		t.Fatal("expected merge when path does not exist")
	}
}

func TestConditionalMerger_InvalidOp(t *testing.T) {
	base := map[string]interface{}{"env": "prod"}
	overlay := map[string]interface{}{}

	cm := NewConditionalMerger([]MergeCondition{
		{Path: "env", Op: ConditionOp("invalid"), Value: "prod"},
	})
	_, _, err := cm.MergeIfSatisfied(base, overlay)
	if err == nil {
		t.Fatal("expected error for invalid op")
	}
}

func TestMergeCondition_String(t *testing.T) {
	c := MergeCondition{Path: "env", Op: OpEquals, Value: "prod"}
	if c.String() == "" {
		t.Error("expected non-empty string")
	}
	c2 := MergeCondition{Path: "feature", Op: OpExists}
	if c2.String() == "" {
		t.Error("expected non-empty string for exists op")
	}
}
