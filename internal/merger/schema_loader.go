package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// schemaFile is the intermediate YAML representation of a schema.
type schemaFile struct {
	Fields map[string]schemaFieldDef `yaml:"fields"`
}

type schemaFieldDef struct {
	Type     string                    `yaml:"type"`
	Required bool                      `yaml:"required"`
	Fields   map[string]schemaFieldDef `yaml:"fields"`
}

// LoadSchema reads a YAML schema definition file and returns a Schema.
func LoadSchema(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file %q: %w", path, err)
	}
	var sf schemaFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("parsing schema file %q: %w", path, err)
	}
	root := &SchemaNode{
		Type:     TypeMap,
		Children: buildChildren(sf.Fields),
	}
	return &Schema{Root: root}, nil
}

func buildChildren(defs map[string]schemaFieldDef) map[string]*SchemaNode {
	if len(defs) == 0 {
		return nil
	}
	children := make(map[string]*SchemaNode, len(defs))
	for name, def := range defs {
		node := &SchemaNode{
			Type:     parseSchemaType(def.Type),
			Required: def.Required,
		}
		if len(def.Fields) > 0 {
			node.Children = buildChildren(def.Fields)
		}
		children[name] = node
	}
	return children
}

func parseSchemaType(s string) SchemaType {
	switch s {
	case "string":
		return TypeString
	case "int":
		return TypeInt
	case "float":
		return TypeFloat
	case "bool":
		return TypeBool
	case "map":
		return TypeMap
	case "list":
		return TypeList
	default:
		return TypeAny
	}
}
