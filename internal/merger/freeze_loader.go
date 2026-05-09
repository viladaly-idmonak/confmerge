package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type freezeRuleSpec struct {
	Path   string `yaml:"path"`
	Reason string `yaml:"reason"`
}

type freezeFileSpec struct {
	Freeze []freezeRuleSpec `yaml:"freeze"`
}

// LoadFreezeRules reads a YAML file defining frozen (immutable) config paths.
// Each rule specifies a dot-separated path and an optional reason.
//
// Example YAML:
//
//	freeze:
//	  - path: database.host
//	    reason: "managed by infra team"
//	  - path: app.secret_key
func LoadFreezeRules(filename string) ([]FreezeRule, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("freeze_loader: read file: %w", err)
	}

	var spec freezeFileSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("freeze_loader: parse yaml: %w", err)
	}

	if len(spec.Freeze) == 0 {
		return nil, nil
	}

	rules := make([]FreezeRule, 0, len(spec.Freeze))
	for i, r := range spec.Freeze {
		if r.Path == "" {
			return nil, fmt.Errorf("freeze_loader: rule %d missing required 'path' field", i)
		}
		rules = append(rules, FreezeRule{
			Path:   r.Path,
			Reason: r.Reason,
		})
	}

	return rules, nil
}
