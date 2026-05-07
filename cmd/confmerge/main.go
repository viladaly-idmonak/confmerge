package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/confmerge/internal/loader"
	"github.com/user/confmerge/internal/merger"
)

func main() {
	var (
		diffFormat = flag.String("diff", "", "diff output format: text or unified")
		outputFile = flag.String("output", "", "write merged result to file (default: stdout)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: confmerge [options] base.yaml override1.yaml [override2.yaml ...]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	base, err := loader.Load(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base file %q: %v\n", args[0], err)
		os.Exit(1)
	}

	tracker := merger.NewOverrideTracker()
	result := base

	for _, path := range args[1:] {
		override, err := loader.Load(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading override file %q: %v\n", path, err)
			os.Exit(1)
		}
		result = merger.MergeTracked(result, override, path, tracker)
	}

	if *diffFormat != "" {
		merger.PrintDiff(tracker, *diffFormat)
	}

	out := os.Stdout
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	if err := writeYAML(out, result); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}
}
