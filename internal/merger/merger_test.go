package merger

import (
	"testing"
)

func TestMerge_BasicOverride(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "port": 5432}
	over := map[string]interface{}{"port": 9999}

	res, err := Merge([]map[string]interface{}{base, over}, []string{"base", "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Data["port"] != 9999 {
		t.Errorf("expected port 9999, got %v", res.Data["port"])
	}
	if res.Data["host"] != "localhost" {
		t.Errorf("expected host localhost, got %v", res.Data["host"])
	}
	if len(res.Overrides) != 1 {
		t.Fatalf("expected 1 override, got %d", len(res.Overrides))
	}
	ov := res.Overrides[0]
	if ov.Key != "port" || ov.BaseVal != 5432 || ov.OverVal != 9999 || ov.Source != "override" {
		t.Errorf("unexpected override record: %+v", ov)
	}
}

func TestMerge_DeepMerge(t *testing.T) {
	base := map[string]interface{}{
		"db": map[string]interface{}{"host": "localhost", "port": 5432},
	}
	over := map[string]interface{}{
		"db": map[string]interface{}{"port": 6543},
	}

	res, err := Merge([]map[string]interface{}{base, over}, []string{"base", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	db, ok := res.Data["db"].(map[string]interface{})
	if !ok {
		t.Fatal("db key should be a map")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected db.host localhost, got %v", db["host"])
	}
	if db["port"] != 6543 {
		t.Errorf("expected db.port 6543, got %v", db["port"])
	}
	if len(res.Overrides) != 1 || res.Overrides[0].Key != "db.port" {
		t.Errorf("expected override for db.port, got %+v", res.Overrides)
	}
}

func TestMerge_NoConflict(t *testing.T) {
	a := map[string]interface{}{"a": 1}
	b := map[string]interface{}{"b": 2}

	res, err := Merge([]map[string]interface{}{a, b}, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overrides) != 0 {
		t.Errorf("expected no overrides, got %d", len(res.Overrides))
	}
	if res.Data["a"] != 1 || res.Data["b"] != 2 {
		t.Errorf("unexpected data: %v", res.Data)
	}
}

func TestMerge_LengthMismatch(t *testing.T) {
	_, err := Merge([]map[string]interface{}{{"a": 1}}, []string{})
	if err == nil {
		t.Error("expected error for mismatched slice lengths")
	}
}
