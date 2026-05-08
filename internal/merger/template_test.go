package merger

import (
	"testing"
)

func TestTemplateExpander_NoTemplates(t *testing.T) {
	te := NewTemplateExpander(map[string]interface{}{"env": "prod"})
	input := map[string]interface{}{"host": "localhost", "port": 8080}
	out, err := te.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" || out["port"] != 8080 {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestTemplateExpander_SimpleVar(t *testing.T) {
	te := NewTemplateExpander(map[string]interface{}{"Env": "production"})
	input := map[string]interface{}{"name": "app-{{.Env}}"}
	out, err := te.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["name"] != "app-production" {
		t.Errorf("expected app-production, got %v", out["name"])
	}
}

func TestTemplateExpander_NestedMap(t *testing.T) {
	te := NewTemplateExpander(map[string]interface{}{"Region": "us-east-1"})
	input := map[string]interface{}{
		"db": map[string]interface{}{
			"endpoint": "db.{{.Region}}.example.com",
		},
	}
	out, err := te.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["db"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested map")
	}
	if db["endpoint"] != "db.us-east-1.example.com" {
		t.Errorf("unexpected endpoint: %v", db["endpoint"])
	}
}

func TestTemplateExpander_ListValues(t *testing.T) {
	te := NewTemplateExpander(map[string]interface{}{"Domain": "example.com"})
	input := map[string]interface{}{
		"hosts": []interface{}{"api.{{.Domain}}", "www.{{.Domain}}"},
	}
	out, err := te.Expand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hosts, ok := out["hosts"].([]interface{})
	if !ok || len(hosts) != 2 {
		t.Fatalf("expected list of 2")
	}
	if hosts[0] != "api.example.com" || hosts[1] != "www.example.com" {
		t.Errorf("unexpected hosts: %v", hosts)
	}
}

func TestTemplateExpander_MissingKey(t *testing.T) {
	te := NewTemplateExpander(map[string]interface{}{})
	input := map[string]interface{}{"name": "app-{{.Missing}}"}
	_, err := te.Expand(input)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestTemplateExpander_InvalidTemplate(t *testing.T) {
	te := NewTemplateExpander(map[string]interface{}{})
	input := map[string]interface{}{"name": "{{unclosed"}
	_, err := te.Expand(input)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}
