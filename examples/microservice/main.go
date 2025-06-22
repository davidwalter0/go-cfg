// Microservice Configuration Example
//
// This example demonstrates complex configuration for a microservice application.
// It shows nested structures, maps, slices, and advanced go-cfg features including
// multiple service endpoints, feature flags, and observability configuration.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davidwalter0/go-cfg"
)

// MicroserviceConfig represents a comprehensive microservice configuration
type MicroserviceConfig struct {
	// Service identification
	ServiceName string `required:"true" json:"service_name"`
	Version     string `default:"unknown" json:"version"`
	Environment string `default:"development" json:"environment"`
	NodeID      string `json:"node_id"`

	// API server configuration
	API struct {
		Host      string        `default:"0.0.0.0" json:"host"`
		Port      int           `default:"8080" json:"port"`
		Timeout   time.Duration `default:"30s" json:"timeout"`
		RateLimit int           `default:"1000" json:"rate_limit"`
		EnableTLS bool          `default:"false" json:"enable_tls"`
	} `json:"api"`

	// Database configuration
	Database struct {
		Host         string        `default:"localhost" json:"host"`
		Port         int           `default:"5432" json:"port"`
		Name         string        `required:"true" json:"name"`
		Username     string        `required:"true" json:"username"`
		Password     string        `required:"true" json:"password"`
		MaxConns     int           `default:"10" json:"max_connections"`
		ConnTimeout  time.Duration `default:"5s" json:"connection_timeout"`
		QueryTimeout time.Duration `default:"30s" json:"query_timeout"`
		SSLMode      string        `default:"prefer" json:"ssl_mode"`
	} `json:"database"`

	// External service endpoints (maps demonstrate complex type support)
	Services map[string]string `json:"services"`

	// Feature flags (boolean map)
	Features map[string]bool `default:"auth:true,metrics:true,tracing:false,caching:true" json:"features"`

	// Service discovery
	Discovery struct {
		Enabled   bool     `default:"false" json:"enabled"`
		Endpoints []string `default:"consul:8500,etcd:2379" json:"endpoints"`
		Namespace string   `default:"default" json:"namespace"`
	} `json:"discovery"`

	// Observability configuration
	Metrics struct {
		Enabled  bool          `default:"true" json:"enabled"`
		Endpoint string        `default:"/metrics" json:"endpoint"`
		Port     int           `default:"9090" json:"port"`
		Interval time.Duration `default:"15s" json:"scrape_interval"`
	} `json:"metrics"`

	Logging struct {
		Level      string `default:"info" json:"level"`
		Format     string `default:"json" json:"format"`
		Output     string `default:"stdout" json:"output"`
		MaxSize    int    `default:"100" json:"max_size_mb"`
		MaxBackups int    `default:"3" json:"max_backups"`
	} `json:"logging"`

	Tracing struct {
		Enabled     bool    `default:"false" json:"enabled"`
		Endpoint    string  `json:"endpoint"`
		SampleRate  float64 `default:"0.1" json:"sample_rate"`
		ServiceName string  `json:"service_name"`
	} `json:"tracing"`

	// Resource limits
	Resources struct {
		MaxMemoryMB int           `default:"512" json:"max_memory_mb"`
		MaxCPU      int           `default:"2" json:"max_cpu_cores"`
		Timeout     time.Duration `default:"5m" json:"timeout"`
	} `json:"resources"`

	// Security settings
	Security struct {
		JWTSecret     string        `json:"jwt_secret"`
		JWTExpiration time.Duration `default:"24h" json:"jwt_expiration"`
		EnableCORS    bool          `default:"true" json:"enable_cors"`
		CORSOrigins   []string      `default:"*" json:"cors_origins"`
		RateLimit     struct {
			Enabled bool `default:"true" json:"enabled"`
			RPS     int  `default:"100" json:"requests_per_second"`
			Burst   int  `default:"200" json:"burst"`
		} `json:"rate_limit"`
	} `json:"security"`
}

