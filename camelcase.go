package cfg

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// ToCamelCase converts strings to CamelCase format.
// Uses the strcase library for reliable conversion from snake_case, kebab-case, and other formats.
// Example: "json_map" -> "JsonMap", "api-endpoint" -> "ApiEndpoint"
func ToCamelCase(str string) string {
	return strcase.ToCamel(str)
}

// ToLowerSnakeCase converts strings to lower_snake_case format.
// Handles CamelCase, kebab-case, and other formats reliably.
// Example: "CamelCase" -> "camel_case", "api-endpoint" -> "api_endpoint"
func ToLowerSnakeCase(str string) string {
	return strcase.ToSnake(str)
}

// ToUpperSnakeCase converts strings to UPPER_SNAKE_CASE format.
// Used primarily for environment variable names.
// Example: "CamelCase" -> "CAMEL_CASE", "api-endpoint" -> "API_ENDPOINT"
func ToUpperSnakeCase(str string) string {
	return strcase.ToScreamingSnake(str)
}

// ToLowerKebabCase converts strings to lower-kebab-case format.
// Used primarily for command-line flag names.
// Example: "CamelCase" -> "camel-case", "api_endpoint" -> "api-endpoint"
func ToLowerKebabCase(str string) string {
	return strcase.ToKebab(str)
}

// ToUpperKebabCase converts strings to UPPER-KEBAB-CASE format.
// Rarely used but provided for completeness.
// Example: "CamelCase" -> "CAMEL-CASE", "api_endpoint" -> "API-ENDPOINT"
func ToUpperKebabCase(str string) string {
	return strcase.ToScreamingKebab(str)
}

// KeyNameFromCamelCase converts the field's KeyName to UPPER_SNAKE_CASE format
// for use as environment variable names. Uses the strcase library for reliable conversion.
func (field *Field) KeyNameFromCamelCase() {
	field.KeyName = ToUpperSnakeCase(field.KeyName)
}

// FlagNameFromCamelCase converts the field's FlagName to lower-kebab-case format
// for use as command-line flag names. Uses the strcase library for reliable conversion.
func (field *Field) FlagNameFromCamelCase() {
	field.FlagName = ToLowerKebabCase(field.FlagName)
}

// ToLowerKebabCase converts the field's FlagName to lower-kebab-case format.
// This method is an alias for FlagNameFromCamelCase for backward compatibility.
func (field *Field) ToLowerKebabCase() {
	field.FlagName = ToLowerKebabCase(field.FlagName)
}

// Capitalize text
func Capitalize(text string) string {
	switch len(text) {
	case 0:
	case 1:
		text = strings.ToUpper(text[0:1])
	default:
		text = strings.ToUpper(text[0:1]) + text[1:]
	}
	return text
}

// Downcase text
func Downcase(text string) string {
	if len(text) > 0 {
		text = strings.ToLower(text)
	}
	return text
}
