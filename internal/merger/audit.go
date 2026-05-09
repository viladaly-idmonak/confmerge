package merger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AuditEntry records a single merge or transform event.
type AuditEntry struct {
	Timestamp time.Time
	Operation string
	Path      string
	OldValue  interface{}
	NewValue  interface{}
	Source    string
}

func (e AuditEntry) String() string {
	return fmt.Sprintf("[%s] %s @ %s: %v -> %v (source: %s)",
		e.Timestamp.Format(time.RFC3339), e.Operation, e.Path, e.OldValue, e.NewValue, e.Source)
}

// AuditLog collects audit entries during a merge session.
type AuditLog struct {
	entries []AuditEntry
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new entry to the log.
func (a *AuditLog) Record(op, path string, oldVal, newVal interface{}, source string) {
	a.entries = append(a.entries, AuditEntry{
		Timestamp: time.Now(),
		Operation: op,
		Path:      path,
		OldValue:  oldVal,
		NewValue:  newVal,
		Source:    source,
	})
}

// Entries returns all recorded audit entries.
func (a *AuditLog) Entries() []AuditEntry {
	return a.entries
}

// Write outputs all audit entries to the given writer.
func (a *AuditLog) Write(w io.Writer) {
	for _, e := range a.entries {
		fmt.Fprintln(w, e.String())
	}
}

// WriteToFile writes the audit log to a file at the given path.
func (a *AuditLog) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("audit: cannot create file %q: %w", path, err)
	}
	defer f.Close()
	a.Write(f)
	return nil
}
