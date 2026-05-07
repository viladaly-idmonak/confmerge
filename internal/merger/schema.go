package merger

import (
	"fmt"
	"strings"
)

// SchemaType represents the expected type for a config value.
type SchemaType int

const (
	TypeAny SchemaType = iota
	TypeString
	TypeInt
	TypeFloat
	TypeBool
	TypeMap
	TypeList
)

func (t SchemaType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBool:
		return "bool"
	case TypeMap:
		return "map"
	case TypeList:
		return "list"
	default:
		return "any"
	}
}

// SchemaNode defines expected type constraints for a key path.
type SchemaNode struct {
	Type     SchemaType
	Children map[string]*SchemaNode
	Required bool
}

// Schema holds the root schema node for validation.
type Schema struct {
	Root *SchemaNode
}

// ValidationError holds all schema violations found during validation.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("schema validation failed:\n  %s", strings.Join(e.Violations, "\n  "))
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Violations) > 0
}

// Validate checks a merged config map against the schema.
func (s *Schema) Validate(data map[string]interface{}) error {
	if s == nil || s.Root == nil {
		return nil
	}
	var violations []string
	validateNode(s.Root, data, "", &violations)
	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}

func validateNode(node *SchemaNode, data map[string]interface{}, prefix string, violations *[]string) {
	for key, child := range node.Children {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		val, exists := data[key]
		if !exists {
			if child.Required {
				*violations = append(*violations, fmt.Sprintf("required key %q is missing", path))
			}
			continue
		}
		checkType(child, val, path, violations)
		if child.Type == TypeMap || child.Type == TypeAny {
			if nested, ok := val.(map[string]interface{}); ok && len(child.Children) > 0 {
				validateNode(child, nested, path, violations)
			}
		}
	}
}

func checkType(node *SchemaNode, val interface{}, path string, violations *[]string) {
	if node.Type == TypeAny {
		return
	}
	var ok bool
	switch node.Type {
	case TypeString:
		_, ok = val.(string)
	case TypeInt:
		switch val.(type) {
		case int, int64, uint64:
			ok = true
		}
	case TypeFloat:
		switch val.(type) {
		case float32, float64:
			ok = true
		}
	case TypeBool:
		_, ok = val.(bool)
	case TypeMap:
		_, ok = val.(map[string]interface{})
	case TypeList:
		_, ok = val.([]interface{})
	}
	if !ok {
		*violations = append(*violations, fmt.Sprintf("key %q: expected type %s, got %T", path, node.Type, val))
	}
}
