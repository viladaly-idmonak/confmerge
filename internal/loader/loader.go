package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents the config file format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatTOML Format = "toml"
	FormatUnknown Format = "unknown"
)

// Load reads a YAML or TOML config file and returns its contents as a
// map[string]interface{}. The format is inferred from the file extension.
func Load(path string) (map[string]interface{}, Format, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, FormatUnknown, fmt.Errorf("loader: reading file %q: %w", path, err)
	}

	fmt := detectFormat(path)
	switch fmt {
	case FormatYAML:
		result, err := parseYAML(data)
		if err != nil {
			return nil, FormatYAML, fmt.Errorf("loader: parsing YAML %q: %w", path, err)
		}
		return result, FormatYAML, nil
	case FormatTOML:
		result, err := parseTOML(data)
		if err != nil {
			return nil, FormatTOML, fmt.Errorf("loader: parsing TOML %q: %w", path, err)
		}
		return result, FormatTOML, nil
	default:
		return nil, FormatUnknown, fmt.Errorf("loader: unsupported file extension for %q", path)
	}
}

func detectFormat(path string) Format {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return FormatYAML
	case ".toml":
		return FormatTOML
	default:
		return FormatUnknown
	}
}

func parseYAML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result == nil {
		return map[string]interface{}{}, nil
	}
	return result, nil
}

func parseTOML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if _, err := toml.Decode(string(data), &result); err != nil {
		return nil, err
	}
	if result == nil {
		return map[string]interface{}{}, nil
	}
	return result, nil
}
