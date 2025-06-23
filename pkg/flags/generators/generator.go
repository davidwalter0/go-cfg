package generators

import (
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

// TypeGenerator generates custom flag types based on templates
type TypeGenerator struct {
	PackageName string
	Imports     []string
	Templates   map[string]*template.Template
}

// NewTypeGenerator creates a new type generator
func NewTypeGenerator(packageName string) *TypeGenerator {
	return &TypeGenerator{
		PackageName: packageName,
		Templates:   make(map[string]*template.Template),
	}
}

// TypeSpec defines the specification for generating a custom type
type TypeSpec struct {
	TypeName        string
	BaseType        string
	ValueMethod     string
	ParseMethod     string
	ValidatorMethod string
	DefaultValue    string
	Usage           string
	JSONMarshal     bool
	YAMLMarshal     bool
	Validators      []ValidatorSpec
}

// ValidatorSpec defines a validation rule
type ValidatorSpec struct {
	Name     string
	Check    string
	ErrorMsg string
}

// RegisterTemplate registers a template for code generation
func (tg *TypeGenerator) RegisterTemplate(name, templateStr string) error {
	tmpl, err := template.New(name).Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	tg.Templates[name] = tmpl
	return nil
}

// GenerateType generates a custom flag type based on the specification
func (tg *TypeGenerator) GenerateType(spec TypeSpec) (string, error) {
	tmpl, exists := tg.Templates["custom_type"]
	if !exists {
		return "", fmt.Errorf("template 'custom_type' not found")
	}

	var buf strings.Builder

	// Generate package header
	buf.WriteString(fmt.Sprintf("package %s\n\n", tg.PackageName))

	// Generate imports
	if len(tg.Imports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range tg.Imports {
			buf.WriteString(fmt.Sprintf("\t%s\n", imp))
		}
		buf.WriteString(")\n\n")
	}

	// Generate type using template
	if err := tmpl.Execute(&buf, spec); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the generated code
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return buf.String(), fmt.Errorf("failed to format generated code: %w", err)
	}

	return string(formatted), nil
}

// GenerateEnumType generates an enum flag type with validation
func (tg *TypeGenerator) GenerateEnumType(typeName string, values []string, defaultValue string) (string, error) {
	spec := TypeSpec{
		TypeName:     typeName,
		BaseType:     "string",
		DefaultValue: defaultValue,
		JSONMarshal:  true,
		YAMLMarshal:  true,
	}

	// Create validators for enum values
	for _, value := range values {
		spec.Validators = append(spec.Validators, ValidatorSpec{
			Name:     fmt.Sprintf("Valid%s", strings.ToUpper(value[:1])+strings.ToLower(value[1:])),
			Check:    fmt.Sprintf(`strings.ToLower(s) == "%s"`, strings.ToLower(value)),
			ErrorMsg: fmt.Sprintf("invalid value, expected one of: %s", strings.Join(values, ", ")),
		})
	}

	return tg.GenerateType(spec)
}

// Standard templates for common patterns

// CustomTypeTemplate is the base template for custom flag types
const CustomTypeTemplate = `
// {{.TypeName}}Value implements pflag.Value for {{.TypeName}}
type {{.TypeName}}Value struct {
	value {{.BaseType}}
	{{if .Validators}}validators []func({{.BaseType}}) error{{end}}
}

// New{{.TypeName}}Value creates a new {{.TypeName}}Value
func New{{.TypeName}}Value(defaultValue {{.BaseType}}) *{{.TypeName}}Value {
	v := &{{.TypeName}}Value{
		value: defaultValue,
	}
	{{if .Validators}}
	v.validators = []func({{.BaseType}}) error{
		{{range .Validators}}v.{{.Name}},{{end}}
	}
	{{end}}
	return v
}

// String returns the string representation (required by pflag.Value)
func (v *{{.TypeName}}Value) String() string {
	{{if eq .BaseType "string"}}
	return v.value
	{{else}}
	return fmt.Sprintf("%v", v.value)
	{{end}}
}

// Set parses and sets the value (required by pflag.Value)
func (v *{{.TypeName}}Value) Set(s string) error {
	{{if .ParseMethod}}
	parsed, err := {{.ParseMethod}}(s)
	if err != nil {
		return fmt.Errorf("invalid {{.TypeName}}: %w", err)
	}
	{{else if eq .BaseType "string"}}
	parsed := s
	{{else if eq .BaseType "int"}}
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid {{.TypeName}}: %w", err)
	}
	{{else if eq .BaseType "bool"}}
	parsed, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("invalid {{.TypeName}}: %w", err)
	}
	{{end}}
	
	{{if .Validators}}
	// Run validators
	for _, validator := range v.validators {
		if err := validator(parsed); err != nil {
			return err
		}
	}
	{{end}}
	
	v.value = parsed
	return nil
}

// Type returns the type name for help text
func (v *{{.TypeName}}Value) Type() string {
	return "{{.TypeName}}"
}

// Get returns the underlying value (required by pflag.Getter)
func (v *{{.TypeName}}Value) Get() interface{} {
	return v.value
}

// GetValue returns the typed value
func (v *{{.TypeName}}Value) GetValue() {{.BaseType}} {
	return v.value
}

{{if .Validators}}
{{range .Validators}}
// {{.Name}} validates {{$.TypeName}} value
func (v *{{$.TypeName}}Value) {{.Name}}(value {{$.BaseType}}) error {
	if !({{.Check}}) {
		return fmt.Errorf("{{.ErrorMsg}}")
	}
	return nil
}
{{end}}
{{end}}

{{if .JSONMarshal}}
// MarshalJSON implements json.Marshaler
func (v {{.TypeName}}Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// UnmarshalJSON implements json.Unmarshaler
func (v *{{.TypeName}}Value) UnmarshalJSON(data []byte) error {
	var value {{.BaseType}}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	{{if eq .BaseType "string"}}
	return v.Set(value)
	{{else}}
	return v.Set(fmt.Sprintf("%v", value))
	{{end}}
}
{{end}}

{{if .YAMLMarshal}}
// MarshalYAML implements yaml.Marshaler
func (v {{.TypeName}}Value) MarshalYAML() (interface{}, error) {
	return v.value, nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (v *{{.TypeName}}Value) UnmarshalYAML(node *yaml.Node) error {
	var value {{.BaseType}}
	if err := node.Decode(&value); err != nil {
		return err
	}
	{{if eq .BaseType "string"}}
	return v.Set(value)
	{{else}}
	return v.Set(fmt.Sprintf("%v", value))
	{{end}}
}
{{end}}
`

