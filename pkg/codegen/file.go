package codegen

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDir ensures that a directory exists, creating it if necessary
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ErrFileExists is returned by WriteFileSafe when the target file already exists.
var ErrFileExists = fmt.Errorf("file already exists")

// WriteFile writes content to a file, creating parent directories if needed.
// It overwrites any existing file at the path.
func WriteFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// WriteFileSafe writes content to a file only if it does not already exist.
// Returns an error wrapping ErrFileExists if a file is already present so
// callers can distinguish a conflict from other I/O errors:
//
//	err := codegen.WriteFileSafe(path, content)
//	if errors.Is(err, codegen.ErrFileExists) { ... }
func WriteFileSafe(path string, content []byte) error {
	if FileExists(path) {
		return fmt.Errorf("%w: %s", ErrFileExists, path)
	}
	return WriteFile(path, content)
}

// ReadFile reads a file and returns its content
func ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return content, nil
}
