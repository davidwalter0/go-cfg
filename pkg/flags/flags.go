package flags

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

// FlagConfig holds configuration for flag registration
type FlagConfig struct {
	FlagSet     *pflag.FlagSet
	EnvPrefix   string
	Required    map[string]bool
	HelpTexts   map[string]string
	Validators  map[string]func(interface{}) error
}

// NewFlagConfig creates a new flag configuration
func NewFlagConfig() *FlagConfig {
	return &FlagConfig{
		FlagSet:    pflag.NewFlagSet("", pflag.ContinueOnError),
		Required:   make(map[string]bool),
		HelpTexts:  make(map[string]string),
		Validators: make(map[string]func(interface{}) error),
	}
}

// RegisterCustomTypes registers all custom flag types with automatic env defaults
func (fc *FlagConfig) RegisterCustomTypes(structPtr interface{}) error {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got %T", structPtr)
	}
	
	structValue := v.Elem()
	structType := structValue.Type()
	nc := &types.NamingConventions{}
	
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)
		
		// Skip unexported fields
		if !field.CanSet() {
			continue
		}
		
		fieldName := fieldType.Name
		flagName := nc.ToKebabCase(fieldName)
		envName := nc.ToUpperSnakeCase(fieldName)
		
		// Get tag information
		usage := fieldType.Tag.Get("usage")
		if usage == "" {
			usage = fieldType.Tag.Get("help")
		}
		if usage == "" {
			usage = fmt.Sprintf("%s configuration", fieldName)
		}
		
		defaultValue := fieldType.Tag.Get("default")
		required := fieldType.Tag.Get("required") == "true"
		
		// Check for environment variable
		envValue := os.Getenv(envName)
		if envValue != "" && defaultValue == "" {
			defaultValue = envValue
		}
		
		// Register based on field type
		if err := fc.registerFieldType(field, fieldName, flagName, usage, defaultValue, required); err != nil {
			return fmt.Errorf("failed to register field %s: %w", fieldName, err)
		}
		
		if required {
			fc.Required[flagName] = true
		}
	}
	
	return nil
}

// registerFieldType registers a specific field type with pflag
func (fc *FlagConfig) registerFieldType(field reflect.Value, fieldName, flagName, usage, defaultValue string, required bool) error {
	fieldType := field.Type()
	
	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
		
		// Initialize the pointer if it's nil
		if field.IsNil() {
			field.Set(reflect.New(fieldType))
		}
		field = field.Elem()
	}
	
	switch {
	// Custom types implementing pflag.Value interface
	case field.Type().String() == "*types.URLValue":
		if field.IsNil() {
			field.Set(reflect.ValueOf(types.NewURLValue(defaultValue)))
		}
		fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
		
	case field.Type().String() == "*types.DurationValue":
		var defaultDuration time.Duration
		if defaultValue != "" {
			if d, err := time.ParseDuration(defaultValue); err == nil {
				defaultDuration = d
			}
		}
		if field.IsNil() {
			field.Set(reflect.ValueOf(types.NewDurationValue(defaultDuration)))
		}
		fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
		
	case field.Type().String() == "*types.EnumValue":
		// For EnumValue, we need the valid values from tags
		// This would need to be initialized with valid values before registration
		if !field.IsNil() {
			fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
		}
		
	case field.Type().String() == "*types.KeyValueValue":
		if field.IsNil() {
			field.Set(reflect.ValueOf(types.NewKeyValueValue()))
		}
		fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
		
	case field.Type().String() == "*types.HostPortValue":
		if field.IsNil() {
			field.Set(reflect.ValueOf(types.NewHostPortValue("localhost", 0)))
		}
		fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
		
	case field.Type().String() == "*types.IPValue":
		if field.IsNil() {
			field.Set(reflect.ValueOf(types.NewIPValue(defaultValue)))
		}
		fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
		
	// Standard types with environment defaults
	case fieldType.Kind() == reflect.String:
		if defaultValue == "" {
			defaultValue = field.String()
		}
		fc.FlagSet.StringVar(field.Addr().Interface().(*string), flagName, defaultValue, usage)
		
	case fieldType.Kind() == reflect.Int:
		defaultInt := 0
		if defaultValue != "" {
			if i, err := strconv.Atoi(defaultValue); err == nil {
				defaultInt = i
			}
		} else {
			defaultInt = int(field.Int())
		}
		fc.FlagSet.IntVar(field.Addr().Interface().(*int), flagName, defaultInt, usage)
		
	case fieldType.Kind() == reflect.Bool:
		defaultBool := false
		if defaultValue != "" {
			if b, err := strconv.ParseBool(defaultValue); err == nil {
				defaultBool = b
			}
		} else {
			defaultBool = field.Bool()
		}
		fc.FlagSet.BoolVar(field.Addr().Interface().(*bool), flagName, defaultBool, usage)
		
	case fieldType == reflect.TypeOf(time.Duration(0)):
		defaultDuration := time.Duration(0)
		if defaultValue != "" {
			if d, err := time.ParseDuration(defaultValue); err == nil {
				defaultDuration = d
			}
		} else {
			defaultDuration = time.Duration(field.Int())
		}
		fc.FlagSet.DurationVar(field.Addr().Interface().(*time.Duration), flagName, defaultDuration, usage)
		
	// Slice types
	case fieldType.Kind() == reflect.Slice:
		elemType := fieldType.Elem()
		switch elemType.Kind() {
		case reflect.String:
			var defaultSlice []string
			if defaultValue != "" {
				defaultSlice = strings.Split(defaultValue, ",")
			}
			fc.FlagSet.StringSliceVar(field.Addr().Interface().(*[]string), flagName, defaultSlice, usage)
			
		case reflect.Int:
			var defaultSlice []int
			if defaultValue != "" {
				parts := strings.Split(defaultValue, ",")
				for _, part := range parts {
					if i, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
						defaultSlice = append(defaultSlice, i)
					}
				}
			}
			fc.FlagSet.IntSliceVar(field.Addr().Interface().(*[]int), flagName, defaultSlice, usage)
			
		default:
			return fmt.Errorf("unsupported slice type: %v", fieldType)
		}
		
	default:
		return fmt.Errorf("unsupported field type: %v", fieldType)
	}
	
	return nil
}

