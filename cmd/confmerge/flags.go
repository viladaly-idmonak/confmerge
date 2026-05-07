package main

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the parsed command-line configuration for confmerge.
type Config struct {
	// InputFiles is the ordered list of config files to merge (base first, overrides after).
	InputFiles []string

	// OutputFile is the path to write the merged result. Empty means stdout.
	OutputFile string

	// DiffFormat controls how override diffs are displayed: "none", "text", or "unified".
	DiffFormat string

	// OutputFormat forces the output format: "yaml" or "toml".
	// If empty, it is inferred from the OutputFile extension or defaults to "yaml".
	OutputFormat string

	// Indent specifies the number of spaces used for YAML indentation.
	Indent int

	// Strict causes the tool to exit with a non-zero status if any overrides are detected.
	Strict bool
}

// parseFlags parses os.Args and returns a populated Config.
// It prints usage and exits on error or when -help is requested.
func parseFlags() Config {
	fs := flag.NewFlagSet("confmerge", flag.ExitOnError)

	output := fs.String("output", "", "Write merged result to `file` instead of stdout")
	diffFormat := fs.String("diff", "text", "Diff format for overrides: none, text, or unified")
	outputFormat := fs.String("format", "", "Force output format: yaml or toml (default: infer from output file or yaml)")
	indent := fs.Int("indent", 2, "Number of spaces for YAML indentation")
	strict := fs.Bool("strict", false, "Exit with non-zero status if any overrides are detected")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: confmerge [options] <base-file> <override-file> [override-file ...]\n\n")
		fmt.Fprintf(os.Stderr, "Deep-merge layered YAML/TOML config files with override tracking and diff output.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  confmerge base.yaml override.yaml\n")
		fmt.Fprintf(os.Stderr, "  confmerge -output merged.yaml -diff unified base.toml prod.toml secrets.toml\n")
		fmt.Fprintf(os.Stderr, "  confmerge -strict -diff none base.yaml ci.yaml\n")
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag.ExitOnError handles this, but be explicit.
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(2)
	}

	args := fs.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "error: at least two input files are required\n\n")
		fs.Usage()
		os.Exit(2)
	}

	if err := validateDiffFormat(*diffFormat); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if err := validateOutputFormat(*outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	return Config{
		InputFiles:   args,
		OutputFile:   *output,
		DiffFormat:   *diffFormat,
		OutputFormat: *outputFormat,
		Indent:       *indent,
		Strict:       *strict,
	}
}

// validateDiffFormat returns an error if the given format string is not supported.
func validateDiffFormat(f string) error {
	switch f {
	case "none", "text", "unified":
		return nil
	default:
		return fmt.Errorf("unsupported diff format %q: must be one of none, text, unified", f)
	}
}

// validateOutputFormat returns an error if the given format string is not supported.
func validateOutputFormat(f string) error {
	switch f {
	case "", "yaml", "toml":
		return nil
	default:
		return fmt.Errorf("unsupported output format %q: must be yaml or toml", f)
	}
}
