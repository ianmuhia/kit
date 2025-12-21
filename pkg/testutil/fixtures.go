package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// LoadFixture loads a test fixture from a JSON file
func LoadFixture(t *testing.T, path string, v interface{}) {
	t.Helper()
	
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("failed to unmarshal fixture %s: %v", path, err)
	}
}

// SaveFixture saves data to a JSON fixture file (useful for updating test fixtures)
func SaveFixture(t *testing.T, path string, v interface{}) {
	t.Helper()
	
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal fixture: %v", err)
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create fixture directory: %v", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write fixture %s: %v", path, err)
	}
}

// TempDir creates a temporary directory for testing
func TempDir(t *testing.T) string {
	t.Helper()
	
	dir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	
	return dir
}

// TempFile creates a temporary file with content for testing
func TempFile(t *testing.T, pattern, content string) string {
	t.Helper()
	
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()
	
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	
	path := f.Name()
	t.Cleanup(func() {
		os.Remove(path)
	})
	
	return path
}
