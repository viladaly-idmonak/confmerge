package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFreezeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFreezeFile: %v", err)
	}
	return path
}

func TestLoadFreezeRules_Basic(t *testing.T) {
	path := writeTempFreezeFile(t, `
freeze:
  - path: database.host
    reason: managed by infra
  - path: app.secret_key
`)
	rules, err := LoadFreezeRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Path != "database.host" {
		t.Errorf("expected path 'database.host', got %q", rules[0].Path)
	}
	if rules[0].Reason != "managed by infra" {
		t.Errorf("expected reason 'managed by infra', got %q", rules[0].Reason)
	}
	if rules[1].Path != "app.secret_key" {
		t.Errorf("expected path 'app.secret_key', got %q", rules[1].Path)
	}
	if rules[1].Reason != "" {
		t.Errorf("expected empty reason, got %q", rules[1].Reason)
	}
}

func TestLoadFreezeRules_Empty(t *testing.T) {
	path := writeTempFreezeFile(t, `freeze: []`)
	rules, err := LoadFreezeRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestLoadFreezeRules_FileNotFound(t *testing.T) {
	_, err := LoadFreezeRules("/nonexistent/freeze.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFreezeRules_InvalidYAML(t *testing.T) {
	path := writeTempFreezeFile(t, `freeze: [invalid: yaml: here`)
	_, err := LoadFreezeRules(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoadFreezeRules_MissingPath(t *testing.T) {
	path := writeTempFreezeFile(t, `
freeze:
  - reason: no path specified
`)
	_, err := LoadFreezeRules(path)
	if err == nil {
		t.Fatal("expected error for rule missing path, got nil")
	}
}
