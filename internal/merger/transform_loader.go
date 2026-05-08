package merger

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// transformRuleSpec is the YAML-serializable form of a transform rule.
type transformRuleSpec struct {
	Path      string `yaml:"path"`
	Operation string `yaml:"op"`
	Value     string `yaml:"value,omitempty"`
}

// LoadTransforms reads a YAML file defining transform rules and returns a Transformer.
// Supported ops: "uppercase", "lowercase", "replace:<old>:<new>", "set".
func LoadTransforms(path string) (*Transformer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading transform file: %w", err)
	}
	var specs []transformRuleSpec
	if err := yaml.Unmarshal(data, &specs); err != nil {
		return nil, fmt.Errorf("parsing transform file: %w", err)
	}
	rules := make([]TransformRule, 0, len(specs))
	for i, spec := range specs {
		if spec.Path == "" {
			return nil, fmt.Errorf("rule at index %d is missing required field \"path\"", i)
		}
		if spec.Operation == "" {
			return nil, fmt.Errorf("rule at index %d (path %q) is missing required field \"op\"", i, spec.Path)
		}
		rule, err := buildTransformRule(spec)
		if err != nil {
			return nil, fmt.Errorf("building rule for path %q: %w", spec.Path, err)
		}
		rules = append(rules, rule)
	}
	return NewTransformer(rules), nil
}

func buildTransformRule(spec transformRuleSpec) (TransformRule, error) {
	var fn TransformFunc
	switch {
	case spec.Operation == "uppercase":
		fn = func(path string, val interface{}) (interface{}, error) {
			s, ok := val.(string)
			if !ok {
				return val, nil
			}
			return strings.ToUpper(s), nil
		}
	case spec.Operation == "lowercase":
		fn = func(path string, val interface{}) (interface{}, error) {
			s, ok := val.(string)
			if !ok {
				return val, nil
			}
			return strings.ToLower(s), nil
		}
	case strings.HasPrefix(spec.Operation, "replace:"):
		parts := strings.SplitN(spec.Operation, ":", 3)
		if len(parts) != 3 {
			return TransformRule{}, fmt.Errorf("invalid replace op %q, expected replace:<old>:<new>", spec.Operation)
		}
		old, newStr := parts[1], parts[2]
		fn = func(path string, val interface{}) (interface{}, error) {
			s, ok := val.(string)
			if !ok {
				return val, nil
			}
			return strings.ReplaceAll(s, old, newStr), nil
		}
	case spec.Operation == "set":
		v := spec.Value
		fn = func(path string, val interface{}) (interface{}, error) {
			return v, nil
		}
	default:
		return TransformRule{}, fmt.Errorf("unknown operation %q", spec.Operation)
	}
	return TransformRule{Path: spec.Path, Transform: fn}, nil
}
