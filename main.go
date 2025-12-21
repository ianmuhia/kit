package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/urfave/cli/v3"
)

//go:embed templates/**/*.tmpl
var templates embed.FS

type Config struct {
	DomainName     string
	OutputDir      string
	WithTests      bool
	WithMessaging  bool
	WithRiver      bool
	WithCQRS       bool
	WithWorkflows  bool
	WithDecorators bool
}

type TemplateData struct {
	DomainTitle string // Capitalized for type names
	DomainLower string // Lowercase for package/file names
}

func main() {
	cmd := &cli.Command{
		Name:  "ddd-lite",
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
			cfg := Config{
				DomainName:     cmd.String("domain"),
				OutputDir:      cmd.String("output"),
				WithTests:      cmd.Bool("with-tests") || cmd.Bool("all"),
				WithMessaging:  cmd.Bool("with-messaging") || cmd.Bool("all"),
				WithRiver:      cmd.Bool("with-river") || cmd.Bool("all"),
				WithCQRS:       cmd.Bool("with-cqrs") || cmd.Bool("all"),
				WithWorkflows:  cmd.Bool("with-workflows") || cmd.Bool("all"),
				WithDecorators: cmd.Bool("with-decorators") || cmd.Bool("all"),
			}

			data := TemplateData{
				DomainTitle: capitalize(cfg.DomainName),
				DomainLower: strings.ToLower(cfg.DomainName),
			}

			if err := generateDomain(cfg, data); err != nil {
				return fmt.Errorf("failed to generate domain: %w", err)
			}

			fmt.Printf("\n✓ SUCCESS: Generated domain '%s' in %s/%s\n",
				data.DomainLower,
				cfg.OutputDir,
				data.DomainLower,
			)
			fmt.Println("\nNext steps:")
			fmt.Printf("  1. Review generated files in %s/%s\n", cfg.OutputDir, data.DomainLower)
			fmt.Printf("  2. Customize domain entity in %s.go\n", data.DomainLower)
			fmt.Println("  3. Add domain-specific repository methods")
			fmt.Println("  4. Implement business logic in app/service.go")
			fmt.Println("  5. Wire up HTTP routes in your application")
			
			if cfg.WithCQRS {
				fmt.Println("  6. Configure Watermill CQRS in cqrs/wiring.go")
			}
			if cfg.WithRiver {
				fmt.Println("  7. Setup River client and run migrations")
			}
			fmt.Println()

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version information",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("ddd-lite version 1.0.0")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func generateDomain(cfg Config, data TemplateData) error {
	fmt.Printf("\nGenerating domain: %s\n", data.DomainTitle)
	fmt.Println(strings.Repeat("-", 50))

	// Create directory structure
	basePath := filepath.Join(cfg.OutputDir, data.DomainLower)
	dirs := []string{
		basePath, // Root domain directory
		filepath.Join(basePath, "app"),
		filepath.Join(basePath, "adapters"),
	}

	if cfg.WithCQRS {
		dirs = append(dirs, filepath.Join(basePath, "cqrs"))
	}

	fmt.Println("\nCreating directories...")
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("  [DIR]  %s\n", dir)
	}

	// Generate files from templates
	files := map[string]string{
		"templates/domain/entity.go.tmpl":     filepath.Join(basePath, data.DomainLower+".go"),
		"templates/domain/repository.go.tmpl": filepath.Join(basePath, "repository.go"),
		"templates/domain/errors.go.tmpl":     filepath.Join(basePath, "errors.go"),
		"templates/domain/events.go.tmpl":     filepath.Join(basePath, "events.go"),
		"templates/domain/validation.go.tmpl": filepath.Join(basePath, "validation.go"),
		"templates/app/service.go.tmpl":       filepath.Join(basePath, "app", "service.go"),
		"templates/adapters/http.go.tmpl":     filepath.Join(basePath, "adapters", data.DomainLower+"_http.go"),
		"templates/adapters/postgres.go.tmpl": filepath.Join(basePath, "adapters", data.DomainLower+"_postgres.go"),
	}

	// Add optional files based on flags
	if cfg.WithTests {
		files["templates/app/service_test.go.tmpl"] = filepath.Join(basePath, "app", "service_test.go")
	}
	if cfg.WithMessaging {
		files["templates/adapters/messaging.go.tmpl"] = filepath.Join(basePath, "adapters", data.DomainLower+"_messaging.go")
	}
	if cfg.WithRiver {
		files["templates/adapters/river.go.tmpl"] = filepath.Join(basePath, "adapters", data.DomainLower+"_river.go")
	}
	if cfg.WithCQRS {
		files["templates/cqrs/commands.go.tmpl"] = filepath.Join(basePath, "cqrs", "commands.go")
		files["templates/cqrs/command_handlers.go.tmpl"] = filepath.Join(basePath, "cqrs", "command_handlers.go")
		files["templates/cqrs/events.go.tmpl"] = filepath.Join(basePath, "cqrs", "events.go")
		files["templates/cqrs/event_handlers.go.tmpl"] = filepath.Join(basePath, "cqrs", "event_handlers.go")
		files["templates/cqrs/wiring.go.tmpl"] = filepath.Join(basePath, "cqrs", "wiring.go")
	}
	if cfg.WithWorkflows {
		files["templates/adapters/temporal.go.tmpl"] = filepath.Join(basePath, "adapters", data.DomainLower+"_temporal.go")
	}
	if cfg.WithDecorators {
		files["templates/app/decorators.go.tmpl"] = filepath.Join(basePath, "app", "decorators.go")
		files["templates/app/wiring_example.go.tmpl"] = filepath.Join(basePath, "app", "wiring_example.go")
	}

	fmt.Println("\nGenerating files from templates...")
	for tmplPath, outputPath := range files {
		if err := generateFile(tmplPath, outputPath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outputPath, err)
		}
		// Show relative path for cleaner output
		relPath, _ := filepath.Rel(cfg.OutputDir, outputPath)
		fmt.Printf("  [FILE] %s\n", relPath)
	}

	return nil
}

func generateFile(tmplPath, outputPath string, data TemplateData) error {
	// Read template from embedded FS
	tmplContent, err := templates.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Execute template
	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
