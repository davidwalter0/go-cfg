package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	API         struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"api"`
	Features map[string]bool   `json:"features"`
	Services map[string]string `json:"services"`
}

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestEngine_FindConfigFiles_EmptyPaths(t *testing.T) {
	engine := NewEngine()

	opts := DiscoveryOptions{
		Mode:        Union | Direct,
		SearchPaths: []string{},
		ProgramName: "test-discovery",
	}

	files, err := engine.FindConfigFiles(opts)
	if err != nil {
		t.Fatalf("FindConfigFiles failed: %v", err)
	}

	// Should return empty list when no files exist
	if len(files) > 0 {
		t.Errorf("Expected empty file list, got %d files", len(files))
	}
}

func TestEngine_FindConfigFiles_Union(t *testing.T) {
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

	engine := NewEngine()
	opts := DiscoveryOptions{
		Mode:        Union | Direct,
		SearchPaths: []string{config1, config2},
	}

	files, err := engine.FindConfigFiles(opts)
	if err != nil {
		t.Fatalf("FindConfigFiles failed: %v", err)
	}

	// Union mode should find both files
	if len(files) != 2 {
		t.Errorf("Expected 2 config files, got %d", len(files))
	}
}

func TestEngine_FindConfigFiles_First(t *testing.T) {
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

	engine := NewEngine()
	opts := DiscoveryOptions{
		Mode:        First | Direct,
		SearchPaths: []string{config1, config2},
	}

	files, err := engine.FindConfigFiles(opts)
	if err != nil {
		t.Fatalf("FindConfigFiles failed: %v", err)
	}

	// First mode should find only the first file
	if len(files) != 1 {
		t.Errorf("Expected 1 config file, got %d", len(files))
	}

	if len(files) > 0 && files[0] != config1 {
		t.Errorf("Expected first file to be %s, got %s", config1, files[0])
	}
}

func TestEngine_FindConfigFiles_EnvironmentOverride(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "env-config.json")

	if err := os.WriteFile(configFile, []byte(`{"service_name": "env-service"}`), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Set environment variable
	envVar := "TEST_CONFIG_OVERRIDE"
	os.Setenv(envVar, configFile)
	defer os.Unsetenv(envVar)

	engine := NewEngine()
	opts := DiscoveryOptions{
		Mode:        Union | Direct,
		EnvOverride: envVar,
	}

	files, err := engine.FindConfigFiles(opts)
	if err != nil {
		t.Fatalf("FindConfigFiles failed: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("Expected at least 1 config file from environment override")
	}

	if files[0] != configFile {
		t.Errorf("Expected first file to be %s, got %s", configFile, files[0])
	}
}

func TestEngine_LoadConfigFile_Direct(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.json")

	// Create test config file
	testConfig := map[string]interface{}{
		"service_name": "test-service",
		"version":      "2.0.0",
		"api": map[string]interface{}{
			"host": "example.com",
			"port": 9090,
		},
		"features": map[string]bool{
			"auth":    true,
			"metrics": false,
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
	engine := NewEngine()
	var config TestConfig

	err = engine.LoadConfigFile(configFile, &config)
	if err != nil {
		t.Fatalf("LoadConfigFile failed: %v", err)
	}

	// Verify loaded values
	if config.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", config.ServiceName)
	}

	if config.Version != "2.0.0" {
		t.Errorf("Expected Version '2.0.0', got '%s'", config.Version)
	}

	if config.API.Host != "example.com" {
		t.Errorf("Expected API.Host 'example.com', got '%s'", config.API.Host)
	}

	if config.API.Port != 9090 {
		t.Errorf("Expected API.Port 9090, got %d", config.API.Port)
	}

	// Test maps
	if !config.Features["auth"] {
		t.Error("Expected Features['auth'] to be true")
	}

	if config.Features["metrics"] {
		t.Error("Expected Features['metrics'] to be false")
	}

	if config.Services["auth_service"] != "https://auth.example.com" {
		t.Errorf("Expected Services['auth_service'] 'https://auth.example.com', got '%s'",
			config.Services["auth_service"])
	}
}

func TestEngine_LoadConfigFile_Indirect(t *testing.T) {
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

	// Test loading
	engine := NewEngine()
	var config TestConfig

	err = engine.LoadConfigFile(indirectConfigFile, &config)
	if err != nil {
		t.Fatalf("LoadConfigFile failed: %v", err)
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

	// Clean up environment
	os.Unsetenv("TEST_VAR")
}

func TestEngine_MergeConfigs(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple config files
	config1 := filepath.Join(tempDir, "config1.json")
	config1Data := map[string]interface{}{
		"service_name": "merged-service",
		"version":      "1.0.0",
		"api": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
	}

	data1, _ := json.Marshal(config1Data)
	if err := os.WriteFile(config1, data1, 0644); err != nil {
		t.Fatalf("Failed to write config1: %v", err)
	}

	config2 := filepath.Join(tempDir, "config2.json")
	config2Data := map[string]interface{}{
		"version": "2.0.0", // Override version
		"api": map[string]interface{}{
			"port": 9090, // Override port
		},
		"features": map[string]bool{
			"auth": true,
		},
	}

	data2, _ := json.Marshal(config2Data)
	if err := os.WriteFile(config2, data2, 0644); err != nil {
		t.Fatalf("Failed to write config2: %v", err)
	}

	// Test merging
	engine := NewEngine()
	var config TestConfig

	err := engine.MergeConfigs([]string{config1, config2}, &config)
	if err != nil {
		t.Fatalf("MergeConfigs failed: %v", err)
	}

	// Verify merged values
	if config.ServiceName != "merged-service" {
		t.Errorf("Expected ServiceName 'merged-service', got '%s'", config.ServiceName)
	}

	// Version should be overridden by config2
	if config.Version != "2.0.0" {
		t.Errorf("Expected Version '2.0.0', got '%s'", config.Version)
	}

	// API host from config1
	if config.API.Host != "localhost" {
		t.Errorf("Expected API.Host 'localhost', got '%s'", config.API.Host)
	}

	// API port overridden by config2
	if config.API.Port != 9090 {
		t.Errorf("Expected API.Port 9090, got %d", config.API.Port)
	}

	// Features from config2
	if !config.Features["auth"] {
		t.Error("Expected Features['auth'] to be true")
	}
}
