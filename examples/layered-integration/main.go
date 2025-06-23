package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-cfg/pkg/integration"
)

// DemoConfig shows layered integration with complex nested structure
type DemoConfig struct {
	// Service identification
	ServiceName string `json:"service_name" required:"true"`
	Version     string `json:"version" default:"1.0.0"`
	Environment string `json:"environment" default:"development"`

	// API server configuration
	API struct {
		Host      string        `json:"host" default:"localhost"`
		Port      int           `json:"port" default:"8080"`
		Timeout   time.Duration `json:"timeout" default:"30s"`
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

	// Observability
	Metrics struct {
		Enabled  bool   `json:"enabled" default:"true"`
		Endpoint string `json:"endpoint" default:"/metrics"`
		Port     int    `json:"port" default:"9090"`
	} `json:"metrics"`
}

func main() {
	var config DemoConfig

	// Enhanced help text explaining layered configuration
	helpText := `
Layered Integration Demo - go-cfg + AutoCfg-style Discovery

This example demonstrates layered configuration loading:

STAGE 1 - File Discovery (AutoCfg-style patterns):
  Searches standard locations for configuration files:
  • $AUTOCFG_FILENAME (environment override)
  • .main.json (current directory, program-name based)
  • ~/.config/main/config.json  
  • /etc/main/config.json

STAGE 2 - Flag Processing (go-cfg FlatNest):
  Generates clean, intuitive flags:
  • --service-name (required)
  • --api-host, --api-port, --api-timeout
  • --database-host, --database-name (required)
  • --metrics-enabled, --metrics-port

PRECEDENCE ORDER (lowest to highest):
  Configuration Files → Environment Variables → Command Flags

REQUIRED FIELDS:
  --service-name, --database-name, --database-username, --database-password

EXAMPLES:
  # Use defaults with required fields
  ./main --service-name=api --database-name=mydb \
         --database-username=user --database-password=secret

  # Override with environment variables
  export API_HOST=0.0.0.0
  export DATABASE_HOST=prod-db.internal
  ./main --service-name=api --database-name=mydb \
         --database-username=user --database-password=secret

  # Use configuration file
  echo '{"service_name":"api","database":{"name":"mydb","username":"user","password":"secret"}}' > .main.json
  ./main

  # Load specific config file
  AUTOCFG_FILENAME=/path/to/config.json ./main
`

	// Configure using layered integration
	if err := integration.LayeredConfigureWithHelp(&config, helpText); err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	// Display loaded configuration
	fmt.Println("=== Layered Integration Demo ===")
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

	// Show metrics config
	if config.Metrics.Enabled {
		fmt.Printf("Metrics: enabled on port %d at %s\n",
			config.Metrics.Port, config.Metrics.Endpoint)
	}

	// Show external services
	if len(config.Services) > 0 {
		fmt.Printf("External services: %v\n", config.Services)
	}

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
	} else if err := stor.Save("layered-demo-current.yaml"); err != nil {
		log.Printf("Could not save configuration: %v", err)
	} else {
		fmt.Println("✓ Current configuration saved to layered-demo-current.yaml")
	}
}
