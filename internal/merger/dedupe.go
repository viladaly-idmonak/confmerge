package merger

import "fmt"

// DedupeStrategy controls how duplicate list items are removed.
type DedupeStrategy string

const (
	DedupeFirst DedupeStrategy = "first" // keep first occurrence
	DedupeLast  DedupeStrategy = "last"  // keep last occurrence
)

// DedupeRule defines a path whose list value should be deduplicated.
type DedupeRule struct {
	Path     string        // dot-separated path to a list key
	Strategy DedupeStrategy // "first" or "last"
}

// String returns a human-readable representation of the rule.
func (r DedupeRule) String() string {
	return fmt.Sprintf("dedupe path=%q strategy=%s", r.Path, r.Strategy)
}

// Deduper removes duplicate scalar entries from list values in a config map.
type Deduper struct {
	rules []DedupeRule
}

// NewDeduper creates a Deduper with the given rules.
func NewDeduper(rules []DedupeRule) *Deduper {
	return &Deduper{rules: rules}
}

// Apply walks each rule's path, finds the list, and deduplicates it.
func (d *Deduper) Apply(cfg map[string]any) error {
	for _, rule := range d.rules {
		parts := splitPath(rule.Path)
		if len(parts) == 0 {
			continue
		}
		parent, key, err := resolveParentMap(cfg, parts)
		if err != nil {
			// path not found — skip silently
			continue
		}
		raw, ok := parent[key]
		if !ok {
			continue
		}
		list, ok := raw.([]any)
		if !ok {
			return fmt.Errorf("dedupe: path %q is not a list", rule.Path)
		}
		parent[key] = dedupeList(list, rule.Strategy)
	}
	return nil
}

// resolveParentMap walks parts[0:n-1] into nested maps and returns the
// penultimate map plus the final key.
func resolveParentMap(cfg map[string]any, parts []string) (map[string]any, string, error) {
	current := cfg
	for _, p := range parts[:len(parts)-1] {
		v, ok := current[p]
		if !ok {
			return nil, "", fmt.Errorf("key %q not found", p)
		}
		next, ok := v.(map[string]any)
		if !ok {
			return nil, "", fmt.Errorf("key %q is not a map", p)
		}
		current = next
	}
	return current, parts[len(parts)-1], nil
}

func dedupeList(list []any, strategy DedupeStrategy) []any {
	seen := make(map[any]struct{}, len(list))
	if strategy == DedupeLast {
		// iterate in reverse, then reverse result
		var result []any
		for i := len(list) - 1; i >= 0; i-- {
			v := list[i]
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
		// reverse
		for l, r := 0, len(result)-1; l < r; l, r = l+1, r-1 {
			result[l], result[r] = result[r], result[l]
		}
		return result
	}
	// default: first
	var result []any
	for _, v := range list {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
