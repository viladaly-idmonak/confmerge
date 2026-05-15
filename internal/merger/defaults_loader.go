package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type defaultRuleYAML struct {
	Path  string      `yaml:"path"`
	Value interface{} `yaml:"value"`
}

type defaultsFileYAML struct {
	Defaults []defaultRuleYAML `yaml:"defaults"`
}

// LoadDefaults reads a YAML file containing default rules and returns a slice
// of DefaultRule values ready for use with NewDefaulter.
func LoadDefaults(path string) ([]DefaultRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("defaults loader: read %q: %w", path, err)
	}

	var raw defaultsFileYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("defaults loader: parse %q: %w", path, err)
	}

	var rules []DefaultRule
	for i, r := range raw.Defaults {
		if r.Path == "" {
			return nil, fmt.Errorf("defaults loader: entry %d missing 'path'", i)
		}
		rules = append(rules, DefaultRule{
			Path:  r.Path,
			Value: r.Value,
		})
	}
	return rules, nil
}