func main() {
	var config MicroserviceConfig

	// Add comprehensive help text
	cfg.HelpText(`
Microservice Configuration Demo

This example demonstrates advanced go-cfg features including:
- Nested configuration structures
- Maps and slices
- Complex types (time.Duration, etc.)
- Feature flags
- Service discovery configuration
- Observability settings

Required Configuration:
  You must provide at minimum:
  --service-name, --database-name, --database-username, --database-password

Examples:
  # Minimal configuration
  ./microservice --service-name=user-service --database-name=users \
                  --database-username=user --database-password=secret

  # Full production configuration
  ./microservice --config-file=production.yaml

  # Development with overrides
  export SERVICE_NAME=user-service
  export DATABASE_PASSWORD=secret
  ./microservice --config-file=development.yaml --debug

Environment Variables:
  All nested fields can be set via environment variables:
  API_HOST, API_PORT, DATABASE_HOST, METRICS_ENABLED, etc.
`)

	// Parse configuration
	if err := cfg.FlatNest(&config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	text, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(text))

	// Validate required fields (temporarily disabled for testing)
	// if !cfg.Ok("service-name") {
	//	log.Fatal("Service name is required")
	// }
	// if !cfg.Ok("database-name") || !cfg.Ok("database-username") || !cfg.Ok("database-password") {
	//	log.Fatal("Database configuration is incomplete")
	// }

	// Set derived values
	if config.Tracing.ServiceName == "" {
		config.Tracing.ServiceName = config.ServiceName
	}

	// Display configuration
	displayConfiguration(config)

	// Simulate service startup
	fmt.Printf("\n🚀 Starting microservice: %s\n", config.ServiceName)
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("API endpoint: http://%s:%d\n", config.API.Host, config.API.Port)

	if config.Metrics.Enabled {
		fmt.Printf("Metrics endpoint: http://%s:%d%s\n",
			config.API.Host, config.Metrics.Port, config.Metrics.Endpoint)
	}

	// Display enabled features
	var enabledFeatures []string
	for feature, enabled := range config.Features {
		if enabled {
			enabledFeatures = append(enabledFeatures, feature)
		}
	}
	if len(enabledFeatures) > 0 {
		fmt.Printf("Enabled features: %s\n", strings.Join(enabledFeatures, ", "))
	}

	// Display external services
	if len(config.Services) > 0 {
		fmt.Printf("External services:\n")
		for name, endpoint := range config.Services {
			fmt.Printf("  %s: %s\n", name, endpoint)
		}
	}

	// Save current configuration for reference
	// Create a Stor from the config and save it
	stor := cfg.NewStor()
	if err := cfg.CopyOut(config, &stor); err != nil {
		log.Printf("Warning: Could not copy configuration: %v", err)
	} else if err := stor.Save("current-config.yaml"); err != nil {
		log.Printf("Warning: Could not save configuration: %v", err)
	} else {
		fmt.Printf("\n💾 Current configuration saved to: current-config.yaml\n")
	}

	fmt.Printf("\n✅ Microservice configuration completed successfully\n")
}

func displayConfiguration(config MicroserviceConfig) {
	fmt.Printf("\n📋 Microservice Configuration Summary:\n")
	fmt.Printf("════════════════════════════════════════\n")

	// Service info
	fmt.Printf("Service: %s v%s (%s)\n", config.ServiceName, config.Version, config.Environment)

	// API configuration
	fmt.Printf("\nAPI Configuration:\n")
	fmt.Printf("  Address: %s:%d\n", config.API.Host, config.API.Port)
	fmt.Printf("  Timeout: %v\n", config.API.Timeout)
	fmt.Printf("  Rate Limit: %d req/s\n", config.API.RateLimit)
	fmt.Printf("  TLS: %v\n", config.API.EnableTLS)

	// Database configuration (mask password)
	fmt.Printf("\nDatabase Configuration:\n")
	fmt.Printf("  Host: %s:%d\n", config.Database.Host, config.Database.Port)
	fmt.Printf("  Database: %s\n", config.Database.Name)
	fmt.Printf("  Username: %s\n", config.Database.Username)
	fmt.Printf("  Password: %s\n", maskPassword(config.Database.Password))
	fmt.Printf("  Max Connections: %d\n", config.Database.MaxConns)
	fmt.Printf("  SSL Mode: %s\n", config.Database.SSLMode)

	// Service discovery
	if config.Discovery.Enabled {
		fmt.Printf("\nService Discovery:\n")
		fmt.Printf("  Endpoints: %v\n", config.Discovery.Endpoints)
		fmt.Printf("  Namespace: %s\n", config.Discovery.Namespace)
	}

	// Security
	fmt.Printf("\nSecurity Configuration:\n")
	fmt.Printf("  JWT Expiration: %v\n", config.Security.JWTExpiration)
	fmt.Printf("  CORS Enabled: %v\n", config.Security.EnableCORS)
	if len(config.Security.CORSOrigins) > 0 {
		fmt.Printf("  CORS Origins: %v\n", config.Security.CORSOrigins)
	}
	fmt.Printf("  Rate Limiting: %v (%d RPS)\n",
		config.Security.RateLimit.Enabled, config.Security.RateLimit.RPS)

	// Resources
	fmt.Printf("\nResource Limits:\n")
	fmt.Printf("  Memory: %d MB\n", config.Resources.MaxMemoryMB)
	fmt.Printf("  CPU Cores: %d\n", config.Resources.MaxCPU)
	fmt.Printf("  Timeout: %v\n", config.Resources.Timeout)

	// Observability
	fmt.Printf("\nObservability:\n")
	fmt.Printf("  Metrics: %v\n", config.Metrics.Enabled)
	fmt.Printf("  Tracing: %v\n", config.Tracing.Enabled)
	fmt.Printf("  Log Level: %s\n", config.Logging.Level)
	fmt.Printf("  Log Format: %s\n", config.Logging.Format)
	cfg.Usage()
}

func maskPassword(password string) string {
	if password == "" {
		return "<not set>"
	}
	if len(password) <= 3 {
		return "***"
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}
