package merger

import (
	"testing"
)

func TestFreezer_NoRules(t *testing.T) {
	fz := NewFreezer(nil)
	base := map[string]interface{}{"key": "value"}
	merged := map[string]interface{}{"key": "changed"}
	errs := fz.CheckViolations(base, merged)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestFreezer_NoViolation(t *testing.T) {
	rules := []FreezeRule{{Path: "server.port"}}
	fz := NewFreezer(rules)
	base := map[string]interface{}{"server": map[string]interface{}{"port": 8080}}
	merged := map[string]interface{}{"server": map[string]interface{}{"port": 8080}}
	errs := fz.CheckViolations(base, merged)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestFreezer_ScalarViolation(t *testing.T) {
	rules := []FreezeRule{{Path: "app.name"}}
	fz := NewFreezer(rules)
	base := map[string]interface{}{"app": map[string]interface{}{"name": "myapp"}}
	merged := map[string]interface{}{"app": map[string]interface{}{"name": "otherapp"}}
	errs := fz.CheckViolations(base, merged)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestFreezer_PathNotInBase(t *testing.T) {
	// If base doesn't have the path, no freeze violation should be raised.
	rules := []FreezeRule{{Path: "missing.key"}}
	fz := NewFreezer(rules)
	base := map[string]interface{}{}
	merged := map[string]interface{}{"missing": map[string]interface{}{"key": "val"}}
	errs := fz.CheckViolations(base, merged)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestFreezer_MultipleRules(t *testing.T) {
	rules := []FreezeRule{{Path: "db.host"}, {Path: "db.port"}}
	fz := NewFreezer(rules)
	base := map[string]interface{}{"db": map[string]interface{}{"host": "localhost", "port": 5432}}
	merged := map[string]interface{}{"db": map[string]interface{}{"host": "remotehost", "port": 5432}}
	errs := fz.CheckViolations(base, merged)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestFreezeRule_String(t *testing.T) {
	r := FreezeRule{Path: "server.host"}
	if r.String() != "freeze:server.host" {
		t.Errorf("unexpected string: %s", r.String())
	}
}
