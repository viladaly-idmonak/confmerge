package merger

import (
	"fmt"
	"regexp"
	"strings"
)

// SecretMasker replaces sensitive values in a config map with a redacted placeholder.
type SecretMasker struct {
	patterns []*regexp.Regexp
	placeholder string
}

// NewSecretMasker creates a SecretMasker that redacts keys matching any of the
// provided glob-style patterns (e.g. "*password*", "*secret*").
func NewSecretMasker(patterns []string, placeholder string) (*SecretMasker, error) {
	if placeholder == "" {
		placeholder = "***REDACTED***"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		// Convert simple glob (*) to regex
		regexStr := "(?i)^" + strings.ReplaceAll(regexp.QuoteMeta(p), `\*`, `.*`) + "$"
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &SecretMasker{patterns: compiled, placeholder: placeholder}, nil
}

// Mask returns a deep copy of m with matching leaf values replaced by the placeholder.
func (s *SecretMasker) Mask(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			out[k] = s.Mask(nested)
		} else if s.matchesAny(k) {
			out[k] = s.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

func (s *SecretMasker) matchesAny(key string) bool {
	for _, re := range s.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
