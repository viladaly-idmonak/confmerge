package merger

import (
	"testing"
)

func TestRenamer_NoRules(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	r := NewRenamer(nil)
	if err := r.Apply(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("expected key to remain unchanged")
	}
}

func TestRenamer_RootRename(t *testing.T) {
	data := map[string]interface{}{"old": 42}
	r := NewRenamer([]RenameRule{{Path: "", From: "old", To: "new"}})
	if err := r.Apply(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := data["old"]; ok {
		t.Error("old key should be removed")
	}
	if data["new"] != 42 {
		t.Errorf("expected new key to hold value 42, got %v", data["new"])
	}
}

func TestRenamer_NestedRename(t *testing.T) {
	data := map[string]interface{}{
		"db": map[string]interface{}{"host": "localhost"},
	}
	r := NewRenamer([]RenameRule{{Path: "db", From: "host", To: "hostname"}})
	if err := r.Apply(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db := data["db"].(map[string]interface{})
	if _, ok := db["host"]; ok {
		t.Error("old key 'host' should be removed")
	}
	if db["hostname"] != "localhost" {
		t.Errorf("expected hostname=localhost, got %v", db["hostname"])
	}
}

func TestRenamer_MissingKey(t *testing.T) {
	data := map[string]interface{}{"other": "val"}
	r := NewRenamer([]RenameRule{{Path: "", From: "missing", To: "new"}})
	if err := r.Apply(data); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestRenamer_InvalidPath(t *testing.T) {
	data := map[string]interface{}{"a": "not-a-map"}
	r := NewRenamer([]RenameRule{{Path: "a", From: "x", To: "y"}})
	if err := r.Apply(data); err == nil {
		t.Error("expected error when path segment is not a map")
	}
}

func TestRenameRule_String(t *testing.T) {
	rule := RenameRule{Path: "db", From: "host", To: "hostname"}
	got := rule.String()
	want := `rename db: "host" -> "hostname"`
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
