package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestConfig for testing
type TestConfig struct {
	ServiceName string `json:"service_name" default:"test-service"`
	Version     string `json:"version" default:"1.0.0"`
	Environment string `json:"environment" default:"development"`

	API struct {
		Host      string        `json:"host" default:"localhost"`
		Port      int           `json:"port" default:"8080"`
		Timeout   time.Duration `json:"timeout" default:"30s"`
		EnableTLS bool          `json:"enable_tls" default:"false"`
	} `json:"api"`

	Database struct {
		Host     string `json:"host" default:"localhost"`
		Port     int    `json:"port" default:"5432"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`

	Features map[string]bool   `json:"features"`
	Services map[string]string `json:"services"`
}

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	// Test that it implements ConfigManager interface
	var _ ConfigManager = manager
}

func TestManager_WithDiscovery(t *testing.T) {
	manager := NewManager()

	newManager := manager.WithDiscovery(First | Indirect)

	// Should return a new instance
	if manager == newManager {
		t.Error("WithDiscovery should return a new instance")
	}
}

func TestManager_WithFlagStyle(t *testing.T) {
	manager := NewManager()

	newManager := manager.WithFlagStyle(Nested)

	// Should return a new instance
	if manager == newManager {
		t.Error("WithFlagStyle should return a new instance")
	}
}

func TestManager_WithValidation(t *testing.T) {
	manager := NewManager()

	newManager := manager.WithValidation(false)

	// Should return a new instance
	if manager == newManager {
		t.Error("WithValidation should return a new instance")
	}
}

func TestManager_WithEnvironmentPrefix(t *testing.T) {
	manager := NewManager()

	newManager := manager.WithEnvironmentPrefix("TEST")

	// Should return a new instance
	if manager == newManager {
		t.Error("WithEnvironmentPrefix should return a new instance")
	}
}

func TestManager_GenerateSample_JSON(t *testing.T) {
	manager := NewManager()

	sample, err := manager.GenerateSample("json")
	if err != nil {
		t.Fatalf("GenerateSample failed: %v", err)
	}

	if len(sample) == 0 {
		t.Error("Generated sample is empty")
	}

	// Should be valid JSON
	if sample[0] != '{' {
		t.Error("Generated sample doesn't look like JSON")
	}
}

func TestManager_GenerateSample_YAML(t *testing.T) {
	manager := NewManager()

	sample, err := manager.GenerateSample("yaml")
	if err != nil {
		t.Fatalf("GenerateSample failed: %v", err)
	}

	if len(sample) == 0 {
		t.Error("Generated sample is empty")
	}
}

func TestManager_GenerateSample_UnsupportedFormat(t *testing.T) {
	manager := NewManager()

	_, err := manager.GenerateSample("xml")
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}

func TestManager_SaveCurrent_JSON(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test-config.json")

	manager := NewManager()
	config := TestConfig{
		ServiceName: "test-save",
		Version:     "2.0.0",
	}

	err := manager.SaveCurrent(&config, filename)
	if err != nil {
		t.Fatalf("SaveCurrent failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Check file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Saved file is empty")
	}
}

func TestManager_SaveCurrent_YAML(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test-config.yaml")

	manager := NewManager()
	config := TestConfig{
		ServiceName: "test-save-yaml",
		Version:     "3.0.0",
	}

	err := manager.SaveCurrent(&config, filename)
	if err != nil {
		t.Fatalf("SaveCurrent failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestLoad_SimpleAPI(t *testing.T) {
	var config TestConfig

	// This should not fail even without any config files
	err := Load(&config)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should have default values
	if config.ServiceName != "test-service" {
		t.Errorf("Expected default ServiceName 'test-service', got '%s'", config.ServiceName)
	}

	if config.Version != "1.0.0" {
		t.Errorf("Expected default Version '1.0.0', got '%s'", config.Version)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.DiscoveryMode != (Union | Direct) {
		t.Errorf("Expected default DiscoveryMode %d, got %d", Union|Direct, opts.DiscoveryMode)
	}

	if opts.FlagStyle != FlatNest {
		t.Errorf("Expected default FlagStyle %d, got %d", FlatNest, opts.FlagStyle)
	}

	if opts.EnvOverride != "AUTOCFG_FILENAME" {
		t.Errorf("Expected default EnvOverride 'AUTOCFG_FILENAME', got '%s'", opts.EnvOverride)
	}

	if !opts.RequiredFields {
		t.Error("Expected RequiredFields to be true by default")
	}

	if !opts.ValidateOnLoad {
		t.Error("Expected ValidateOnLoad to be true by default")
	}

	if opts.SampleFormat != "json" {
		t.Errorf("Expected default SampleFormat 'json', got '%s'", opts.SampleFormat)
	}
}
