package merger

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditLog_Empty(t *testing.T) {
	log := NewAuditLog()
	if len(log.Entries()) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(log.Entries()))
	}
}

func TestAuditLog_Record(t *testing.T) {
	log := NewAuditLog()
	log.Record("override", "db.host", "localhost", "prod.host", "prod.yaml")
	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Operation != "override" {
		t.Errorf("expected operation 'override', got %q", e.Operation)
	}
	if e.Path != "db.host" {
		t.Errorf("expected path 'db.host', got %q", e.Path)
	}
	if e.Source != "prod.yaml" {
		t.Errorf("expected source 'prod.yaml', got %q", e.Source)
	}
}

func TestAuditEntry_String(t *testing.T) {
	log := NewAuditLog()
	log.Record("add", "app.name", nil, "myapp", "base.yaml")
	s := log.Entries()[0].String()
	if !strings.Contains(s, "add") || !strings.Contains(s, "app.name") {
		t.Errorf("unexpected String() output: %s", s)
	}
}

func TestAuditLog_Write(t *testing.T) {
	log := NewAuditLog()
	log.Record("override", "x", 1, 2, "s")
	var buf bytes.Buffer
	log.Write(&buf)
	if !strings.Contains(buf.String(), "override") {
		t.Errorf("expected 'override' in output, got: %s", buf.String())
	}
}

func TestAuditLog_WriteToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	log := NewAuditLog()
	log.Record("add", "key", nil, "val", "src")
	if err := log.WriteToFile(path); err != nil {
		t.Fatalf("WriteToFile error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "add") {
		t.Errorf("expected 'add' in file, got: %s", string(data))
	}
}
