package merger

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateExpander applies Go text/template expressions to string values in a map.
type TemplateExpander struct {
	vars map[string]interface{}
}

// NewTemplateExpander creates a TemplateExpander with the provided template variables.
func NewTemplateExpander(vars map[string]interface{}) *TemplateExpander {
	return &TemplateExpander{vars: vars}
}

// Expand walks the config map and renders any string values as Go templates.
func (te *TemplateExpander) Expand(data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		expanded, err := te.expandValue(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

func (te *TemplateExpander) expandValue(v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case string:
		return te.renderTemplate(val)
	case map[string]interface{}:
		return te.Expand(val)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			expanded, err := te.expandValue(item)
			if err != nil {
				return nil, fmt.Errorf("index %d: %w", i, err)
			}
			result[i] = expanded
		}
		return result, nil
	default:
		return v, nil
	}
}

func (te *TemplateExpander) renderTemplate(s string) (string, error) {
	tmpl, err := template.New("").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", s, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, te.vars); err != nil {
		return "", fmt.Errorf("execute template %q: %w", s, err)
	}
	return buf.String(), nil
}
