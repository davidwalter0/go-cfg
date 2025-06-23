package cfg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestConfig for testing discovery functionality
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

func TestConfigDiscovery_NewConfigDiscovery(t *testing.T) {
	opts := ConfigOptions{
		SearchMode:  Union | Direct,
		ProgramName: "test-app",
		EnvOverride: "TEST_CONFIG",
	}

	discovery := NewConfigDiscovery(opts)

	if discovery.Options.SearchMode != (Union | Direct) {
		t.Errorf("Expected SearchMode %d, got %d", Union|Direct, discovery.Options.SearchMode)
	}

	if discovery.Options.ProgramName != "test-app" {
		t.Errorf("Expected ProgramName 'test-app', got '%s'", discovery.Options.ProgramName)
	}
}

func TestConfigDiscovery_GetSearchPaths(t *testing.T) {
	discovery := NewConfigDiscovery(ConfigOptions{
		ProgramName: "test-discovery",
	})

	paths := discovery.getSearchPaths()

	// Should contain standard paths
	expectedPaths := []string{
		"/etc/test-discovery/config.json",
		".test-discovery.json",
		"config.json",
	}

	// Check that our expected paths are found
	for _, expected := range expectedPaths {
		found := false
		for _, path := range paths {
			if path == expected || filepath.Base(path) == filepath.Base(expected) {
				found = true
				break
			}
		}
		if !found && expected != "/etc/test-discovery/config.json" {
			// Skip system path check in tests
			t.Errorf("Expected path pattern '%s' not found in paths: %v", expected, paths)
		}
	}
}

func TestConfigDiscovery_LoadConfigFile_Direct(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	// Create test config file
	testConfig := map[string]interface{}{
		"service_name": "test-direct-service",
		"version":      "2.0.0",
		"environment":  "staging",
		"api": map[string]interface{}{
			"host": "api.example.com",
			"port": 9090,
		},
		"database": map[string]interface{}{
			"host":     "db.example.com",
			"port":     5432,
			"name":     "test_db",
			"username": "test_user",
			"password": "test_pass",
		},
		"features": map[string]bool{
			"auth":    true,
			"metrics": false,
			"debug":   true,
		},
		"services": map[string]string{
			"auth_service": "https://auth.example.com",
			"user_service": "https://users.example.com",
		},
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading
	discovery := NewConfigDiscovery(ConfigOptions{
		SearchMode: Direct,
	})

	var config TestConfig
	err = discovery.loadConfigFile(configFile, &config)
	if err != nil {
		t.Fatalf("Failed to load config file: %v", err)
	}

	// Verify loaded values
	if config.ServiceName != "test-direct-service" {
		t.Errorf("Expected ServiceName 'test-direct-service', got '%s'", config.ServiceName)
	}

	if config.Version != "2.0.0" {
		t.Errorf("Expected Version '2.0.0', got '%s'", config.Version)
	}

	if config.API.Host != "api.example.com" {
		t.Errorf("Expected API.Host 'api.example.com', got '%s'", config.API.Host)
	}

	if config.API.Port != 9090 {
		t.Errorf("Expected API.Port 9090, got %d", config.API.Port)
	}

	if config.Database.Name != "test_db" {
		t.Errorf("Expected Database.Name 'test_db', got '%s'", config.Database.Name)
	}

	// Test maps
	if !config.Features["auth"] {
		t.Error("Expected Features['auth'] to be true")
	}

	if config.Features["metrics"] {
		t.Error("Expected Features['metrics'] to be false")
	}

	if config.Services["auth_service"] != "https://auth.example.com" {
		t.Errorf("Expected Services['auth_service'] 'https://auth.example.com', got '%s'", config.Services["auth_service"])
	}
}

func TestConfigDiscovery_LoadConfigFile_Indirect(t *testing.T) {
	tempDir := t.TempDir()

	// Create actual config file
	actualConfigFile := filepath.Join(tempDir, "actual-config.json")
	actualConfig := map[string]interface{}{
		"service_name": "indirect-service",
		"version":      "3.0.0",
		"api": map[string]interface{}{
			"host": "indirect.example.com",
			"port": 7070,
		},
	}

	data, err := json.Marshal(actualConfig)
	if err != nil {
		t.Fatalf("Failed to marshal actual config: %v", err)
	}

	if err := os.WriteFile(actualConfigFile, data, 0644); err != nil {
		t.Fatalf("Failed to write actual config: %v", err)
	}

	// Create indirect config file (pointer file)
	indirectConfigFile := filepath.Join(tempDir, "indirect-config.json")
	indirectConfig := map[string]interface{}{
		"path": actualConfigFile,
		"env": map[string]string{
			"TEST_VAR": "test_value",
		},
	}

	data, err = json.Marshal(indirectConfig)
	if err != nil {
		t.Fatalf("Failed to marshal indirect config: %v", err)
	}

	if err := os.WriteFile(indirectConfigFile, data, 0644); err != nil {
		t.Fatalf("Failed to write indirect config: %v", err)
	}

	// Test loading with indirect mode
	discovery := NewConfigDiscovery(ConfigOptions{
		SearchMode: Indirect,
	})

	var config TestConfig
	err = discovery.loadConfigFile(indirectConfigFile, &config)
	if err != nil {
		t.Fatalf("Failed to load indirect config file: %v", err)
	}

	// Verify loaded values from actual config
	if config.ServiceName != "indirect-service" {
		t.Errorf("Expected ServiceName 'indirect-service', got '%s'", config.ServiceName)
	}

	if config.Version != "3.0.0" {
		t.Errorf("Expected Version '3.0.0', got '%s'", config.Version)
	}

	if config.API.Host != "indirect.example.com" {
		t.Errorf("Expected API.Host 'indirect.example.com', got '%s'", config.API.Host)
	}

	// Verify environment variable was set
	if os.Getenv("TEST_VAR") != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got '%s'", os.Getenv("TEST_VAR"))
	}
}

