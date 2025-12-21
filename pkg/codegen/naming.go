// Package codegen provides utilities for code generation
package codegen

import (
	"unicode"
)

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize converts the first letter of a string to lowercase
func Uncapitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
