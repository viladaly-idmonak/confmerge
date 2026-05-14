package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempRenameFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "renames.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp rename file: %v", err)
	}
	return p
}

func TestLoadRenameRules_Basic(t *testing.T) {
	p := writeTempRenameFile(t, `
renames:
  - path: "database"
    from: "host"
    to: "hostname"
`)
	rules, err := LoadRenameRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Path != "database" || rules[0].From != "host" || rules[0].To != "hostname" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestLoadRenameRules_Empty(t *testing.T) {
	p := writeTempRenameFile(t, "renames: []\n")
	rules, err := LoadRenameRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestLoadRenameRules_FileNotFound(t *testing.T) {
	_, err := LoadRenameRules("/nonexistent/renames.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRenameRules_InvalidYAML(t *testing.T) {
	p := writeTempRenameFile(t, ": invalid: yaml: [")
	_, err := LoadRenameRules(p)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadRenameRules_MissingFrom(t *testing.T) {
	p := writeTempRenameFile(t, `
renames:
  - path: "db"
    to: "hostname"
`)
	_, err := LoadRenameRules(p)
	if err == nil {
		t.Error("expected error for missing 'from' field")
	}
}

func TestLoadRenameRules_MissingTo(t *testing.T) {
	p := writeTempRenameFile(t, `
renames:
  - path: "db"
    from: "host"
`)
	_, err := LoadRenameRules(p)
	if err == nil {
		t.Error("expected error for missing 'to' field")
	}
}
