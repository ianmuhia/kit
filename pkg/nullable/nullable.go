// Package nullable provides utility functions for working with nullable types.
//
// This package offers bidirectional conversion functions between Go pointer types
// and the nullable.Nullable[T] type from github.com/oapi-codegen/nullable.
// This is particularly useful when working with OpenAPI specifications and JSON APIs
// where the distinction between null, undefined, and zero values is important.
//
// Key concepts:
//   - nil pointer: represents an absent/undefined value
//   - nullable.Nullable[T]{}: represents an unspecified value (not sent in request)
//   - nullable.Nullable[T] with value: represents an explicitly set value (including null)
//
// Example usage:
//
//	// Converting to Nullable for API responses
//	name := "John"
//	nullableName := nullable.ToNullableString(&name)
//
//	// Converting from Nullable for database operations
//	var dbName *string = nullable.FromNullableString(nullableName)
//
//	// Handling null values
//	var nilName *string = nil
//	unspecified := nullable.ToNullableString(nilName) // Creates unspecified Nullable
package nullable

import (
	"time"

	"github.com/oapi-codegen/nullable"
)

// String conversions

// ToNullableString converts a *string to nullable.Nullable[string].
// Returns an unspecified Nullable if the pointer is nil.
//
// Example:
//
//	name := "Alice"
//	n := ToNullableString(&name) // Specified with value "Alice"
//	n = ToNullableString(nil)    // Unspecified
func ToNullableString(s *string) nullable.Nullable[string] {
	if s == nil {
		return nullable.Nullable[string]{}
	}
	return nullable.NewNullableWithValue(*s)
}

