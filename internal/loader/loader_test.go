package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoad_YAML(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "server:\n  host: localhost\n  port: 8080\n")

	result, fmt, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != FormatYAML {
		t.Errorf("expected FormatYAML, got %v", fmt)
	}
	server, ok := result["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'server' key to be a map")
	}
	if server["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", server["host"])
	}
}

func TestLoad_TOML(t *testing.T) {
	path := writeTempFile(t, "config.toml", "[server]\nhost = \"localhost\"\nport = 8080\n")

	result, fmt, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != FormatTOML {
		t.Errorf("expected FormatTOML, got %v", fmt)
	}
	server, ok := result["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'server' key to be a map")
	}
	if server["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", server["host"])
	}
}

func TestLoad_UnsupportedExtension(t *testing.T) {
	path := writeTempFile(t, "config.json", `{}`)
	_, _, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, _, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyYAML(t *testing.T) {
	path := writeTempFile(t, "empty.yaml", "")
	result, _, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestDetectFormat(t *testing.T) {
	cases := []struct {
		path   string
		wanted Format
	}{
		{"config.yaml", FormatYAML},
		{"config.yml", FormatYAML},
		{"config.YAML", FormatYAML},
		{"config.toml", FormatTOML},
		{"config.json", FormatUnknown},
		{"config", FormatUnknown},
	}
	for _, c := range cases {
		got := detectFormat(c.path)
		if got != c.wanted {
			t.Errorf("detectFormat(%q) = %v, want %v", c.path, got, c.wanted)
		}
	}
}
