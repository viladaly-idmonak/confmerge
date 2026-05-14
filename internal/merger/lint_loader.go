package merger

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type lintFileRule struct {
	Pattern     string `yaml:"pattern"`
	Description string `yaml:"description"`
	Severity    string `yaml:"severity"`
}

type lintFile struct {
	Rules []lintFileRule `yaml:"rules"`
}

// LoadLintRules reads a YAML lint rules file and returns a slice of LintRule.
func LoadLintRules(path string) ([]LintRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lint_loader: read file %q: %w", path, err)
	}

	var lf lintFile
	if err := yaml.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("lint_loader: parse YAML %q: %w", path, err)
	}

	var rules []LintRule
	for i, r := range lf.Rules {
		if r.Pattern == "" {
			return nil, fmt.Errorf("lint_loader: rule[%d] missing pattern", i)
		}
		if r.Description == "" {
			return nil, fmt.Errorf("lint_loader: rule[%d] missing description", i)
		}
		sev := r.Severity
		if sev == "" {
			sev = "warn"
		}
		if sev != "warn" && sev != "error" {
			return nil, fmt.Errorf("lint_loader: rule[%d] invalid severity %q", i, sev)
		}
		rules = append(rules, LintRule{
			Pattern:     r.Pattern,
			Description: r.Description,
			Severity:    sev,
		})
	}
	return rules, nil
}
