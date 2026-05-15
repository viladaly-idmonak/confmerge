package merger

import (
	"testing"
)

func TestDefaulter_NoRules(t *testing.T) {
	d := NewDefaulter(nil)
	dst := map[string]interface{}{"key": "val"}
	if err := d.Apply(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["key"] != "val" {
		t.Errorf("expected val, got %v", dst["key"])
	}
}

func TestDefaulter_AbsentKey(t *testing.T) {
	d := NewDefaulter([]DefaultRule{{Path: "timeout", Value: 30}})
	dst := map[string]interface{}{}
	if err := d.Apply(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["timeout"] != 30 {
		t.Errorf("expected 30, got %v", dst["timeout"])
	}
}

func TestDefaulter_PresentKeyNotOverwritten(t *testing.T) {
	d := NewDefaulter([]DefaultRule{{Path: "timeout", Value: 30}})
	dst := map[string]interface{}{"timeout": 60}
	if err := d.Apply(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["timeout"] != 60 {
		t.Errorf("expected 60, got %v", dst["timeout"])
	}
}

func TestDefaulter_NestedAbsent(t *testing.T) {
	d := NewDefaulter([]DefaultRule{{Path: "db.port", Value: 5432}})
	dst := map[string]interface{}{"db": map[string]interface{}{}}
	if err := d.Apply(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db := dst["db"].(map[string]interface{})
	if db["port"] != 5432 {
		t.Errorf("expected 5432, got %v", db["port"])
	}
}

func TestDefaulter_NestedPresent(t *testing.T) {
	d := NewDefaulter([]DefaultRule{{Path: "db.port", Value: 5432}})
	dst := map[string]interface{}{"db": map[string]interface{}{"port": 3306}}
	if err := d.Apply(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db := dst["db"].(map[string]interface{})
	if db["port"] != 3306 {
		t.Errorf("expected 3306, got %v", db["port"])
	}
}

func TestDefaultRule_String(t *testing.T) {
	r := DefaultRule{Path: "foo.bar", Value: "baz"}
	got := r.String()
	if got != "default(foo.bar=baz)" {
		t.Errorf("unexpected string: %q", got)
	}
}
