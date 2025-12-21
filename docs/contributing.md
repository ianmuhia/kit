# Contributing to Kit

Thank you for your interest in contributing to Kit! This document provides guidelines and instructions for contributing.

## 🌟 Ways to Contribute

- **Report Bugs**: Submit detailed bug reports via GitHub Issues
- **Suggest Features**: Propose new generators, utilities, or improvements
- **Submit Code**: Fix bugs, add features, or improve documentation
- **Improve Documentation**: Help make the docs clearer and more comprehensive
- **Share Feedback**: Tell us what works and what doesn't

## 🚀 Getting Started

### Prerequisites

- Go 1.25 or higher
- Git
- Make (optional, but recommended)
- golangci-lint (for linting)

### Setup Development Environment

1. **Fork the repository** on GitHub

2. **Clone your fork**:

```bash
git clone https://github.com/YOUR_USERNAME/kit.git
cd kit
```

1. **Add upstream remote**:

```bash
git remote add upstream https://github.com/ianmuhia/kit.git
```

1. **Install dependencies**:

```bash
go mod download
```

1. **Verify setup**:

```bash
make test
make build
```

## 📝 Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

Branch naming conventions:

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions/improvements

### 2. Make Your Changes

- Write clean, idiomatic Go code
- Follow existing code style and patterns
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run specific package tests
go test ./pkg/slices/...

# Run with coverage
make test-coverage

# Lint your code
make lint

# Format code
make fmt
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add Map function to slices package"
```

Commit message format:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:

- Clear title describing the change
- Detailed description of what and why
- Reference any related issues
- Screenshots/examples if applicable

## 🏗️ Project Structure

```
kit/
├── cmd/              # Generators and CLI tools
├── pkg/              # Public packages (importable)
├── internal/         # Private packages
├── docs/             # Documentation
├── tools/            # Build tools
└── Makefile          # Build automation
```

### Adding a New Package

1. Create directory in `pkg/`:

```bash
mkdir -p pkg/newpackage
```

1. Add package files with documentation:

```go
// Package newpackage provides utilities for...
package newpackage

// Function does something useful
func Function() {
    // implementation
}
```

1. Add tests:

```go
package newpackage

import "testing"

func TestFunction(t *testing.T) {
    // test implementation
}
```

1. Update main README with package documentation

### Adding a New Generator

1. Create directory in `cmd/`:

```bash
mkdir -p cmd/new-gen
```

1. Create `main.go`:

```go
package main

import (
    "github.com/ianmuhia/kit/internal/newgen"
)

func main() {
    // CLI setup
}
```

1. Create implementation in `internal/`:

```bash
mkdir -p internal/newgen
```

1. Add documentation in `docs/`

## ✅ Code Quality Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting (run `make fmt`)
- Use meaningful variable names
- Keep functions small and focused
- Add comments for exported functions

### Testing Requirements

- Unit tests for all new functionality
- Table-driven tests for multiple scenarios
- Test error cases
- Aim for >80% code coverage
- Use descriptive test names

Example:

```go
func TestMap(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        fn       func(int) int
        expected []int
    }{
        {
            name:     "double numbers",
            input:    []int{1, 2, 3},
            fn:       func(n int) int { return n * 2 },
            expected: []int{2, 4, 6},
        },
        // more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Map(tt.input, tt.fn)
            // assertions...
        })
    }
}
```

### Documentation Requirements

- Add package-level documentation
- Document all exported functions, types, and constants
- Include usage examples in comments
- Update README for new packages/features
- Add entries to relevant docs/ files

Example:

```go
// ToPascalCase converts a string to PascalCase.
// It handles various input formats including snake_case,
// kebab-case, and space-separated words.
//
// Example:
//   ToPascalCase("hello_world") // Returns: "HelloWorld"
//   ToPascalCase("hello-world") // Returns: "HelloWorld"
func ToPascalCase(s string) string {
    // implementation
}
```

## 🐛 Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the bug
2. **Steps to Reproduce**: Minimal steps to reproduce
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happens
5. **Environment**: Go version, OS, etc.
6. **Code Sample**: Minimal code that demonstrates the issue

## 💡 Suggesting Features

When suggesting features, please:

1. Check if the feature already exists or is planned
2. Provide a clear use case
3. Describe the proposed solution
4. Consider backward compatibility
5. Be open to discussion and alternatives

## 📋 Pull Request Checklist

Before submitting a PR, ensure:

- [ ] Code follows project style guidelines
- [ ] All tests pass (`make test`)
- [ ] Code is properly formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] New code has tests
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] PR description explains the change

## 🤝 Code Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge
4. Your contribution will be in the next release!

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 💬 Questions?

- Open a [GitHub Discussion](https://github.com/ianmuhia/kit/discussions)
- Comment on related issues
- Reach out to maintainers

## 🙏 Thank You

Your contributions make Kit better for everyone!
