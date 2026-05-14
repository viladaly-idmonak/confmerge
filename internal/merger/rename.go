package merger

import "fmt"

// RenameRule defines a key rename operation at a given path.
type RenameRule struct {
	Path string // dot-separated path to the parent map
	From string // old key name
	To   string // new key name
}

func (r RenameRule) String() string {
	return fmt.Sprintf("rename %s: %q -> %q", r.Path, r.From, r.To)
}

// Renamer applies key rename rules to a config map.
type Renamer struct {
	rules []RenameRule
}

// NewRenamer creates a Renamer with the given rules.
func NewRenamer(rules []RenameRule) *Renamer {
	return &Renamer{rules: rules}
}

// Apply walks the config map and renames keys according to the rules.
// It returns an error if a source key is not found at the specified path.
func (r *Renamer) Apply(data map[string]interface{}) error {
	for _, rule := range r.rules {
		target, err := resolveParent(data, rule.Path)
		if err != nil {
			return fmt.Errorf("rename rule %s: %w", rule, err)
		}
		val, ok := target[rule.From]
		if !ok {
			return fmt.Errorf("rename rule %s: key %q not found", rule, rule.From)
		}
		delete(target, rule.From)
		target[rule.To] = val
	}
	return nil
}

// resolveParent navigates dot-separated path segments to return the nested map.
// An empty path returns the root map.
func resolveParent(data map[string]interface{}, path string) (map[string]interface{}, error) {
	if path == "" {
		return data, nil
	}
	parts := splitPath(path)
	current := data
	for i, part := range parts {
		val, ok := current[part]
		if !ok {
			return nil, fmt.Errorf("path segment %q not found (index %d)", part, i)
		}
		nested, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path segment %q is not a map", part)
		}
		current = nested
	}
	return current, nil
}
