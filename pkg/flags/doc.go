// Package flags provides enhanced flag handling with custom types, automatic environment
// variable defaults, and consistent naming conventions.
//
// # Overview
//
// This package extends pflag with:
//   - Custom flag types implementing pflag.Value interface
//   - Automatic environment variable defaults using UPPER_SNAKE_CASE naming
//   - JSON/YAML friendly kebab-case flag names
//   - Type-safe validation and parsing
//   - Code generation utilities for common patterns
//
// # Basic Usage
//
//	type Config struct {
//		APIEndpoint *types.URLValue      `usage:"API endpoint URL"`
//		Timeout     *types.DurationValue `usage:"Request timeout"`
//		LogLevel    *types.LogLevelValue `usage:"Log level"`
//	}
//
//	config := &Config{
//		APIEndpoint: types.NewURLValue("https://api.example.com"),
//		Timeout:     types.NewDurationValue(30 * time.Second),
//		LogLevel:    types.NewLogLevelValue("info"),
//	}
//
//	flagConfig := flags.NewFlagConfig()
//	flagConfig.RegisterCustomTypes(config)
//	flagConfig.Parse(os.Args[1:])
//
// # Environment Variable Integration
//
// Field names are automatically converted to environment variables:
//   - APIEndpoint -> API_ENDPOINT
//   - DatabaseURL -> DATABASE_URL
//   - HTTPSProxy -> HTTPS_PROXY
//
// Flag names use kebab-case:
//   - APIEndpoint -> --api-endpoint
//   - DatabaseURL -> --database-url
//   - HTTPSProxy -> --https-proxy
//
// # Custom Types
//
// All custom types implement the pflag.Value interface:
//
//	type Value interface {
//		String() string
//		Set(string) error
//		Type() string
//	}
//
// Available types:
//   - URLValue: URL with validation
//   - DurationValue: Enhanced duration parsing
//   - HumanDurationValue: Human-friendly duration formats
//   - EnumValue: String enum with validation
//   - LogLevelValue: Specialized log level enum
//   - EnvironmentValue: Deployment environment enum
//   - KeyValueValue: Key-value pairs (key=value,...)
//   - HostPortValue: Host:port combinations
//   - IPValue: IP address validation
//   - CIDRValue: CIDR notation support
//   - PortRangeValue: Port range parsing
//
// # Manual Registration
//
// For more control, register flags manually:
//
//	flagSet := pflag.NewFlagSet("myapp", pflag.ExitOnError)
//
//	var serverURL *types.URLValue
//	flags.RegisterURLFlag(flagSet, "ServerURL", &serverURL, "Server endpoint")
//
//	var logLevel *types.LogLevelValue
//	validLevels := []string{"debug", "info", "warn", "error"}
//	flags.RegisterEnumFlag(flagSet, "LogLevel", &logLevel, validLevels, "info", "Log level")
//
// # Code Generation
//
// Generate custom types using templates:
//
//	tg := generators.NewTypeGenerator("types")
//	tg.RegisterDefaultTemplates()
//
//	code, err := tg.GenerateEnumType("Priority",
//		[]string{"low", "medium", "high"}, "medium")
//
// # Validation
//
// Add custom validators for any field:
//
//	flagConfig.AddValidator("Port", func(value interface{}) error {
//		if v, ok := value.(*int); ok && *v <= 0 {
//			return fmt.Errorf("port must be positive")
//		}
//		return nil
//	})
//
// # JSON/YAML Marshaling
//
// Most custom types support JSON and YAML marshaling:
//
//	type Config struct {
//		Timeout *types.DurationValue `json:"timeout" yaml:"timeout"`
//		Server  *types.URLValue      `json:"server" yaml:"server"`
//	}
//
// Duration types marshal differently in JSON vs YAML:
//   - JSON: numeric seconds
//   - YAML: human-readable strings
//
// # Thread Safety
//
// All types are designed to be thread-safe for concurrent read access.
// Concurrent writes should be synchronized by the caller.
//
// # Error Handling
//
// Types provide detailed error messages with suggestions:
//
//	invalid log level 'verbos', valid options: debug, info, warn, error, fatal
//	invalid URL 'not-a-url': parse "not-a-url": invalid URI for request
//	port number out of range (1-65535): 70000
//
// # Examples
//
// See the examples/ directory for complete working examples:
//   - basic_example.go: Automatic registration with struct tags
//   - manual_registration.go: Manual flag registration with validation
//   - generator_example.go: Code generation utilities
//   - test_types.go: Comprehensive type testing
package flags
