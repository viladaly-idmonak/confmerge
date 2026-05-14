package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type renameFileEntry struct {
	Path string `yaml:"path"`
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type renameFile struct {
	Renames []renameFileEntry `yaml:"renames"`
}

// LoadRenameRules reads a YAML file and returns a slice of RenameRule.
//
// Expected format:
//
//	renames:
//	  - path: "database"
//	    from: "host"
//	    to: "hostname"
func LoadRenameRules(path string) ([]RenameRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("rename loader: read file: %w", err)
	}

	var rf renameFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("rename loader: parse yaml: %w", err)
	}

	if len(rf.Renames) == 0 {
		return nil, nil
	}

	rules := make([]RenameRule, 0, len(rf.Renames))
	for i, entry := range rf.Renames {
		if entry.From == "" {
			return nil, fmt.Errorf("rename loader: entry %d missing 'from'", i)
		}
		if entry.To == "" {
			return nil, fmt.Errorf("rename loader: entry %d missing 'to'", i)
		}
		rules = append(rules, RenameRule{
			Path: entry.Path,
			From: entry.From,
			To:   entry.To,
		})
	}
	return rules, nil
}
