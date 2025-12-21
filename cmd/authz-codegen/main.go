package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ianmuhia/kit/pkg/authzgen"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	var (
		outputDir  = flag.String("output", ".", "Output directory for generated code")
		schemaFile = flag.String("schema", "", "Path to the AuthZed schema file")
	)
	flag.Parse()

	// Handle positional arguments for backward compatibility
	args := flag.Args()
	if *schemaFile == "" && len(args) > 0 {
		*schemaFile = args[0]
	}
	if len(args) > 1 {
		*outputDir = args[1]
	}

	if *schemaFile == "" {
		slog.Error("Schema file is required but not provided")
		fmt.Fprintf(os.Stderr, "Error: schema file is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <schema-file> [output-dir]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create generator with functional options
	generator, err := authzgen.NewGenerator(
		authzgen.WithSchemaFile(*schemaFile),
		authzgen.WithOutputDir(*outputDir),
		authzgen.WithLogger(logger),
	)
	if err != nil {
		slog.Error("Failed to create generator", "error", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Generate code
	if err := generator.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Code generation completed successfully!\n")
}
