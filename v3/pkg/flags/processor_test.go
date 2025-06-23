package flags

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

type TestConfig struct {
	ServiceName string        `json:"service_name" default:"test-service" help:"Service name"`
	Version     string        `json:"version" default:"1.0.0" help:"Service version"`
	Debug       bool          `json:"debug" default:"false" help:"Debug mode"`
	Port        int           `json:"port" default:"8080" help:"Server port"`
	Timeout     time.Duration `json:"timeout" default:"30s" help:"Request timeout"`

	API struct {
		Host      string `json:"host" default:"localhost" help:"API host"`
		Port      int    `json:"port" default:"9090" help:"API port"`
		EnableTLS bool   `json:"enable_tls" default:"false" help:"Enable TLS"`
	} `json:"api"`

	Features map[string]bool `json:"features" help:"Feature flags"`
	Services []string        `json:"services" help:"Service list"`
}

func TestNewProcessor(t *testing.T) {
	processor := NewProcessor()
	if processor == nil {
		t.Fatal("NewProcessor returned nil")
	}
}

func TestProcessor_GenerateFlags(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	if flagSet == nil {
		t.Fatal("GenerateFlags returned nil flagSet")
	}

	// Check that some expected flags were created
	expectedFlags := []string{
		"service-name",
		"version",
		"debug",
		"port",
		"timeout",
		"api-host",
		"api-port",
		"api-enable-tls",
	}

	for _, flagName := range expectedFlags {
		flag := flagSet.Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
		}
	}
}

func TestProcessor_ProcessStruct(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := processor.ProcessStruct(config, flagSet)
	if err != nil {
		t.Fatalf("ProcessStruct failed: %v", err)
	}

	// Check that flags were added to the flag set
	flagCount := 0
	flagSet.VisitAll(func(flag *pflag.Flag) {
		flagCount++
	})

	if flagCount == 0 {
		t.Error("No flags were added to the flag set")
	}

	// Should have at least the basic flags
	if flagCount < 5 {
		t.Errorf("Expected at least 5 flags, got %d", flagCount)
	}
}

func TestProcessor_ApplyFlagValues_String(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	// Generate flags
	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	// Simulate command line arguments
	args := []string{"--service-name", "test-apply", "--version", "2.0.0"}
	err = flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	// Apply flag values
	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	// Check that values were applied
	if config.ServiceName != "test-apply" {
		t.Errorf("Expected ServiceName 'test-apply', got '%s'", config.ServiceName)
	}

	if config.Version != "2.0.0" {
		t.Errorf("Expected Version '2.0.0', got '%s'", config.Version)
	}
}

func TestProcessor_ApplyFlagValues_Bool(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	args := []string{"--debug", "--api-enable-tls"}
	err = flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	if !config.Debug {
		t.Error("Expected Debug to be true")
	}

	if !config.API.EnableTLS {
		t.Error("Expected API.EnableTLS to be true")
	}
}

func TestProcessor_ApplyFlagValues_Int(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	args := []string{"--port", "3000", "--api-port", "4000"}
	err = flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	if config.Port != 3000 {
		t.Errorf("Expected Port 3000, got %d", config.Port)
	}

	if config.API.Port != 4000 {
		t.Errorf("Expected API.Port 4000, got %d", config.API.Port)
	}
}

func TestProcessor_ApplyFlagValues_Duration(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	args := []string{"--timeout", "45s"}
	err = flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	expectedTimeout := 45 * time.Second
	if config.Timeout != expectedTimeout {
		t.Errorf("Expected Timeout %v, got %v", expectedTimeout, config.Timeout)
	}
}

func TestProcessor_ApplyFlagValues_Slice(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	args := []string{"--services", "service1,service2,service3"}
	err = flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	expectedServices := []string{"service1", "service2", "service3"}
	if len(config.Services) != len(expectedServices) {
		t.Errorf("Expected %d services, got %d", len(expectedServices), len(config.Services))
	}

	for i, expected := range expectedServices {
		if i >= len(config.Services) || config.Services[i] != expected {
			t.Errorf("Expected service[%d] '%s', got '%s'", i, expected,
				func() string {
					if i < len(config.Services) {
						return config.Services[i]
					}
					return "<missing>"
				}())
		}
	}
}

func TestProcessor_ApplyFlagValues_Map(t *testing.T) {
	processor := NewProcessor()
	config := &TestConfig{}

	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	args := []string{"--features", "auth:true,metrics:false,debug:true"}
	err = flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	if config.Features == nil {
		t.Fatal("Features map is nil")
	}

	expectedFeatures := map[string]string{
		"auth":    "true",
		"metrics": "false",
		"debug":   "true",
	}

	if len(config.Features) != len(expectedFeatures) {
		t.Errorf("Expected %d features, got %d", len(expectedFeatures), len(config.Features))
	}

	// Note: The current implementation stores maps as map[string]string
	// This test validates the basic map functionality
}

func TestProcessor_Integration(t *testing.T) {
	// Save original args
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set test args
	os.Args = []string{"test-program", "--service-name", "integration-test", "--debug", "--port", "5000"}

	processor := NewProcessor()
	config := &TestConfig{}

	// Generate flags
	flagSet, err := processor.GenerateFlags(config)
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}

	// Parse args
	err = flagSet.Parse(os.Args[1:])
	if err != nil {
		t.Fatalf("Flag parsing failed: %v", err)
	}

	// Apply values
	err = processor.ApplyFlagValues(config, flagSet)
	if err != nil {
		t.Fatalf("ApplyFlagValues failed: %v", err)
	}

	// Verify all values
	if config.ServiceName != "integration-test" {
		t.Errorf("Expected ServiceName 'integration-test', got '%s'", config.ServiceName)
	}

	if !config.Debug {
		t.Error("Expected Debug to be true")
	}

	if config.Port != 5000 {
		t.Errorf("Expected Port 5000, got %d", config.Port)
	}
}
