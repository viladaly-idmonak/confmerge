package merger

import "fmt"

// Override records a single value override during a merge.
type Override struct {
	Key      string
	BaseVal  interface{}
	OverVal  interface{}
	Source   string
}

// Result holds the merged map and the list of overrides that occurred.
type Result struct {
	Data      map[string]interface{}
	Overrides []Override
}

// Merge deep-merges a slice of maps in order (later entries override earlier ones).
// Each source label in sources corresponds to the map at the same index.
func Merge(maps []map[string]interface{}, sources []string) (*Result, error) {
	if len(maps) != len(sources) {
		return nil, fmt.Errorf("merger: maps and sources slices must have equal length")
	}

	result := &Result{
		Data:      make(map[string]interface{}),
		Overrides: []Override{},
	}

	for i, m := range maps {
		mergeInto(result.Data, m, sources[i], "", &result.Overrides)
	}

	return result, nil
}

// mergeInto recursively merges src into dst, recording overrides.
func mergeInto(dst, src map[string]interface{}, source, prefix string, overrides *[]Override) {
	for k, srcVal := range src {
		qualified := k
		if prefix != "" {
			qualified = prefix + "." + k
		}

		dstVal, exists := dst[k]
		if !exists {
			dst[k] = srcVal
			continue
		}

		// Both sides are maps — recurse.
		srcMap, srcIsMap := toMap(srcVal)
		dstMap, dstIsMap := toMap(dstVal)
		if srcIsMap && dstIsMap {
			mergeInto(dstMap, srcMap, source, qualified, overrides)
			dst[k] = dstMap
			continue
		}

		// Scalar override.
		*overrides = append(*overrides, Override{
			Key:     qualified,
			BaseVal: dstVal,
			OverVal: srcVal,
			Source:  source,
		})
		dst[k] = srcVal
	}
}

func toMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}
