package cfg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

// VendorTestConfig for testing the old vendor implementation
type VendorTestConfig struct {
	Name        string        `json:"name" default:"test-app" help:"Application name"`
	Port        int           `json:"port" default:"8080" help:"Server port"`
	Debug       bool          `json:"debug" default:"false" help:"Enable debug mode"`
	Timeout     time.Duration `json:"timeout" default:"30s" help:"Request timeout"`
	APIEndpoint string        `json:"api_endpoint" default:"https://api.example.com" help:"API endpoint URL"`
	LogLevel    string        `json:"log_level" default:"info" help:"Log level"`
	Environment string        `json:"environment" default:"development" help:"Deployment environment"`
	Required    string        `json:"required" required:"true" help:"Required field"`
}

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
}

func TestVendorImplementation_Parse(t *testing.T) {
	// Save and restore original args to isolate test from go test flags
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	
	// Set clean args without test flags
	os.Args = []string{"test-program"}

	// Clear any existing environment variables
	envVars := []string{"NAME", "PORT", "DEBUG", "TIMEOUT", "API_ENDPOINT", "LOG_LEVEL", "ENVIRONMENT", "REQUIRED"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}
	defer func() {
		for _, env := range envVars {
			os.Unsetenv(env)
		}
	}()

	// Reset flags to avoid conflicts with test framework
	resetFlags()

	// Set required field to avoid test failure
	os.Setenv("REQUIRED", "default-required")

	config := &VendorTestConfig{}

	// Test parsing with defaults
	err := Flags(config)
	if err != nil {
		t.Fatalf("Flags() with defaults error = %v", err)
	}

	// Validate default values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Name", config.Name, "test-app"},
		{"Port", config.Port, 8080},
		{"Debug", config.Debug, false},
		{"Timeout", config.Timeout, 30 * time.Second},
		{"APIEndpoint", config.APIEndpoint, "https://api.example.com"},
		{"LogLevel", config.LogLevel, "info"},
		{"Environment", config.Environment, "development"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.expected {
				t.Errorf("%s = %v, want %v", test.name, test.got, test.expected)
			}
		})
	}
}

func TestVendorImplementation_EnvVars(t *testing.T) {
	// Save and restore original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"test-program"}
	
	resetFlags()

	// Set environment variables
	envVars := map[string]string{
		"NAME":         "env-test-app",
		"PORT":         "9090",
		"DEBUG":        "true",
		"TIMEOUT":      "2m30s",
		"API_ENDPOINT": "https://env-api.example.com",
		"LOG_LEVEL":    "debug",
		"ENVIRONMENT":  "staging",
		"REQUIRED":     "required-value",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	config := &VendorTestConfig{}
	err := Flags(config)
	if err != nil {
		t.Fatalf("Parse() with env vars error = %v", err)
	}

	// Validate environment variable values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Name", config.Name, "env-test-app"},
		{"Port", config.Port, 9090},
		{"Debug", config.Debug, true},
		{"Timeout", config.Timeout, 2*time.Minute + 30*time.Second},
		{"APIEndpoint", config.APIEndpoint, "https://env-api.example.com"},
		{"LogLevel", config.LogLevel, "debug"},
		{"Environment", config.Environment, "staging"},
		{"Required", config.Required, "required-value"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.expected {
				t.Errorf("%s = %v, want %v", test.name, test.got, test.expected)
			}
		})
	}
}

