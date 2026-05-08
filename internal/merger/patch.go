package merger

import (
	"fmt"
	"strings"
)

// PatchOp represents a single patch operation type.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
)

// Patch represents a single key-path mutation to apply to a config map.
type Patch struct {
	Op    PatchOp
	Path  []string
	Value interface{}
}

// String returns a human-readable representation of the patch.
func (p Patch) String() string {
	path := strings.Join(p.Path, ".")
	if p.Op == PatchDelete {
		return fmt.Sprintf("delete %s", path)
	}
	return fmt.Sprintf("set %s = %v", path, p.Value)
}

// ApplyPatches applies a slice of Patch operations to the given config map,
// returning a new map with mutations applied. The original map is not modified.
func ApplyPatches(base map[string]interface{}, patches []Patch) (map[string]interface{}, error) {
	result := deepCopyMap(base)
	for _, p := range patches {
		if len(p.Path) == 0 {
			return nil, fmt.Errorf("patch has empty path")
		}
		switch p.Op {
		case PatchSet:
			if err := setNested(result, p.Path, p.Value); err != nil {
				return nil, fmt.Errorf("patch set %s: %w", strings.Join(p.Path, "."), err)
			}
		case PatchDelete:
			deleteNested(result, p.Path)
		default:
			return nil, fmt.Errorf("unknown patch op: %q", p.Op)
		}
	}
	return result, nil
}

func setNested(m map[string]interface{}, path []string, value interface{}) error {
	if len(path) == 1 {
		m[path[0]] = value
		return nil
	}
	child, ok := m[path[0]]
	if !ok {
		child = map[string]interface{}{}
		m[path[0]] = child
	}
	childMap, ok := child.(map[string]interface{})
	if !ok {
		return fmt.Errorf("key %q is not a map", path[0])
	}
	return setNested(childMap, path[1:], value)
}

func deleteNested(m map[string]interface{}, path []string) {
	if len(path) == 1 {
		delete(m, path[0])
		return
	}
	child, ok := m[path[0]]
	if !ok {
		return
	}
	childMap, ok := child.(map[string]interface{})
	if !ok {
		return
	}
	deleteNested(childMap, path[1:])
}

func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{}, len(m))
	for k, v := range m {
		if nested, ok := v.(map[string]interface{}); ok {
			copy[k] = deepCopyMap(nested)
		} else {
			copy[k] = v
		}
	}
	return copy
}
