package merger

import "fmt"

// FreezeRule defines a key path that should not be overridden after a base layer.
type FreezeRule struct {
	Path string
}

func (f FreezeRule) String() string {
	return fmt.Sprintf("freeze:%s", f.Path)
}

// Freezer holds a set of frozen key paths and enforces them during merge.
type Freezer struct {
	rules []FreezeRule
}

// NewFreezer creates a Freezer from a list of FreezeRules.
func NewFreezer(rules []FreezeRule) *Freezer {
	return &Freezer{rules: rules}
}

// CheckViolations inspects a merged map against a base map and returns errors
// for any frozen paths whose values differ between base and merged.
func (fz *Freezer) CheckViolations(base, merged map[string]interface{}) []error {
	var errs []error
	for _, rule := range fz.rules {
		parts := splitPath(rule.Path)
		baseVal, baseOk := getNestedValue(base, parts)
		mergedVal, mergedOk := getNestedValue(merged, parts)
		if !baseOk {
			// path not present in base — nothing to freeze
			continue
		}
		if !mergedOk || fmt.Sprintf("%v", baseVal) != fmt.Sprintf("%v", mergedVal) {
			errs = append(errs, fmt.Errorf("frozen key %q was modified: base=%v merged=%v", rule.Path, baseVal, mergedVal))
		}
	}
	return errs
}

// getNestedValue traverses a map using the provided path parts and returns the value.
func getNestedValue(m map[string]interface{}, parts []string) (interface{}, bool) {
	if len(parts) == 0 {
		return nil, false
	}
	val, ok := m[parts[0]]
	if !ok {
		return nil, false
	}
	if len(parts) == 1 {
		return val, true
	}
	child, ok := val.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return getNestedValue(child, parts[1:])
}
