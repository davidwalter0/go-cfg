package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davidwalter0/go-cfg/v3/pkg/config"
)

// BasicConfig demonstrates simple configuration
type BasicConfig struct {
	ServiceName string        `json:"service_name" required:"true" help:"Name of the service"`
	Version     string        `json:"version" default:"1.0.0" help:"Service version"`
	Debug       bool          `json:"debug" default:"false" help:"Enable debug mode"`
	Port        int           `json:"port" default:"8080" help:"Server port"`
	Timeout     time.Duration `json:"timeout" default:"30s" help:"Request timeout"`
}

func main() {
	var cfg BasicConfig

	fmt.Println("=== Basic Configuration Example ===")

	// Simple API - uses all defaults
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("  Service: %s v%s\n", cfg.ServiceName, cfg.Version)
	fmt.Printf("  Debug: %v\n", cfg.Debug)
	fmt.Printf("  Port: %d\n", cfg.Port)
	fmt.Printf("  Timeout: %v\n", cfg.Timeout)

	// Advanced API with custom options
	manager := config.NewManager().
		WithDiscovery(config.Union | config.Direct).
		WithFlagStyle(config.FlatNest).
		WithValidation(true)

	// Generate sample configuration
	if sample, err := manager.GenerateSample("json"); err == nil {
		fmt.Printf("\nSample JSON configuration:\n%s\n", sample)
	}

	// Save current configuration
	if err := manager.SaveCurrent(&cfg, "current-config.yaml"); err != nil {
		log.Printf("Failed to save config: %v", err)
	} else {
		fmt.Println("Current configuration saved to current-config.yaml")
	}
}
