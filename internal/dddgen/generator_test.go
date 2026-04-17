package dddgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDomainName(t *testing.T) {
	valid := []string{"booking", "user", "Order", "my_domain", "domain2"}
	for _, n := range valid {
		assert.NoError(t, validateDomainName(n), "expected valid: %q", n)
	}

	invalid := []string{"", "1invalid", "has-hyphen", "has space", "has.dot"}
	for _, n := range invalid {
		assert.Error(t, validateDomainName(n), "expected invalid: %q", n)
	}
}

func TestNew_missingModulePath(t *testing.T) {
	_, err := New(Config{DomainName: "booking", OutputDir: t.TempDir()})
	require.ErrorContains(t, err, "module path is required")
}

func TestNew_invalidDomainName(t *testing.T) {
	_, err := New(Config{DomainName: "1bad", ModulePath: "github.com/x/y", OutputDir: t.TempDir()})
	require.Error(t, err)
}

func TestNew_existingDomainRejected(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "booking"), 0755))
	_, err := New(Config{DomainName: "booking", ModulePath: "github.com/x/y", OutputDir: dir})
	require.ErrorContains(t, err, "already exists")
}

func TestNew_success(t *testing.T) {
	g, err := New(Config{
		DomainName: "booking",
		ModulePath: "github.com/x/y",
		OutputDir:  t.TempDir(),
	})
	require.NoError(t, err)
	assert.Equal(t, "Booking", g.data.DomainTitle)
	assert.Equal(t, "booking", g.data.DomainLower)
	assert.Equal(t, "github.com/x/y", g.data.ModulePath)
}

func TestGenerate_createsFiles(t *testing.T) {
	dir := t.TempDir()
	g, err := New(Config{
		DomainName: "order",
		ModulePath: "github.com/x/y",
		OutputDir:  dir,
	})
	require.NoError(t, err)
	require.NoError(t, g.Generate())

	expected := []string{
		filepath.Join(dir, "order", "order.go"),
		filepath.Join(dir, "order", "repository.go"),
		filepath.Join(dir, "order", "errors.go"),
		filepath.Join(dir, "order", "app", "service.go"),
		filepath.Join(dir, "order", "adapters", "order_http.go"),
	}
	for _, f := range expected {
		assert.FileExists(t, f)
	}
}
