package merger

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a value at a given path.
type TransformFunc func(path string, value interface{}) (interface{}, error)

// TransformRule defines a transformation to apply to a specific key path.
type TransformRule struct {
	Path      string
	Transform TransformFunc
}

// Transformer applies a set of transform rules to a config map.
type Transformer struct {
	rules []TransformRule
}

// NewTransformer creates a new Transformer with the given rules.
func NewTransformer(rules []TransformRule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply walks the config map and applies matching transform rules.
func (t *Transformer) Apply(data map[string]interface{}) (map[string]interface{}, error) {
	result := deepCopyMap(data)
	for _, rule := range t.rules {
		parts := strings.Split(rule.Path, ".")
		if err := applyTransform(result, parts, rule.Path, rule.Transform); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func applyTransform(data map[string]interface{}, parts []string, fullPath string, fn TransformFunc) error {
	if len(parts) == 0 {
		return nil
	}
	key := parts[0]
	if len(parts) == 1 {
		val, ok := data[key]
		if !ok {
			return nil
		}
		newVal, err := fn(fullPath, val)
		if err != nil {
			return fmt.Errorf("transform at %q: %w", fullPath, err)
		}
		data[key] = newVal
		return nil
	}
	nested, ok := data[key]
	if !ok {
		return nil
	}
	nestMap, ok := nested.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map at %q, got %T", key, nested)
	}
	return applyTransform(nestMap, parts[1:], fullPath, fn)
}
