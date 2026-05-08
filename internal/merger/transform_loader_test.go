package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempTransform(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "transforms.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp transform file: %v", err)
	}
	return p
}

func TestLoadTransforms_Uppercase(t *testing.T) {
	p := writeTempTransform(t, "- path: env\n  op: uppercase\n")
	tr, err := LoadTransforms(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := tr.Apply(map[string]interface{}{"env": "development"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if result["env"] != "DEVELOPMENT" {
		t.Errorf("expected DEVELOPMENT, got %v", result["env"])
	}
}

func TestLoadTransforms_Lowercase(t *testing.T) {
	p := writeTempTransform(t, "- path: region\n  op: lowercase\n")
	tr, err := LoadTransforms(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := tr.Apply(map[string]interface{}{"region": "US-EAST-1"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if result["region"] != "us-east-1" {
		t.Errorf("expected us-east-1, got %v", result["region"])
	}
}

func TestLoadTransforms_Replace(t *testing.T) {
	p := writeTempTransform(t, "- path: url\n  op: \"replace:localhost:prod.example.com\"\n")
	tr, err := LoadTransforms(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := tr.Apply(map[string]interface{}{"url": "http://localhost:8080"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if result["url"] != "http://prod.example.com:8080" {
		t.Errorf("unexpected url: %v", result["url"])
	}
}

func TestLoadTransforms_Set(t *testing.T) {
	p := writeTempTransform(t, "- path: mode\n  op: set\n  value: production\n")
	tr, err := LoadTransforms(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := tr.Apply(map[string]interface{}{"mode": "debug"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if result["mode"] != "production" {
		t.Errorf("expected production, got %v", result["mode"])
	}
}

func TestLoadTransforms_FileNotFound(t *testing.T) {
	_, err := LoadTransforms("/nonexistent/transforms.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadTransforms_UnknownOp(t *testing.T) {
	p := writeTempTransform(t, "- path: key\n  op: explode\n")
	_, err := LoadTransforms(p)
	if err == nil {
		t.Fatal("expected error for unknown operation")
	}
}
