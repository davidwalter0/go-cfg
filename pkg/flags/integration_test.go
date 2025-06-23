package flags

import (
	"os"
	"testing"
	"time"

	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

// TestConfig represents a typical configuration structure
type TestConfig struct {
	Name        string                  `usage:"Application name" default:"test-app"`
	Port        int                     `usage:"Server port" default:"8080"`
	Debug       bool                    `usage:"Enable debug mode" default:"false"`
	Timeout     *types.DurationValue    `usage:"Request timeout"`
	APIEndpoint *types.URLValue         `usage:"API endpoint URL"`
	LogLevel    *types.LogLevelValue    `usage:"Log level"`
	Environment *types.EnvironmentValue `usage:"Deployment environment"`
	Labels      *types.KeyValueValue    `usage:"Application labels"`
	ListenAddr  *types.HostPortValue    `usage:"Listen address"`
}

func TestFlagConfig_RegisterCustomTypes_Integration(t *testing.T) {
	// Set up environment variables
	testEnvVars := map[string]string{
		"NAME":         "integration-test-app",
		"PORT":         "9090",
		"DEBUG":        "true",
		"TIMEOUT":      "2m30s",
		"API_ENDPOINT": "https://api-test.example.com/v2",
		"LOG_LEVEL":    "debug",
		"ENVIRONMENT":  "staging",
		"LABELS":       "service=test,version=v1.0.0",
		"LISTEN_ADDR":  "0.0.0.0:9090",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	// Initialize configuration
	config := &TestConfig{
		Timeout:     types.NewDurationValue(30 * time.Second),
		APIEndpoint: types.NewURLValue("https://default-api.example.com"),
		LogLevel:    types.NewLogLevelValue("info"),
		Environment: types.NewEnvironmentValue("development"),
		Labels:      types.NewKeyValueValue(),
		ListenAddr:  types.NewHostPortValue("localhost", 8080),
	}

	// Create flag configuration
	flagConfig := NewFlagConfig()

	// Register custom types
	err := flagConfig.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterCustomTypes() error = %v", err)
	}

	// Parse with empty args (should use environment defaults)
	err = flagConfig.Parse([]string{})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Validate environment variables were applied
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Name", config.Name, "integration-test-app"},
		{"Port", config.Port, 9090},
		{"Debug", config.Debug, true},
		{"Timeout", config.Timeout.GetDuration(), 2*time.Minute + 30*time.Second},
		{"APIEndpoint", config.APIEndpoint.String(), "https://api-test.example.com/v2"},
		{"LogLevel", config.LogLevel.GetLevel(), "debug"},
		{"Environment", config.Environment.GetEnvironment(), "staging"},
		{"Labels service", func() string {
			if val, exists := config.Labels.GetValue("service"); exists {
				return val
			}
			return ""
		}(), "test"},
		{"Labels version", func() string {
			if val, exists := config.Labels.GetValue("version"); exists {
				return val
			}
			return ""
		}(), "v1.0.0"},
		{"ListenAddr Host", config.ListenAddr.GetHost(), "0.0.0.0"},
		{"ListenAddr Port", config.ListenAddr.GetPort(), 9090},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.expected {
				t.Errorf("%s = %v, want %v", test.name, test.got, test.expected)
			}
		})
	}
}

func TestFlagConfig_CommandLineOverride_Integration(t *testing.T) {
	// Initialize configuration with defaults
	config := &TestConfig{
		Timeout:     types.NewDurationValue(30 * time.Second),
		APIEndpoint: types.NewURLValue("https://default-api.example.com"),
		LogLevel:    types.NewLogLevelValue("info"),
		Environment: types.NewEnvironmentValue("development"),
		Labels:      types.NewKeyValueValue(),
		ListenAddr:  types.NewHostPortValue("localhost", 8080),
	}

	// Create flag configuration
	flagConfig := NewFlagConfig()

	// Register custom types
	err := flagConfig.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterCustomTypes() error = %v", err)
	}

	// Parse with command line arguments
	args := []string{
		"--name=cli-override-app",
		"--port=7070",
		"--debug=true",
		"--timeout=5m",
		"--api-endpoint=https://cli-api.example.com",
		"--log-level=warn",
		"--environment=production",
		"--labels=env=prod,tier=backend",
		"--listen-addr=127.0.0.1:7070",
	}

	err = flagConfig.Parse(args)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Validate command line arguments were applied
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Name", config.Name, "cli-override-app"},
		{"Port", config.Port, 7070},
		{"Debug", config.Debug, true},
		{"Timeout", config.Timeout.GetDuration(), 5 * time.Minute},
		{"APIEndpoint", config.APIEndpoint.String(), "https://cli-api.example.com"},
		{"LogLevel", config.LogLevel.GetLevel(), "warn"},
		{"Environment", config.Environment.GetEnvironment(), "production"},
		{"Environment IsProduction", config.Environment.IsProduction(), true},
		{"Labels env", func() string {
			if val, exists := config.Labels.GetValue("env"); exists {
				return val
			}
			return ""
		}(), "prod"},
		{"Labels tier", func() string {
			if val, exists := config.Labels.GetValue("tier"); exists {
				return val
			}
			return ""
		}(), "backend"},
		{"ListenAddr Host", config.ListenAddr.GetHost(), "127.0.0.1"},
		{"ListenAddr Port", config.ListenAddr.GetPort(), 7070},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.expected {
				t.Errorf("%s = %v, want %v", test.name, test.got, test.expected)
			}
		})
	}
}