// FromNullableString converts nullable.Nullable[string] to *string.
// Returns nil if the Nullable is unspecified.
//
// Example:
//
//	n := nullable.NewNullableWithValue("Bob")
//	s := FromNullableString(n) // *string pointing to "Bob"
func FromNullableString(n nullable.Nullable[string]) *string {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// StringValue returns the string value or the provided default if unspecified.
//
// Example:
//
//	n := nullable.Nullable[string]{}
//	s := StringValue(n, "default") // Returns "default"
func StringValue(n nullable.Nullable[string], defaultVal string) string {
	if !n.IsSpecified() {
		return defaultVal
	}
	return n.MustGet()
}

// StringOrEmpty returns the string value or empty string if unspecified.
//
// Example:
//
//	n := nullable.Nullable[string]{}
//	s := StringOrEmpty(n) // Returns ""
func StringOrEmpty(n nullable.Nullable[string]) string {
	return StringValue(n, "")
}

// Integer conversions

// ToNullableInt converts a *int to nullable.Nullable[int].
// Returns an unspecified Nullable if the pointer is nil.
//
// Example:
//
//	age := 25
//	n := ToNullableInt(&age) // Specified with value 25
//	n = ToNullableInt(nil)   // Unspecified
func ToNullableInt(i *int) nullable.Nullable[int] {
	if i == nil {
		return nullable.Nullable[int]{}
	}
	return nullable.NewNullableWithValue(*i)
}

// Boolean conversions

// ToNullableBool converts a *bool to nullable.Nullable[bool].
// Returns an unspecified Nullable if the pointer is nil.
//
// Example:
//
//	active := true
//	n := ToNullableBool(&active) // Specified with value true
//	n = ToNullableBool(nil)      // Unspecified
func ToNullableBool(b *bool) nullable.Nullable[bool] {
	if b == nil {
		return nullable.Nullable[bool]{}
	}
	return nullable.NewNullableWithValue(*b)
}

// FromNullableBool converts nullable.Nullable[bool] to *bool.
// Returns nil if the Nullable is unspecified.
//
// Example:
//
//	n := nullable.NewNullableWithValue(false)
//	b := FromNullableBool(n) // *bool pointing to false
func FromNullableBool(n nullable.Nullable[bool]) *bool {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// BoolValue returns the bool value or the provided default if unspecified.
//
// Example:
//
//	n := nullable.Nullable[bool]{}
//	b := BoolValue(n, false) // Returns false
func BoolValue(n nullable.Nullable[bool], defaultVal bool) bool {
	if !n.IsSpecified() {
		return defaultVal
	}
	return n.MustGet()
}

// Integer conversions

// ToNullableInt converts a *int to nullable.Nullable[int].
//
// Example:
//
//	n := nullable.Nullable[int]{}
//	i := IntValue(n, 0) // Returns 0
func IntValue(n nullable.Nullable[int], defaultVal int) int {
	if !n.IsSpecified() {
		return defaultVal
	}
	return n.MustGet()
}

// ToNullableInt32 converts a *int32 to nullable.Nullable[int32].
func ToNullableInt32(i *int32) nullable.Nullable[int32] {
	if i == nil {
		return nullable.Nullable[int32]{}
	}
	return nullable.NewNullableWithValue(*i)
}

// FromNullableInt32 converts nullable.Nullable[int32] to *int32.
func FromNullableInt32(n nullable.Nullable[int32]) *int32 {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// ToNullableInt64 converts a *int64 to nullable.Nullable[int64].
func ToNullableInt64(i *int64) nullable.Nullable[int64] {
	if i == nil {
		return nullable.Nullable[int64]{}
	}
	return nullable.NewNullableWithValue(*i)
}

// FromNullableInt64 converts nullable.Nullable[int64] to *int64.
func FromNullableInt64(n nullable.Nullable[int64]) *int64 {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// Float conversions

// ToNullableFloat32 converts a *float32 to nullable.Nullable[float32].
func ToNullableFloat32(f *float32) nullable.Nullable[float32] {
	if f == nil {
		return nullable.Nullable[float32]{}
	}
	return nullable.NewNullableWithValue(*f)
}

// FromNullableFloat32 converts nullable.Nullable[float32] to *float32.
func FromNullableFloat32(n nullable.Nullable[float32]) *float32 {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// ToNullableFloat64 converts a *float64 to nullable.Nullable[float64].
// Returns an unspecified Nullable if the pointer is nil.
//
// Example:
//
//	price := 99.99
//	n := ToNullableFloat64(&price) // Specified with value 99.99
//	n = ToNullableFloat64(nil)     // Unspecified
func ToNullableFloat64(f *float64) nullable.Nullable[float64] {
	if f == nil {
		return nullable.Nullable[float64]{}
	}
	return nullable.NewNullableWithValue(*f)
}

// FromNullableFloat64 converts nullable.Nullable[float64] to *float64.
// Returns nil if the Nullable is unspecified.
//
// Example:
//
//	n := nullable.NewNullableWithValue(3.14)
//	f := FromNullableFloat64(n) // *float64 pointing to 3.14
func FromNullableFloat64(n nullable.Nullable[float64]) *float64 {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// Float64Value returns the float64 value or the provided default if unspecified.
//
// Example:
//
//	n := nullable.Nullable[float64]{}
//	f := Float64Value(n, 0.0) // Returns 0.0
func Float64Value(n nullable.Nullable[float64], defaultVal float64) float64 {
	if !n.IsSpecified() {
		return defaultVal
	}
	return n.MustGet()
}

// Time conversions

// ToNullableTime converts a *time.Time to nullable.Nullable[time.Time].
// Returns an unspecified Nullable if the pointer is nil.
//
// Example:
//
//	now := time.Now()
//	n := ToNullableTime(&now) // Specified with current time
//	n = ToNullableTime(nil)   // Unspecified
func ToNullableTime(t *time.Time) nullable.Nullable[time.Time] {
	if t == nil {
		return nullable.Nullable[time.Time]{}
	}
	return nullable.NewNullableWithValue(*t)
}

// FromNullableTime converts nullable.Nullable[time.Time] to *time.Time.
// Returns nil if the Nullable is unspecified.
//
// Example:
//
//	n := nullable.NewNullableWithValue(time.Now())
//	t := FromNullableTime(n) // *time.Time
func FromNullableTime(n nullable.Nullable[time.Time]) *time.Time {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// TimeValue returns the time.Time value or the provided default if unspecified.
//
// Example:
//
//	n := nullable.Nullable[time.Time]{}
//	t := TimeValue(n, time.Now()) // Returns current time
func TimeValue(n nullable.Nullable[time.Time], defaultVal time.Time) time.Time {
	if !n.IsSpecified() {
		return defaultVal
	}
	return n.MustGet()
}

// Generic helpers

// IsSpecified returns true if the Nullable has a specified value.
// This is a convenience wrapper around nullable.Nullable.IsSpecified().
//
// Example:
//
//	n := nullable.NewNullableWithValue("test")
//	if IsSpecified(n) {
//	    // Value is specified (even if it's the zero value)
//	}
func IsSpecified[T any](n nullable.Nullable[T]) bool {
	return n.IsSpecified()
}

// ValueOr returns the value if specified, otherwise returns the default.
// This is a generic version of the type-specific value functions.
//
// Example:
//
//	n := nullable.Nullable[string]{}
//	s := ValueOr(n, "default") // Returns "default"
func ValueOr[T any](n nullable.Nullable[T], defaultVal T) T {
	if !n.IsSpecified() {
		return defaultVal
	}
	return n.MustGet()
}

// Ptr is a helper function that returns a pointer to the given value.
// Useful for creating pointers to literals.
//
// Example:
//
//	n := ToNullableString(Ptr("hello"))
//	age := ToNullableInt(Ptr(25))
func Ptr[T any](v T) *T {
	return &v
}
