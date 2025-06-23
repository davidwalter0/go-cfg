package flags

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// Processor implements FlagProcessor interface
type Processor struct {
	naming NamingStrategy
}

// NewProcessor creates a new flag processor
func NewProcessor() *Processor {
	return &Processor{
		naming: NewFlatNestStrategy(),
	}
}

// ProcessStruct processes a struct and generates flags
func (p *Processor) ProcessStruct(config interface{}, flagSet *pflag.FlagSet) error {
	return p.processValue(reflect.ValueOf(config), reflect.TypeOf(config), []string{}, flagSet)
}

// GenerateFlags generates a flag set from a struct
func (p *Processor) GenerateFlags(config interface{}) (*pflag.FlagSet, error) {
	flagSet := pflag.NewFlagSet("config", pflag.ContinueOnError)

	// First apply defaults
	if err := p.applyDefaults(reflect.ValueOf(config), reflect.TypeOf(config), []string{}); err != nil {
		return nil, err
	}

	if err := p.ProcessStruct(config, flagSet); err != nil {
		return nil, err
	}

	return flagSet, nil
}

// ApplyFlagValues applies flag values back to the config struct
func (p *Processor) ApplyFlagValues(config interface{}, flagSet *pflag.FlagSet) error {
	return p.applyValues(reflect.ValueOf(config), reflect.TypeOf(config), []string{}, flagSet)
}

// applyDefaults applies default values from struct tags
func (p *Processor) applyDefaults(val reflect.Value, typ reflect.Type, path []string) error {
	// Handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() || !field.CanSet() {
			continue
		}

		fieldPath := append(path, fieldType.Name)

		// Handle nested structs
		if field.Kind() == reflect.Struct && fieldType.Type != reflect.TypeOf(time.Duration(0)) {
			if err := p.applyDefaults(field, fieldType.Type, fieldPath); err != nil {
				return err
			}
			continue
		}

		// Apply default value if present
		if defaultValue := fieldType.Tag.Get("default"); defaultValue != "" {
			if err := p.setDefaultValue(field, fieldType.Type, defaultValue); err != nil {
				return err
			}
		}
	}

	return nil
}

// setDefaultValue sets a default value on a field
func (p *Processor) setDefaultValue(field reflect.Value, fieldType reflect.Type, defaultValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(defaultValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldType == reflect.TypeOf(time.Duration(0)) {
			if d, err := time.ParseDuration(defaultValue); err == nil {
				field.SetInt(int64(d))
			}
		} else {
			if i, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
				field.SetInt(i)
			}
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(defaultValue); err == nil {
			field.SetBool(b)
		}
	}
	return nil
}

// processValue recursively processes struct fields to generate flags
func (p *Processor) processValue(val reflect.Value, typ reflect.Type, path []string, flagSet *pflag.FlagSet) error {
	// Handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		fieldPath := append(path, fieldType.Name)

		// Handle nested structs
		if field.Kind() == reflect.Struct && fieldType.Type != reflect.TypeOf(time.Duration(0)) {
			if err := p.processValue(field, fieldType.Type, fieldPath, flagSet); err != nil {
				return err
			}
			continue
		}

		// Generate flag for this field
		if err := p.addFlag(field, fieldType, fieldPath, flagSet); err != nil {
			return err
		}
	}

	return nil
}

