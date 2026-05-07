package main

import (
	"fmt"
	"io"
	"sort"
)

// writeYAML writes a map as simple YAML to the given writer.
// It handles nested maps and scalar values with basic indentation.
func writeYAML(w io.Writer, data map[string]interface{}) error {
	return writeYAMLIndented(w, data, 0)
}

func writeYAMLIndented(w io.Writer, data map[string]interface{}, depth int) error {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		switch val := v.(type) {
		case map[string]interface{}:
			if _, err := fmt.Fprintf(w, "%s%s:\n", indent, k); err != nil {
				return err
			}
			if err := writeYAMLIndented(w, val, depth+1); err != nil {
				return err
			}
		case []interface{}:
			if _, err := fmt.Fprintf(w, "%s%s:\n", indent, k); err != nil {
				return err
			}
			for _, item := range val {
				if _, err := fmt.Fprintf(w, "%s  - %v\n", indent, item); err != nil {
					return err
				}
			}
		default:
			if _, err := fmt.Fprintf(w, "%s%s: %v\n", indent, k, val); err != nil {
				return err
			}
		}
	}
	return nil
}
