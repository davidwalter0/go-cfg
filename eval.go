/*
Package cfg provides modern pflag-based configuration management for Go applications.
It uses reflection to process struct fields as configuration elements, supporting
environment variables, command-line flags, and default values through struct tags.

Configuration Modes:

The library offers several configuration processing modes:

  - FlatNest(ptrs...): Clean nested flags without top-level struct prefix (recommended)
  - Eval(ptrs...): Standard mode with struct name prefixing
  - Init(ptrs...): Evaluation with automatic freeze after processing
  - Simple(ptrs...): No prefixing, environment variables only
  - SimpleFlags(ptrs...): No prefixing with environment variables and flags
  - Nest(ptrs...): Hierarchical prefixing using struct field names
  - Wrap(prefix, ptrs...): Custom prefix for all flags and environment variables

Struct Tag Support:

Configuration behavior is controlled through struct tags:

  - json:"name" - Override field name for flags and environment variables
  - default:"value" - Set default value if not provided via env or flags
  - doc:"description" - Set field documentation for help text
  - required:"true" - Mark field as required (validation enforced)
  - ignore:"true" - Skip field during configuration processing

Type Support:

All basic Go types are supported along with complex types:

  - Primitive types: string, int, bool, float64, time.Duration
  - Aggregate types: []string (comma-delimited), map[string]T (key:value,key:value)
  - Nested structs: Processed recursively with appropriate prefixing

Example Usage:

	type Config struct {
	    ServiceName string `json:"service_name" default:"myapp" doc:"Service identifier"`
	    Port        int    `json:"port" default:"8080" doc:"HTTP server port"`
	    Debug       bool   `default:"false" doc:"Enable debug logging"`
	    API         struct {
	        Host string `json:"host" default:"localhost" doc:"API host"`
	        Port int    `json:"port" default:"9090" doc:"API port"`
	    } `json:"api"`
	    Features map[string]bool `doc:"Feature toggles (name:true,other:false)"`
	}

	var config Config
	if err := cfg.FlatNest(&config); err != nil {
	    log.Fatal(err)
	}

Generated Flags:

# This generates clean flags like --service-name, --port, --debug, --api-host, --api-port

Environment Variables:

Environment variables use UPPER_SNAKE_CASE naming:
SERVICE_NAME, PORT, DEBUG, API_HOST, API_PORT

Precedence Order:

Configuration values are applied in this order (later values override earlier):
 1. Struct default tags
 2. Environment variables
 3. Command-line flags

Migration Note:

This version uses modern pflag internally while maintaining full API compatibility
with previous go-flag based versions.
*/
package cfg

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

// decorate var setting, disable prefixes when undecorated
var decorate bool = true

// Arg is the settings passed for recursion
type Arg struct {
	Depth          int
	Prefix         string
	Prefixed       bool
	UseFlags       bool
	SkipRootPrefix bool // Skip top-level struct name in prefixing
}

// NewArg sets values in Arg for cfg to process structs
func NewArg(name string) *Arg {
	var prefixed bool
	var prefix string
	var text string
	var ok bool
	if text, ok = LookupEnv(cfgDecorate); ok {
		if v, err := strconv.ParseBool(text); err != nil {
			log.Println(err)
		} else {
			decorate = v
		}
	}
	if decorate {
		prefix, ok = LookupEnv(cfgEnvKeyPrefix)
		if ok && len(prefix) > 0 {
			prefix = prefix + "_"
			prefixed = true
		}
		if len(name) > 0 {
			prefixed = true
			prefix = prefix + name
		}
	}
	return &Arg{Depth: 0, Prefixed: prefixed, Prefix: prefix, UseFlags: true}
}

// Undecorate structs with prefix
func Undecorate() bool {
	decorate = false
	os.Setenv(cfgDecorate, "false")
	return decorate
}

// Decorate structs with prefix
func Decorate() bool {
	decorate = true
	os.Setenv(cfgDecorate, "true")
	return decorate
}

func Unprefix() {
	os.Unsetenv(cfgEnvKeyPrefix)
}

// Eval processes one or more configuration structures using the standard mode.
// This mode includes struct name prefixing for flags and environment variables,
// making it suitable when you need clear namespace separation between different
// configuration structs.
//
// Flag Naming (with struct name prefix):
//   - For struct "Config" with field "Port": --config-port
//   - Environment variable: CONFIG_PORT
//
// This function does not freeze flag parsing, allowing additional configuration
// structs to be processed in subsequent calls. Use Init() for automatic freezing.
//
// Example:
//
//	type ServerConfig struct {
//	    Port int `default:"8080"`
//	}
//	var config ServerConfig
//	if err := cfg.Eval(&config); err != nil {
//	    return err
//	}
//	// Generates: --server-config-port, SERVER_CONFIG_PORT
//
// Parameters:
//
//	ptrs: One or more pointers to structs to be configured
//
// Returns:
//
//	error: nil on success, error if struct processing fails
func Eval(ptrs ...interface{}) error {
	args := NewArg("")
	return Run(args, ptrs...)
}

