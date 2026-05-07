package merger_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confmerge/internal/loader"
	"github.com/yourorg/confmerge/internal/merger"
)

// writeTempIntegrationFile writes content to a temp file with the given extension.
func writeTempIntegrationFile(t *testing.T, ext, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*"+ext)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

// writeTempSchemaFile writes a YAML schema to a temp file.
func writeTempSchemaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write schema file: %v", err)
	}
	return path
}

// TestIntegration_MergeAndValidate loads two YAML configs, merges them,
// then validates the result against a schema.
func TestIntegration_MergeAndValidate(t *testing.T) {
	base := writeTempIntegrationFile(t, ".yaml", `
server:
  host: localhost
  port: 8080
debug: false
`)
	override := writeTempIntegrationFile(t, ".yaml", `
server:
  port: 9090
debug: true
`)
	schemaPath := writeTempSchemaFile(t, `
server:
  type: object
  required: true
  children:
    host:
      type: string
      required: true
    port:
      type: int
      required: true
debug:
  type: bool
  required: false
`)

	baseMap, err := loader.Load(base)
	if err != nil {
		t.Fatalf("failed to load base: %v", err)
	}
	overrideMap, err := loader.Load(override)
	if err != nil {
		t.Fatalf("failed to load override: %v", err)
	}

	merged := merger.Merge(baseMap, overrideMap)

	schema, err := merger.LoadSchema(schemaPath)
	if err != nil {
		t.Fatalf("failed to load schema: %v", err)
	}

	if errs := merger.ValidateSchema(merged, schema); len(errs) != 0 {
		t.Errorf("expected no validation errors, got: %v", errs)
	}

	// Verify merged values
	server, ok := merged["server"].(map[string]interface{})
	if !ok {
		t.Fatal("expected server to be a map")
	}
	if server["port"] != 9090 {
		t.Errorf("expected port 9090, got %v", server["port"])
	}
	if server["host"] != "localhost" {
		t.Errorf("expected host localhost, got %v", server["host"])
	}
	if merged["debug"] != true {
		t.Errorf("expected debug true, got %v", merged["debug"])
	}
}

// TestIntegration_MergeAndValidate_MissingRequired verifies that validation
// catches a missing required field after merging configs.
func TestIntegration_MergeAndValidate_MissingRequired(t *testing.T) {
	base := writeTempIntegrationFile(t, ".yaml", `
server:
  port: 8080
`)
	override := writeTempIntegrationFile(t, ".yaml", `
debug: true
`)
	schemaPath := writeTempSchemaFile(t, `
server:
  type: object
  required: true
  children:
    host:
      type: string
      required: true
    port:
      type: int
      required: true
`)

	baseMap, err := loader.Load(base)
	if err != nil {
		t.Fatalf("failed to load base: %v", err)
	}
	overrideMap, err := loader.Load(override)
	if err != nil {
		t.Fatalf("failed to load override: %v", err)
	}

	merged := merger.Merge(baseMap, overrideMap)

	schema, err := merger.LoadSchema(schemaPath)
	if err != nil {
		t.Fatalf("failed to load schema: %v", err)
	}

	errs := merger.ValidateSchema(merged, schema)
	if len(errs) == 0 {
		t.Error("expected validation errors for missing required field 'host', got none")
	}
}
