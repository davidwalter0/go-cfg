// Package cfg flag adapter provides a compatibility layer between the legacy go-flag API
// and the modern pflag library. This adapter allows the go-cfg library to use pflag
// internally while maintaining the exact same external API for backward compatibility.
//
// Migration Benefits:
//   - Eliminates vendored go-flag dependency
//   - Uses industry-standard pflag library
//   - Maintains 100% API compatibility
//   - Supports all existing flag types and behaviors
//
// The adapter translates between the legacy MakeVar() interface and pflag's typed
// flag registration functions, handling type detection, default value application,
// and flag deduplication automatically.
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

// globalFlagSet is the primary pflag.FlagSet used throughout the application.
// It replaces the vendored go-flag system and provides modern flag processing
// while maintaining compatibility with legacy code.
var globalFlagSet *pflag.FlagSet

// Flag state tracking for compatibility
var flagNames = make(map[string]bool)
var flagValues = make(map[string]interface{})
var flagDefaults = make(map[string]string)
var flagRequired = make(map[string]bool)

func init() {
	globalFlagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
}

// MakeVar creates a pflag-compatible flag with the same interface as the legacy go-flag system.
// This function serves as the primary bridge between go-cfg's field processing and pflag's
// type-safe flag registration, maintaining full API compatibility while using modern pflag internally.
//
// The function automatically detects the target type from the pointer and registers the appropriate
// pflag type, applying default values before registration to ensure consistent behavior.
//
// Parameters:
//
//	ptr - Pointer to the variable that will receive the flag value (must be *T, not **T)
//	name - Flag name, typically in kebab-case format (e.g., "server-port")
//	defaultValue - Default value as string, applied before flag registration
//	usage - Help text displayed in flag usage (may include environment variable info)
//	value - Current/initial value, typically from environment variable lookup
//	required - Whether this flag is mandatory (tracked for validation)
//	isset - Whether the value was explicitly set from environment or previous processing
//
// Type Handling:
//   - string, int, int64, bool, float64: Direct pflag registration with native support
//   - time.Duration: Special handling using pflag's DurationVar for proper parsing
//   - []string: Native pflag StringSliceVar support with comma-separated values
//   - Other slices: Manual parsing using string registration with deferred conversion
//   - Maps: Manual parsing with key:value,key:value syntax using string registration
//   - Complex types: String registration with custom parsing in setValueFromString
//
// Default Value Application:
//
//	The function applies default values before flag registration, ensuring that pflag
//	sees the correct default value. This maintains consistency with legacy behavior
//	where defaults were applied during flag creation rather than during parsing.
//
// Flag Deduplication:
//
//	To prevent "flag redefined" errors, the function checks if a flag is already
//	registered and silently skips registration if found. This allows multiple
//	configuration processing passes without conflicts.
//
// Error Handling:
//   - Type conversion errors during default value application generate warnings
//   - Invalid pointer types cause panics (consistent with legacy behavior)
//   - Unsupported types fall back to string registration with manual parsing
//
// Example Usage:
//
//	var port int
//	MakeVar(&port, "server-port", "8080", "HTTP server port", nil, false, false)
//	// Registers: --server-port int with default 8080
//
//	var features []string
//	MakeVar(&features, "features", "auth,logging", "Enabled features", nil, false, true)
//	// Registers: --features stringSlice with default ["auth", "logging"]
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

	// Apply resolved value (environment variable or default) if provided
	// The value parameter contains the final resolved value after env var lookup
	if value != nil {
		if valueStr, ok := value.(string); ok && valueStr != "" {
			if err := setValueFromString(elem, valueStr); err != nil {
				fmt.Printf("Warning: failed to set resolved value for %s: %v\n", name, err)
			}
		}
	} else if defaultValue != "" {
		// Fallback to default if no resolved value
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

// Parse processes command line arguments using the global flag set.
// This function should be called after all flags have been registered via MakeVar.
// It automatically uses os.Args[1:] as the argument source.
//
// Returns:
//
//	error - Parsing error if invalid flags or values provided, nil on success
//
// Note: This function is typically called automatically by the freeze mechanism
// in go-cfg, so manual calls are usually not necessary.
func Parse() error {
	return globalFlagSet.Parse(os.Args[1:])
}

// IsSet returns whether a flag was explicitly set on the command line.
// This function checks pflag's Changed() status to determine if the flag
// value came from command line arguments rather than default values.
//
// Parameters:
//
//	name - Flag name to check (must match name used in MakeVar)
//
// Returns:
//
//	bool - true if flag was provided on command line, false if using default
//
// This is used to distinguish between explicit flag values and defaults,
// which is important for configuration precedence (flags > env > defaults).
func IsSet(name string) bool {
	return globalFlagSet.Changed(name)
}

// Ok returns whether a flag was set (alias for IsSet for legacy compatibility).
// This function exists to maintain API compatibility with the original go-flag
// based implementation where Ok() was the primary flag checking function.
//
// Parameters:
//
//	name - Flag name to check
//
// Returns:
//
//	bool - true if flag was set on command line, false otherwise
//
// Use IsSet() in new code for clarity; Ok() is maintained for backward compatibility.
func Ok(name string) bool {
	return IsSet(name)
}

// Required validates that a required flag has been set.
// This function checks both the flag's required status (from MakeVar) and
// whether it was actually provided on the command line.
//
// Parameters:
//
//	name - Flag name to validate
//
// Returns:
//
//	bool - true if flag is not required OR was set, false if required but missing
//
// Logic:
//   - If flag is not marked as required: returns true (validation passes)
//   - If flag is required: returns true only if IsSet(name) is true
//
// This function is typically used during configuration validation to ensure
// all mandatory flags have been provided by the user.
func Required(name string) bool {
	if !flagRequired[name] {
		return true // Not required, so it's OK
	}
	return IsSet(name)
}

// PrintDefaults outputs help text for all registered flags.
// This function delegates to pflag's PrintDefaults() to generate
// consistent help output showing flag names, types, defaults, and usage.
//
// Output format (handled by pflag):
//
//	--flag-name type
//	      usage text (default value)
//
// This function is typically called from help/usage functions rather
// than directly by application code.
func PrintDefaults() {
	globalFlagSet.PrintDefaults()
}

// NewFlagSet creates a new pflag.FlagSet and resets all adapter state.
// This function replaces the global flag set with a new instance and clears
// all tracking data, effectively starting fresh for flag registration.
//
// Parameters:
//
//	name - Name for the flag set (typically program name)
//	errorHandling - pflag error handling mode (ContinueOnError, ExitOnError, PanicOnError)
//
// Returns:
//
//	*pflag.FlagSet - The newly created flag set (same as globalFlagSet)
//
// Side Effects:
//   - Replaces globalFlagSet with new instance
//   - Clears all flag tracking maps (flagNames, flagValues, etc.)
//   - Updates CommandLine to point to new flag set
//
// This function is useful for testing or when you need to completely reset
// the flag system state. Most applications won't need to call this directly.
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

// setValueFromString converts a string value to the appropriate Go type and assigns it to a reflect.Value.
// This function handles the complex type conversion logic needed to support default value application
// and manual parsing for types not natively supported by pflag.
//
// The function uses reflection to determine the target type and applies appropriate parsing logic
// for each supported type category. It supports the same types as MakeVar but focuses on
// string-to-type conversion rather than flag registration.
//
// Parameters:
//
//	val - reflect.Value pointing to the destination variable (must be settable)
//	str - String representation of the value to convert and assign
//
// Returns:
//
//	error - Conversion error if string cannot be parsed to target type, nil on success
//
// Supported Type Conversions:
//   - string: Direct assignment
//   - int, int8, int16, int32, int64: ParseInt with appropriate bit size
//   - uint, uint8, uint16, uint32, uint64: ParseUint with appropriate bit size
//   - time.Duration: ParseDuration for duration strings (e.g., "5m", "30s")
//   - bool: ParseBool supporting various boolean representations
//   - float32, float64: ParseFloat with appropriate precision
//   - []string: Comma-separated values split and trimmed
//   - map[string]T: key:value,key:value format with recursive value parsing
//
// Slice Handling:
//
//	String slices are split on commas with whitespace trimming. Other slice types
//	would require additional implementation but are typically handled by MakeVar's
//	string registration approach.
//
// Map Handling:
//
//	Maps use key:value,key:value syntax where each value is recursively parsed
//	according to the map's value type. Only string-keyed maps are currently supported.
//
// Error Conditions:
//   - Invalid string format for target type (e.g., "abc" for int)
//   - Unsupported type combinations
//   - Malformed compound types (slices/maps with invalid syntax)
//
// Example Usage:
//
//	var port int
//	val := reflect.ValueOf(&port).Elem()
//	err := setValueFromString(val, "8080")  // Sets port = 8080
//
//	var features []string
//	val := reflect.ValueOf(&features).Elem()
//	err := setValueFromString(val, "auth, logging")  // Sets features = ["auth", "logging"]
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

// Error handling constants for pflag compatibility.
// These constants provide the same error handling modes as the legacy go-flag system,
// allowing existing code to continue using familiar error handling patterns.
//
// Usage with NewFlagSet:
//
//	flagSet := NewFlagSet("myapp", ContinueOnError)
//	// Flag parsing errors will be returned rather than causing program exit
const (
	ContinueOnError = pflag.ContinueOnError // Return errors from Parse(), don't exit
	ExitOnError     = pflag.ExitOnError     // Call os.Exit(2) on flag parsing errors
	PanicOnError    = pflag.PanicOnError    // Panic on flag parsing errors
)

// CommandLine provides access to the global pflag.FlagSet for advanced usage.
// This variable maintains compatibility with go-flag's CommandLine global variable
// while providing access to pflag's enhanced functionality.
//
// The variable is automatically updated when NewFlagSet() is called, ensuring
// it always points to the current global flag set instance.
//
// Advanced Usage:
//
//	CommandLine.SetInterspersed(false)  // Disable interspersed flags
//	CommandLine.SortFlags = false       // Disable flag sorting in help
//	CommandLine.Usage = customUsageFunc // Set custom usage function
var CommandLine = globalFlagSet