// Init processes configuration structures and automatically freezes flag parsing.
// This is a convenience function that combines Eval() with Freeze(), making it
// suitable for one-time configuration initialization where no further flag
// processing is needed.
//
// The function uses the same struct name prefixing as Eval() but ensures that
// flag parsing is finalized after processing, preventing accidental modification
// of the configuration system.
//
// Use this function when:
//   - You have all configuration structs ready to process at once
//   - You don't need to add more configuration structs later
//   - You want automatic flag parsing finalization
//
// Example:
//
//	type Config struct {
//	    Port int `default:"8080"`
//	    Host string `default:"localhost"`
//	}
//	var config Config
//	if err := cfg.Init(&config); err != nil {
//	    log.Fatal(err)
//	}
//	// Configuration is now frozen and ready to use
//
// Parameters:
//
//	ptrs: One or more pointers to structs to be configured
//
// Returns:
//
//	error: nil on success, error if struct processing or flag parsing fails
func Init(ptrs ...interface{}) error {
	args := NewArg("")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// EvalName one or more configuration structures overriding the name
func EvalName(name string, ptrs ...interface{}) error {
	args := NewArg(name)
	err := Run(args, ptrs...)
	Freeze()
	return err
}

func Run(args *Arg, ptrs ...interface{}) error {
	var err error
	err = CheckArgs(ptrs...)
	if err != nil {

		os.Exit(1)
	}
	for _, ptr := range ptrs {
		err = Enter(args, ptr)
		if err != nil {
			return err
		}
	}
	return err
}

// SimpleFlags don't prefix with base struct prefix, and they are
// undecorated.
//
// type T struct {
//    I int
// }
// var ptr = &T{}
// Default
// Eval(ptr) is decorated/prefixed with the struct name
// type T struct {
// ...
//    I int
// }
// type T struct {
// Eval
// EvalName("D", ptr) is decorated/prefixed with flag "--d-i"
// EvalName("D", ptr) is decorated/prefixed with env var "D_I"
//
// export CFG_KEY_PREFIX=KP
// Eval(p) is decorated/prefixed with flag "--kp-d-i"
// EvalName("D", p) is decorated/prefixed with env var "KP_D_I"
//
// export CFG_KEY_PREFIX=KP
// export CFG_DECORATE=false
// ignore prefix and struct name prefixing
// Eval(p) is decorated/prefixed with flag "--kp-d-i"
// EvalName("D", p) is decorated/prefixed with env var "KP_D_I"

// Simple processes configuration structures without any prefixing or decoration.
// This mode creates clean, unprefixed environment variables and is suitable for
// simple applications where you want direct field-to-environment mapping without
// namespace complexity.
//
// Flag Naming (no prefixes):
//   - Field "Port" becomes environment variable: PORT
//   - No command-line flags are created in this mode
//   - No struct name prefixing applied
//
// This function automatically disables prefixing and decoration, making it ideal
// for scenarios where you want the simplest possible environment variable mapping.
// Use SimpleFlags() if you also want command-line flag support.
//
// Example:
//
//	type Config struct {
//	    Port int    `default:"8080"`
//	    Host string `default:"localhost"`
//	}
//	var config Config
//	if err := cfg.Simple(&config); err != nil {
//	    return err
//	}
//	// Creates environment variables: PORT, HOST
//	// No command-line flags created
//
// Parameters:
//
//	ptrs: One or more pointers to structs to be configured
//
// Returns:
//
//	error: nil on success, error if struct processing fails
func Simple(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	args := NewArg("")
	return Run(args, ptrs...)
}

// SimpleFlags create env vars and flags without prefices
func SimpleFlags(ptrs ...interface{}) error {
	err := Simple(ptrs...)
	Freeze()
	return err
}

// Unwrap alias of Simple create env vars without prefices
func Unwrap(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	args := NewArg("")
	err := Run(args, ptrs...)
	return err
}

// Final freezes calling flag.Parse, no more additions to the
// configuration after Final
func Final() {
	Freeze()
}

// Bare alias of Simple create env vars without prefices
func Bare(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	args := NewArg("")
	err := Run(args, ptrs...)
	return err
}

// Wrap (optionally) with prefix and struct names create env vars without prefices
func Wrap(name string, ptrs ...interface{}) error {
	args := NewArg(name)
	err := Run(args, ptrs...)
	return err
}

// Add alias of Eval
func Add(ptrs ...interface{}) error {
	args := NewArg("")
	return Run(args, ptrs...)
}

var Package = func() string {
	type Empty struct{}
	return reflect.TypeOf(Empty{}).PkgPath()
}()

// Flags alias of Eval
func Flags(ptrs ...interface{}) error {
	args := NewArg("")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// Nest uses struct names recursively to prefix each flag or
// environment variable name as it descends the configuration object
// heirarchy. For exqmple, type B struct{ I int }; type A struct{ B }
// will create a flag with prefixes a and b, i.e. --a-b-i
// --a-b-i int
//   	 Env _A_B_I                           : (I) (int)

func Nest(ptrs ...interface{}) error {
	args := NewArg("")
	args.Prefixed = true
	err := Run(args, ptrs...)
	return err
}

// NestWrap objects retaining object hierarchy with prefix
func NestWrap(prefix string, ptrs ...interface{}) error {
	// fmt.Println("NestWrap", prefix)
	// fmt.Println(ptrs...)
	args := NewArg(prefix)
	args.Prefixed = true
	err := Run(args, ptrs...)
	return err
}

// FlatNest provides clean nested flags without the top-level struct name prefix.
// This is the recommended configuration mode for most applications as it generates
// intuitive flag names while maintaining structure for nested fields.
//
// Flag Naming:
//   - Top-level fields: --field-name
//   - Nested fields: --parent-field-name
//   - No top-level struct prefix
//
// Example:
//
//	type Config struct {
//	    ServiceName string `json:"service_name"`
//	    API struct {
//	        Host string `json:"host"`
//	        Port int    `json:"port"`
//	    } `json:"api"`
//	}
//
// Generates flags: --service-name, --api-host, --api-port
// Environment vars: SERVICE_NAME, API_HOST, API_PORT
//
// The function automatically freezes flag parsing after processing, preventing
// further configuration changes. Use Thaw() to allow additional modifications.
//
// Parameters:
//
//	ptrs: One or more pointers to structs to be configured
//
// Returns:
//
//	error: nil on success, error if struct processing fails or flag conflicts occur
func FlatNest(ptrs ...interface{}) error {
	args := NewArg("")
	args.Prefixed = true
	args.SkipRootPrefix = true // Skip the top-level struct name in prefixing
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// // https://graphemica.com/%E2%9C%93
// // https://graphemica.com/%E2%9C%97
// log.Println("\n\u2713 text\n\u2716 text\n")

// CheckArgs validate that the pointers are pointers to struct
func CheckArgs(ptrs ...interface{}) error {
	checkArgs := []string{}
	ok := true
	for i, ptr := range ptrs {
		typeOf := reflect.Indirect(reflect.ValueOf(ptr)).Type()
		kind := reflect.TypeOf(ptr).Kind()
		switch kind {
		case reflect.Ptr:
			elem := reflect.TypeOf(ptr).Elem()
			isPtr := reflect.Ptr == elem.Kind()
			if isPtr {
				ok = false
				checkArgs = append(checkArgs, fmt.Sprintf("%s !ok %2d [%v] is a pointer to %v.", Red("\u2716"), i, reflect.TypeOf(ptr), typeOf))
			} else {
				checkArgs = append(checkArgs, fmt.Sprintf("%s  ok %2d [%v] is a pointer to %v. expected.", Green("\u2713"), i, reflect.TypeOf(ptr), typeOf))
			}
		default:
			ok = false
			name := reflect.TypeOf(ptr).Name()
			checkArgs = append(checkArgs, fmt.Sprintf("%s !ok %2d [%v] is not a pointer.", Red("\u2716"), i, name))
		}
	}
	var err error
	var str string
	if !ok {
		for _, text := range checkArgs {
			str += text
		}
		str += fmt.Sprintf("\n%s\n\nThe argument list ptrs []interface{} is a slice of struct pointers (*struct) \n\n", Red(ErrInvalidArgPointerRequired))
		// fmt.Println(CallerAndArgs(ptrs))
		str += fmt.Sprintln(Caller())
		err = fmt.Errorf("%s", str)
	}
	return err
}

// Caller get the calling frame info
func Caller() string {
	text := ""
	{
		var pcs [10]uintptr
		n := runtime.Callers(2, pcs[:])
		frames := runtime.CallersFrames(pcs[:n])
		var frame runtime.Frame
		var more bool = true
		for i := 0; i < n-2 && more; i++ {
			frame, more = frames.Next()
			text += fmt.Sprintf("%s:%d:%s\n", frame.File, frame.Line, frame.Func.Name())
		}
	}
	return text
}
