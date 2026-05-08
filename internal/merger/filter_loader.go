package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// filterFile is the YAML structure for a filter definition file.
type filterFile struct {
	Filters []filterEntry `yaml:"filters"`
}

type filterEntry struct {
	Path string `yaml:"path"`
	Mode string `yaml:"mode"`
}

// LoadFilters reads a YAML filter definition file and returns a slice of FilterRules.
//
// Example YAML:
//
//	filters:
//	  - path: database.password
//	    mode: exclude
//	  - path: debug
//	    mode: exclude
func LoadFilters(path string) ([]FilterRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("filter loader: read file: %w", err)
	}

	var ff filterFile
	if err := yaml.Unmarshal(data, &ff); err != nil {
		return nil, fmt.Errorf("filter loader: parse YAML: %w", err)
	}

	var rules []FilterRule
	for i, entry := range ff.Filters {
		if entry.Path == "" {
			return nil, fmt.Errorf("filter loader: entry %d missing 'path'", i)
		}
		mode := entry.Mode
		if mode == "" {
			mode = "exclude"
		}
		if mode != "exclude" && mode != "include" {
			return nil, fmt.Errorf("filter loader: entry %d invalid mode %q", i, mode)
		}
		rules = append(rules, FilterRule{Path: entry.Path, Mode: mode})
	}
	return rules, nil
}
