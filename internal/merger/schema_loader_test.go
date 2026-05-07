package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp schema: %v", err)
	}
	return path
}

func TestLoadSchema_Basic(t *testing.T) {
	content := `
fields:
  host:
    type: string
    required: true
  port:
    type: int
    required: true
  debug:
    type: bool
`
	path := writeTempSchema(t, content)
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Root == nil {
		t.Fatal("expected non-nil root")
	}
	if len(s.Root.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(s.Root.Children))
	}
	if !s.Root.Children["host"].Required {
		t.Error("expected host to be required")
	}
	if s.Root.Children["host"].Type != TypeString {
		t.Errorf("expected host type string, got %v", s.Root.Children["host"].Type)
	}
}

func TestLoadSchema_Nested(t *testing.T) {
	content := `
fields:
  database:
    type: map
    fields:
      name:
        type: string
        required: true
      pool:
        type: int
`
	path := writeTempSchema(t, content)
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := s.Root.Children["database"]
	if !ok {
		t.Fatal("expected database key")
	}
	if db.Type != TypeMap {
		t.Errorf("expected database type map, got %v", db.Type)
	}
	if len(db.Children) != 2 {
		t.Errorf("expected 2 nested children, got %d", len(db.Children))
	}
}

func TestLoadSchema_FileNotFound(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSchema_InvalidYAML(t *testing.T) {
	path := writeTempSchema(t, "fields: [invalid: yaml: content")
	_, err := LoadSchema(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParseSchemaType_Defaults(t *testing.T) {
	if got := parseSchemaType("unknown"); got != TypeAny {
		t.Errorf("expected TypeAny for unknown type, got %v", got)
	}
	if got := parseSchemaType(""); got != TypeAny {
		t.Errorf("expected TypeAny for empty type, got %v", got)
	}
}
