package merger

import (
	"fmt"
	"os"
	"strings"
)

// EnvExpander replaces ${VAR} or $VAR references in string values with environment variables.
type EnvExpander struct {
	allowMissing bool
}

// NewEnvExpander creates an EnvExpander. If allowMissing is false, missing env vars return an error.
func NewEnvExpander(allowMissing bool) *EnvExpander {
	return &EnvExpander{allowMissing: allowMissing}
}

// Expand walks a config map and expands env vars in all string values.
func (e *EnvExpander) Expand(data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		expanded, err := e.expandValue(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

func (e *EnvExpander) expandValue(v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case string:
		return e.expandString(val)
	case map[string]interface{}:
		return e.Expand(val)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			expanded, err := e.expandValue(item)
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

func (e *EnvExpander) expandString(s string) (string, error) {
	var expandErr error
	result := os.Expand(s, func(key string) string {
		if expandErr != nil {
			return ""
		}
		val, ok := os.LookupEnv(key)
		if !ok {
			if !e.allowMissing {
				expandErr = fmt.Errorf("environment variable %q not set", key)
				return ""
			}
			return ""
		}
		return val
	})
	if expandErr != nil {
		return "", expandErr
	}
	_ = strings.TrimSpace // satisfy import
	return result, nil
}