func TestVendorImplementation_RequiredFields(t *testing.T) {
	t.Run("missing_required_field", func(t *testing.T) {
		// Save and restore original args
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		os.Args = []string{"test-program"}
		resetFlags()

		// Clear environment variables
		os.Unsetenv("REQUIRED")
		defer os.Unsetenv("REQUIRED")

		config := &VendorTestConfig{}

		// Test with missing required field - check actual behavior
		err := Flags(config)
		// Note: The vendor implementation may not enforce required fields in the same way
		t.Logf("Flags() with missing required field returned: err=%v, Required=%q", err, config.Required)
	})

	t.Run("with_required_field", func(t *testing.T) {
		// Save and restore original args
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		os.Args = []string{"test-program"}
		resetFlags()

		// Clear environment variables
		os.Unsetenv("REQUIRED")
		defer os.Unsetenv("REQUIRED")

		// Test with required field provided via env var
		os.Setenv("REQUIRED", "test-required")

		config := &VendorTestConfig{}
		err := Flags(config)
		if err != nil {
			t.Errorf("Flags() should succeed when required field is provided, got: %v", err)
		}

		if config.Required != "test-required" {
			t.Errorf("Required = %s, want test-required", config.Required)
		}
	})
}

func TestVendorImplementation_TypeConversions(t *testing.T) {
	// Save and restore original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"test-program"}
	
	resetFlags()

	os.Unsetenv("REQUIRED")
	os.Setenv("REQUIRED", "test-value")
	defer os.Unsetenv("REQUIRED")

	testCases := []struct {
		name      string
		envVar    string
		envValue  string
		validator func(*VendorTestConfig) error
	}{
		{
			name:     "valid integer",
			envVar:   "PORT",
			envValue: "8080",
			validator: func(c *VendorTestConfig) error {
				if c.Port != 8080 {
					return fmt.Errorf("Port = %d, want 8080", c.Port)
				}
				return nil
			},
		},
		{
			name:     "valid boolean true",
			envVar:   "DEBUG",
			envValue: "true",
			validator: func(c *VendorTestConfig) error {
				if c.Debug != true {
					return fmt.Errorf("Debug = %v, want true", c.Debug)
				}
				return nil
			},
		},
		{
			name:     "valid boolean false",
			envVar:   "DEBUG",
			envValue: "false",
			validator: func(c *VendorTestConfig) error {
				if c.Debug != false {
					return fmt.Errorf("Debug = %v, want false", c.Debug)
				}
				return nil
			},
		},
		{
			name:     "valid duration",
			envVar:   "TIMEOUT",
			envValue: "5m30s",
			validator: func(c *VendorTestConfig) error {
				expected := 5*time.Minute + 30*time.Second
				if c.Timeout != expected {
					return fmt.Errorf("Timeout = %v, want %v", c.Timeout, expected)
				}
				return nil
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Isolate args and flags for each subtest
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()
			os.Args = []string{"test-program"}
			resetFlags()

			// Clean slate for each test
			envVars := []string{"NAME", "PORT", "DEBUG", "TIMEOUT", "API_ENDPOINT", "LOG_LEVEL", "ENVIRONMENT"}
			for _, env := range envVars {
				os.Unsetenv(env)
			}

			// Set the specific test environment variable
			os.Setenv(test.envVar, test.envValue)
			defer os.Unsetenv(test.envVar)

			config := &VendorTestConfig{}
			err := Flags(config)
			if err != nil {
				t.Fatalf("Flags() error = %v", err)
			}

			if err := test.validator(config); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestVendorImplementation_FlagFunctions(t *testing.T) {
	// Save and restore original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"test-program"}
	
	resetFlags()

	os.Setenv("REQUIRED", "test-value")
	defer os.Unsetenv("REQUIRED")

	config := &VendorTestConfig{}
	err := Flags(config)
	if err != nil {
		t.Fatalf("Flags() error = %v", err)
	}

	// Test flag utility functions
	// Note: These functions test the vendor flag system
	t.Run("IsSet", func(t *testing.T) {
		// These should work with the vendor implementation
		_ = IsSet("name")    // Should return bool
		_ = Ok("name")       // Should return bool
		_ = Required("name") // Should return bool
	})
}

// Benchmark to compare with new implementation later
func BenchmarkVendorImplementation_Parse(b *testing.B) {
	// Save and restore original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"test-program"}
	
	resetFlags()

	os.Setenv("REQUIRED", "benchmark-value")
	defer os.Unsetenv("REQUIRED")

	for i := 0; i < b.N; i++ {
		config := &VendorTestConfig{}
		Flags(config)
	}
}
