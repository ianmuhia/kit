package errorgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsUpperSnakeCase(t *testing.T) {
	valid := []string{"NOT_FOUND", "INVALID_INPUT", "A", "CODE_123"}
	for _, s := range valid {
		assert.True(t, isUpperSnakeCase(s), "expected true for %q", s)
	}

	invalid := []string{"", "_LEADING", "TRAILING_", "lower_case", "Mixed", "HAS SPACE"}
	for _, s := range invalid {
		assert.False(t, isUpperSnakeCase(s), "expected false for %q", s)
	}
}

func TestIsValidHTTPStatus(t *testing.T) {
	assert.True(t, isValidHTTPStatus(200))
	assert.True(t, isValidHTTPStatus(100))
	assert.True(t, isValidHTTPStatus(599))
	assert.False(t, isValidHTTPStatus(0))
	assert.False(t, isValidHTTPStatus(99))
	assert.False(t, isValidHTTPStatus(600))
}

func TestIsValidSeverity(t *testing.T) {
	for _, s := range []string{"critical", "high", "medium", "low", "CRITICAL", "High"} {
		assert.True(t, isValidSeverity(s), "expected true for %q", s)
	}
	assert.False(t, isValidSeverity(""))
	assert.False(t, isValidSeverity("info"))
	assert.False(t, isValidSeverity("severe"))
}

func TestValidate(t *testing.T) {
	t.Run("missing package", func(t *testing.T) {
		c := &ErrorConfig{Errors: []ErrorDefinition{{Name: "ErrFoo", Code: "FOO", Message: "foo"}}}
		require.ErrorContains(t, c.validate(), "package name is required")
	})

	t.Run("empty errors list", func(t *testing.T) {
		c := &ErrorConfig{Package: "errs"}
		require.ErrorContains(t, c.validate(), "errors list must not be empty")
	})

	t.Run("missing required fields", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors:  []ErrorDefinition{{Name: "ErrFoo"}},
		}
		require.Error(t, c.validate())
	})

	t.Run("duplicate code", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors: []ErrorDefinition{
				{Name: "ErrFoo", Code: "FOO", Message: "foo"},
				{Name: "ErrBar", Code: "FOO", Message: "bar"},
			},
		}
		require.ErrorContains(t, c.validate(), "duplicate error code")
	})

	t.Run("duplicate name", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors: []ErrorDefinition{
				{Name: "ErrFoo", Code: "FOO", Message: "foo"},
				{Name: "ErrFoo", Code: "BAR", Message: "bar"},
			},
		}
		require.ErrorContains(t, c.validate(), "duplicate error name")
	})

	t.Run("invalid code format", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors:  []ErrorDefinition{{Name: "ErrFoo", Code: "not_upper", Message: "foo"}},
		}
		require.ErrorContains(t, c.validate(), "UPPER_SNAKE_CASE")
	})

	t.Run("invalid http status", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors:  []ErrorDefinition{{Name: "ErrFoo", Code: "FOO", Message: "foo", HTTPStatus: 999}},
		}
		require.ErrorContains(t, c.validate(), "invalid HTTP status")
	})

	t.Run("invalid severity", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors:  []ErrorDefinition{{Name: "ErrFoo", Code: "FOO", Message: "foo", Severity: "fatal"}},
		}
		require.ErrorContains(t, c.validate(), "invalid severity")
	})

	t.Run("parameter not in message", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors: []ErrorDefinition{
				{Name: "ErrFoo", Code: "FOO", Message: "no placeholder", Parameters: []string{"id"}},
			},
		}
		require.ErrorContains(t, c.validate(), "not found in message")
	})

	t.Run("valid config", func(t *testing.T) {
		c := &ErrorConfig{
			Package: "errs",
			Errors: []ErrorDefinition{
				{Name: "ErrFoo", Code: "FOO", Message: "item {id} not found", Parameters: []string{"id"}, HTTPStatus: 404, Severity: "medium"},
			},
		}
		require.NoError(t, c.validate())
	})
}

func TestNewGenerator(t *testing.T) {
	t.Run("defaults applied", func(t *testing.T) {
		g, err := NewGenerator(WithInputFile("errors.cue"), WithOutputFile("errors.go"))
		require.NoError(t, err)
		assert.Equal(t, "errors.cue", g.config.inputFile)
		assert.Equal(t, "errors.go", g.config.outputFile)
	})

	t.Run("missing input file", func(t *testing.T) {
		_, err := NewGenerator(WithInputFile(""), WithOutputFile("out.go"))
		require.ErrorContains(t, err, "input file is required")
	})

	t.Run("missing output file", func(t *testing.T) {
		_, err := NewGenerator(WithInputFile("in.cue"), WithOutputFile(""))
		require.ErrorContains(t, err, "output file is required")
	})

	t.Run("package name override", func(t *testing.T) {
		g, err := NewGenerator(WithInputFile("in.cue"), WithOutputFile("out.go"), WithPackageName("myerrs"))
		require.NoError(t, err)
		assert.Equal(t, "myerrs", g.config.packageName)
	})
}
