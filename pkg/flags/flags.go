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

// StringSliceValue implements pflag.Value for string slices that preserve empty strings
type StringSliceValue struct {
	value []string
}

func NewStringSliceValue(defaultValue []string) *StringSliceValue {
	return &StringSliceValue{value: defaultValue}
}

func (ssv *StringSliceValue) String() string {
	return strings.Join(ssv.value, ",")
}

func (ssv *StringSliceValue) Set(s string) error {
	// If the string is empty, create a slice with one empty string
	if s == "" {
		ssv.value = []string{""}
		return nil
	}
	
	// Split by comma and preserve empty strings
	ssv.value = strings.Split(s, ",")
	return nil
}

func (ssv *StringSliceValue) Type() string {
	return "stringSlice"
}

func (ssv *StringSliceValue) GetSlice() []string {
	return ssv.value
}

// FlagConfig holds configuration for flag registration
type FlagConfig struct {
	FlagSet     *pflag.FlagSet
	EnvPrefix   string
	Required    map[string]bool
	HelpTexts   map[string]string
	Validators  map[string]func(interface{}) error
	sliceUpdaters map[string]func() // Functions to update slice fields after parsing
}

// NewFlagConfig creates a new flag configuration
func NewFlagConfig() *FlagConfig {
	return &FlagConfig{
		FlagSet:    pflag.NewFlagSet("", pflag.ContinueOnError),
		Required:   make(map[string]bool),
		HelpTexts:  make(map[string]string),
		Validators: make(map[string]func(interface{}) error),
		sliceUpdaters: make(map[string]func()),
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
		
		// Check for environment variable (env vars take precedence over defaults)
		envValue := os.Getenv(envName)
		if envValue != "" {
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
	originalFieldType := field.Type()
	
	// Get the actual type string for debugging (use original type before pointer handling)
	typeStr := originalFieldType.String()
	
	// Handle pointer types
	if originalFieldType.Kind() == reflect.Ptr {
		// Initialize the pointer if it's nil
		if field.IsNil() {
			// Create instances based on type
			switch typeStr {
			case "*types.URLValue":
				field.Set(reflect.ValueOf(types.NewURLValue(defaultValue)))
			case "*types.DurationValue":
				var defaultDuration time.Duration
				if defaultValue != "" {
					if d, err := time.ParseDuration(defaultValue); err == nil {
						defaultDuration = d
					}
				}
				field.Set(reflect.ValueOf(types.NewDurationValue(defaultDuration)))
			case "*types.HumanDurationValue":
				var defaultDuration time.Duration
				if defaultValue != "" {
					if d, err := time.ParseDuration(defaultValue); err == nil {
						defaultDuration = d
					}
				}
				field.Set(reflect.ValueOf(types.NewHumanDurationValue(defaultDuration)))
			case "*types.LogLevelValue":
				defaultLevel := "info"
				if defaultValue != "" {
					defaultLevel = defaultValue
				}
				field.Set(reflect.ValueOf(types.NewLogLevelValue(defaultLevel)))
			case "*types.EnvironmentValue":
				defaultEnv := "development"
				if defaultValue != "" {
					defaultEnv = defaultValue
				}
				field.Set(reflect.ValueOf(types.NewEnvironmentValue(defaultEnv)))
			case "*types.KeyValueValue":
				kv := types.NewKeyValueValue()
				if defaultValue != "" {
					kv.Set(defaultValue)
				}
				field.Set(reflect.ValueOf(kv))
			case "*types.LabelsValue":
				lv := types.NewLabelsValue()
				if defaultValue != "" {
					lv.Set(defaultValue)
				}
				field.Set(reflect.ValueOf(lv))
			case "*types.HostPortValue":
				hp := types.NewHostPortValue("localhost", 0)
				if defaultValue != "" {
					hp.Set(defaultValue)
				}
				field.Set(reflect.ValueOf(hp))
			case "*types.IPValue":
				field.Set(reflect.ValueOf(types.NewIPValue(defaultValue)))
			case "*types.CIDRValue":
				field.Set(reflect.ValueOf(types.NewCIDRValue(defaultValue)))
			case "*types.PortRangeValue":
				pr := types.NewPortRangeValue(0, 0)
				if defaultValue != "" {
					pr.Set(defaultValue)
				}
				field.Set(reflect.ValueOf(pr))
			case "*types.URLSliceValue":
				us := types.NewURLSliceValue()
				if defaultValue != "" {
					us.Set(defaultValue)
				}
				field.Set(reflect.ValueOf(us))
			default:
				// For other pointer types, create a new instance
				elemType := originalFieldType.Elem()
				field.Set(reflect.New(elemType))
			}
		}
		
		// Check if the pointer type implements pflag.Value
		if originalFieldType.Implements(reflect.TypeOf((*pflag.Value)(nil)).Elem()) {
			// Apply default value if provided (including from environment variables)
			if defaultValue != "" {
				if setValue, ok := field.Interface().(interface{ Set(string) error }); ok {
					setValue.Set(defaultValue)
				}
			}
			fc.FlagSet.Var(field.Interface().(pflag.Value), flagName, usage)
			return nil
		}
		
		// Dereference for further processing
		field = field.Elem()
	}
	
	// Handle non-pointer types or types that don't implement pflag.Value
	fieldType := field.Type()
	switch {
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
			// Use custom StringSliceValue to preserve empty strings
			stringSliceValue := NewStringSliceValue(defaultSlice)
			fc.FlagSet.Var(stringSliceValue, flagName, usage)
			
			// Store reference to update the original field after parsing
			fc.storeSliceUpdater(flagName, field, stringSliceValue)
			
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
			// Check if it's a time.Duration slice
			if fieldType.Elem() == reflect.TypeOf(time.Duration(0)) {
				var defaultSlice []time.Duration
				if defaultValue != "" {
					parts := strings.Split(defaultValue, ",")
					for _, part := range parts {
						if d, err := time.ParseDuration(strings.TrimSpace(part)); err == nil {
							defaultSlice = append(defaultSlice, d)
						}
					}
				}
				fc.FlagSet.DurationSliceVar(field.Addr().Interface().(*[]time.Duration), flagName, defaultSlice, usage)
			} else {
				return fmt.Errorf("unsupported slice type: %v", fieldType)
			}
		}
		
	// Handle nested structs
	case fieldType.Kind() == reflect.Struct:
		// Recursively register nested struct fields
		return fc.registerNestedStruct(field, fieldName, flagName, usage)
		
	default:
		return fmt.Errorf("unsupported field type: %v", fieldType)
	}
	
	return nil
}

// registerNestedStruct recursively registers fields from nested structs
func (fc *FlagConfig) registerNestedStruct(structValue reflect.Value, parentName, parentFlagName, parentUsage string) error {
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
		// Create nested flag name: parent-field becomes parent-field-subfield
		nestedFlagName := parentFlagName + "-" + nc.ToKebabCase(fieldName)
		// Create nested env name: PARENT_FIELD_SUBFIELD
		nestedEnvName := nc.ToUpperSnakeCase(parentName + "_" + fieldName)
		
		// Get tag information
		usage := fieldType.Tag.Get("usage")
		if usage == "" {
			usage = fieldType.Tag.Get("help")
		}
		if usage == "" {
			usage = fmt.Sprintf("%s %s configuration", parentName, fieldName)
		}
		
		defaultValue := fieldType.Tag.Get("default")
		required := fieldType.Tag.Get("required") == "true"
		
		// Check for environment variable (env vars take precedence over defaults)
		envValue := os.Getenv(nestedEnvName)
		if envValue != "" {
			defaultValue = envValue
		}
		
		// Register the nested field using the same logic as top-level fields
		if err := fc.registerFieldType(field, parentName+"_"+fieldName, nestedFlagName, usage, defaultValue, required); err != nil {
			return fmt.Errorf("failed to register nested field %s.%s: %w", parentName, fieldName, err)
		}
		
		if required {
			fc.Required[nestedFlagName] = true
		}
	}
	
	return nil
}

// Parse parses command line arguments and validates required fields
func (fc *FlagConfig) Parse(args []string) error {
	if err := fc.FlagSet.Parse(args); err != nil {
		return err
	}
	
	// Update slice fields with parsed values
	for _, updater := range fc.sliceUpdaters {
		updater()
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

// storeSliceUpdater stores a function to update a slice field after parsing
func (fc *FlagConfig) storeSliceUpdater(flagName string, field reflect.Value, stringSliceValue *StringSliceValue) {
	fc.sliceUpdaters[flagName] = func() {
		// Update the original slice field with the parsed values
		field.Set(reflect.ValueOf(stringSliceValue.GetSlice()))
	}
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