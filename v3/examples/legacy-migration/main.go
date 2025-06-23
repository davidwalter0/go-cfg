package main

import (
	"fmt"
	"log"

	// Legacy API (compatible)
	legacy "github.com/davidwalter0/go-cfg/v3/legacy"

	// Modern v3 API
	"github.com/davidwalter0/go-cfg/v3/pkg/config"
)

// Example configuration struct
type AppConfig struct {
	ServiceName string `json:"service_name" default:"my-service"`
	Version     string `json:"version" default:"1.0.0"`
	API         struct {
		Host string `json:"host" default:"localhost"`
		Port int    `json:"port" default:"8080"`
	} `json:"api"`
	Features map[string]bool `json:"features"`
}

// Example showing legacy API compatibility
func ExampleLegacyCompatibility() {
	var cfg AppConfig

	// Legacy API calls work exactly as before
	if err := legacy.Init(&cfg); err != nil {
		log.Fatalf("Legacy configuration failed: %v", err)
	}

	fmt.Printf("Legacy API loaded: %s v%s\n", cfg.ServiceName, cfg.Version)
}

// Example showing migration to modern v3 API
func ExampleModernAPI() {
	var cfg AppConfig

	// Modern v3 API - simple usage
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("Modern configuration failed: %v", err)
	}

	fmt.Printf("Modern API loaded: %s v%s\n", cfg.ServiceName, cfg.Version)
}

// Example showing enhanced v3 API with advanced features
func ExampleAdvancedV3API() {
	var cfg AppConfig

	// Advanced v3 API with custom configuration
	manager := config.NewManager().
		WithDiscovery(config.Union | config.Direct).
		WithFlagStyle(config.FlatNest).
		WithValidation(true)

	if err := manager.Load(&cfg); err != nil {
		log.Fatalf("Advanced configuration failed: %v", err)
	}

	// Generate sample configuration
	sample, err := manager.GenerateSample("json")
	if err != nil {
		log.Printf("Failed to generate sample: %v", err)
	} else {
		fmt.Printf("Sample config:\n%s\n", sample)
	}

	fmt.Printf("Advanced API loaded: %s v%s\n", cfg.ServiceName, cfg.Version)
}

// Example showing enhanced FlatNest API
func ExampleFlatNestAPI() {
	var cfg AppConfig

	// Enhanced FlatNest with discovery (new capability)
	if err := legacy.FlatNestWithDiscovery(config.Union|config.Direct, &cfg); err != nil {
		log.Fatalf("FlatNest configuration failed: %v", err)
	}

	fmt.Printf("FlatNest API loaded: %s v%s\n", cfg.ServiceName, cfg.Version)
}

// Migration guide function showing equivalent APIs
func MigrationGuide() {
	var _ AppConfig // Placeholder for documentation

	// Old: legacy.Init(&cfg)
	// New: config.Load(&cfg)
	// Or:  legacy.Init(&cfg) (still works!)

	// Old: legacy.Simple(&cfg)
	// New: config.NewManager().WithFlagStyle(config.FlatNest).Load(&cfg)
	// Or:  legacy.Simple(&cfg) (still works!)

	// Old: legacy.Nest(&cfg)
	// New: config.NewManager().WithFlagStyle(config.Nested).Load(&cfg)
	// Or:  legacy.Nest(&cfg) (still works!)

	// New capabilities in v3:
	// - AutoCfg-style discovery
	// - Multiple naming strategies
	// - Configuration validation
	// - Sample generation
	// - Help generation
	// - Configuration saving

	fmt.Println("Migration complete - all legacy APIs preserved, new features available!")
}