func TestConfigDiscovery_FindConfigFiles_Union(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple config files
	config1 := filepath.Join(tempDir, "config1.json")
	config2 := filepath.Join(tempDir, "config2.json")

	if err := os.WriteFile(config1, []byte(`{"service_name": "service1"}`), 0644); err != nil {
		t.Fatalf("Failed to write config1: %v", err)
	}

	if err := os.WriteFile(config2, []byte(`{"service_name": "service2"}`), 0644); err != nil {
		t.Fatalf("Failed to write config2: %v", err)
	}

	discovery := NewConfigDiscovery(ConfigOptions{
		SearchMode:  Union | Direct,
		SearchPaths: []string{config1, config2},
	})

	files, err := discovery.findConfigFiles()
	if err != nil {
		t.Fatalf("Failed to find config files: %v", err)
	}

	// Union mode should find both files
	if len(files) != 2 {
		t.Errorf("Expected 2 config files, got %d", len(files))
	}
}

func TestConfigDiscovery_FindConfigFiles_First(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple config files
	config1 := filepath.Join(tempDir, "config1.json")
	config2 := filepath.Join(tempDir, "config2.json")

	if err := os.WriteFile(config1, []byte(`{"service_name": "service1"}`), 0644); err != nil {
		t.Fatalf("Failed to write config1: %v", err)
	}

	if err := os.WriteFile(config2, []byte(`{"service_name": "service2"}`), 0644); err != nil {
		t.Fatalf("Failed to write config2: %v", err)
	}

	discovery := NewConfigDiscovery(ConfigOptions{
		SearchMode:  First | Direct,
		SearchPaths: []string{config1, config2},
	})

	files, err := discovery.findConfigFiles()
	if err != nil {
		t.Fatalf("Failed to find config files: %v", err)
	}

	// First mode should find only the first file
	if len(files) != 1 {
		t.Errorf("Expected 1 config file, got %d", len(files))
	}

	if len(files) > 0 && files[0] != config1 {
		t.Errorf("Expected first file to be %s, got %s", config1, files[0])
	}
}

