package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/confmerge/internal/merger"
)

// envFlags holds CLI flags related to environment variable expansion.
type envFlags struct {
	envOverrideFile string
	allowMissingEnv bool
	expandEnvVars   bool
}

func addEnvFlags(fs *flag.FlagSet, ef *envFlags) {
	fs.StringVar(&ef.envOverrideFile, "env-overrides", "", "YAML file mapping config keys to env var names")
	fs.BoolVar(&ef.allowMissingEnv, "allow-missing-env", false, "silently ignore missing env vars during expansion")
	fs.BoolVar(&ef.expandEnvVars, "expand-env", false, "expand ${VAR} references in string config values")
}

// applyEnvFeatures applies env expansion and/or env-override-file logic to data.
func applyEnvFeatures(data map[string]interface{}, ef envFlags) (map[string]interface{}, error) {
	if ef.envOverrideFile != "" {
		if _, err := os.Stat(ef.envOverrideFile); err != nil {
			return nil, fmt.Errorf("env-overrides file not accessible: %w", err)
		}
		overrides, err := merger.LoadEnvOverrides(ef.envOverrideFile)
		if err != nil {
			return nil, fmt.Errorf("loading env overrides: %w", err)
		}
		if err := merger.ApplyEnvOverrides(data, overrides); err != nil {
			return nil, fmt.Errorf("applying env overrides: %w", err)
		}
	}
	if ef.expandEnvVars {
		expander := merger.NewEnvExpander(ef.allowMissingEnv)
		expanded, err := expander.Expand(data)
		if err != nil {
			return nil, fmt.Errorf("expanding env vars: %w", err)
		}
		return expanded, nil
	}
	return data, nil
}
