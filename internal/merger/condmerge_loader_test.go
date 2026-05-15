package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempCondFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "conditions.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp cond file: %v", err)
	}
	return p
}

func TestLoadConditions_Basic(t *testing.T) {
	p := writeTempCondFile(t, `
conditions:
  - path: env
    op: eq
    value: production
  - path: debug
    op: not_exists
`)
	conds, err := LoadConditions(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conds) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(conds))
	}
	if conds[0].Op != OpEquals || conds[0].Value != "production" {
		t.Errorf("unexpected first condition: %+v", conds[0])
	}
	if conds[1].Op != OpNotExists {
		t.Errorf("unexpected second condition: %+v", conds[1])
	}
}

func TestLoadConditions_InvalidOp(t *testing.T) {
	p := writeTempCondFile(t, `
conditions:
  - path: env
    op: badop
    value: x
`)
	_, err := LoadConditions(p)
	if err == nil {
		t.Fatal("expected error for invalid op")
	}
}

func TestLoadConditions_MissingPath(t *testing.T) {
	p := writeTempCondFile(t, `
conditions:
  - op: eq
    value: production
`)
	_, err := LoadConditions(p)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestLoadConditions_FileNotFound(t *testing.T) {
	_, err := LoadConditions("/nonexistent/conditions.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConditions_InvalidYAML(t *testing.T) {
	p := writeTempCondFile(t, `:::invalid yaml:::`)
	_, err := LoadConditions(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
