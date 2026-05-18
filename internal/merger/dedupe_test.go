package merger

import (
	"testing"
)

func TestDeduper_NoRules(t *testing.T) {
	cfg := map[string]any{"tags": []any{"a", "b", "a"}}
	d := NewDeduper(nil)
	if err := d.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// unchanged
	list := cfg["tags"].([]any)
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestDeduper_FirstStrategy(t *testing.T) {
	cfg := map[string]any{"tags": []any{"a", "b", "a", "c", "b"}}
	d := NewDeduper([]DedupeRule{{Path: "tags", Strategy: DedupeFirst}})
	if err := d.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list := cfg["tags"].([]any)
	want := []any{"a", "b", "c"}
	if len(list) != len(want) {
		t.Fatalf("expected %v, got %v", want, list)
	}
	for i, v := range want {
		if list[i] != v {
			t.Errorf("index %d: expected %v, got %v", i, v, list[i])
		}
	}
}

func TestDeduper_LastStrategy(t *testing.T) {
	cfg := map[string]any{"tags": []any{"a", "b", "a", "c", "b"}}
	d := NewDeduper([]DedupeRule{{Path: "tags", Strategy: DedupeLast}})
	if err := d.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list := cfg["tags"].([]any)
	want := []any{"a", "c", "b"}
	if len(list) != len(want) {
		t.Fatalf("expected %v, got %v", want, list)
	}
	for i, v := range want {
		if list[i] != v {
			t.Errorf("index %d: expected %v, got %v", i, v, list[i])
		}
	}
}

func TestDeduper_NestedPath(t *testing.T) {
	cfg := map[string]any{
		"server": map[string]any{
			"hosts": []any{"h1", "h2", "h1"},
		},
	}
	d := NewDeduper([]DedupeRule{{Path: "server.hosts", Strategy: DedupeFirst}})
	if err := d.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hosts := cfg["server"].(map[string]any)["hosts"].([]any)
	if len(hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestDeduper_PathNotFound(t *testing.T) {
	cfg := map[string]any{"other": "value"}
	d := NewDeduper([]DedupeRule{{Path: "missing.path", Strategy: DedupeFirst}})
	if err := d.Apply(cfg); err != nil {
		t.Fatalf("expected no error for missing path, got: %v", err)
	}
}

func TestDeduper_NotAList(t *testing.T) {
	cfg := map[string]any{"name": "alice"}
	d := NewDeduper([]DedupeRule{{Path: "name", Strategy: DedupeFirst}})
	if err := d.Apply(cfg); err == nil {
		t.Fatal("expected error for non-list value")
	}
}

func TestDedupeRule_String(t *testing.T) {
	r := DedupeRule{Path: "tags", Strategy: DedupeFirst}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
