package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ianmuhia/kit/pkg/errorgen"
)

func main() {
	// Parse command-line flags
	inputFile := flag.String("input", "errors.cue", "Input CUE file or directory")
	outputFile := flag.String("output", "errors.go", "Output Go file")
	templateFile := flag.String("template", "", "Custom error template file (optional)")
	packageName := flag.String("package", "", "Override package name (optional)")
	flag.Parse()

	// Build options
	opts := []errorgen.GeneratorOption{
		errorgen.WithInputFile(*inputFile),
		errorgen.WithOutputFile(*outputFile),
	}

	if *templateFile != "" {
		opts = append(opts, errorgen.WithTemplateFile(*templateFile))
	}

	if *packageName != "" {
		opts = append(opts, errorgen.WithPackageName(*packageName))
	}

	// Create generator and run
	generator, err := errorgen.NewGenerator(opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating generator: %v\n", err)
		os.Exit(1)
	}

	if err := generator.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Error code generated successfully in %s\n", *outputFile)
}
