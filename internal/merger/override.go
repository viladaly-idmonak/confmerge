package merger

import "fmt"

// Override records a single value replacement during a merge operation.
type Override struct {
	Key      string
	BaseVal  interface{}
	OverVal  interface{}
}

// String returns a human-readable representation of the override.
func (o Override) String() string {
	return fmt.Sprintf("%s: %v -> %v", o.Key, o.BaseVal, o.OverVal)
}

// OverrideTracker collects overrides that occur during a merge.
type OverrideTracker struct {
	Overrides []Override
}

// NewOverrideTracker creates an empty OverrideTracker.
func NewOverrideTracker() *OverrideTracker {
	return &OverrideTracker{}
}

// Record stores an override entry.
func (t *OverrideTracker) Record(key string, baseVal, overVal interface{}) {
	t.Overrides = append(t.Overrides, Override{
		Key:     key,
		BaseVal: baseVal,
		OverVal: overVal,
	})
}

// HasOverrides returns true when at least one override was recorded.
func (t *OverrideTracker) HasOverrides() bool {
	return len(t.Overrides) > 0
}

// MergeTracked performs a deep merge of over into base, recording every
// scalar value replacement in the supplied tracker. The merged result is
// returned as a new map.
func MergeTracked(base, over map[string]interface{}, tracker *OverrideTracker) map[string]interface{} {
	return mergeTrackedInto("", base, over, tracker)
}

func mergeTrackedInto(prefix string, base, over map[string]interface{}, tracker *OverrideTracker) map[string]interface{} {
	result := make(map[string]interface{}, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, overVal := range over {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		baseVal, exists := result[k]
		if !exists {
			result[k] = overVal
			continue
		}

		baseMap, baseIsMap := toMap(baseVal)
		overMap, overIsMap := toMap(overVal)

		if baseIsMap && overIsMap {
			result[k] = mergeTrackedInto(fullKey, baseMap, overMap, tracker)
		} else {
			if fmt.Sprintf("%v", baseVal) != fmt.Sprintf("%v", overVal) {
				tracker.Record(fullKey, baseVal, overVal)
			}
			result[k] = overVal
		}
	}
	return result
}
