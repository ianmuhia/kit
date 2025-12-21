package authzgen

import (
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

// Generator handles AuthZed schema code generation
type Generator struct {
	schemaFile string
	outputDir  string
	logger     *slog.Logger
}

// Option is a functional option for configuring the Generator
type Option func(*Generator)

// WithSchemaFile sets the schema file path
func WithSchemaFile(path string) Option {
	return func(g *Generator) {
		g.schemaFile = path
	}
}

// WithOutputDir sets the output directory
func WithOutputDir(dir string) Option {
	return func(g *Generator) {
		g.outputDir = dir
	}
}

// WithLogger sets the logger
func WithLogger(logger *slog.Logger) Option {
	return func(g *Generator) {
		g.logger = logger
	}
}

// NewGenerator creates a new AuthZed code generator with the given options
func NewGenerator(opts ...Option) (*Generator, error) {
	g := &Generator{
		outputDir: ".",
		logger:    slog.Default(),
	}

	for _, opt := range opts {
		opt(g)
	}

	if g.schemaFile == "" {
		return nil, fmt.Errorf("schema file is required")
	}

	return g, nil
}

// Generate parses the schema and generates the code
func (g *Generator) Generate() error {
	g.logger.Info("Starting schema parsing", "file", g.schemaFile)

	schema, err := g.parseSchema()
	if err != nil {
		g.logger.Error("Schema parsing failed", "file", g.schemaFile, "error", err)
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Use a single package name for all definitions
	packageName := "authz"
	if len(schema.Definitions) > 0 {
		if schema.Definitions[0].Package != "" {
			packageName = schema.Definitions[0].Package
		}
	}

	g.logger.Info("Generating code for single package", "package", packageName, "definitions_count", len(schema.Definitions))

	if err := g.generateCode(packageName, schema.Definitions); err != nil {
		g.logger.Error("Code generation failed", "package", packageName, "output_dir", g.outputDir, "error", err)
		return fmt.Errorf("failed to generate code for package %s: %w", packageName, err)
	}

	g.logger.Info("Code generation completed successfully", "package", packageName, "output", filepath.Join(g.outputDir, packageName+".gen.go"))
	return nil
}

func (g *Generator) parseSchema() (*Schema, error) {
	content, err := os.ReadFile(g.schemaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	g.logger.Debug("Schema file read successfully", "size_bytes", len(content))

	lexer := NewLexer(string(content))
	tokens := lexer.TokenizeAll()

	g.logger.Debug("Lexical analysis complete", "token_count", len(tokens))

	parser := NewParser(tokens)
	astDefinitions, err := parser.ParseDefinitions()
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	g.logger.Info("AST parsing successful", "definitions_count", len(astDefinitions))

	// Convert AST definitions to our internal format
	var schema Schema
	for _, astDef := range astDefinitions {
		pkg := astDef.ObjectType.Prefix
		if pkg == "" {
			pkg = "authz"
		}

		def := Definition{
			Package: pkg,
			Name:    astDef.ObjectType.Name,
		}

		// Convert relations
		for _, astRel := range astDef.Relations {
			relation := Relation{
				Name:  astRel.Name,
				Types: extractTypesFromRelationExpression(astRel.Expression),
			}
			if len(relation.Types) > 1 {
				relation.IsUnion = true
			}
			def.Relations = append(def.Relations, relation)
		}

		// Convert permissions
		for _, astPerm := range astDef.Permissions {
			permission := Permission{
				Name:       astPerm.Name,
				Expression: astPerm.Expression.String(),
			}
			def.Permissions = append(def.Permissions, permission)
		}

		schema.Definitions = append(schema.Definitions, def)
	}

	return &schema, nil
}

func (g *Generator) generateCode(packageName string, definitions []Definition) error {
	if err := os.MkdirAll(g.outputDir, 0o755); err != nil {
		return err
	}

	tmpl := template.New("code").Funcs(template.FuncMap{
		"camelcase": ToPascalCase,
		"lower":     strings.ToLower,
		"extractType": func(fullType string) string {
			parts := strings.Split(fullType, "/")
			typeName := fullType
			if len(parts) > 1 {
				typeName = parts[1]
			}
			if idx := strings.Index(typeName, "#"); idx != -1 {
				typeName = typeName[:idx]
			}
			return typeName
		},
	})

	tmpl, err := tmpl.Parse(codeTemplate)
	if err != nil {
		return err
	}

	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].Name < definitions[j].Name
	})

	data := struct {
		Package     string
		Definitions []Definition
	}{
		Package:     packageName,
		Definitions: definitions,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		formatted = []byte(buf.String())
	}

	filename := filepath.Join(g.outputDir, packageName+".gen.go")
	return os.WriteFile(filename, formatted, 0o644)
}

func extractTypesFromRelationExpression(expr RelationExpressionNode) []string {
	var types []string

	switch node := expr.(type) {
	case *SingleRelationNode:
		if node.Fragment != "" {
			types = append(types, fmt.Sprintf("%s#%s", node.Value, node.Fragment))
		} else {
			types = append(types, node.Value)
		}
	case *UnionRelationNode:
		types = append(types, extractTypesFromRelationExpression(node.Left)...)
		types = append(types, extractTypesFromRelationExpression(node.Right)...)
	}

	return types
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	var result strings.Builder
	s = strings.TrimSpace(s)

	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})

	for _, word := range words {
		if word == "" {
			continue
		}
		runes := []rune(word)
		for i, r := range runes {
			if i == 0 {
				result.WriteRune(unicode.ToUpper(r))
			} else {
				result.WriteRune(unicode.ToLower(r))
			}
		}
	}

	return result.String()
}

// Schema represents the parsed AuthZed schema
type Schema struct {
	Definitions []Definition
}

// Definition represents a definition in the schema
type Definition struct {
	Name        string
	Package     string
	Relations   []Relation
	Permissions []Permission
}

// Relation represents a relation in a definition
type Relation struct {
	Name    string
	Types   []string
	IsUnion bool
}

// Permission represents a permission in a definition
type Permission struct {
	Name       string
	Expression string
}