// SliceTypeTemplate generates slice flag types
const SliceTypeTemplate = `
// {{.TypeName}}SliceValue implements pflag.Value for []{{.BaseType}}
type {{.TypeName}}SliceValue struct {
	values []{{.BaseType}}
	separator string
}

// New{{.TypeName}}SliceValue creates a new {{.TypeName}}SliceValue
func New{{.TypeName}}SliceValue(defaultValues []{{.BaseType}}, separator string) *{{.TypeName}}SliceValue {
	if separator == "" {
		separator = ","
	}
	return &{{.TypeName}}SliceValue{
		values: append([]{{.BaseType}}(nil), defaultValues...),
		separator: separator,
	}
}

// String returns the string representation
func (v *{{.TypeName}}SliceValue) String() string {
	if len(v.values) == 0 {
		return ""
	}
	
	var strs []string
	for _, val := range v.values {
		{{if eq .BaseType "string"}}
		strs = append(strs, val)
		{{else}}
		strs = append(strs, fmt.Sprintf("%v", val))
		{{end}}
	}
	return strings.Join(strs, v.separator)
}

// Set parses and appends values
func (v *{{.TypeName}}SliceValue) Set(s string) error {
	if s == "" {
		return nil
	}
	
	parts := strings.Split(s, v.separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		{{if eq .BaseType "string"}}
		v.values = append(v.values, part)
		{{else if eq .BaseType "int"}}
		parsed, err := strconv.Atoi(part)
		if err != nil {
			return fmt.Errorf("invalid integer: %s", part)
		}
		v.values = append(v.values, parsed)
		{{else if eq .BaseType "bool"}}
		parsed, err := strconv.ParseBool(part)
		if err != nil {
			return fmt.Errorf("invalid boolean: %s", part)
		}
		v.values = append(v.values, parsed)
		{{end}}
	}
	
	return nil
}

// Type returns the type name
func (v *{{.TypeName}}SliceValue) Type() string {
	return "{{.TypeName}}Slice"
}

// Get returns the underlying slice
func (v *{{.TypeName}}SliceValue) Get() interface{} {
	return append([]{{.BaseType}}(nil), v.values...)
}

// GetValues returns the typed slice
func (v *{{.TypeName}}SliceValue) GetValues() []{{.BaseType}} {
	return append([]{{.BaseType}}(nil), v.values...)
}

// Len returns the number of values
func (v *{{.TypeName}}SliceValue) Len() int {
	return len(v.values)
}

// Clear removes all values
func (v *{{.TypeName}}SliceValue) Clear() {
	v.values = nil
}
`

// RegisterDefaultTemplates registers the standard templates
func (tg *TypeGenerator) RegisterDefaultTemplates() error {
	if err := tg.RegisterTemplate("custom_type", CustomTypeTemplate); err != nil {
		return err
	}

	if err := tg.RegisterTemplate("slice_type", SliceTypeTemplate); err != nil {
		return err
	}

	return nil
}

// Common generator functions

// GenerateLogLevelType generates a log level enum type
func GenerateLogLevelType() (string, error) {
	tg := NewTypeGenerator("types")
	tg.Imports = []string{`"encoding/json"`, `"fmt"`, `"strings"`, `"gopkg.in/yaml.v3"`}

	if err := tg.RegisterDefaultTemplates(); err != nil {
		return "", err
	}

	return tg.GenerateEnumType("LogLevel",
		[]string{"debug", "info", "warn", "error", "fatal"},
		"info")
}

// GenerateEnvironmentType generates an environment enum type
func GenerateEnvironmentType() (string, error) {
	tg := NewTypeGenerator("types")
	tg.Imports = []string{`"encoding/json"`, `"fmt"`, `"strings"`, `"gopkg.in/yaml.v3"`}

	if err := tg.RegisterDefaultTemplates(); err != nil {
		return "", err
	}

	return tg.GenerateEnumType("Environment",
		[]string{"development", "staging", "production"},
		"development")
}

// GenerateStringSliceType generates a string slice type
func GenerateStringSliceType() (string, error) {
	tg := NewTypeGenerator("types")
	tg.Imports = []string{`"fmt"`, `"strings"`}

	if err := tg.RegisterDefaultTemplates(); err != nil {
		return "", err
	}

	tmpl := tg.Templates["slice_type"]

	spec := TypeSpec{
		TypeName: "String",
		BaseType: "string",
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, spec); err != nil {
		return "", err
	}

	return buf.String(), nil
}
