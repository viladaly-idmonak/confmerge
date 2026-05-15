package merger

import "fmt"

// DefaultRule specifies a default value to apply at a given path
// if the key is absent in the merged config.
type DefaultRule struct {
	Path  string
	Value interface{}
}

func (d DefaultRule) String() string {
	return fmt.Sprintf("default(%s=%v)", d.Path, d.Value)
}

// Defaulter applies default values to a config map.
type Defaulter struct {
	rules []DefaultRule
}

// NewDefaulter creates a Defaulter with the given rules.
func NewDefaulter(rules []DefaultRule) *Defaulter {
	return &Defaulter{rules: rules}
}

// Apply sets default values in dst for any path that is not already present.
func (d *Defaulter) Apply(dst map[string]interface{}) error {
	for _, rule := range d.rules {
		parts, err := splitPath(rule.Path)
		if err != nil {
			return fmt.Errorf("defaults: invalid path %q: %w", rule.Path, err)
		}
		if !pathExists(dst, parts) {
			if err := setNested(dst, parts, rule.Value); err != nil {
				return fmt.Errorf("defaults: set %q: %w", rule.Path, err)
			}
		}
	}
	return nil
}

// pathExists returns true if the dot-separated path exists in m.
func pathExists(m map[string]interface{}, parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	v, ok := m[parts[0]]
	if !ok {
		return false
	}
	if len(parts) == 1 {
		return true
	}
	child, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	return pathExists(child, parts[1:])
}