// Parse parses command line arguments and validates required fields
func (fc *FlagConfig) Parse(args []string) error {
	if err := fc.FlagSet.Parse(args); err != nil {
		return err
	}
	
	// Validate required fields
	for flagName, isRequired := range fc.Required {
		if isRequired {
			flag := fc.FlagSet.Lookup(flagName)
			if flag == nil {
				return fmt.Errorf("required flag not found: %s", flagName)
			}
			
			if !flag.Changed {
				// Check if there's a default value or env var
				if flag.DefValue == "" {
					return fmt.Errorf("required flag not provided: --%s", flagName)
				}
			}
		}
	}
	
	// Run custom validators
	for flagName, validator := range fc.Validators {
		flag := fc.FlagSet.Lookup(flagName)
		if flag != nil && flag.Changed {
			if err := validator(flag.Value); err != nil {
				return fmt.Errorf("validation failed for flag --%s: %w", flagName, err)
			}
		}
	}
	
	return nil
}

// AddValidator adds a custom validator for a specific field
func (fc *FlagConfig) AddValidator(fieldName string, validator func(interface{}) error) {
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)
	fc.Validators[flagName] = validator
}

// SetRequired marks a field as required
func (fc *FlagConfig) SetRequired(fieldName string) {
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)
	fc.Required[flagName] = true
}

// Usage returns the usage string for all registered flags
func (fc *FlagConfig) Usage() string {
	return fc.FlagSet.FlagUsages()
}

// Helper functions for common flag registration patterns

// RegisterURLFlag registers a URL flag with environment default
func RegisterURLFlag(flagSet *pflag.FlagSet, fieldName string, target **types.URLValue, usage string) {
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)
	
	if *target == nil {
		*target = types.NewURLValue("")
	}
	
	// Try to get default from environment
	if envDefault := types.GetEnvDefault(fieldName); envDefault != "" {
		(*target).Set(envDefault)
	}
	
	flagSet.Var(*target, flagName, usage)
}

// RegisterEnumFlag registers an enum flag with validation
func RegisterEnumFlag(flagSet *pflag.FlagSet, fieldName string, target **types.EnumValue, validValues []string, defaultValue, usage string) {
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)
	
	if *target == nil {
		*target = types.NewEnumValue(validValues, defaultValue, false)
	}
	
	// Try to get default from environment
	if envDefault := types.GetEnvDefault(fieldName); envDefault != "" {
		(*target).Set(envDefault)
	}
	
	flagSet.Var(*target, flagName, usage)
}

// RegisterDurationFlag registers a duration flag with environment default
func RegisterDurationFlag(flagSet *pflag.FlagSet, fieldName string, target **types.DurationValue, defaultValue time.Duration, usage string) {
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)
	
	if *target == nil {
		*target = types.NewDurationValue(defaultValue)
	}
	
	// Try to get default from environment
	if envDefault := types.GetEnvDefault(fieldName); envDefault != "" {
		(*target).Set(envDefault)
	}
	
	flagSet.Var(*target, flagName, usage)
}

// RegisterKeyValueFlag registers a key-value flag with environment default
func RegisterKeyValueFlag(flagSet *pflag.FlagSet, fieldName string, target **types.KeyValueValue, usage string) {
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(fieldName)
	
	if *target == nil {
		*target = types.NewKeyValueValue()
	}
	
	// Try to get default from environment
	if envDefault := types.GetEnvDefault(fieldName); envDefault != "" {
		(*target).Set(envDefault)
	}
	
	flagSet.Var(*target, flagName, usage)
}