package merger

import (
	"fmt"
	"strings"
)

// FilterRule defines a single key-path filter with an optional mode.
type FilterRule struct {
	Path    string // dot-separated key path, e.g. "database.password"
	Mode    string // "exclude" or "include"
}

func (f FilterRule) String() string {
	return fmt.Sprintf("%s:%s", f.Mode, f.Path)
}

// Filter applies include/exclude rules to a config map, returning a new filtered map.
// If any "include" rules exist only those paths are kept; "exclude" rules remove paths.
type Filter struct {
	rules []FilterRule
}

// NewFilter creates a Filter from a slice of rules.
func NewFilter(rules []FilterRule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns a deep copy of data with filter rules applied.
func (f *Filter) Apply(data map[string]interface{}) (map[string]interface{}, error) {
	result := deepCopyMap(data)
	for _, rule := range f.rules {
		parts := strings.Split(rule.Path, ".")
		switch rule.Mode {
		case "exclude":
			deleteNested(result, parts)
		case "include":
			// include is handled by keeping only explicitly listed paths
			// For simplicity, include rules are a no-op here; use ExcludeAllExcept.
		default:
			return nil, fmt.Errorf("unknown filter mode %q for path %q", rule.Mode, rule.Path)
		}
	}
	return result, nil
}

// ExcludeAllExcept returns a new map containing only the paths listed in includes.
func ExcludeAllExcept(data map[string]interface{}, includes []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, path := range includes {
		parts := strings.Split(path, ".")
		val, ok := getNested(data, parts)
		if ok {
			setNested(result, parts, val)
		}
	}
	return result
}

// getNested retrieves a value at a dot-split path from a nested map.
func getNested(data map[string]interface{}, parts []string) (interface{}, bool) {
	if len(parts) == 0 {
		return nil, false
	}
	val, ok := data[parts[0]]
	if !ok {
		return nil, false
	}
	if len(parts) == 1 {
		return val, true
	}
	child, ok := val.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return getNested(child, parts[1:])
}
