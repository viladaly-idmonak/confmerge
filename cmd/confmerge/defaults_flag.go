package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/confmerge/internal/merger"
)

// addDefaultsFlags registers the --defaults flag on the given FlagSet.
func addDefaultsFlags(fs *flag.FlagSet, defaultsFile *string) {
	fs.StringVar(defaultsFile, "defaults", "", "path to YAML file defining default values")
}

// applyDefaults loads default rules from defaultsFile (if non-empty) and
// applies them to dst. Missing keys are filled with their declared defaults.
func applyDefaults(defaultsFile string, dst map[string]interface{}) error {
	if defaultsFile == "" {
		return nil
	}

	rules, err := merger.LoadDefaults(defaultsFile)
	if err != nil {
		return fmt.Errorf("--defaults: %w", err)
	}

	defaulter := merger.NewDefaulter(rules)
	if err := defaulter.Apply(dst); err != nil {
		return fmt.Errorf("--defaults: apply: %w", err)
	}

	fmt.Fprintf(os.Stderr, "[defaults] applied %d rule(s) from %s\n", len(rules), defaultsFile)
	return nil
}
