package merger

import "testing"

func TestMergeWithAudit_Addition(t *testing.T) {
	base := map[string]interface{}{"a": 1}
	override := map[string]interface{}{"b": 2}
	log := NewAuditLog()

	result := MergeWithAudit(base, override, "override.yaml", log)

	if result["a"] != 1 || result["b"] != 2 {
		t.Errorf("unexpected result: %v", result)
	}
	entries := log.Entries()
	if len(entries) != 1 || entries[0].Operation != "add" {
		t.Errorf("expected 1 add entry, got: %v", entries)
	}
}

func TestMergeWithAudit_Override(t *testing.T) {
	base := map[string]interface{}{"host": "localhost"}
	override := map[string]interface{}{"host": "prod.example.com"}
	log := NewAuditLog()

	result := MergeWithAudit(base, override, "prod.yaml", log)

	if result["host"] != "prod.example.com" {
		t.Errorf("expected overridden host, got: %v", result["host"])
	}
	entries := log.Entries()
	if len(entries) != 1 || entries[0].Operation != "override" {
		t.Errorf("expected 1 override entry, got: %v", entries)
	}
	if entries[0].OldValue != "localhost" {
		t.Errorf("expected OldValue 'localhost', got: %v", entries[0].OldValue)
	}
}

func TestMergeWithAudit_DeepMerge(t *testing.T) {
	base := map[string]interface{}{
		"db": map[string]interface{}{"host": "localhost", "port": 5432},
	}
	override := map[string]interface{}{
		"db": map[string]interface{}{"host": "db.prod"},
	}
	log := NewAuditLog()

	result := MergeWithAudit(base, override, "prod.yaml", log)

	db, _ := toMap(result["db"])
	if db["host"] != "db.prod" {
		t.Errorf("expected db.host 'db.prod', got: %v", db["host"])
	}
	if db["port"] != 5432 {
		t.Errorf("expected db.port 5432, got: %v", db["port"])
	}
	for _, e := range log.Entries() {
		if e.Path == "db.host" && e.Operation == "override" {
			return
		}
	}
	t.Errorf("expected override entry for db.host, entries: %v", log.Entries())
}

func TestMergeWithAudit_NoMutation(t *testing.T) {
	base := map[string]interface{}{"x": 10}
	override := map[string]interface{}{"x": 20}
	log := NewAuditLog()

	MergeWithAudit(base, override, "s", log)

	if base["x"] != 10 {
		t.Errorf("base should not be mutated, got: %v", base["x"])
	}
}
