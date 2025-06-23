package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davidwalter0/go-cfg/v3/pkg/config"
)

// AdvancedConfig demonstrates complex configuration with nested structures
type AdvancedConfig struct {
	ServiceName string `json:"service_name" required:"true" help:"Service identifier"`
	Version     string `json:"version" default:"1.0.0" help:"Service version"`
	Environment string `json:"environment" default:"development" help:"Deployment environment"`

	// API configuration
	API struct {
		Host      string        `json:"host" default:"localhost" help:"API server host"`
		Port      int           `json:"port" default:"8080" help:"API server port"`
		Timeout   time.Duration `json:"timeout" default:"30s" help:"Request timeout"`
		EnableTLS bool          `json:"enable_tls" default:"false" help:"Enable TLS encryption"`
	} `json:"api" help:"API server configuration"`

	// Database configuration
	Database struct {
		Host     string `json:"host" default:"localhost" help:"Database host"`
		Port     int    `json:"port" default:"5432" help:"Database port"`
		Name     string `json:"name" required:"true" help:"Database name"`
		Username string `json:"username" required:"true" help:"Database username"`
		Password string `json:"password" required:"true" help:"Database password"`
	} `json:"database" help:"Database configuration"`

	// Feature flags
	Features map[string]bool `json:"features" help:"Feature toggles"`

	// External services
	Services map[string]string `json:"services" help:"External service URLs"`

	// Logging configuration
	Logging struct {
		Level  string `json:"level" default:"info" help:"Log level (debug, info, warn, error)"`
		Format string `json:"format" default:"json" help:"Log format (json, text)"`
		Output string `json:"output" default:"stdout" help:"Log output (stdout, file, syslog)"`
	} `json:"logging" help:"Logging configuration"`
}

func main() {
	var cfg AdvancedConfig

	fmt.Println("=== Advanced Configuration Example ===")

	// Create manager with custom settings
	manager := config.NewManager().
		WithDiscovery(config.Union | config.Direct).
		WithFlagStyle(config.FlatNest).
		WithValidation(true).
		WithEnvironmentPrefix("MYAPP")

	// Load configuration
	if err := manager.Load(&cfg); err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	// Display loaded configuration
	fmt.Printf("Service: %s v%s (%s)\n", cfg.ServiceName, cfg.Version, cfg.Environment)
	fmt.Printf("API: %s:%d (TLS: %v, Timeout: %v)\n",
		cfg.API.Host, cfg.API.Port, cfg.API.EnableTLS, cfg.API.Timeout)
	fmt.Printf("Database: %s@%s:%d/%s\n",
		cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	// Display features
	if len(cfg.Features) > 0 {
		fmt.Printf("Features enabled: ")
		for name, enabled := range cfg.Features {
			if enabled {
				fmt.Printf("%s ", name)
			}
		}
		fmt.Println()
	}

	// Display services
	if len(cfg.Services) > 0 {
		fmt.Printf("External services:\n")
		for name, url := range cfg.Services {
			fmt.Printf("  %s: %s\n", name, url)
		}
	}

	fmt.Printf("Logging: %s level, %s format, output to %s\n",
		cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output)

	// Generate documentation
	fmt.Println("\n=== Configuration Documentation ===")
	fmt.Println("Use --help to see all available flags")

	// Generate sample configurations
	if jsonSample, err := manager.GenerateSample("json"); err == nil {
		fmt.Printf("\n=== Sample JSON Configuration ===\n%s\n", jsonSample)
	}

	if yamlSample, err := manager.GenerateSample("yaml"); err == nil {
		fmt.Printf("\n=== Sample YAML Configuration ===\n%s\n", yamlSample)
	}

	// Save current configuration
	if err := manager.SaveCurrent(&cfg, "advanced-config.yaml"); err != nil {
		log.Printf("Failed to save config: %v", err)
	} else {
		fmt.Println("Current configuration saved to advanced-config.yaml")
	}
}