func TestFlagConfig_Validation_Integration(t *testing.T) {
	config := &TestConfig{
		Timeout:     types.NewDurationValue(30 * time.Second),
		APIEndpoint: types.NewURLValue(""),
		LogLevel:    types.NewLogLevelValue("info"),
		Environment: types.NewEnvironmentValue("development"),
		Labels:      types.NewKeyValueValue(),
		ListenAddr:  types.NewHostPortValue("localhost", 8080),
	}

	flagConfig := NewFlagConfig()

	// Mark API endpoint as required
	flagConfig.SetRequired("APIEndpoint")

	// Add custom validators
	flagConfig.AddValidator("Port", func(value interface{}) error {
		// This validator should receive the pflag.Value, not the actual value
		return nil // Simple test - just ensure no panic
	})

	err := flagConfig.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterCustomTypes() error = %v", err)
	}

	// Test missing required field
	err = flagConfig.Parse([]string{})
	if err == nil {
		t.Errorf("Parse() should fail when required field is missing")
	}

	// Test with required field provided
	err = flagConfig.Parse([]string{"--api-endpoint=https://required-api.example.com"})
	if err != nil {
		t.Errorf("Parse() should succeed when required field is provided, got: %v", err)
	}

	if config.APIEndpoint.String() != "https://required-api.example.com" {
		t.Errorf("APIEndpoint = %s, want https://required-api.example.com", config.APIEndpoint.String())
	}
}

func TestFlagConfig_TypeSafety_Integration(t *testing.T) {
	config := &TestConfig{
		Timeout:     types.NewDurationValue(30 * time.Second),
		APIEndpoint: types.NewURLValue("https://api.example.com"),
		LogLevel:    types.NewLogLevelValue("info"),
		Environment: types.NewEnvironmentValue("development"),
		Labels:      types.NewKeyValueValue(),
		ListenAddr:  types.NewHostPortValue("localhost", 8080),
	}

	flagConfig := NewFlagConfig()
	err := flagConfig.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterCustomTypes() error = %v", err)
	}

	// Test various type safety scenarios
	testCases := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "valid log level",
			args:      []string{"--log-level=debug"},
			wantError: false,
		},
		{
			name:      "invalid log level",
			args:      []string{"--log-level=invalid"},
			wantError: true,
		},
		{
			name:      "valid environment",
			args:      []string{"--environment=production"},
			wantError: false,
		},
		{
			name:      "invalid environment",
			args:      []string{"--environment=invalid"},
			wantError: true,
		},
		{
			name:      "valid URL",
			args:      []string{"--api-endpoint=https://valid.example.com"},
			wantError: false,
		},
		{
			name:      "invalid URL scheme",
			args:      []string{"--api-endpoint=ftp://invalid.example.com"},
			wantError: false, // URL parsing is lenient
		},
		{
			name:      "valid duration",
			args:      []string{"--timeout=1h30m"},
			wantError: false,
		},
		{
			name:      "invalid duration",
			args:      []string{"--timeout=invalid"},
			wantError: true,
		},
		{
			name:      "valid host:port",
			args:      []string{"--listen-addr=0.0.0.0:9090"},
			wantError: false,
		},
		{
			name:      "valid labels",
			args:      []string{"--labels=key1=value1,key2=value2"},
			wantError: false,
		},
		{
			name:      "invalid labels format",
			args:      []string{"--labels=invalid-format"},
			wantError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Reset config to defaults
			_ = config.LogLevel.Set("info")
			_ = config.Environment.Set("development")
			_ = config.APIEndpoint.Set("https://api.example.com")
			_ = config.Timeout.Set("30s")
			_ = config.ListenAddr.Set("localhost:8080")
			config.Labels = types.NewKeyValueValue()

			err := flagConfig.Parse(test.args)
			if (err != nil) != test.wantError {
				t.Errorf("Parse(%v) error = %v, wantError = %v", test.args, err, test.wantError)
			}
		})
	}
}

// Benchmark tests to ensure performance is acceptable
func BenchmarkFlagConfig_RegisterCustomTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config := &TestConfig{
			Timeout:     types.NewDurationValue(30 * time.Second),
			APIEndpoint: types.NewURLValue("https://api.example.com"),
			LogLevel:    types.NewLogLevelValue("info"),
			Environment: types.NewEnvironmentValue("development"),
			Labels:      types.NewKeyValueValue(),
			ListenAddr:  types.NewHostPortValue("localhost", 8080),
		}

		flagConfig := NewFlagConfig()
		_ = flagConfig.RegisterCustomTypes(config)
	}
}

func BenchmarkFlagConfig_Parse(b *testing.B) {
	config := &TestConfig{
		Timeout:     types.NewDurationValue(30 * time.Second),
		APIEndpoint: types.NewURLValue("https://api.example.com"),
		LogLevel:    types.NewLogLevelValue("info"),
		Environment: types.NewEnvironmentValue("development"),
		Labels:      types.NewKeyValueValue(),
		ListenAddr:  types.NewHostPortValue("localhost", 8080),
	}

	flagConfig := NewFlagConfig()
	_ = flagConfig.RegisterCustomTypes(config)

	args := []string{
		"--name=benchmark-app",
		"--port=8080",
		"--timeout=1m",
		"--api-endpoint=https://api.example.com",
		"--log-level=info",
		"--environment=development",
		"--labels=service=test,version=v1.0.0",
		"--listen-addr=localhost:8080",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flagConfig.Parse(args)
	}
}
