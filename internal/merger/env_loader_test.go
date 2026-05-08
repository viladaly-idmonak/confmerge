package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "env_overrides.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env file: %v", err)
	}
	return p
}

func TestLoadEnvOverrides_Basic(t *testing.T) {
	t.Setenv("MY_HOST", "prod.example.com")
	p := writeTempEnvFile(t, "mappings:\n  server.host: MY_HOST\n")
	resolved, err := LoadEnvOverrides(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved["server.host"] != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %v", resolved["server.host"])
	}
}

func TestLoadEnvOverrides_MissingEnvVar(t *testing.T) {
	os.Unsetenv("TOTALLY_MISSING_VAR")
	p := writeTempEnvFile(t, "mappings:\n  db.pass: TOTALLY_MISSING_VAR\n")
	_, err := LoadEnvOverrides(p)
	if err == nil {
		t.Fatal("expected error for missing env var")
	}
}

func TestLoadEnvOverrides_EmptyMappings(t *testing.T) {
	p := writeTempEnvFile(t, "mappings: {}\n")
	resolved, err := LoadEnvOverrides(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 0 {
		t.Errorf("expected empty map, got %v", resolved)
	}
}

func TestLoadEnvOverrides_FileNotFound(t *testing.T) {
	_, err := LoadEnvOverrides("/nonexistent/path/env.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestApplyEnvOverrides_Nested(t *testing.T) {
	data := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
	}
	overrides := map[string]string{
		"server.host": "prod.example.com",
	}
	if err := ApplyEnvOverrides(data, overrides); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	server := data["server"].(map[string]interface{})
	if server["host"] != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %v", server["host"])
	}
}

func TestSplitPath(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"a.b.c", []string{"a", "b", "c"}},
		{"simple", []string{"simple"}},
		{"x.y", []string{"x", "y"}},
	}
	for _, tc := range cases {
		got := splitPath(tc.input)
		if len(got) != len(tc.expected) {
			t.Errorf("splitPath(%q): got %v, want %v", tc.input, got, tc.expected)
			continue
		}
		for i := range got {
			if got[i] != tc.expected[i] {
				t.Errorf("splitPath(%q)[%d]: got %q, want %q", tc.input, i, got[i], tc.expected[i])
			}
		}
	}
}
