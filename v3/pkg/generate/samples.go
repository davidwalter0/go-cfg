package generate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SampleGenerator creates sample configuration files
type SampleGenerator struct{}

// NewSampleGenerator creates a new sample generator
func NewSampleGenerator() *SampleGenerator {
	return &SampleGenerator{}
}

// GenerateJSON creates a sample JSON configuration
func (g *SampleGenerator) GenerateJSON(config interface{}) ([]byte, error) {
	sample := g.generateSampleData(config)
	return json.MarshalIndent(sample, "", "  ")
}

// GenerateYAML creates a sample YAML configuration
func (g *SampleGenerator) GenerateYAML(config interface{}) ([]byte, error) {
	sample := g.generateSampleData(config)
	return yaml.Marshal(sample)
}

// generateSampleData recursively generates sample data from struct
func (g *SampleGenerator) generateSampleData(config interface{}) interface{} {
	return g.processValue(reflect.ValueOf(config), reflect.TypeOf(config))
}

// processValue processes a reflect.Value to generate sample data
func (g *SampleGenerator) processValue(val reflect.Value, typ reflect.Type) interface{} {
	// Handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			// Create a new instance
			newVal := reflect.New(typ.Elem())
			return g.processValue(newVal.Elem(), typ.Elem())
		}
		return g.processValue(val.Elem(), typ.Elem())
	}

	switch val.Kind() {
	case reflect.Struct:
		return g.processStruct(val, typ)
	case reflect.String:
		return g.getSampleString(typ)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return g.getSampleInt(typ)
	case reflect.Bool:
		return g.getSampleBool(typ)
	case reflect.Map:
		return g.getSampleMap(typ)
	case reflect.Slice:
		return g.getSampleSlice(typ)
	default:
		return nil
	}
}

// processStruct processes a struct to generate sample data
func (g *SampleGenerator) processStruct(val reflect.Value, typ reflect.Type) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Get JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Remove omitempty and other options
		if idx := strings.Index(jsonTag, ","); idx >= 0 {
			jsonTag = jsonTag[:idx]
		}

		// Generate sample value
		sampleValue := g.processValue(field, fieldType.Type)
		if sampleValue != nil {
			result[jsonTag] = sampleValue
		}
	}

	return result
}

// getSampleString generates sample string values
func (g *SampleGenerator) getSampleString(typ reflect.Type) string {
	// Check if this is a special type
	if typ == reflect.TypeOf(time.Duration(0)) {
		return "30s"
	}

	return "example-value"
}

// getSampleInt generates sample integer values
func (g *SampleGenerator) getSampleInt(typ reflect.Type) int64 {
	// Special handling for duration
	if typ == reflect.TypeOf(time.Duration(0)) {
		return int64(30 * time.Second)
	}

	return 8080
}

// getSampleBool generates sample boolean values
func (g *SampleGenerator) getSampleBool(typ reflect.Type) bool {
	return false
}

// getSampleMap generates sample map values
func (g *SampleGenerator) getSampleMap(typ reflect.Type) map[string]interface{} {
	return map[string]interface{}{
		"example_key": "example_value",
	}
}

// getSampleSlice generates sample slice values
func (g *SampleGenerator) getSampleSlice(typ reflect.Type) []interface{} {
	return []interface{}{"example1", "example2"}
}

// DocumentationGenerator creates documentation from struct tags
type DocumentationGenerator struct{}

// NewDocumentationGenerator creates a new documentation generator
func NewDocumentationGenerator() *DocumentationGenerator {
	return &DocumentationGenerator{}
}

// GenerateMarkdown creates markdown documentation
func (g *DocumentationGenerator) GenerateMarkdown(config interface{}) string {
	var doc strings.Builder

	doc.WriteString("# Configuration Documentation\n\n")
	doc.WriteString("This document describes all available configuration options.\n\n")

	g.processStructForDocs(reflect.ValueOf(config), reflect.TypeOf(config), &doc, "")

	return doc.String()
}

// processStructForDocs processes a struct to generate documentation
func (g *DocumentationGenerator) processStructForDocs(val reflect.Value, typ reflect.Type, doc *strings.Builder, prefix string) {
	// Handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			newVal := reflect.New(typ.Elem())
			g.processStructForDocs(newVal.Elem(), typ.Elem(), doc, prefix)
			return
		}
		g.processStructForDocs(val.Elem(), typ.Elem(), doc, prefix)
		return
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Remove omitempty and other options
		if idx := strings.Index(jsonTag, ","); idx >= 0 {
			jsonTag = jsonTag[:idx]
		}

		fieldPath := jsonTag
		if prefix != "" {
			fieldPath = prefix + "." + jsonTag
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			doc.WriteString(fmt.Sprintf("## %s\n\n", fieldPath))
			helpText := fieldType.Tag.Get("help")
			if helpText != "" {
				doc.WriteString(fmt.Sprintf("%s\n\n", helpText))
			}
			g.processStructForDocs(field, fieldType.Type, doc, fieldPath)
			continue
		}

		// Document this field
		doc.WriteString(fmt.Sprintf("### %s\n\n", fieldPath))
		doc.WriteString(fmt.Sprintf("- **Type**: %s\n", fieldType.Type.String()))

		if defaultVal := fieldType.Tag.Get("default"); defaultVal != "" {
			doc.WriteString(fmt.Sprintf("- **Default**: `%s`\n", defaultVal))
		}

		if required := fieldType.Tag.Get("required"); required == "true" {
			doc.WriteString("- **Required**: Yes\n")
		}

		if helpText := fieldType.Tag.Get("help"); helpText != "" {
			doc.WriteString(fmt.Sprintf("- **Description**: %s\n", helpText))
		}

		doc.WriteString("\n")
	}
}
