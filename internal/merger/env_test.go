package merger

import (
	"testing"
)

func TestEnvExpander_NoVars(t *testing.T) {
	e := NewEnvExpander(false)
	input := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}
	out, err := e.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" {
		t.Errorf("expected localhost, got %v", out["host"])
	}
	if out["port"] != 8080 {
		t.Errorf("expected 8080, got %v", out["port"])
	}
}

func TestEnvExpander_SimpleVar(t *testing.T) {
	t.Setenv("APP_HOST", "example.com")
	e := NewEnvExpander(false)
	input := map[string]interface{}{
		"host": "${APP_HOST}",
	}
	out, err := e.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "example.com" {
		t.Errorf("expected example.com, got %v", out["host"])
	}
}

func TestEnvExpander_MissingVar_NotAllowed(t *testing.T) {
	e := NewEnvExpander(false)
	input := map[string]interface{}{
		"secret": "${UNDEFINED_VAR_XYZ}",
	}
	_, err := e.Expand(input)
	if err == nil {
		t.Fatal("expected error for missing env var, got nil")
	}
}

func TestEnvExpander_MissingVar_Allowed(t *testing.T) {
	e := NewEnvExpander(true)
	input := map[string]interface{}{
		"secret": "${UNDEFINED_VAR_XYZ}",
	}
	out, err := e.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["secret"] != "" {
		t.Errorf("expected empty string, got %v", out["secret"])
	}
}

func TestEnvExpander_NestedMap(t *testing.T) {
	t.Setenv("DB_PASS", "s3cr3t")
	e := NewEnvExpander(false)
	input := map[string]interface{}{
		"database": map[string]interface{}{
			"password": "${DB_PASS}",
			"port":     5432,
		},
	}
	out, err := e.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["password"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %v", db["password"])
	}
}

func TestEnvExpander_SliceValues(t *testing.T) {
	t.Setenv("ITEM_ONE", "alpha")
	e := NewEnvExpander(false)
	input := map[string]interface{}{
		"items": []interface{}{"${ITEM_ONE}", "beta"},
	}
	out, err := e.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	items, ok := out["items"].([]interface{})
	if !ok {
		t.Fatal("expected slice")
	}
	if items[0] != "alpha" {
		t.Errorf("expected alpha, got %v", items[0])
	}
	if items[1] != "beta" {
		t.Errorf("expected beta, got %v", items[1])
	}
}