// addFlag adds a single flag to the flag set
func (p *Processor) addFlag(field reflect.Value, fieldType reflect.StructField, path []string, flagSet *pflag.FlagSet) error {
	flagName := p.naming.FlagName(path)
	envName := p.naming.EnvVarName(path)

	// Get default value and help text from struct tags
	defaultValue := fieldType.Tag.Get("default")
	helpText := fieldType.Tag.Get("help")
	if helpText == "" {
		helpText = fmt.Sprintf("(%s)", strings.ToLower(fieldType.Name))
	}

	// Add environment variable info to help
	if envName != "" {
		helpText = fmt.Sprintf("Env %s : %s", envName, helpText)
	}

	// Add flag based on field type
	switch field.Kind() {
	case reflect.String:
		flagSet.String(flagName, defaultValue, helpText)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldType.Type == reflect.TypeOf(time.Duration(0)) {
			// Handle time.Duration specially
			var defaultDuration time.Duration
			if defaultValue != "" {
				if d, err := time.ParseDuration(defaultValue); err == nil {
					defaultDuration = d
				}
			}
			flagSet.Duration(flagName, defaultDuration, helpText)
		} else {
			defaultInt := int64(0)
			if defaultValue != "" {
				if i, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
					defaultInt = i
				}
			}
			flagSet.Int64(flagName, defaultInt, helpText)
		}
	case reflect.Bool:
		defaultBool := false
		if defaultValue == "true" {
			defaultBool = true
		}
		flagSet.Bool(flagName, defaultBool, helpText)
	case reflect.Map:
		// Handle maps as string slices for now
		flagSet.StringSlice(flagName, nil, helpText+" (key:value pairs)")
	case reflect.Slice:
		flagSet.StringSlice(flagName, nil, helpText+" (comma-separated)")
	default:
		// Skip unsupported types
		return nil
	}

	return nil
}

// applyValues applies flag values back to struct fields
func (p *Processor) applyValues(val reflect.Value, typ reflect.Type, path []string, flagSet *pflag.FlagSet) error {
	// Handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() || !field.CanSet() {
			continue
		}

		fieldPath := append(path, fieldType.Name)

		// Handle nested structs
		if field.Kind() == reflect.Struct && fieldType.Type != reflect.TypeOf(time.Duration(0)) {
			if err := p.applyValues(field, fieldType.Type, fieldPath, flagSet); err != nil {
				return err
			}
			continue
		}

		// Apply flag value to this field
		if err := p.applyFlag(field, fieldType, fieldPath, flagSet); err != nil {
			return err
		}
	}

	return nil
}

// applyFlag applies a single flag value to a struct field
func (p *Processor) applyFlag(field reflect.Value, fieldType reflect.StructField, path []string, flagSet *pflag.FlagSet) error {
	flagName := p.naming.FlagName(path)

	// Check if flag was set
	flag := flagSet.Lookup(flagName)
	if flag == nil || !flag.Changed {
		return nil
	}

	// Apply value based on field type
	switch field.Kind() {
	case reflect.String:
		if val, err := flagSet.GetString(flagName); err == nil {
			field.SetString(val)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldType.Type == reflect.TypeOf(time.Duration(0)) {
			if val, err := flagSet.GetDuration(flagName); err == nil {
				field.SetInt(int64(val))
			}
		} else {
			if val, err := flagSet.GetInt64(flagName); err == nil {
				field.SetInt(val)
			}
		}
	case reflect.Bool:
		if val, err := flagSet.GetBool(flagName); err == nil {
			field.SetBool(val)
		}
	case reflect.Map:
		if vals, err := flagSet.GetStringSlice(flagName); err == nil && len(vals) > 0 {
			// Create map from key:value pairs
			mapVal := reflect.MakeMap(field.Type())
			for _, kv := range vals {
				parts := strings.SplitN(kv, ":", 2)
				if len(parts) == 2 {
					key := reflect.ValueOf(parts[0])
					// Handle different value types
					var val reflect.Value
					switch field.Type().Elem().Kind() {
					case reflect.String:
						val = reflect.ValueOf(parts[1])
					case reflect.Bool:
						if b, err := strconv.ParseBool(parts[1]); err == nil {
							val = reflect.ValueOf(b)
						} else {
							continue
						}
					default:
						val = reflect.ValueOf(parts[1])
					}
					mapVal.SetMapIndex(key, val)
				}
			}
			field.Set(mapVal)
		}
	case reflect.Slice:
		if vals, err := flagSet.GetStringSlice(flagName); err == nil && len(vals) > 0 {
			// Create slice
			sliceVal := reflect.MakeSlice(field.Type(), len(vals), len(vals))
			for i, val := range vals {
				sliceVal.Index(i).SetString(val)
			}
			field.Set(sliceVal)
		}
	}

	return nil
}
