package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLintFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "lint.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write lint file: %v", err)
	}
	return p
}

func TestLoadLintRules_Basic(t *testing.T) {
	p := writeTempLintFile(t, `
rules:
  - pattern: "database.password"
    description: "avoid plaintext passwords"
    severity: error
  - pattern: "secrets.*"
    description: "secrets namespace"
    severity: warn
`)
	rules, err := LoadLintRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Pattern != "database.password" {
		t.Errorf("unexpected pattern: %s", rules[0].Pattern)
	}
	if rules[1].Severity != "warn" {
		t.Errorf("unexpected severity: %s", rules[1].Severity)
	}
}

func TestLoadLintRules_DefaultSeverity(t *testing.T) {
	p := writeTempLintFile(t, `
rules:
  - pattern: "app.debug"
    description: "debug mode enabled"
`)
	rules, err := LoadLintRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Severity != "warn" {
		t.Errorf("expected default severity 'warn', got %q", rules[0].Severity)
	}
}

func TestLoadLintRules_InvalidSeverity(t *testing.T) {
	p := writeTempLintFile(t, `
rules:
  - pattern: "x"
    description: "bad"
    severity: critical
`)
	_, err := LoadLintRules(p)
	if err == nil {
		t.Error("expected error for invalid severity")
	}
}

func TestLoadLintRules_FileNotFound(t *testing.T) {
	_, err := LoadLintRules("/nonexistent/lint.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadLintRules_InvalidYAML(t *testing.T) {
	p := writeTempLintFile(t, ":::invalid yaml:::")
	_, err := LoadLintRules(p)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadLintRules_MissingPattern(t *testing.T) {
	p := writeTempLintFile(t, `
rules:
  - description: "no pattern here"
    severity: warn
`)
	_, err := LoadLintRules(p)
	if err == nil {
		t.Error("expected error for missing pattern")
	}
}