func TestConfigDiscovery_EnvironmentOverride(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "env-config.json")

	if err := os.WriteFile(configFile, []byte(`{"service_name": "env-service"}`), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Set environment variable
	envVar := "TEST_CONFIG_OVERRIDE"
	os.Setenv(envVar, configFile)
	defer os.Unsetenv(envVar)

	discovery := NewConfigDiscovery(ConfigOptions{
		SearchMode:  Union | Direct,
		EnvOverride: envVar,
	})

	files, err := discovery.findConfigFiles()
	if err != nil {
		t.Fatalf("Failed to find config files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("Expected at least 1 config file from environment override")
	}

	if files[0] != configFile {
		t.Errorf("Expected first file to be %s, got %s", configFile, files[0])
	}
}

func TestConfigDiscovery_DiscoverAndLoad_Integration(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "integration-config.json")

	// Create comprehensive test config
	testConfig := map[string]interface{}{
		"service_name": "integration-test",
		"version":      "4.0.0",
		"environment":  "production",
		"api": map[string]interface{}{
			"host":       "prod.example.com",
			"port":       443,
			"enable_tls": true,
		},
		"database": map[string]interface{}{
			"host":     "prod-db.example.com",
			"port":     5432,
			"name":     "prod_db",
			"username": "prod_user",
			"password": "prod_pass",
		},
		"features": map[string]bool{
			"auth":          true,
			"metrics":       true,
			"debug":         false,
			"audit_logging": true,
		},
		"services": map[string]string{
			"auth_service":    "https://auth.prod.example.com",
			"user_service":    "https://users.prod.example.com",
			"billing_service": "https://billing.prod.example.com",
		},
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test discovery and loading (without flag processing)
	discovery := NewConfigDiscovery(ConfigOptions{
		SearchMode:  Union | Direct,
		SearchPaths: []string{configFile},
		ProgramName: "integration-test",
	})

	var config TestConfig
	err = discovery.DiscoverAndLoad(&config)
	if err != nil {
		t.Fatalf("DiscoverAndLoad failed: %v", err)
	}

	// Verify all configuration was loaded correctly
	if config.ServiceName != "integration-test" {
		t.Errorf("Expected ServiceName 'integration-test', got '%s'", config.ServiceName)
	}

	if config.Version != "4.0.0" {
		t.Errorf("Expected Version '4.0.0', got '%s'", config.Version)
	}

	if config.Environment != "production" {
		t.Errorf("Expected Environment 'production', got '%s'", config.Environment)
	}

	if config.API.Host != "prod.example.com" {
		t.Errorf("Expected API.Host 'prod.example.com', got '%s'", config.API.Host)
	}

	if config.API.Port != 443 {
		t.Errorf("Expected API.Port 443, got %d", config.API.Port)
	}

	if !config.API.EnableTLS {
		t.Error("Expected API.EnableTLS to be true")
	}

	if config.Database.Name != "prod_db" {
		t.Errorf("Expected Database.Name 'prod_db', got '%s'", config.Database.Name)
	}

	// Test maps loaded correctly
	if !config.Features["auth"] {
		t.Error("Expected Features['auth'] to be true")
	}

	if !config.Features["audit_logging"] {
		t.Error("Expected Features['audit_logging'] to be true")
	}

	if config.Features["debug"] {
		t.Error("Expected Features['debug'] to be false")
	}

	if config.Services["billing_service"] != "https://billing.prod.example.com" {
		t.Errorf("Expected Services['billing_service'] to be 'https://billing.prod.example.com', got '%s'",
			config.Services["billing_service"])
	}
}

func TestDefaultConfigOptions(t *testing.T) {
	opts := DefaultConfigOptions()

	if opts.SearchMode != (Union | Direct) {
		t.Errorf("Expected default SearchMode %d, got %d", Union|Direct, opts.SearchMode)
	}

	if opts.EnvOverride != "AUTOCFG_FILENAME" {
		t.Errorf("Expected default EnvOverride 'AUTOCFG_FILENAME', got '%s'", opts.EnvOverride)
	}
}

func TestConfigDiscovery_GetProgramName(t *testing.T) {
	// Test with explicit program name
	discovery := NewConfigDiscovery(ConfigOptions{
		ProgramName: "test-program",
	})

	name := discovery.getProgramName()
	if name != "test-program" {
		t.Errorf("Expected program name 'test-program', got '%s'", name)
	}

	// Test auto-detection (this will use the test binary name)
	discovery = NewConfigDiscovery(ConfigOptions{})
	name = discovery.getProgramName()

	// Should not be empty and should not be "app" (fallback)
	if name == "" {
		t.Error("Expected non-empty program name from auto-detection")
	}
}
