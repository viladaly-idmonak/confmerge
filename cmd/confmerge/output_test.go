package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteYAML_Scalar(t *testing.T) {
	data := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}
	var buf bytes.Buffer
	if err := writeYAML(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "host: localhost") {
		t.Errorf("expected 'host: localhost' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "port: 8080") {
		t.Errorf("expected 'port: 8080' in output, got:\n%s", out)
	}
}

func TestWriteYAML_Nested(t *testing.T) {
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "db.local",
			"port": 5432,
		},
	}
	var buf bytes.Buffer
	if err := writeYAML(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "database:") {
		t.Errorf("expected 'database:' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "  host: db.local") {
		t.Errorf("expected indented 'host: db.local', got:\n%s", out)
	}
}

func TestWriteYAML_List(t *testing.T) {
	data := map[string]interface{}{
		"tags": []interface{}{"a", "b", "c"},
	}
	var buf bytes.Buffer
	if err := writeYAML(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "  - a") {
		t.Errorf("expected list items in output, got:\n%s", out)
	}
}

func TestWriteYAML_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := writeYAML(&buf, map[string]interface{}{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty map, got: %q", buf.String())
	}
}
