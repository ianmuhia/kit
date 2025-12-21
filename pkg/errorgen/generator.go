package errorgen

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

//go:embed templates/*.tmpl
var Templates embed.FS

// ErrorDefinition represents a single error definition.
type ErrorDefinition struct {
	Name        string
	Code        string
	Message     string
	Category    string
	HTTPStatus  int
	Severity    string
	Description string
	Parameters  []string
}

// ErrorConfig holds all error definitions.
type ErrorConfig struct {
	Package string
	Errors  []ErrorDefinition
}

// GeneratorConfig holds configuration for the error generator.
type GeneratorConfig struct {
	inputFile    string
	outputFile   string
	templateFile string
	packageName  string
}

// GeneratorOption is a functional option for configuring the generator.
type GeneratorOption func(*GeneratorConfig)

// WithInputFile sets the input CUE file or directory.
func WithInputFile(path string) GeneratorOption {
	return func(c *GeneratorConfig) {
		c.inputFile = path
	}
}

// WithOutputFile sets the output Go file path.
func WithOutputFile(path string) GeneratorOption {
	return func(c *GeneratorConfig) {
		c.outputFile = path
	}
}

// WithTemplateFile sets a custom template file.
func WithTemplateFile(path string) GeneratorOption {
	return func(c *GeneratorConfig) {
		c.templateFile = path
	}
}

// WithPackageName overrides the package name.
func WithPackageName(name string) GeneratorOption {
	return func(c *GeneratorConfig) {
		c.packageName = name
	}
}

// defaultGeneratorConfig returns sensible defaults.
func defaultGeneratorConfig() *GeneratorConfig {
	return &GeneratorConfig{
		inputFile:  "errors.cue",
		outputFile: "errors.go",
	}
}

// Generator handles error code generation.
type Generator struct {
	config *GeneratorConfig
}

// NewGenerator creates a new error generator.
func NewGenerator(opts ...GeneratorOption) (*Generator, error) {
	config := defaultGeneratorConfig()

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	if config.inputFile == "" {
		return nil, fmt.Errorf("input file is required")
	}

	if config.outputFile == "" {
		return nil, fmt.Errorf("output file is required")
	}

	return &Generator{config: config}, nil
}

// Generate generates error code from CUE definitions.
func (g *Generator) Generate() error {
	// Load CUE configuration
	errorConfig, err := g.loadCUEConfig()
	if err != nil {
		return fmt.Errorf("failed to load CUE config: %w", err)
	}

	// Override package name if specified
	if g.config.packageName != "" {
		errorConfig.Package = g.config.packageName
	}

	// Validate config
	if err := errorConfig.validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Generate code from template
	if err := g.generateCode(errorConfig); err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}

// loadCUEConfig loads error definitions from a CUE file.
func (g *Generator) loadCUEConfig() (*ErrorConfig, error) {
	inputPath := g.config.inputFile
	if !filepath.IsAbs(inputPath) {
		wd, _ := os.Getwd()
		inputPath = filepath.Join(wd, inputPath)
	}

	// Create CUE context
	ctx := cuecontext.New()

	// Determine if input is a file or directory
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat input path: %w", err)
	}

	var value cue.Value

	if fileInfo.IsDir() {
		// Load as a package directory
		buildInstances := load.Instances([]string{inputPath}, nil)
		if len(buildInstances) == 0 {
			return nil, fmt.Errorf("no CUE instances found in %s", inputPath)
		}
		if buildInstances[0].Err != nil {
			return nil, fmt.Errorf("failed to load CUE package: %w", buildInstances[0].Err)
		}
		value = ctx.BuildInstance(buildInstances[0])
	} else {
		// Load as a single file
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CUE file: %w", err)
		}
		value = ctx.CompileBytes(data, cue.Filename(inputPath))
	}

	if err := value.Err(); err != nil {
		return nil, fmt.Errorf("CUE compilation error: %w", err)
	}

	// Build config by extracting concrete values
	config := &ErrorConfig{}

	// Get package name
	packageValue := value.LookupPath(cue.ParsePath("package"))
	if packageValue.Exists() {
		pkgStr, _ := packageValue.String()
		config.Package = pkgStr
	}
	if config.Package == "" {
		config.Package = "errors" // default
	}

	// Get errors array
	errorsValue := value.LookupPath(cue.ParsePath("errors"))
	if !errorsValue.Exists() {
		return nil, fmt.Errorf("errors field not found in CUE file")
	}

	// Iterate through errors array
	iter, err := errorsValue.List()
	if err != nil {
		return nil, fmt.Errorf("errors must be a list: %w", err)
	}

	for iter.Next() {
		errorDef := ErrorDefinition{}
		errVal := iter.Value()

		// Extract each field
		if code := errVal.LookupPath(cue.ParsePath("code")); code.Exists() {
			errorDef.Code, _ = code.String()
		}
		if name := errVal.LookupPath(cue.ParsePath("name")); name.Exists() {
			errorDef.Name, _ = name.String()
		}
		if message := errVal.LookupPath(cue.ParsePath("message")); message.Exists() {
			errorDef.Message, _ = message.String()
		}
		if category := errVal.LookupPath(cue.ParsePath("category")); category.Exists() {
			errorDef.Category, _ = category.String()
		}
		if httpStatus := errVal.LookupPath(cue.ParsePath("httpStatus")); httpStatus.Exists() {
			if status, err := httpStatus.Int64(); err == nil {
				errorDef.HTTPStatus = int(status)
			}
		}
		if severity := errVal.LookupPath(cue.ParsePath("severity")); severity.Exists() {
			errorDef.Severity, _ = severity.String()
		}
		if description := errVal.LookupPath(cue.ParsePath("description")); description.Exists() {
			errorDef.Description, _ = description.String()
		}
		if parameters := errVal.LookupPath(cue.ParsePath("parameters")); parameters.Exists() {
			paramIter, _ := parameters.List()
			for paramIter.Next() {
				if param, err := paramIter.Value().String(); err == nil {
					errorDef.Parameters = append(errorDef.Parameters, param)
				}
			}
		}

		config.Errors = append(config.Errors, errorDef)
	}

	return config, nil
}

