package merger

import (
	"testing"
)

func TestFilter_ExcludeSingleKey(t *testing.T) {
	data := map[string]interface{}{
		"host": "localhost",
		"port": 5432,
		"password": "secret",
	}
	f := NewFilter([]FilterRule{{Path: "password", Mode: "exclude"}})
	out, err := f.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["password"]; ok {
		t.Error("expected 'password' to be excluded")
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", out["host"])
	}
}

func TestFilter_ExcludeNestedKey(t *testing.T) {
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host":     "localhost",
			"password": "secret",
		},
	}
	f := NewFilter([]FilterRule{{Path: "database.password", Mode: "exclude"}})
	out, err := f.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, _ := out["database"].(map[string]interface{})
	if _, ok := db["password"]; ok {
		t.Error("expected nested 'password' to be excluded")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", db["host"])
	}
}

func TestFilter_InvalidMode(t *testing.T) {
	data := map[string]interface{}{"key": "val"}
	f := NewFilter([]FilterRule{{Path: "key", Mode: "unknown"}})
	_, err := f.Apply(data)
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestExcludeAllExcept(t *testing.T) {
	data := map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"password": "secret",
	}
	out := ExcludeAllExcept(data, []string{"host", "port"})
	if _, ok := out["password"]; ok {
		t.Error("expected 'password' to be absent")
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", out["host"])
	}
}

func TestFilterRule_String(t *testing.T) {
	r := FilterRule{Path: "db.pass", Mode: "exclude"}
	if r.String() != "exclude:db.pass" {
		t.Errorf("unexpected string: %s", r.String())
	}
}
