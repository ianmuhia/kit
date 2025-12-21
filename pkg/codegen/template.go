package codegen

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateEngine handles template processing
type TemplateEngine struct {
	fs embed.FS
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(fs embed.FS) *TemplateEngine {
	return &TemplateEngine{fs: fs}
}

// Execute reads a template file, parses it, and writes the result to outputPath
func (te *TemplateEngine) Execute(templatePath, outputPath string, data any) error {
	// Read template from embedded FS
	tmplContent, err := te.fs.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplContent))
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

// ExecuteString executes a template string and returns the result
func ExecuteString(tmplStr string, data any) (string, error) {
	tmpl, err := template.New("inline").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}
