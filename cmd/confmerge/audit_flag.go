package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/confmerge/internal/merger"
)

// auditFlags holds audit-related CLI options.
type auditFlags struct {
	enabled  bool
	filePath string
}

// addAuditFlags registers audit flags onto the provided FlagSet.
func addAuditFlags(fs *flag.FlagSet, af *auditFlags) {
	fs.BoolVar(&af.enabled, "audit", false, "enable audit logging of merge operations")
	fs.StringVar(&af.filePath, "audit-file", "", "write audit log to this file (default: stdout)")
}

// applyAuditLog writes the audit log according to the configured flags.
// It prints to stdout if no file path is given, or writes to the specified file.
func applyAuditLog(af auditFlags, log *merger.AuditLog) error {
	if !af.enabled {
		return nil
	}
	if af.filePath != "" {
		if err := log.WriteToFile(af.filePath); err != nil {
			return fmt.Errorf("audit: %w", err)
		}
		return nil
	}
	log.Write(os.Stdout)
	return nil
}
