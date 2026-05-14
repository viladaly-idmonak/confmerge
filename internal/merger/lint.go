package merger

import (
	"fmt"
	"strings"
)

// LintRule represents a single linting rule for config keys.
type LintRule struct {
	Pattern     string // key pattern to match (supports wildcard *)
	Description string // human-readable description of the rule
	Severity    string // "error" or "warn"
}

// LintResult holds the outcome of a lint check.
type LintResult struct {
	Path     string
	Message  string
	Severity string
}

func (r LintResult) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(r.Severity), r.Path, r.Message)
}

// Linter applies lint rules to a merged config map.
type Linter struct {
	rules []LintRule
}

// NewLinter creates a Linter with the given rules.
func NewLinter(rules []LintRule) *Linter {
	return &Linter{rules: rules}
}

// Lint walks the config map and returns any lint results.
func (l *Linter) Lint(data map[string]interface{}) []LintResult {
	var results []LintResult
	l.walk(data, "", &results)
	return results
}

func (l *Linter) walk(node interface{}, path string, results *[]LintResult) {
	switch v := node.(type) {
	case map[string]interface{}:
		for k, val := range v {
			child := k
			if path != "" {
				child = path + "." + k
			}
			for _, rule := range l.rules {
				if matchPattern(rule.Pattern, child) {
					*results = append(*results, LintResult{
						Path:     child,
						Message:  rule.Description,
						Severity: rule.Severity,
					})
				}
			}
			l.walk(val, child, results)
		}
	}
}

// matchPattern matches a dot-separated path against a pattern supporting trailing *.
func matchPattern(pattern, path string) bool {
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(path, prefix+".")
	}
	return pattern == path
}
