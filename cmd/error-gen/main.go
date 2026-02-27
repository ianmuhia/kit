package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ianmuhia/kit/pkg/errorgen"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "error-gen",
		Usage:   "Generate strongly-typed error codes from CUE definitions",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input CUE file or directory",
				Value:   "errors.cue",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output Go file path",
				Value:   "errors.go",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "Custom error template file (optional)",
			},
			&cli.StringFlag{
				Name:    "package",
				Aliases: []string{"p"},
				Usage:   "Override package name (optional)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			opts := []errorgen.GeneratorOption{
				errorgen.WithInputFile(cmd.String("input")),
				errorgen.WithOutputFile(cmd.String("output")),
			}

			if t := cmd.String("template"); t != "" {
				opts = append(opts, errorgen.WithTemplateFile(t))
			}
			if p := cmd.String("package"); p != "" {
				opts = append(opts, errorgen.WithPackageName(p))
			}

			generator, err := errorgen.NewGenerator(opts...)
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			if err := generator.Generate(); err != nil {
				return fmt.Errorf("failed to generate code: %w", err)
			}

			fmt.Printf("✓ Error code generated successfully in %s\n", cmd.String("output"))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
