# Kit

> A production-ready Go toolkit providing code generators, reusable packages, and best-practice implementations for building scalable applications.

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ianmuhia/kit)](https://goreportcard.com/report/github.com/ianmuhia/kit)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](docs/contributing.md)

## Overview

**Kit** is a comprehensive Go toolkit designed to accelerate development while maintaining code quality and consistency:

- **🏗️ Code Generators**: CLI tools for scaffolding DDD domains, AuthZed permissions, error types, and APIs
- **📦 Reusable Packages**: Battle-tested utilities for common tasks (slices, strings, HTTP, messaging, testing)
- **🎯 Best Practices**: Implements clean architecture, CQRS, event-driven patterns, and functional options
- **⚡ Production-Ready**: Type-safe, well-tested, and designed for real-world applications

## Table of Contents

- [Quick Start](#quick-start)
- [Code Generators](#code-generators)
  - [DDD Generator](#-ddd-generator-ddd-gen)
  - [AuthZ Code Generator](#-authz-code-generator-authz-codegen)
  - [Error Generator](#-error-generator-error-gen)
- [Reusable Packages](#reusable-packages)
  - [Code Generation](#-code-generation-pkgcodegen)
  - [Slices](#-slices-pkgslices)
  - [Strings](#-strings-pkgstringutil)
  - [HTTP Utilities](#-http-utilities-pkghttputil)
  - [Messaging](#-messaging-pkgmessaging)
  - [Error Generation](#-error-generation-pkgerrorgen)
  - [AuthZ Generation](#-authz-generation-pkgauthzgen)
  - [Testing](#-testing-pkgtestutil)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Contributing](#-contributing)
- [License](#-license)

## Quick Start

### Installation

Install all generators:

```bash
go install github.com/ianmuhia/kit/cmd/...@latest
```

Or install specific tools:

```bash
# DDD domain generator
go install github.com/ianmuhia/kit/cmd/ddd-gen@latest

# AuthZed code generator
go install github.com/ianmuhia/kit/cmd/authz-codegen@latest

# Error type generator
go install github.com/ianmuhia/kit/cmd/error-gen@latest
```

### Use as Library

```bash
go get github.com/ianmuhia/kit@latest
```

```go
import (
    "github.com/ianmuhia/kit/pkg/slices"
    "github.com/ianmuhia/kit/pkg/stringutil"
    "github.com/ianmuhia/kit/pkg/httputil"
    "github.com/ianmuhia/kit/pkg/messaging"
)
```

---

## Code Generators

### 🏗️ DDD Generator (`ddd-gen`)

Generate complete Domain-Driven Design modules with hexagonal architecture, CQRS, and event sourcing.

**Installation:**

```bash
go install github.com/ianmuhia/kit/cmd/ddd-gen@latest
```

**Usage:**

```bash
# Basic domain generation
ddd-gen --domain=user --output=./internal

# With all features (CQRS, messaging, workflows, etc.)
ddd-gen --domain=order --all

# Selective features
ddd-gen --domain=booking --with-cqrs --with-messaging --with-tests
```

**Generated Structure:**

```
internal/user/
├── user.go              # Domain entity
├── repository.go        # Repository interface
├── errors.go            # Domain-specific errors
├── events.go            # Domain events
├── validation.go        # Validation logic
├── app/
│   └── service.go       # Application service
└── adapters/
    ├── user_http.go     # Huma v2 HTTP handlers
    └── user_postgres.go # PostgreSQL repository
```

**Features:**

| Feature | Flag | Description |
|---------|------|-------------|
| HTTP API | (default) | Type-safe REST API with Huma v2 and auto-generated OpenAPI |
| PostgreSQL | (default) | Repository implementation with pgx/v5 |
| CQRS | `--with-cqrs` | Command/Query separation with Watermill |
| Messaging | `--with-messaging` | Event pub/sub with Watermill + NATS |
| Workflows | `--with-workflows` | Temporal workflow integration |
| Job Queues | `--with-river` | Background jobs with River |
| Decorators | `--with-decorators` | Service decorators (auth, audit, cache, metrics) |
| Tests | `--with-tests` | Unit and integration test scaffolds |
| All Features | `--all` | Enable everything above |

**Key Technologies:**

- **HTTP**: [Huma v2](https://huma.rocks/) - Type-safe REST with auto-generated OpenAPI
- **Database**: [pgx/v5](https://github.com/jackc/pgx) - High-performance PostgreSQL driver
- **Messaging**: [Watermill](https://watermill.io/) + [NATS](https://nats.io/) - Event-driven architecture
- **Workflows**: [Temporal](https://temporal.io/) - Durable execution engine
- **Jobs**: [River](https://riverqueue.com/) - Fast and reliable background jobs

📖 **[Full Documentation →](docs/ddd-generator.md)**

---

### 🔐 AuthZ Code Generator (`authz-codegen`)

Generate type-safe Go client code from AuthZed schema definitions for SpiceDB/Zanzibar-style permissions.

**Installation:**

```bash
go install github.com/ianmuhia/kit/cmd/authz-codegen@latest
```

**Usage:**

```bash
# Generate from schema file
authz-codegen --schema=schema.zed --output=./authz

# Using positional arguments
authz-codegen schema.zed ./authz
```

**Input Schema Example:**

```zed
definition user {}

definition document {
    relation owner: user
    relation viewer: user
    
    permission view = viewer + owner
    permission edit = owner
}
```

**Generated Code:**

```go
// Type-safe resource types
doc := authz.NewDocument("doc-123")
user := authz.NewUser("user-456")

// Configure client (singleton pattern)
authz.SetClientConfig("localhost", "50051", "your-token")

// Create relations
err := doc.CreateOwnerRelations(ctx, authz.DocumentOwnerObjects{
    User: []authz.User{user},
})

// Check permissions
hasEdit, err := doc.CheckEdit(ctx, authz.CheckDocumentEditInputs{
    User: []authz.User{user},
})

// Lookup resources
docs, err := authz.LookupDocumentViewResourcesForUser(ctx, user)

// Read relations
owners, err := doc.ReadOwnerRelations(ctx)
```

**Features:**

- ✅ Type-safe API (no string-based resource types or relations)
- ✅ Singleton client pattern with functional options
- ✅ Auto-generated methods for all relations and permissions
- ✅ Support for union relations and computed permissions
- ✅ Context-aware operations with proper error handling
- ✅ Fully formatted, linted, and production-ready code

**Key Technologies:**

- **AuthZed SDK**: Official Go client for SpiceDB
- **Schema Parsing**: Custom lexer/parser for `.zed` files
- **Code Generation**: Template-based with `text/template`

📖 **[AuthZed Documentation](https://authzed.com/docs)**

---

### ⚠️ Error Generator (`error-gen`)

Generate strongly-typed error types from CUE definitions with consistent codes, messages, and HTTP status mappings.

**Installation:**

```bash
go install github.com/ianmuhia/kit/cmd/error-gen@latest
```

**Usage:**

```bash
# Generate from CUE file
error-gen --input=errors.cue --output=errors.gen.go --package=myapp

# Using defaults (errors.cue -> errors.go)
error-gen
```

**Input CUE Schema:**

```cue
package errors

errors: [
    {
        code: "USER_NOT_FOUND"
        message: "User not found"
        httpStatus: 404
    },
    {
        code: "INVALID_INPUT"
        message: "Invalid input provided"
        httpStatus: 400
    },
    {
        code: "INTERNAL_ERROR"
        message: "Internal server error"
        httpStatus: 500
    },
]
```

**Generated Code:**

```go
// Strongly-typed error codes
const (
    ErrCodeUserNotFound  ErrorCode = "USER_NOT_FOUND"
    ErrCodeInvalidInput  ErrorCode = "INVALID_INPUT"
    ErrCodeInternalError ErrorCode = "INTERNAL_ERROR"
)

// Constructor functions
func NewUserNotFoundError(details ...string) *AppError {
    return &AppError{
        Code:       ErrCodeUserNotFound,
        Message:    "User not found",
        HTTPStatus: 404,
        Details:    details,
    }
}

// HTTP status mapping
func (e *AppError) HTTPStatusCode() int {
    return e.HTTPStatus
}

// JSON serialization
func (e *AppError) MarshalJSON() ([]byte, error) {
    // ...
}
```

**Features:**

- ✅ CUE-based schema validation
- ✅ Type-safe error codes (no magic strings)
- ✅ Automatic HTTP status code mapping
- ✅ Constructor functions for each error type
- ✅ JSON serialization support
- ✅ Integration with `pkg/httputil` for consistent API responses
- ✅ Supports custom templates

**Key Technologies:**

- **CUE**: Configuration language with strong typing
- **Code Generation**: Template-based with functional options

---

## Reusable Packages

### 📝 Code Generation (`pkg/codegen`)

Utilities for template processing, file operations, and naming conventions used by code generators.

```go
import "github.com/ianmuhia/kit/pkg/codegen"

// Naming conventions
codegen.Capitalize("hello")        // "Hello"
codegen.Uncapitalize("World")      // "world"
codegen.ToPascalCase("user_name")  // "UserName"
codegen.ToCamelCase("user_name")   // "userName"

// File operations
codegen.EnsureDir("/path/to/dir")
codegen.WriteFile("/path/file.txt", []byte("content"))
codegen.FileExists("config.yaml")  // bool

// Template utilities
codegen.ProcessTemplate(templateStr, data)
```

---

### 🔀 Slices (`pkg/slices`)

Type-safe generic slice operations (requires Go 1.18+).

```go
import "github.com/ianmuhia/kit/pkg/slices"

numbers := []int{1, 2, 3, 4, 5}

// Functional operations
doubled := slices.Map(numbers, func(n int) int { return n * 2 })
// [2, 4, 6, 8, 10]

evens := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
// [2, 4]

sum := slices.Reduce(numbers, 0, func(acc, n int) int { return acc + n })
// 15

// Utilities
slices.Contains(numbers, 3)            // true
slices.Unique([]int{1, 1, 2, 3})      // [1, 2, 3]
slices.Chunk(numbers, 2)               // [[1, 2], [3, 4], [5]]
slices.Reverse(numbers)                // [5, 4, 3, 2, 1]
slices.Flatten([][]int{{1,2},{3,4}})  // [1, 2, 3, 4]

// Set operations
slices.Union([]int{1,2}, []int{2,3})       // [1, 2, 3]
slices.Intersection([]int{1,2}, []int{2,3}) // [2]
slices.Difference([]int{1,2,3}, []int{2})   // [1, 3]
```

**Features:**

- Full generic support (works with any type)
- Functional programming patterns (map, filter, reduce)
- Common utilities (contains, unique, chunk)
- Set operations (union, intersection, difference)
- Zero allocations for most operations

---

### 📋 Strings (`pkg/stringutil`)

Comprehensive string manipulation and validation utilities.

```go
import "github.com/ianmuhia/kit/pkg/stringutil"

// Case conversions
stringutil.ToPascalCase("hello_world")   // "HelloWorld"
stringutil.ToCamelCase("hello_world")    // "helloWorld"
stringutil.ToSnakeCase("HelloWorld")     // "hello_world"
stringutil.ToKebabCase("HelloWorld")     // "hello-world"
stringutil.ToTitleCase("hello world")    // "Hello World"

// Validation
stringutil.IsEmail("test@example.com")   // true
stringutil.IsURL("https://example.com")  // true
stringutil.IsEmpty("   ")                // true
stringutil.IsAlpha("Hello")              // true
stringutil.IsNumeric("12345")            // true
stringutil.IsAlphanumeric("Hello123")    // true

// Utilities
stringutil.Truncate("long text here", 10, "...")  // "long te..."
stringutil.Reverse("hello")                        // "olleh"
stringutil.Slugify("Hello World!")                 // "hello-world"
stringutil.RandomString(16)                        // "aB3dE5fG7hI9jK1l"
stringutil.Pad("Go", 5, "*", "both")              // "*Go**"
```

**Features:**

- Multiple case conversion formats
- Comprehensive validation (email, URL, alpha, numeric)
- String manipulation (truncate, reverse, slugify)
- Random string generation
- Padding and alignment

---

### 🌐 HTTP Utilities (`pkg/httputil`)

Production-ready HTTP helpers for responses, errors, and middleware.

```go
import "github.com/ianmuhia/kit/pkg/httputil"

func Handler(w http.ResponseWriter, r *http.Request) {
    // Success responses
    httputil.Success(w, map[string]string{"status": "ok"})
    httputil.Created(w, user)
    httputil.NoContent(w)
    
    // Error responses
    httputil.BadRequest(w, "Invalid input")
    httputil.NotFound(w, "Resource not found")
    httputil.Unauthorized(w, "Invalid credentials")
    httputil.Forbidden(w, "Access denied")
    httputil.InternalServerError(w, "Something went wrong")
    
    // Custom errors with proper HTTP status
    err := httputil.NewHTTPError(422, "VALIDATION_ERROR", "Invalid fields")
    httputil.WriteError(w, err)
}

// Middleware
handler := httputil.LoggingMiddleware(myHandler)
handler = httputil.RecoveryMiddleware(handler)
handler = httputil.CORSMiddleware([]string{"https://example.com"})(handler)
handler = httputil.RateLimitMiddleware(100, time.Minute)(handler)

// Request utilities
body, err := httputil.ReadJSON(r, &user)
query := httputil.GetQueryParam(r, "page", "1")
```

**Features:**

- Consistent JSON response format
- Standard HTTP error helpers
- Production-ready middleware (CORS, rate limiting, logging, recovery)
- Request parsing utilities
- Integration with standard `net/http`

---

### 📨 Messaging (`pkg/messaging`)

High-level NATS publisher and subscriber with Watermill integration, using functional options pattern.

```go
import "github.com/ianmuhia/kit/pkg/messaging"

// Create publisher
publisher, err := messaging.NewPublisher(
    messaging.WithURL("nats://localhost:4222"),
    messaging.WithLogger(logger),
    messaging.WithMarshaler(myMarshaler),
)

// Publish message
err = publisher.Publish(ctx, "user.created", event)

// Create subscriber
subscriber, err := messaging.NewSubscriber(
    messaging.WithSubscriberURL("nats://localhost:4222"),
    messaging.WithDurablePrefix("myapp"),
    messaging.WithLogger(logger),
)

// Subscribe to topic
messages, err := subscriber.Subscribe(ctx, "user.created")
for msg := range messages {
    // Process message
    msg.Ack()
}
```

**Features:**

- Functional options pattern for clean configuration
- NATS JetStream integration for durability
- Automatic reconnection and error handling
- Message acknowledgment support
- Structured logging
- Custom marshalers/unmarshalers

**Key Technologies:**

- **Watermill**: Event-driven messaging library
- **NATS**: Cloud-native messaging system

---

### ⚠️ Error Generation (`pkg/errorgen`)

Library for generating typed errors from CUE definitions (used by `error-gen` CLI).

```go
import "github.com/ianmuhia/kit/pkg/errorgen"

// Create generator
gen, err := errorgen.NewGenerator(
    errorgen.WithInputFile("errors.cue"),
    errorgen.WithOutputFile("errors.gen.go"),
    errorgen.WithPackageName("myapp"),
)

// Generate code
err = gen.Generate()
```

**Features:**

- Functional options pattern
- CUE schema validation
- Template-based code generation
- Structured logging

---

### 🔐 AuthZ Generation (`pkg/authzgen`)

Library for generating AuthZed client code from schema files (used by `authz-codegen` CLI).

```go
import "github.com/ianmuhia/kit/pkg/authzgen"

// Create generator
gen, err := authzgen.NewGenerator(
    authzgen.WithSchemaFile("schema.zed"),
    authzgen.WithOutputDir("./authz"),
    authzgen.WithLogger(logger),
)

// Generate code
err = gen.Generate()
```

**Features:**

- Full lexer/parser for AuthZed schema language
- AST-based code generation
- Type-safe API generation
- Functional options pattern

---

### 🧪 Testing (`pkg/testutil`)

Comprehensive testing utilities and assertions for unit and integration tests.

```go
import "github.com/ianmuhia/kit/pkg/testutil"

func TestMyHandler(t *testing.T) {
    // HTTP testing
    rec := testutil.DoHTTPRequest(t, handler, "GET", "/users", nil)
    testutil.AssertStatus(t, rec, 200)
    testutil.AssertJSON(t, rec, expectedResponse)
    testutil.AssertHeader(t, rec, "Content-Type", "application/json")
    
    // JSON assertions
    var user User
    testutil.AssertJSONMatch(t, rec.Body.Bytes(), &user)
    testutil.AssertEqual(t, user.Name, "John Doe")
    
    // Fixtures
    var testUser User
    testutil.LoadFixture(t, "testdata/user.json", &testUser)
    
    // Temporary files (auto-cleanup)
    dir := testutil.TempDir(t)
    file := testutil.TempFile(t, "*.txt", "test content")
    
    // Database testing
    db := testutil.SetupTestDB(t, "postgres://...")
    testutil.TeardownTestDB(t, db)
}
```

**Features:**

- HTTP request/response testing
- JSON comparison and assertions
- Fixture loading from files
- Automatic temp file cleanup
- Database test helpers
- Custom assertions

---

## 🏗️ Project Structure

```
kit/
├── cmd/                       # CLI tools and generators
│   ├── ddd-gen/              # DDD domain generator
│   ├── authz-codegen/        # AuthZed code generator
│   ├── error-gen/            # Error type generator
│   ├── api-gen/              # API generator (future)
│   └── db-migrate/           # Database migration tool (future)
│
├── pkg/                       # Public, reusable packages
│   ├── codegen/              # Code generation utilities
│   ├── slices/               # Generic slice operations
│   ├── stringutil/           # String manipulation & validation
│   ├── httputil/             # HTTP helpers & middleware
│   ├── messaging/            # NATS pub/sub with Watermill
│   ├── errorgen/             # Error generation library
│   ├── authzgen/             # AuthZed generation library
│   └── testutil/             # Testing utilities
│
├── internal/                  # Private packages
│   ├── dddgen/               # DDD generator implementation
│   │   ├── config.go
│   │   ├── generator.go
│   │   └── templates/        # Embedded templates
│   └── shared/               # Shared internal utilities
│
├── templates/                 # Template files
│   ├── domain/               # Domain layer templates
│   ├── app/                  # Application layer templates
│   └── adapters/             # Infrastructure templates
│
├── docs/                      # Documentation
│   ├── ddd-generator.md      # DDD generator guide
│   ├── huma-adapter.md       # Huma integration guide
│   └── contributing.md       # Contributing guidelines
│
├── examples/                  # Example usage
│   ├── basic/
│   └── advanced/
│
├── tools/                     # Tool dependencies
├── .github/                   # GitHub workflows
├── Makefile                   # Build automation
├── go.mod                     # Go module definition
└── README.md                  # This file
```

---

## 🛠️ Development

### Prerequisites

- Go 1.25 or higher
- Make (for build automation)
- Docker (optional, for testing with PostgreSQL/NATS)

### Clone and Build

```bash
# Clone repository
git clone https://github.com/ianmuhia/kit.git
cd kit

# Install dependencies
go mod download

# Build all tools
make build

# Build specific tool
make build-ddd-gen
make build-authz-codegen
make build-error-gen

# Install to GOPATH/bin
make install
make install-ddd-gen
```

### Available Make Targets

```bash
make help                 # Show all available targets
make build                # Build all binaries
make build-<tool>         # Build specific tool
make install              # Install all tools to GOPATH/bin
make install-<tool>       # Install specific tool
make test                 # Run all tests
make test-coverage        # Run tests with coverage report
make lint                 # Run linters (golangci-lint)
make fmt                  # Format code with gofmt
make clean                # Remove build artifacts
make docs                 # Generate documentation
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test -v ./pkg/slices/...
go test -v ./internal/dddgen/...

# Run integration tests
go test -v -tags=integration ./...
```

### Code Quality

```bash
# Format code
make fmt
gofmt -w .

# Run linters
make lint
golangci-lint run

# Static analysis
go vet ./...
staticcheck ./...
```

---

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [DDD Generator Guide](docs/ddd-generator.md) | Complete guide to using the DDD generator |
| [Huma HTTP Adapter](docs/huma-adapter.md) | Huma v2 integration and best practices |
| [Contributing Guide](docs/contributing.md) | How to contribute to the project |
| [Migration Guide](MIGRATION.md) | Migrating from ddd-lite to kit |

---

## 🤝 Contributing

Contributions are welcome! We follow standard GitHub workflow and code quality standards.

### How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes following our coding standards
4. **Write/update** tests for your changes
5. **Run** tests and linters (`make test lint`)
6. **Commit** with clear, descriptive messages (`git commit -m 'Add amazing feature'`)
7. **Push** to your branch (`git push origin feature/amazing-feature`)
8. **Open** a Pull Request with detailed description

### Coding Standards

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Write comprehensive tests (aim for >80% coverage)
- Document exported functions, types, and packages
- Use functional options pattern for configuration
- Handle errors explicitly (no silent failures)

### Pull Request Guidelines

- **Title**: Clear, concise summary (e.g., "Add rate limiting to HTTP middleware")
- **Description**: What, why, and how of your changes
- **Tests**: Include tests for new functionality
- **Documentation**: Update relevant docs
- **Breaking Changes**: Clearly mark and explain in PR description

### Reporting Issues

When reporting bugs, please include:

- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

📖 **Full guide: [CONTRIBUTING.md](docs/contributing.md)**

---

## 🗺️ Roadmap

### Current (v1.0)

- ✅ DDD domain generator with Huma v2
- ✅ AuthZed code generator
- ✅ Error type generator
- ✅ Core utility packages (slices, strings, HTTP, messaging, testing)
- ✅ Functional options pattern
- ✅ Comprehensive documentation

### Planned (v1.1)

- [ ] API generator (REST/gRPC scaffolding)
- [ ] Database migration tool (schema versioning)
- [ ] GraphQL generator
- [ ] Event sourcing templates
- [ ] Saga pattern support
- [ ] OpenTelemetry integration

### Future (v2.0)

- [ ] CLI framework for building custom generators
- [ ] Plugin system for extending generators
- [ ] Web UI for generator configuration
- [ ] Project templates (monorepo, microservices)
- [ ] CI/CD pipeline templates

---

## 📊 Benchmarks

### Slice Operations

```
BenchmarkMap-8           5000000    245 ns/op      0 B/op    0 allocs/op
BenchmarkFilter-8        3000000    412 ns/op      0 B/op    0 allocs/op
BenchmarkReduce-8       10000000    156 ns/op      0 B/op    0 allocs/op
```

### String Operations

```
BenchmarkToPascalCase-8  2000000    658 ns/op    128 B/op    3 allocs/op
BenchmarkToSnakeCase-8   3000000    523 ns/op     96 B/op    2 allocs/op
BenchmarkSlugify-8       1000000   1024 ns/op    256 B/op    5 allocs/op
```

### HTTP Utilities

```
BenchmarkSuccessResponse-8     500000   2345 ns/op   512 B/op   8 allocs/op
BenchmarkErrorResponse-8       500000   2567 ns/op   544 B/op   9 allocs/op
BenchmarkMiddlewareChain-8    1000000   1234 ns/op   256 B/op   4 allocs/op
```

---

## 🙏 Acknowledgments

This toolkit builds on the shoulders of giants:

- **[Huma](https://huma.rocks/)** - Type-safe REST API framework
- **[Watermill](https://watermill.io/)** - Event-driven architecture
- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL driver
- **[AuthZed](https://authzed.com/)** - Permissions infrastructure
- **[CUE](https://cuelang.org/)** - Configuration language
- **[Temporal](https://temporal.io/)** - Workflow engine
- **[River](https://riverqueue.com/)** - Background jobs

Special thanks to all [contributors](https://github.com/ianmuhia/kit/graphs/contributors)!

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

**TL;DR**: You can use this toolkit freely in personal and commercial projects. Attribution appreciated but not required.

---

## 👤 Author

**Ian Muhia**

- GitHub: [@Ianmuhia](https://github.com/Ianmuhia)
- Twitter: [@ianmuhia](https://twitter.com/ianmuhia)
- Email: <ian@example.com>

---

## ⭐ Show Your Support

If this project helps you build better Go applications:

- Give it a ⭐️ on [GitHub](https://github.com/ianmuhia/kit)
- Share it with your team
- [Sponsor the project](https://github.com/sponsors/Ianmuhia)
- Contribute back with PRs
- Report issues and suggest features

---

## 📈 Stats

![GitHub stars](https://img.shields.io/github/stars/Ianmuhia/kit?style=social)
![GitHub forks](https://img.shields.io/github/forks/Ianmuhia/kit?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/Ianmuhia/kit?style=social)
![GitHub contributors](https://img.shields.io/github/contributors/Ianmuhia/kit)
![GitHub issues](https://img.shields.io/github/issues/Ianmuhia/kit)
![GitHub pull requests](https://img.shields.io/github/issues-pr/Ianmuhia/kit)

---

<div align="center">

**Built with ❤️ using Go**

[⬆ Back to Top](#kit)

</div>
