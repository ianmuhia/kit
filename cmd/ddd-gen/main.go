package main

import (
	"context"
	"log"
	"os"

	"github.com/ianmuhia/kit/internal/dddgen"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "ddd-gen",
		Usage: "Generate DDD domain modules for Go projects",
		Authors: []any{
			"Ian Muhia <https://github.com/Ianmuhia>",
		},
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "domain",
				Aliases:  []string{"d"},
				Usage:    "Domain name (e.g., 'booking', 'user', 'order')",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output directory for generated code",
				Value:   "./internal",
			},
			&cli.BoolFlag{
				Name:    "with-tests",
				Aliases: []string{"t"},
				Usage:   "Generate test files",
			},
			&cli.BoolFlag{
				Name:    "with-messaging",
				Aliases: []string{"m"},
				Usage:   "Generate messaging adapter (Watermill pub/sub)",
			},
			&cli.BoolFlag{
				Name:    "with-river",
				Aliases: []string{"r"},
				Usage:   "Generate River job queue adapter",
			},
			&cli.BoolFlag{
				Name:    "with-cqrs",
				Aliases: []string{"c"},
				Usage:   "Generate CQRS components (Watermill commands, events, handlers)",
			},
			&cli.BoolFlag{
				Name:    "with-workflows",
				Aliases: []string{"w"},
				Usage:   "Generate Temporal workflow adapter",
			},
			&cli.BoolFlag{
				Name:  "with-decorators",
				Usage: "Generate service decorators (permissions, audit, cache, metrics)",
			},
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Generate all optional components",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := dddgen.Config{
				DomainName:     cmd.String("domain"),
				OutputDir:      cmd.String("output"),
				WithTests:      cmd.Bool("with-tests") || cmd.Bool("all"),
				WithMessaging:  cmd.Bool("with-messaging") || cmd.Bool("all"),
				WithRiver:      cmd.Bool("with-river") || cmd.Bool("all"),
				WithCQRS:       cmd.Bool("with-cqrs") || cmd.Bool("all"),
				WithWorkflows:  cmd.Bool("with-workflows") || cmd.Bool("all"),
				WithDecorators: cmd.Bool("with-decorators") || cmd.Bool("all"),
			}

			generator := dddgen.New(cfg)
			return generator.Generate()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
