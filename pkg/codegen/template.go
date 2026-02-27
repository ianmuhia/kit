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
	fs      embed.FS
	funcMap template.FuncMap
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(fs embed.FS) *TemplateEngine {
	return &TemplateEngine{fs: fs}
}

// WithFuncMap registers additional template functions available in all templates
// executed by this engine. It returns the engine for chaining.
func (te *TemplateEngine) WithFuncMap(fm template.FuncMap) *TemplateEngine {
	if te.funcMap == nil {
		te.funcMap = make(template.FuncMap)
	}
	for k, v := range fm {
		te.funcMap[k] = v
	}
	return te
}

// Execute reads a template file, parses it, and writes the result to outputPath.
// Any FuncMap registered via WithFuncMap is available inside the template.
func (te *TemplateEngine) Execute(templatePath, outputPath string, data any) error {
	tmplContent, err := te.fs.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl := template.New(filepath.Base(templatePath))
	if len(te.funcMap) > 0 {
		tmpl = tmpl.Funcs(te.funcMap)
	}

	tmpl, err = tmpl.Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
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
