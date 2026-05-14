package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/confmerge/internal/merger"
)

// addLintFlags registers lint-related CLI flags onto the given FlagSet.
func addLintFlags(fs *flag.FlagSet) *string {
	return fs.String("lint", "", "path to lint rules YAML file")
}

// applyLint loads lint rules and runs the linter against the merged config.
// It prints results to stderr and returns true if any errors were found.
func applyLint(lintFile string, data map[string]interface{}) (hasErrors bool) {
	if lintFile == "" {
		return false
	}

	rules, err := merger.LoadLintRules(lintFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint: failed to load rules: %v\n", err)
		return true
	}

	linter := merger.NewLinter(rules)
	results := linter.Lint(data)

	for _, r := range results {
		fmt.Fprintln(os.Stderr, r.String())
		if r.Severity == "error" {
			hasErrors = true
		}
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "lint: no issues found")
	}
	return hasErrors
}
