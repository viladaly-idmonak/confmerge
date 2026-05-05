package merger

import (
	"fmt"
	"io"
	"strings"
)

// DiffFormat controls the output style of PrintDiff.
type DiffFormat int

const (
	DiffFormatText DiffFormat = iota
	DiffFormatUnified
)

// PrintDiff writes a human-readable diff of the recorded overrides to w.
func PrintDiff(w io.Writer, overrides []Override, format DiffFormat) {
	if len(overrides) == 0 {
		fmt.Fprintln(w, "(no overrides)")
		return
	}

	switch format {
	case DiffFormatUnified:
		printUnified(w, overrides)
	default:
		printText(w, overrides)
	}
}

func printText(w io.Writer, overrides []Override) {
	for _, ov := range overrides {
		fmt.Fprintf(w, "key: %s\n", ov.Key)
		fmt.Fprintf(w, "  base:  %v\n", ov.BaseVal)
		fmt.Fprintf(w, "  override (%s): %v\n", ov.Source, ov.OverVal)
	}
}

func printUnified(w io.Writer, overrides []Override) {
	for _, ov := range overrides {
		sep := strings.Repeat("-", 40)
		fmt.Fprintln(w, sep)
		fmt.Fprintf(w, "--- base\t%s\n", ov.Key)
		fmt.Fprintf(w, "+++ %s\t%s\n", ov.Source, ov.Key)
		fmt.Fprintf(w, "- %v\n", ov.BaseVal)
		fmt.Fprintf(w, "+ %v\n", ov.OverVal)
	}
}
