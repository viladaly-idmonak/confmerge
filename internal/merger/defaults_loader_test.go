package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempDefaultsFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "defaults.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp defaults: %v", err)
	}
	return p
}

func TestLoadDefaults_Basic(t *testing.T) {
	p := writeTempDefaultsFile(t, `
defaults:
  - path: timeout
    value: 30
  - path: retries
    value: 3
`)
	rules, err := LoadDefaults(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Path != "timeout" || rules[0].Value != 30 {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestLoadDefaults_Empty(t *testing.T) {
	p := writeTempDefaultsFile(t, "defaults: []\n")
	rules, err := LoadDefaults(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestLoadDefaults_FileNotFound(t *testing.T) {
	_, err := LoadDefaults("/nonexistent/defaults.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadDefaults_InvalidYAML(t *testing.T) {
	p := writeTempDefaultsFile(t, ": bad: yaml: [")
	_, err := LoadDefaults(p)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadDefaults_MissingPath(t *testing.T) {
	p := writeTempDefaultsFile(t, `
defaults:
  - value: 42
`)
	_, err := LoadDefaults(p)
	if err == nil {
		t.Error("expected error for missing path field")
	}
}
