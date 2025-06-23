package types

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

// Value interface that all custom types must implement
type Value interface {
	pflag.Value
	Get() interface{}
}

// NamingConventions handles field name transformations
type NamingConventions struct{}

// ToUpperSnakeCase converts CamelCase to UPPER_SNAKE_CASE for env vars
func (nc *NamingConventions) ToUpperSnakeCase(s string) string {
	// Insert underscores before capital letters (except first)
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)

	// Handle sequences of capitals (e.g., HTTPSProxy -> HTTPS_Proxy)
	re2 := regexp.MustCompile(`([A-Z])([A-Z][a-z])`)
	snake = re2.ReplaceAllString(snake, `${1}_${2}`)

	return strings.ToUpper(snake)
}

// ToKebabCase converts CamelCase to kebab-case for JSON/YAML
func (nc *NamingConventions) ToKebabCase(s string) string {
	// Insert hyphens before capital letters (except first)
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	kebab := re.ReplaceAllString(s, `${1}-${2}`)

	// Handle sequences of capitals
	re2 := regexp.MustCompile(`([A-Z])([A-Z][a-z])`)
	kebab = re2.ReplaceAllString(kebab, `${1}-${2}`)

	return strings.ToLower(kebab)
}

// GetEnvDefault retrieves environment variable with UPPER_SNAKE_CASE naming
func GetEnvDefault(fieldName string) string {
	nc := &NamingConventions{}
	envName := nc.ToUpperSnakeCase(fieldName)
	return os.Getenv(envName)
}

// RegisterValue registers a custom value with pflag using naming conventions
func RegisterValue(flagSet *pflag.FlagSet, value Value, fieldName, usage string) {
	nc := &NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)

	// Try to get default from environment
	if envDefault := GetEnvDefault(fieldName); envDefault != "" {
		if err := value.Set(envDefault); err != nil {
			// Log error but don't fail - use type's default instead
			fmt.Printf("Warning: failed to set %s from env var %s: %v\n",
				flagName, nc.ToUpperSnakeCase(fieldName), err)
		}
	}

	flagSet.Var(value, flagName, usage)
}

// TypeValidator provides common validation patterns
type TypeValidator struct{}

// ValidateRequired checks if a value is set when required
func (tv *TypeValidator) ValidateRequired(value interface{}, fieldName string) error {
	if value == nil {
		return fmt.Errorf("required field %s is not set", fieldName)
	}

	// Handle common zero values
	switch typedValue := value.(type) {
	case string:
		if typedValue == "" {
			return fmt.Errorf("required field %s cannot be empty", fieldName)
		}
	case int, int8, int16, int32, int64:
		// Zero is valid for integers
	case float32, float64:
		// Zero is valid for floats
	case bool:
		// false is valid for booleans
	}

	return nil
}

// SafeParser provides common parsing utilities with error handling
type SafeParser struct{}

// ParseWithFallback attempts to parse a value, falling back to default on error
func (sp *SafeParser) ParseWithFallback(input string, parser func(string) error, fallback func()) error {
	if err := parser(input); err != nil {
		fallback()
		return fmt.Errorf("failed to parse '%s', using default: %w", input, err)
	}
	return nil
}
