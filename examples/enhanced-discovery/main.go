package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidwalter0/go-cfg"
)

// Duration is a custom type that can unmarshal from both JSON strings and nanoseconds
type Duration time.Duration

// UnmarshalJSON implements json.Unmarshaler for Duration
func (d *Duration) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first (e.g., "30s")
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		parsed, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		*d = Duration(parsed)
		return nil
	}
	
	// Try to unmarshal as nanoseconds (backwards compatibility)
	var ns int64
	if err := json.Unmarshal(data, &ns); err != nil {
		return err
	}
	*d = Duration(time.Duration(ns))
	return nil
}

// String implements fmt.Stringer
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Set implements flag.Value interface for command line flag parsing
func (d *Duration) Set(s string) error {
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}

// Type returns the type name for help text
func (d *Duration) Type() string {
	return "duration"
}

// ServiceConfig demonstrates enhanced discovery with nested configuration
type ServiceConfig struct {
	// Service identification
	ServiceName string `json:"service_name" required:"true"`
	Version     string `json:"version" default:"1.0.0"`
	Environment string `json:"environment" default:"development"`

	// API configuration
	API struct {
		Host      string        `json:"host" default:"localhost"`
		Port      int           `json:"port" default:"8080"`
		Timeout   Duration `json:"timeout" default:"30s"`
		EnableTLS bool          `json:"enable_tls" default:"false"`
	} `json:"api"`

	// Database configuration
	Database struct {
		Host     string `json:"host" default:"localhost"`
		Port     int    `json:"port" default:"5432"`
		Name     string `json:"name" required:"true"`
		Username string `json:"username" required:"true"`
		Password string `json:"password" required:"true"`
	} `json:"database"`

	// Feature flags
	Features map[string]bool `json:"features" default:"auth:true,metrics:true,debug:false"`

	// External services
	Services map[string]string `json:"services"`

	// Logging configuration
	Logging struct {
		Level  string `json:"level" default:"info"`
		Format string `json:"format" default:"json"`
		Output string `json:"output" default:"stdout"`
	} `json:"logging"`
}

func main() {
	var config ServiceConfig

	// Enhanced help text explaining discovery modes
	cfg.HelpText(`
Enhanced Discovery Demo - go-cfg with AutoCfg-style Configuration Discovery

This example demonstrates advanced configuration discovery patterns:

DISCOVERY MODES:
  Union Mode (default):  Merges all found configuration files
  First Mode:           Stops at first found configuration file
  Direct Mode:          Loads configuration files directly  
  Indirect Mode:        Supports AutoCfg pointer files

SEARCH LOCATIONS (Union mode order):
  1. /etc/enhanced-discovery/config.json
  2. ~/.config/enhanced-discovery/config.json
  3. .enhanced-discovery.json (current directory)
  4. config.json (current directory)
  5. $AUTOCFG_FILENAME (environment override)

CONFIGURATION PRECEDENCE (lowest to highest):
  System Config → User Config → Local Config → Environment Override → Environment Variables → Command Flags

EXAMPLES:

  # Use default discovery (Union + Direct mode)
  ./enhanced-discovery --service-name=api --database-name=mydb \
                       --database-username=user --database-password=secret

  # Use First mode (stop on first found config)
  DISCOVERY_MODE=first ./enhanced-discovery --service-name=api

  # Use custom config file
  AUTOCFG_FILENAME=/path/to/config.json ./enhanced-discovery

  # Use indirect configuration (AutoCfg pointer file)
  echo '{"path":"/secure/config.json","env":{"DB_HOST":"prod-db"}}' > .enhanced-discovery.json
  DISCOVERY_MODE=indirect ./enhanced-discovery
`)

	// Determine discovery mode from environment
	discoveryMode := getDiscoveryMode()

	// Configure using enhanced discovery
	if err := cfg.FlatNestWithDiscovery(discoveryMode, &config); err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	// Display loaded configuration
	displayConfiguration(config)
}

// getDiscoveryMode determines the discovery mode based on environment variables
func getDiscoveryMode() cfg.ConfigOptions {
	mode := cfg.DefaultConfigOptions() // Start with defaults

	// Allow customization via environment
	switch os.Getenv("DISCOVERY_MODE") {
	case "first":
		mode.SearchMode = cfg.First | cfg.Direct
		fmt.Println("🔍 Using First mode: will stop at first found config file")
	case "indirect":
		mode.SearchMode = cfg.Union | cfg.Indirect
		fmt.Println("🔍 Using Indirect mode: will process AutoCfg pointer files")
	case "union":
		mode.SearchMode = cfg.Union | cfg.Direct
		fmt.Println("🔍 Using Union mode: will merge all found config files")
	default:
		fmt.Println("🔍 Using default Union + Direct mode")
	}

	return mode
}

func displayConfiguration(config ServiceConfig) {
	fmt.Println("\n=== Enhanced Discovery Demo Results ===")
	fmt.Printf("Service: %s v%s (%s)\n", config.ServiceName, config.Version, config.Environment)
	fmt.Printf("API: %s:%d (TLS: %v, Timeout: %v)\n",
		config.API.Host, config.API.Port, config.API.EnableTLS, config.API.Timeout)
	fmt.Printf("Database: %s@%s:%d/%s\n",
		config.Database.Username, config.Database.Host, config.Database.Port, config.Database.Name)

	// Show features
	if len(config.Features) > 0 {
		fmt.Printf("Features: ")
		for feature, enabled := range config.Features {
			if enabled {
				fmt.Printf("%s ", feature)
			}
		}
		fmt.Println()
	}

	// Show external services
	if len(config.Services) > 0 {
		fmt.Printf("External services: %v\n", config.Services)
	}

	// Show logging config
	fmt.Printf("Logging: %s level, %s format, output to %s\n",
		config.Logging.Level, config.Logging.Format, config.Logging.Output)

	// JSON dump for verification
	fmt.Println("\n=== Configuration JSON ===")
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	// Save current configuration
	fmt.Println("\n=== Configuration Persistence ===")
	stor := cfg.NewStor()
	if err := cfg.CopyOut(config, &stor); err != nil {
		log.Printf("Could not copy configuration: %v", err)
	} else if err := stor.Save("enhanced-discovery-current.yaml"); err != nil {
		log.Printf("Could not save configuration: %v", err)
	} else {
		fmt.Println("✓ Current configuration saved to enhanced-discovery-current.yaml")
	}
}