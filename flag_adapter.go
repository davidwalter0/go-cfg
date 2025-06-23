package cfg

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// Global pflag.FlagSet to replace the vendored go-flag
var globalFlagSet *pflag.FlagSet

// Flag state tracking for compatibility
var flagNames = make(map[string]bool)
var flagValues = make(map[string]interface{})
var flagDefaults = make(map[string]string)
var flagRequired = make(map[string]bool)

func init() {
	globalFlagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
}

// MakeVar creates a flag with the same interface as the legacy go-flag
func MakeVar(ptr interface{}, name, defaultValue, usage string, value interface{}, required, isset bool) {
	// Store flag metadata
	flagNames[name] = true
	flagDefaults[name] = defaultValue
	flagRequired[name] = required

	// Get the type of the pointer target
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() != reflect.Ptr {
		panic("MakeVar requires a pointer")
	}

	elem := ptrVal.Elem()
	elemType := elem.Type()

	// Apply default value if provided and not already set
	if defaultValue != "" {
		if err := setValueFromString(elem, defaultValue); err != nil {
			fmt.Printf("Warning: failed to set default value for %s: %v\n", name, err)
		}
	}

	// Skip if flag already registered
	if globalFlagSet.Lookup(name) != nil {
		return
	}

	// Register the flag with pflag based on type
	switch elemType.Kind() {
	case reflect.String:
		stringPtr := ptr.(*string)
		globalFlagSet.StringVar(stringPtr, name, *stringPtr, usage)
		
	case reflect.Int:
		intPtr := ptr.(*int)
		globalFlagSet.IntVar(intPtr, name, *intPtr, usage)
		
	case reflect.Int64:
		// Handle time.Duration specially
		if elemType == reflect.TypeOf(time.Duration(0)) {
			durationPtr := ptr.(*time.Duration)
			globalFlagSet.DurationVar(durationPtr, name, *durationPtr, usage)
		} else {
			int64Ptr := ptr.(*int64)
			globalFlagSet.Int64Var(int64Ptr, name, *int64Ptr, usage)
		}
		
	case reflect.Bool:
		boolPtr := ptr.(*bool)
		globalFlagSet.BoolVar(boolPtr, name, *boolPtr, usage)
		
	case reflect.Float64:
		float64Ptr := ptr.(*float64)
		globalFlagSet.Float64Var(float64Ptr, name, *float64Ptr, usage)
		
	case reflect.Slice:
		// Handle string slices
		if elemType.Elem().Kind() == reflect.String {
			slicePtr := ptr.(*[]string)
			globalFlagSet.StringSliceVar(slicePtr, name, *slicePtr, usage)
		} else {
			// For other slice types, use string and parse manually
			stringPtr := new(string)
			globalFlagSet.StringVar(stringPtr, name, "", usage)
			// Store reference for later parsing
			flagValues[name] = ptr
		}
		
	case reflect.Map:
		// For maps, use string and parse manually
		stringPtr := new(string)
		globalFlagSet.StringVar(stringPtr, name, "", usage)
		// Store reference for later parsing
		flagValues[name] = ptr
		
	default:
		// For other types, use string and parse manually
		stringPtr := new(string)
		globalFlagSet.StringVar(stringPtr, name, defaultValue, usage)
		// Store reference for later parsing
		flagValues[name] = ptr
	}

	// Store the current value
	flagValues[name] = value
}

// Parse parses command line flags
func Parse() error {
	return globalFlagSet.Parse(os.Args[1:])
}

// IsSet returns whether a flag was set on the command line
func IsSet(name string) bool {
	return globalFlagSet.Changed(name)
}

// Ok returns whether a flag was set (alias for IsSet for compatibility)
func Ok(name string) bool {
	return IsSet(name)
}

// Required returns whether a required flag was set
func Required(name string) bool {
	if !flagRequired[name] {
		return true // Not required, so it's OK
	}
	return IsSet(name)
}

// PrintDefaults prints default values of all defined flags
func PrintDefaults() {
	globalFlagSet.PrintDefaults()
}

// NewFlagSet creates a new flag set with given name and error handling
func NewFlagSet(name string, errorHandling pflag.ErrorHandling) *pflag.FlagSet {
	globalFlagSet = pflag.NewFlagSet(name, errorHandling)
	// Clear state for new flag set
	flagNames = make(map[string]bool)
	flagValues = make(map[string]interface{})
	flagDefaults = make(map[string]string)
	flagRequired = make(map[string]bool)
	CommandLine = globalFlagSet
	return globalFlagSet
}

// Usage function type for compatibility
var Usage func()

// setValueFromString sets a reflect.Value from a string representation
func setValueFromString(val reflect.Value, str string) error {
	switch val.Kind() {
	case reflect.String:
		val.SetString(str)
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Type() == reflect.TypeOf(time.Duration(0)) {
			// Handle time.Duration
			if d, err := time.ParseDuration(str); err == nil {
				val.Set(reflect.ValueOf(d))
			} else {
				return err
			}
		} else {
			if i, err := strconv.ParseInt(str, 0, val.Type().Bits()); err == nil {
				val.SetInt(i)
			} else {
				return err
			}
		}
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i, err := strconv.ParseUint(str, 0, val.Type().Bits()); err == nil {
			val.SetUint(i)
		} else {
			return err
		}
		
	case reflect.Bool:
		if b, err := strconv.ParseBool(str); err == nil {
			val.SetBool(b)
		} else {
			return err
		}
		
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(str, val.Type().Bits()); err == nil {
			val.SetFloat(f)
		} else {
			return err
		}
		
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.String {
			// String slice
			parts := strings.Split(str, ",")
			slice := reflect.MakeSlice(val.Type(), len(parts), len(parts))
			for i, part := range parts {
				slice.Index(i).SetString(strings.TrimSpace(part))
			}
			val.Set(slice)
		}
		
	case reflect.Map:
		if val.Type().Key().Kind() == reflect.String {
			// String-keyed map
			mapVal := reflect.MakeMap(val.Type())
			if str != "" {
				pairs := strings.Split(str, ",")
				for _, pair := range pairs {
					kv := strings.SplitN(pair, ":", 2)
					if len(kv) == 2 {
						key := reflect.ValueOf(strings.TrimSpace(kv[0]))
						
						// Set value based on map value type
						valueType := val.Type().Elem()
						value := reflect.New(valueType).Elem()
						if err := setValueFromString(value, strings.TrimSpace(kv[1])); err == nil {
							mapVal.SetMapIndex(key, value)
						}
					}
				}
			}
			val.Set(mapVal)
		}
		
	default:
		return fmt.Errorf("unsupported type: %v", val.Type())
	}
	
	return nil
}

// Error handling constants for compatibility
const (
	ContinueOnError = pflag.ContinueOnError
	ExitOnError     = pflag.ExitOnError
	PanicOnError    = pflag.PanicOnError
)

// CommandLine provides access to the global flag set
var CommandLine = globalFlagSet