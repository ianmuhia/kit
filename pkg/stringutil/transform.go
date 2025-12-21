// Package stringutil provides string manipulation utilities
package stringutil

import (
	"regexp"
	"strings"
	"unicode"
)

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	words[0] = strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}
	return strings.Join(words, "")
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "-")
}

// splitWords splits a string into words, handling various naming conventions
func splitWords(s string) []string {
	// Replace common separators with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	
	// Split on transitions from lowercase to uppercase
	re := regexp.MustCompile("([a-z])([A-Z])")
	s = re.ReplaceAllString(s, "$1 $2")
	
	// Split on spaces and filter empty strings
	words := strings.Fields(s)
	return words
}

// Truncate truncates a string to the specified length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// TruncateWithEllipsis truncates a string and adds "..." if truncated
func TruncateWithEllipsis(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// Reverse reverses a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsAlphanumeric checks if a string contains only alphanumeric characters
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
