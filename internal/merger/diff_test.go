package merger

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintDiff_NoOverrides(t *testing.T) {
	var buf bytes.Buffer
	PrintDiff(&buf, []Override{}, DiffFormatText)
	if !strings.Contains(buf.String(), "no overrides") {
		t.Errorf("expected 'no overrides' message, got: %q", buf.String())
	}
}

func TestPrintDiff_TextFormat(t *testing.T) {
	ovs := []Override{
		{Key: "db.port", BaseVal: 5432, OverVal: 9999, Source: "prod"},
	}
	var buf bytes.Buffer
	PrintDiff(&buf, ovs, DiffFormatText)
	out := buf.String()

	if !strings.Contains(out, "db.port") {
		t.Error("expected key db.port in output")
	}
	if !strings.Contains(out, "5432") {
		t.Error("expected base value 5432 in output")
	}
	if !strings.Contains(out, "9999") {
		t.Error("expected override value 9999 in output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected source 'prod' in output")
	}
}

func TestPrintDiff_UnifiedFormat(t *testing.T) {
	ovs := []Override{
		{Key: "host", BaseVal: "localhost", OverVal: "prod.example.com", Source: "prod"},
	}
	var buf bytes.Buffer
	PrintDiff(&buf, ovs, DiffFormatUnified)
	out := buf.String()

	if !strings.Contains(out, "--- base") {
		t.Error("expected '--- base' marker")
	}
	if !strings.Contains(out, "+++ prod") {
		t.Error("expected '+++ prod' marker")
	}
	if !strings.Contains(out, "- localhost") {
		t.Error("expected '- localhost' line")
	}
	if !strings.Contains(out, "+ prod.example.com") {
		t.Error("expected '+ prod.example.com' line")
	}
}
