package merger

import (
	"testing"
)

func TestLinter_NoRules(t *testing.T) {
	l := NewLinter(nil)
	data := map[string]interface{}{"key": "value"}
	results := l.Lint(data)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestLinter_ExactMatch(t *testing.T) {
	rules := []LintRule{
		{Pattern: "database.password", Description: "avoid plaintext passwords", Severity: "error"},
	}
	l := NewLinter(rules)
	data := map[string]interface{}{
		"database": map[string]interface{}{"password": "secret"},
	}
	results := l.Lint(data)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Path != "database.password" {
		t.Errorf("unexpected path: %s", results[0].Path)
	}
	if results[0].Severity != "error" {
		t.Errorf("unexpected severity: %s", results[0].Severity)
	}
}

func TestLinter_WildcardMatch(t *testing.T) {
	rules := []LintRule{
		{Pattern: "secrets.*", Description: "secrets namespace detected", Severity: "warn"},
	}
	l := NewLinter(rules)
	data := map[string]interface{}{
		"secrets": map[string]interface{}{"token": "abc", "key": "xyz"},
	}
	results := l.Lint(data)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestLinter_NoMatch(t *testing.T) {
	rules := []LintRule{
		{Pattern: "database.password", Description: "plaintext password", Severity: "error"},
	}
	l := NewLinter(rules)
	data := map[string]interface{}{"app": map[string]interface{}{"name": "myapp"}}
	results := l.Lint(data)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestLintResult_String(t *testing.T) {
	r := LintResult{Path: "db.pass", Message: "plaintext", Severity: "error"}
	got := r.String()
	expected := "[ERROR] db.pass: plaintext"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestMatchPattern_Exact(t *testing.T) {
	if !matchPattern("a.b.c", "a.b.c") {
		t.Error("expected match")
	}
	if matchPattern("a.b.c", "a.b") {
		t.Error("expected no match")
	}
}

func TestMatchPattern_Wildcard(t *testing.T) {
	if !matchPattern("secrets.*", "secrets.token") {
		t.Error("expected match")
	}
	if matchPattern("secrets.*", "secrets") {
		t.Error("expected no match for prefix itself")
	}
}
