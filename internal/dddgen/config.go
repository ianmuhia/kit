package dddgen

// Config holds the configuration for domain generation
type Config struct {
	DomainName     string
	OutputDir      string
	ModulePath     string // The Go module path (e.g., "github.com/user/project" or "ibnb")
	WithTests      bool
	WithMessaging  bool
	WithRiver      bool
	WithCQRS       bool
	WithWorkflows  bool
	WithDecorators bool
}

// TemplateData holds data passed to templates
type TemplateData struct {
	DomainTitle string // Capitalized for type names
	DomainLower string // Lowercase for package/file names
	ModulePath  string // The Go module path for imports
}
