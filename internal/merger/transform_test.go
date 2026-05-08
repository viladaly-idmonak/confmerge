package merger

import (
	"fmt"
	"strings"
	"testing"
)

func TestTransformer_NoRules(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	tr := NewTransformer(nil)
	result, err := tr.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected 'value', got %v", result["key"])
	}
}

func TestTransformer_ScalarUppercase(t *testing.T) {
	data := map[string]interface{}{"name": "hello"}
	rules := []TransformRule{
		{
			Path: "name",
			Transform: func(path string, val interface{}) (interface{}, error) {
				s, ok := val.(string)
				if !ok {
					return val, nil
				}
				return strings.ToUpper(s), nil
			},
		},
	}
	tr := NewTransformer(rules)
	result, err := tr.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "HELLO" {
		t.Errorf("expected 'HELLO', got %v", result["name"])
	}
}

func TestTransformer_NestedPath(t *testing.T) {
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"port": 5432,
		},
	}
	rules := []TransformRule{
		{
			Path: "database.port",
			Transform: func(path string, val interface{}) (interface{}, error) {
				return 9999, nil
			},
		},
	}
	tr := NewTransformer(rules)
	result, err := tr.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db := result["database"].(map[string]interface{})
	if db["port"] != 9999 {
		t.Errorf("expected 9999, got %v", db["port"])
	}
}

func TestTransformer_MissingPath(t *testing.T) {
	data := map[string]interface{}{"other": "val"}
	rules := []TransformRule{
		{
			Path: "missing.key",
			Transform: func(path string, val interface{}) (interface{}, error) {
				return "replaced", nil
			},
		},
	}
	tr := NewTransformer(rules)
	result, err := tr.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["other"] != "val" {
		t.Errorf("original data mutated unexpectedly")
	}
}

func TestTransformer_ErrorPropagates(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	rules := []TransformRule{
		{
			Path: "key",
			Transform: func(path string, val interface{}) (interface{}, error) {
				return nil, fmt.Errorf("intentional error")
			},
		},
	}
	tr := NewTransformer(rules)
	_, err := tr.Apply(data)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTransformer_OriginalUnchanged(t *testing.T) {
	data := map[string]interface{}{"env": "dev"}
	rules := []TransformRule{
		{
			Path: "env",
			Transform: func(path string, val interface{}) (interface{}, error) {
				return "prod", nil
			},
		},
	}
	tr := NewTransformer(rules)
	_, err := tr.Apply(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["env"] != "dev" {
		t.Errorf("original data was mutated")
	}
}
