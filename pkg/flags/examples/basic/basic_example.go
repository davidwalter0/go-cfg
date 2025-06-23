package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidwalter0/go-cfg/pkg/flags"
	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

// AppConfig demonstrates various custom flag types
type AppConfig struct {
	// Basic types with env defaults
	Name  string `usage:"Application name" default:"myapp"`
	Port  int    `usage:"Server port" default:"8080"`
	Debug bool   `usage:"Enable debug mode" default:"false"`

	// Custom duration with human-friendly parsing
	Timeout   *types.DurationValue      `usage:"Request timeout"`
	KeepAlive *types.HumanDurationValue `usage:"Keep-alive duration"`

	// URL with validation
	APIEndpoint *types.URLValue `usage:"API endpoint URL"`
	DatabaseURL *types.URLValue `usage:"Database connection URL"`

	// Enum types
	LogLevel    *types.LogLevelValue    `usage:"Application log level"`
	Environment *types.EnvironmentValue `usage:"Deployment environment"`

	// Network configuration
	ListenAddr *types.HostPortValue `usage:"Listen address"`
	BindIP     *types.IPValue       `usage:"Bind IP address"`

	// Key-value pairs
	Labels *types.KeyValueValue `usage:"Application labels (key=value,...)"`
	Config *types.KeyValueValue `usage:"Additional config (key=value,...)"`

	// URL collections
	Mirrors *types.URLSliceValue `usage:"Mirror URLs"`
}

func main() {
	config := &AppConfig{
		// Initialize custom types with defaults
		Timeout:     types.NewDurationValue(30 * time.Second),
		KeepAlive:   types.NewHumanDurationValue(5 * time.Minute),
		APIEndpoint: types.NewURLValue("https://api.example.com"),
		DatabaseURL: types.NewURLValue("postgres://localhost:5432/myapp"),
		LogLevel:    types.NewLogLevelValue("info"),
		Environment: types.NewEnvironmentValue("development"),
		ListenAddr:  types.NewHostPortValue("0.0.0.0", 8080),
		BindIP:      types.NewIPValue("127.0.0.1"),
		Labels:      types.NewKeyValueValue(),
		Config:      types.NewKeyValueValue(),
		Mirrors:     types.NewURLSliceValue(),
	}

	// Create flag configuration
	flagConfig := flags.NewFlagConfig()

	// Register all types automatically
	if err := flagConfig.RegisterCustomTypes(config); err != nil {
		log.Fatalf("Failed to register flags: %v", err)
	}

	// Mark some fields as required
	flagConfig.SetRequired("DatabaseURL")

	// Add custom validators
	flagConfig.AddValidator("Port", func(value interface{}) error {
		if v, ok := value.(*int); ok && *v <= 0 {
			return fmt.Errorf("port must be positive")
		}
		return nil
	})

	// Parse command line arguments
	if err := flagConfig.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUsage:\n%s", flagConfig.Usage())
		os.Exit(1)
	}

	// Display configuration
	fmt.Println("=== Application Configuration ===")
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Printf("Debug: %t\n", config.Debug)
	fmt.Printf("Timeout: %v\n", config.Timeout.GetDuration())
	fmt.Printf("Keep-Alive: %v\n", config.KeepAlive.GetDuration())

	if config.APIEndpoint.IsSet() {
		fmt.Printf("API Endpoint: %s (absolute: %t)\n",
			config.APIEndpoint.String(), config.APIEndpoint.IsAbsolute())
	}

	if config.DatabaseURL.IsSet() {
		fmt.Printf("Database URL: %s\n", config.DatabaseURL.String())
	}

	fmt.Printf("Log Level: %s (debug enabled: %t)\n",
		config.LogLevel.GetLevel(), config.LogLevel.IsDebugEnabled())

	fmt.Printf("Environment: %s (production: %t)\n",
		config.Environment.GetEnvironment(), config.Environment.IsProduction())

	if config.ListenAddr != nil {
		fmt.Printf("Listen: %s (address: %s)\n",
			config.ListenAddr.String(), config.ListenAddr.GetAddress())
	}

	if config.BindIP.IsSet() {
		fmt.Printf("Bind IP: %s (IPv4: %t, IPv6: %t)\n",
			config.BindIP.String(), config.BindIP.IsIPv4(), config.BindIP.IsIPv6())
	}

	if config.Labels.Len() > 0 {
		fmt.Printf("Labels: %s\n", config.Labels.String())
		for _, key := range config.Labels.Keys() {
			if value, exists := config.Labels.GetValue(key); exists {
				fmt.Printf("  %s = %s\n", key, value)
			}
		}
	}

	if config.Config.Len() > 0 {
		fmt.Printf("Config: %s\n", config.Config.String())
	}

	if config.Mirrors.Len() > 0 {
		fmt.Printf("Mirrors (%d):\n", config.Mirrors.Len())
		for i, url := range config.Mirrors.GetURLs() {
			fmt.Printf("  %d: %s\n", i+1, url.String())
		}
	}

	// Environment variable examples
	fmt.Println("\n=== Environment Variable Examples ===")
	fmt.Println("Set these environment variables to provide defaults:")
	fmt.Println("  export NAME='my-service'")
	fmt.Println("  export PORT=9090")
	fmt.Println("  export DEBUG=true")
	fmt.Println("  export TIMEOUT='1m30s'")
	fmt.Println("  export KEEP_ALIVE='10 minutes'")
	fmt.Println("  export API_ENDPOINT='https://prod-api.example.com/v2'")
	fmt.Println("  export DATABASE_URL='postgres://user:pass@db:5432/prod'")
	fmt.Println("  export LOG_LEVEL=warn")
	fmt.Println("  export ENVIRONMENT=production")
	fmt.Println("  export LISTEN_ADDR='0.0.0.0:9090'")
	fmt.Println("  export BIND_IP='10.0.0.1'")
	fmt.Println("  export LABELS='service=api,version=v1.2.3,env=prod'")
	fmt.Println("  export MIRRORS='https://mirror1.example.com,https://mirror2.example.com'")
}

// Example usage:
//
// Basic usage:
//   ./basic_example --name myservice --port 9090 --database-url postgres://localhost/myapp
//
// With custom types:
//   ./basic_example \
//     --timeout 1m30s \
//     --keep-alive "10 minutes" \
//     --api-endpoint https://api.example.com/v2 \
//     --log-level debug \
//     --environment production \
//     --listen-addr 0.0.0.0:9090 \
//     --bind-ip 10.0.0.1 \
//     --labels service=api,version=v1.0.0 \
//     --mirrors https://mirror1.com,https://mirror2.com
//
// With environment variables:
//   export DATABASE_URL=postgres://localhost/myapp
//   export LOG_LEVEL=warn
//   export ENVIRONMENT=staging
//   ./basic_example --port 9090
