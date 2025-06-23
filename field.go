package cfg

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

var debug bool

// Debug set the debug test var to true
func Debug(args ...bool) {
	if len(args) == 0 {
		debug = true
		return
	}
	debug = args[0]
}

// FieldPtr for the struct field field
type FieldPtr interface{}

// Field holds comprehensive metadata for a single struct field during configuration processing.
// This struct contains all the information needed to create command-line flags, environment
// variables, and handle type conversion for a single configuration field.
//
// The Field is created during reflection processing of configuration structs and contains
// both the original struct field information and the processed configuration metadata.
type Field struct {
	reflect.StructField        // Embedded struct field information from reflection
	FieldPtr                   // Pointer to the actual field value for modification
	UseFlags            bool   // Whether this field should generate command-line flags
	StructName          string // Name of the owning struct (used for prefixing)
	Name                string // Field name (from struct field or json tag override)
	KeyName             string // Environment variable name (UPPER_SNAKE_CASE format)
	Default             string // Default value from struct tag (empty if none specified)
	EnvText             string // Current environment variable value (empty if not set)
	Short               string // Short flag name (single character, rarely used)
	Doc                 string // Documentation/help text from struct tag
	FlagName            string // Command-line flag name (kebab-case format)
	Value               string // Final resolved value (env var > default > type default)
	Omit                bool   // Whether to omit empty values (from json:",omitempty")
	Required            bool   // Whether this field must have a non-empty value
	Depth               int    // Nesting depth in struct hierarchy (0 = top level)
	Ignore              bool   // Whether to skip this field entirely during processing
	Error               error  // Any error encountered during field processing
	Type                string // Go type name as string (for debugging/logging)
	Prefix              string // Prefix applied to flag/env names (includes hierarchy)
}

// Get a tag from the struct tags
func (field *Field) Get(name string) string {
	text := field.StructField.Tag.Get(name)
	if len(text) > 0 {
		return text
	}
	return ""
}

// SetField from the struct tags, env, or interpolated values
func (field *Field) SetField() {
	defer func() {
		if err := recover(); err != nil {
			// Log and store the error instead of panicking
			field.Error = fmt.Errorf("field processing panic: %v", err)
		}
	}()

	field.SetDefault()
	field.SetName()
	field.SetIgnore()
	field.SetDoc()
	//	field.SetShort()
	field.SetOmit()
	field.SetRequired()
	field.SetKeyName()
	field.SetFlagName()
	field.SetValueFromEnv()
	field.SetType()
	// Must run last if referencing any info from field settings
	if field.UseFlags {
		field.AddFlag()
	}
}

// String formats args as yaml string
func (field *Field) String() string {
	text, err := yaml.Marshal(*field)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	return string(text)
}

// SetOmit read tag omit option and set when enabled, via ,omit
// ,omitempty or '-' the hyphen option
func (field *Field) SetOmit() {
	json := field.Get("json")
	field.Omit = json == "-" || strings.Contains(json, ",omitempty") || strings.Contains(json, ",omit")
}

// SetDefault read tag default option and save the text
func (field *Field) SetDefault() {
	field.Default = field.Get("default")
	field.Value = field.Default
}

// SetIgnore read tag ignore option and save the text
func (field *Field) SetIgnore() {
	text := field.Get("ignore")
	if v, err := strconv.ParseBool(text); err == nil {
		field.Ignore = v
	}
}

// SetDoc read tag doc option and save the text
func (field *Field) SetDoc() {
	field.Doc = field.Get("doc")
	if len(field.Doc) == 0 {
		field.Doc = field.Get("help")
	}
	if len(field.Doc) == 0 {
		field.Doc = field.Get("usage")
	}
}

// SetShort read tag short option and save the text
func (field *Field) SetShort() {
	field.Short = field.Get("short")
}

// // SetPrefix read tag prefix option and save the text
// func (field *Field) SetPrefix() {
// 	// field.Prefix = field.Get("prefix")
// }

// SetName read tag name option and save the text
func (field *Field) SetName() {
	name := field.Get("json")
	i := strings.Index(name, ",")
	if i > 0 {
		name = name[0:i]
	}
	if len(name) > 0 {
		field.Name = name
	}
}

// SetIgnore read tag required option and save the text
func (field *Field) SetRequired() {
	text := field.Get("required")
	if v, err := strconv.ParseBool(text); err == nil {
		field.Required = v
	}
}

// SetValueFromEnv uses the value from the environment for this
// structure tag replacing the default tag value
func (field *Field) SetValueFromEnv() {
	if len(field.KeyName) > 0 {
		field.EnvText, _ = LookupEnv(field.KeyName)
	}
	if len(field.EnvText) != 0 {
		field.Value = field.EnvText
	}
}

// SetKeyName read tag keyword option and save the text
func (field *Field) SetKeyName() {
	if len(field.KeyName) == 0 {
		panic("len(field.KeyName) == 0")
	}
	field.KeyNameFromCamelCase()
	field.KeyName = strings.Replace(field.KeyName, "-", "_", -1)
}

// SetFlagName read tag keyword option and save the text
func (field *Field) SetFlagName() {
	if len(field.FlagName) == 0 {
		panic("len(field.FlagName) == 0")
	}
	field.FlagName = Capitalize(field.FlagName)
	field.FlagNameFromCamelCase()
}

// SetType read tag keyword option and save the text
func (field *Field) SetType() {
	// field.Type = field.StructField.Name
}
