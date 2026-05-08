package merger

import (
	"os"
	"testing"
)

func writeTempFilterFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "filter-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFilters_Basic(t *testing.T) {
	path := writeTempFilterFile(t, `
filters:
  - path: database.password
    mode: exclude
  - path: debug
    mode: exclude
`)
	rules, err := LoadFilters(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Path != "database.password" || rules[0].Mode != "exclude" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestLoadFilters_DefaultMode(t *testing.T) {
	path := writeTempFilterFile(t, `
filters:
  - path: secret
`)
	rules, err := LoadFilters(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Mode != "exclude" {
		t.Errorf("expected default mode 'exclude', got %q", rules[0].Mode)
	}
}

func TestLoadFilters_InvalidMode(t *testing.T) {
	path := writeTempFilterFile(t, `
filters:
  - path: key
    mode: redact
`)
	_, err := LoadFilters(path)
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestLoadFilters_FileNotFound(t *testing.T) {
	_, err := LoadFilters("/nonexistent/filter.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFilters_InvalidYAML(t *testing.T) {
	path := writeTempFilterFile(t, `:::invalid yaml:::`)
	_, err := LoadFilters(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