// generateCode generates Go code from the error config.
func (g *Generator) generateCode(config *ErrorConfig) error {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"default": func(def any, val any) any {
			if val == nil || val == "" || val == 0 {
				return def
			}
			return val
		},
		"codeConstName": func(name string) string {
			return "Code" + strings.TrimPrefix(name, "Err")
		},
		"paramName": func(param string) string {
			return strings.ToLower(param)
		},
		"sanitizeName": func(name string) string {
			return strings.ReplaceAll(strings.ReplaceAll(name, " ", "_"), "-", "_")
		},
		"getUniqueCategories": func(errors []ErrorDefinition) []string {
			seen := make(map[string]bool)
			var categories []string
			for _, e := range errors {
				if e.Category != "" && !seen[e.Category] {
					categories = append(categories, e.Category)
					seen[e.Category] = true
				}
			}
			return categories
		},
	}

	// Parse template
	var tmpl *template.Template
	var err error

	if g.config.templateFile != "" {
		// Use custom template file
		tmpl, err = template.New(filepath.Base(g.config.templateFile)).Funcs(funcMap).ParseFiles(g.config.templateFile)
	} else {
		// Use embedded template
		tmplContent, readErr := Templates.ReadFile("templates/error.go.tmpl")
		if readErr != nil {
			return fmt.Errorf("failed to read embedded template: %w", readErr)
		}
		tmpl, err = template.New("error.go.tmpl").Funcs(funcMap).Parse(string(tmplContent))
	}

	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	outputPath := g.config.outputFile
	if !filepath.IsAbs(outputPath) {
		wd, _ := os.Getwd()
		outputPath = filepath.Join(wd, outputPath)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Execute template
	if err := tmpl.Execute(outFile, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// validate ensures the error config is valid.
func (c *ErrorConfig) validate() error {
	if c.Package == "" {
		return fmt.Errorf("package name is required")
	}

	seenCodes := make(map[string]bool)
	seenNames := make(map[string]bool)

	for _, e := range c.Errors {
		if e.Name == "" || e.Code == "" || e.Message == "" {
			return fmt.Errorf("error definition missing required fields: name=%s, code=%s, message=%s",
				e.Name, e.Code, e.Message)
		}

		if seenCodes[e.Code] {
			return fmt.Errorf("duplicate error code: %s", e.Code)
		}
		if seenNames[e.Name] {
			return fmt.Errorf("duplicate error name: %s", e.Name)
		}

		if e.Severity != "" && !isValidSeverity(e.Severity) {
			return fmt.Errorf("invalid severity %s for error %s; must be one of: critical, high, medium, low",
				e.Severity, e.Name)
		}

		if len(e.Parameters) > 0 {
			for _, param := range e.Parameters {
				if !strings.Contains(e.Message, "{"+param+"}") {
					return fmt.Errorf("parameter %s in error %s not found in message", param, e.Name)
				}
			}
		}

		seenCodes[e.Code] = true
		seenNames[e.Name] = true
	}

	return nil
}

func isValidSeverity(severity string) bool {
	validSeverities := map[string]bool{
		"critical": true,
		"high":     true,
		"medium":   true,
		"low":      true,
	}
	return validSeverities[strings.ToLower(severity)]
}
