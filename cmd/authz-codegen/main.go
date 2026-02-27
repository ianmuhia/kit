package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ianmuhia/kit/pkg/authzgen"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "authz-codegen",
		Usage:   "Generate type-safe Go client code from AuthZed permission schemas",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "schema",
				Aliases:  []string{"s"},
				Usage:    "Path to the AuthZed schema (.zed) file",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output directory for generated code",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Log level (debug, info, warn, error)",
				Value: "info",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			level := slog.LevelInfo
			switch cmd.String("log-level") {
			case "debug":
				level = slog.LevelDebug
			case "warn":
				level = slog.LevelWarn
			case "error":
				level = slog.LevelError
			}

			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
			slog.SetDefault(logger)

			generator, err := authzgen.NewGenerator(
				authzgen.WithSchemaFile(cmd.String("schema")),
				authzgen.WithOutputDir(cmd.String("output")),
				authzgen.WithLogger(logger),
			)
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			if err := generator.Generate(); err != nil {
				return fmt.Errorf("code generation failed: %w", err)
			}

			fmt.Println("Code generation completed successfully!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
