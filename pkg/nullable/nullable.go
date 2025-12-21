// Package nullable provides utility functions for working with nullable types
package nullable

import "github.com/oapi-codegen/nullable"

// ToNullableString converts *string to nullable.Nullable[string]
func ToNullableString(s *string) nullable.Nullable[string] {
	if s == nil {
		return nullable.Nullable[string]{}
	}
	return nullable.NewNullableWithValue(*s)
}

// FromNullableString converts nullable.Nullable[string] to *string
func FromNullableString(n nullable.Nullable[string]) *string {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// ToNullableInt converts *int to nullable.Nullable[int]
func ToNullableInt(i *int) nullable.Nullable[int] {
	if i == nil {
		return nullable.Nullable[int]{}
	}
	return nullable.NewNullableWithValue(*i)
}

// FromNullableInt converts nullable.Nullable[int] to *int
func FromNullableInt(n nullable.Nullable[int]) *int {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// ToNullableBool converts *bool to nullable.Nullable[bool]
func ToNullableBool(b *bool) nullable.Nullable[bool] {
	if b == nil {
		return nullable.Nullable[bool]{}
	}
	return nullable.NewNullableWithValue(*b)
}

// FromNullableBool converts nullable.Nullable[bool] to *bool
func FromNullableBool(n nullable.Nullable[bool]) *bool {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}

// ToNullableFloat64 converts *float64 to nullable.Nullable[float64]
func ToNullableFloat64(f *float64) nullable.Nullable[float64] {
	if f == nil {
		return nullable.Nullable[float64]{}
	}
	return nullable.NewNullableWithValue(*f)
}

// FromNullableFloat64 converts nullable.Nullable[float64] to *float64
func FromNullableFloat64(n nullable.Nullable[float64]) *float64 {
	if !n.IsSpecified() {
		return nil
	}
	val := n.MustGet()
	return &val
}
