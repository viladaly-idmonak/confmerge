package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvOverrideFile represents a YAML file that maps config keys to env var names.
// Example:
//   database.password: DB_PASSWORD
//   server.host: APP_HOST

type envOverrideFile struct {
	Mappings map[string]string `yaml:"mappings"`
}

// LoadEnvOverrides reads a YAML file defining key->envVar mappings and returns
// a flat map of config path -> resolved env var value.
func LoadEnvOverrides(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("env override: read file %q: %w", path, err)
	}
	var f envOverrideFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("env override: parse %q: %w", path, err)
	}
	if f.Mappings == nil {
		return map[string]string{}, nil
	}
	resolved := make(map[string]string, len(f.Mappings))
	for configKey, envVar := range f.Mappings {
		val, ok := os.LookupEnv(envVar)
		if !ok {
			return nil, fmt.Errorf("env override: env var %q (for key %q) is not set", envVar, configKey)
		}
		resolved[configKey] = val
	}
	return resolved, nil
}

// ApplyEnvOverrides sets dotted-path keys in data to the provided string values.
// Existing nested maps are traversed; missing intermediate maps are created.
func ApplyEnvOverrides(data map[string]interface{}, overrides map[string]string) error {
	for path, val := range overrides {
		if err := setNested(data, splitPath(path), val); err != nil {
			return fmt.Errorf("env override: set %q: %w", path, err)
		}
	}
	return nil
}

func splitPath(path string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}
