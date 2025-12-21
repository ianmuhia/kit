package dddgen

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ianmuhia/kit/pkg/codegen"
)

//go:embed templates/**/*.tmpl
var Templates embed.FS

// Generator handles DDD domain generation
type Generator struct {
	config Config
	data   TemplateData
}

// New creates a new Generator instance
func New(cfg Config) *Generator {
	return &Generator{
		config: cfg,
		data: TemplateData{
			DomainTitle: codegen.Capitalize(cfg.DomainName),
			DomainLower: strings.ToLower(cfg.DomainName),
		},
	}
}

// Generate creates the domain structure and files
func (g *Generator) Generate() error {
	fmt.Printf("\nGenerating domain: %s\n", g.data.DomainTitle)
	fmt.Println(strings.Repeat("-", 50))

	// Create directory structure
	if err := g.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files from templates
	if err := g.generateFiles(); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Print success message
	g.printSuccess()

	return nil
}

func (g *Generator) createDirectories() error {
	basePath := filepath.Join(g.config.OutputDir, g.data.DomainLower)
	dirs := []string{
		basePath,
		filepath.Join(basePath, "app"),
		filepath.Join(basePath, "adapters"),
	}

	if g.config.WithCQRS {
		dirs = append(dirs, filepath.Join(basePath, "cqrs"))
	}

	fmt.Println("\nCreating directories...")
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("  [DIR]  %s\n", dir)
	}

	return nil
}

func (g *Generator) generateFiles() error {
	files := g.getFileMapping()

	fmt.Println("\nGenerating files from templates...")
	for tmplPath, outputPath := range files {
		if err := g.generateFile(tmplPath, outputPath); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outputPath, err)
		}
		relPath, _ := filepath.Rel(g.config.OutputDir, outputPath)
		fmt.Printf("  [FILE] %s\n", relPath)
	}

	return nil
}

func (g *Generator) getFileMapping() map[string]string {
	basePath := filepath.Join(g.config.OutputDir, g.data.DomainLower)

	files := map[string]string{
		"templates/domain/entity.go.tmpl":     filepath.Join(basePath, g.data.DomainLower+".go"),
		"templates/domain/repository.go.tmpl": filepath.Join(basePath, "repository.go"),
		"templates/domain/errors.go.tmpl":     filepath.Join(basePath, "errors.go"),
		"templates/domain/events.go.tmpl":     filepath.Join(basePath, "events.go"),
		"templates/domain/validation.go.tmpl": filepath.Join(basePath, "validation.go"),
		"templates/app/service.go.tmpl":       filepath.Join(basePath, "app", "service.go"),
		"templates/adapters/http.go.tmpl":     filepath.Join(basePath, "adapters", g.data.DomainLower+"_http.go"),
		"templates/adapters/postgres.go.tmpl": filepath.Join(basePath, "adapters", g.data.DomainLower+"_postgres.go"),
	}

	// Add optional files based on flags
	if g.config.WithTests {
		files["templates/app/service_test.go.tmpl"] = filepath.Join(basePath, "app", "service_test.go")
	}
	if g.config.WithMessaging {
		files["templates/adapters/messaging.go.tmpl"] = filepath.Join(basePath, "adapters", g.data.DomainLower+"_messaging.go")
	}
	if g.config.WithRiver {
		files["templates/adapters/river.go.tmpl"] = filepath.Join(basePath, "adapters", g.data.DomainLower+"_river.go")
	}
	if g.config.WithCQRS {
		files["templates/cqrs/commands.go.tmpl"] = filepath.Join(basePath, "cqrs", "commands.go")
		files["templates/cqrs/command_handlers.go.tmpl"] = filepath.Join(basePath, "cqrs", "command_handlers.go")
		files["templates/cqrs/events.go.tmpl"] = filepath.Join(basePath, "cqrs", "events.go")
		files["templates/cqrs/event_handlers.go.tmpl"] = filepath.Join(basePath, "cqrs", "event_handlers.go")
		files["templates/cqrs/wiring.go.tmpl"] = filepath.Join(basePath, "cqrs", "wiring.go")
	}
	if g.config.WithWorkflows {
		files["templates/adapters/temporal.go.tmpl"] = filepath.Join(basePath, "adapters", g.data.DomainLower+"_temporal.go")
	}
	if g.config.WithDecorators {
		files["templates/app/decorators.go.tmpl"] = filepath.Join(basePath, "app", "decorators.go")
		files["templates/app/wiring_example.go.tmpl"] = filepath.Join(basePath, "app", "wiring_example.go")
	}

	return files
}

func (g *Generator) generateFile(tmplPath, outputPath string) error {
	// Read template from embedded FS
	tmplContent, err := Templates.ReadFile(tmplPath)
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
	if err := tmpl.Execute(outFile, g.data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func (g *Generator) printSuccess() {
	fmt.Printf("\n✓ SUCCESS: Generated domain '%s' in %s/%s\n",
		g.data.DomainLower,
		g.config.OutputDir,
		g.data.DomainLower,
	)
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. Review generated files in %s/%s\n", g.config.OutputDir, g.data.DomainLower)
	fmt.Printf("  2. Customize domain entity in %s.go\n", g.data.DomainLower)
	fmt.Println("  3. Add domain-specific repository methods")
	fmt.Println("  4. Implement business logic in app/service.go")
	fmt.Println("  5. Wire up HTTP routes in your application")

	if g.config.WithCQRS {
		fmt.Println("  6. Configure Watermill CQRS in cqrs/wiring.go")
	}
	if g.config.WithRiver {
		fmt.Println("  7. Setup River client and run migrations")
	}
	fmt.Println()
}
