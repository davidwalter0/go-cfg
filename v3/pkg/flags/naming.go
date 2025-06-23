package flags

import (
	"strings"
)

// NamingStrategy defines how struct fields map to flag names
type NamingStrategy interface {
	FlagName(fieldPath []string) string
	EnvVarName(fieldPath []string) string
	HelpText(field FieldInfo) string
}

// FieldInfo contains metadata about a struct field
type FieldInfo struct {
	Name         string
	Type         string
	JSONTag      string
	DefaultValue string
	HelpText     string
	Required     bool
	FieldPath    []string
}

// FlatNestStrategy implements clean flat naming (--api-host)
type FlatNestStrategy struct{}

// NewFlatNestStrategy creates a new flat nest naming strategy
func NewFlatNestStrategy() NamingStrategy {
	return &FlatNestStrategy{}
}

// FlagName generates clean flat flag names
func (s *FlatNestStrategy) FlagName(fieldPath []string) string {
	// Convert CamelCase to kebab-case and join with hyphens
	var parts []string
	for _, part := range fieldPath {
		parts = append(parts, camelToKebab(part))
	}
	return strings.Join(parts, "-")
}

// EnvVarName generates environment variable names
func (s *FlatNestStrategy) EnvVarName(fieldPath []string) string {
	// Convert to UPPER_SNAKE_CASE
	var parts []string
	for _, part := range fieldPath {
		parts = append(parts, camelToUpperSnake(part))
	}
	return strings.Join(parts, "_")
}

// HelpText generates help text for a field
func (s *FlatNestStrategy) HelpText(field FieldInfo) string {
	if field.HelpText != "" {
		return field.HelpText
	}
	return field.Name
}

// NestedStrategy implements hierarchical naming (--api.host)
type NestedStrategy struct{}

// NewNestedStrategy creates a new nested naming strategy
func NewNestedStrategy() NamingStrategy {
	return &NestedStrategy{}
}

// FlagName generates nested flag names
func (s *NestedStrategy) FlagName(fieldPath []string) string {
	var parts []string
	for _, part := range fieldPath {
		parts = append(parts, camelToKebab(part))
	}
	return strings.Join(parts, ".")
}

// EnvVarName generates environment variable names
func (s *NestedStrategy) EnvVarName(fieldPath []string) string {
	var parts []string
	for _, part := range fieldPath {
		parts = append(parts, camelToUpperSnake(part))
	}
	return strings.Join(parts, "_")
}

// HelpText generates help text for a field
func (s *NestedStrategy) HelpText(field FieldInfo) string {
	if field.HelpText != "" {
		return field.HelpText
	}
	return field.Name
}

// PrefixedStrategy implements prefixed naming (--config-api-host)
type PrefixedStrategy struct {
	prefix string
}

// NewPrefixedStrategy creates a new prefixed naming strategy
func NewPrefixedStrategy() NamingStrategy {
	return &PrefixedStrategy{prefix: "config"}
}

// FlagName generates prefixed flag names
func (s *PrefixedStrategy) FlagName(fieldPath []string) string {
	var parts []string
	parts = append(parts, s.prefix)
	for _, part := range fieldPath {
		parts = append(parts, camelToKebab(part))
	}
	return strings.Join(parts, "-")
}

// EnvVarName generates environment variable names
func (s *PrefixedStrategy) EnvVarName(fieldPath []string) string {
	var parts []string
	parts = append(parts, strings.ToUpper(s.prefix))
	for _, part := range fieldPath {
		parts = append(parts, camelToUpperSnake(part))
	}
	return strings.Join(parts, "_")
}

// HelpText generates help text for a field
func (s *PrefixedStrategy) HelpText(field FieldInfo) string {
	if field.HelpText != "" {
		return field.HelpText
	}
	return field.Name
}

// Utility functions for case conversion

// camelToKebab converts CamelCase to kebab-case, handling acronyms properly
func camelToKebab(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		// Add hyphen before uppercase letter if:
		// 1. Not at start
		// 2. Previous char was lowercase, OR
		// 3. This is uppercase and next char is lowercase (end of acronym)
		if i > 0 && 'A' <= r && r <= 'Z' {
			prevLower := i > 0 && 'a' <= runes[i-1] && runes[i-1] <= 'z'
			nextLower := i < len(runes)-1 && 'a' <= runes[i+1] && runes[i+1] <= 'z'

			if prevLower || nextLower {
				result.WriteRune('-')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// camelToUpperSnake converts CamelCase to UPPER_SNAKE_CASE, handling acronyms properly
func camelToUpperSnake(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		// Add underscore before uppercase letter if:
		// 1. Not at start
		// 2. Previous char was lowercase, OR
		// 3. This is uppercase and next char is lowercase (end of acronym)
		if i > 0 && 'A' <= r && r <= 'Z' {
			prevLower := i > 0 && 'a' <= runes[i-1] && runes[i-1] <= 'z'
			nextLower := i < len(runes)-1 && 'a' <= runes[i+1] && runes[i+1] <= 'z'

			if prevLower || nextLower {
				result.WriteRune('_')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToUpper(result.String())
}
