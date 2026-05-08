package merger

import (
	"testing"
)

func TestApplyPatches_Set(t *testing.T) {
	base := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}
	patches := []Patch{
		{Op: PatchSet, Path: []string{"host"}, Value: "example.com"},
		{Op: PatchSet, Path: []string{"port"}, Value: 443},
	}
	result, err := ApplyPatches(base, patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["host"] != "example.com" {
		t.Errorf("expected host=example.com, got %v", result["host"])
	}
	if result["port"] != 443 {
		t.Errorf("expected port=443, got %v", result["port"])
	}
	// original must not be mutated
	if base["host"] != "localhost" {
		t.Errorf("base was mutated")
	}
}

func TestApplyPatches_Delete(t *testing.T) {
	base := map[string]interface{}{
		"debug": true,
		"level": "info",
	}
	patches := []Patch{
		{Op: PatchDelete, Path: []string{"debug"}},
	}
	result, err := ApplyPatches(base, patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := result["debug"]; exists {
		t.Errorf("expected debug key to be deleted")
	}
	if result["level"] != "info" {
		t.Errorf("expected level=info to remain")
	}
}

func TestApplyPatches_NestedSet(t *testing.T) {
	base := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}
	patches := []Patch{
		{Op: PatchSet, Path: []string{"database", "host"}, Value: "db.prod"},
		{Op: PatchSet, Path: []string{"database", "name"}, Value: "mydb"},
	}
	result, err := ApplyPatches(base, patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db := result["database"].(map[string]interface{})
	if db["host"] != "db.prod" {
		t.Errorf("expected db.host=db.prod, got %v", db["host"])
	}
	if db["name"] != "mydb" {
		t.Errorf("expected db.name=mydb, got %v", db["name"])
	}
}

func TestApplyPatches_EmptyPath(t *testing.T) {
	base := map[string]interface{}{}
	patches := []Patch{{Op: PatchSet, Path: []string{}, Value: "x"}}
	_, err := ApplyPatches(base, patches)
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestPatch_String(t *testing.T) {
	p := Patch{Op: PatchSet, Path: []string{"a", "b"}, Value: 42}
	if p.String() != "set a.b = 42" {
		t.Errorf("unexpected string: %s", p.String())
	}
	d := Patch{Op: PatchDelete, Path: []string{"x", "y"}}
	if d.String() != "delete x.y" {
		t.Errorf("unexpected string: %s", d.String())
	}
}
