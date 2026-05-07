package merger

import (
	"testing"
)

func makeSchema() *Schema {
	return &Schema{
		Root: &SchemaNode{
			Type: TypeMap,
			Children: map[string]*SchemaNode{
				"host": {Type: TypeString, Required: true},
				"port": {Type: TypeInt, Required: true},
				"debug": {Type: TypeBool},
				"database": {
					Type: TypeMap,
					Children: map[string]*SchemaNode{
						"name": {Type: TypeString, Required: true},
						"pool": {Type: TypeInt},
					},
				},
			},
		},
	}
}

func TestSchema_Valid(t *testing.T) {
	s := makeSchema()
	data := map[string]interface{}{
		"host":  "localhost",
		"port":  8080,
		"debug": true,
		"database": map[string]interface{}{
			"name": "mydb",
			"pool": 5,
		},
	}
	if err := s.Validate(data); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSchema_MissingRequired(t *testing.T) {
	s := makeSchema()
	data := map[string]interface{}{
		"host": "localhost",
		// port is missing
	}
	err := s.Validate(data)
	if err == nil {
		t.Fatal("expected validation error for missing required key")
	}
	ve, ok := err.(*ValidationError)
	if !ok || !ve.HasErrors() {
		t.Fatalf("expected ValidationError, got: %v", err)
	}
}

func TestSchema_WrongType(t *testing.T) {
	s := makeSchema()
	data := map[string]interface{}{
		"host": 1234, // should be string
		"port": 8080,
	}
	err := s.Validate(data)
	if err == nil {
		t.Fatal("expected validation error for wrong type")
	}
}

func TestSchema_NestedMissingRequired(t *testing.T) {
	s := makeSchema()
	data := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
		"database": map[string]interface{}{
			// name is missing
			"pool": 5,
		},
	}
	err := s.Validate(data)
	if err == nil {
		t.Fatal("expected validation error for nested missing required key")
	}
}

func TestSchema_NilSchema(t *testing.T) {
	var s *Schema
	data := map[string]interface{}{"any": "value"}
	if err := s.Validate(data); err != nil {
		t.Fatalf("nil schema should not error, got: %v", err)
	}
}

func TestSchemaType_String(t *testing.T) {
	cases := map[SchemaType]string{
		TypeString: "string",
		TypeInt:    "int",
		TypeFloat:  "float",
		TypeBool:   "bool",
		TypeMap:    "map",
		TypeList:   "list",
		TypeAny:    "any",
	}
	for typ, expected := range cases {
		if got := typ.String(); got != expected {
			t.Errorf("SchemaType(%d).String() = %q, want %q", typ, got, expected)
		}
	}
}
