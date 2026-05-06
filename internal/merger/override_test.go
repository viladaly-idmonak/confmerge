package merger

import (
	"testing"
)

func TestOverrideTracker_Empty(t *testing.T) {
	tracker := NewOverrideTracker()
	if tracker.HasOverrides() {
		t.Fatal("expected no overrides on a fresh tracker")
	}
}

func TestOverrideTracker_Record(t *testing.T) {
	tracker := NewOverrideTracker()
	tracker.Record("host", "localhost", "prod.example.com")

	if !tracker.HasOverrides() {
		t.Fatal("expected HasOverrides to return true after Record")
	}
	if len(tracker.Overrides) != 1 {
		t.Fatalf("expected 1 override, got %d", len(tracker.Overrides))
	}
	o := tracker.Overrides[0]
	if o.Key != "host" || o.BaseVal != "localhost" || o.OverVal != "prod.example.com" {
		t.Errorf("unexpected override contents: %+v", o)
	}
}

func TestOverride_String(t *testing.T) {
	o := Override{Key: "port", BaseVal: 8080, OverVal: 9090}
	got := o.String()
	want := "port: 8080 -> 9090"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMergeTracked_ScalarOverride(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "port": 8080}
	over := map[string]interface{}{"host": "prod.example.com"}

	tracker := NewOverrideTracker()
	result := MergeTracked(base, over, tracker)

	if result["host"] != "prod.example.com" {
		t.Errorf("expected host to be overridden, got %v", result["host"])
	}
	if result["port"] != 8080 {
		t.Errorf("expected port to remain 8080, got %v", result["port"])
	}
	if len(tracker.Overrides) != 1 || tracker.Overrides[0].Key != "host" {
		t.Errorf("expected one override for 'host', got %+v", tracker.Overrides)
	}
}

func TestMergeTracked_DeepOverride(t *testing.T) {
	base := map[string]interface{}{
		"database": map[string]interface{}{"host": "localhost", "port": 5432},
	}
	over := map[string]interface{}{
		"database": map[string]interface{}{"host": "db.prod"},
	}

	tracker := NewOverrideTracker()
	result := MergeTracked(base, over, tracker)

	db, ok := result["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected database to be a map")
	}
	if db["host"] != "db.prod" {
		t.Errorf("expected db.host = db.prod, got %v", db["host"])
	}
	if db["port"] != 5432 {
		t.Errorf("expected db.port = 5432, got %v", db["port"])
	}
	if len(tracker.Overrides) != 1 || tracker.Overrides[0].Key != "database.host" {
		t.Errorf("expected override for 'database.host', got %+v", tracker.Overrides)
	}
}

func TestMergeTracked_NoChangeNotRecorded(t *testing.T) {
	base := map[string]interface{}{"level": "info"}
	over := map[string]interface{}{"level": "info"}

	tracker := NewOverrideTracker()
	MergeTracked(base, over, tracker)

	if tracker.HasOverrides() {
		t.Errorf("expected no overrides when values are identical, got %+v", tracker.Overrides)
	}
}
