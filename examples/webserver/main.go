// Basic Web Server Configuration Example
//
// This example demonstrates a typical web server configuration using go-cfg.
// It shows how to configure common web server settings like host, port, timeouts,
// and TLS certificates using environment variables, command-line flags, and defaults.
//
// Run examples:
//   go run main.go --help
//   go run main.go --host=0.0.0.0 --port=9000
//   export HOST=127.0.0.1 && go run main.go
//   go run main.go --config=server.yaml

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidwalter0/go-cfg"
)

// ServerConfig defines all configuration options for the web server
type ServerConfig struct {
	// Network settings
	Host string `default:"localhost" json:"host"`
	Port int    `default:"8080" json:"port"`

	// TLS configuration
	TLSCertFile string `json:"tls_cert_file"`
	TLSKeyFile  string `json:"tls_key_file"`
	EnableTLS   bool   `default:"false" json:"enable_tls"`

	// Timeout settings
	ReadTimeout  time.Duration `default:"30s" json:"read_timeout"`
	WriteTimeout time.Duration `default:"30s" json:"write_timeout"`
	IdleTimeout  time.Duration `default:"120s" json:"idle_timeout"`

	// Feature flags
	EnableMetrics bool `default:"true" json:"enable_metrics"`
	EnableLogging bool `default:"true" json:"enable_logging"`
	EnableCORS    bool `default:"false" json:"enable_cors"`

	// Operational settings
	LogLevel    string `default:"info" json:"log_level"`
	Environment string `default:"development" json:"environment"`
	Debug       bool   `default:"false" json:"debug"`

	// Optional configuration file
	ConfigFile string `json:"config_file"`
}

func main() {
	var config ServerConfig

	// Add custom help text
	cfg.HelpText(`
WebServer is a configurable HTTP server demonstration.

Environment Variables:
  All configuration can be set via environment variables using UPPER_SNAKE_CASE:
  HOST, PORT, ENABLE_TLS, READ_TIMEOUT, etc.

Configuration File:
  Use --config-file to load settings from JSON or YAML file.
  Command-line flags and environment variables override file settings.

Examples:
  # Start with defaults
  ./webserver

  # Custom host and port
  ./webserver --host=0.0.0.0 --port=9000

  # Enable TLS
  ./webserver --enable-tls --tls-cert-file=cert.pem --tls-key-file=key.pem

  # Production mode with config file
  ./webserver --config-file=production.yaml --environment=production
`)

	// Parse configuration from all sources
	if err := cfg.Flags(&config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Load configuration file if specified
	if config.ConfigFile != "" {
		stor := cfg.NewStor()
		if err := stor.Load(config.ConfigFile); err != nil {
			log.Printf("Warning: Could not load config file '%s': %v", config.ConfigFile, err)
		} else if err := cfg.CopyOut(stor, &config); err != nil {
			log.Printf("Warning: Could not copy loaded configuration: %v", err)
		} else {
			fmt.Printf("Loaded configuration from: %s\n", config.ConfigFile)
		}
	}

	// Validate required configuration for TLS
	if config.EnableTLS {
		if config.TLSCertFile == "" || config.TLSKeyFile == "" {
			log.Fatal("TLS enabled but certificate files not specified. Use --tls-cert-file and --tls-key-file")
		}
	}

	// Display configuration
	fmt.Printf("Server Configuration:\n")
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  TLS Enabled: %v\n", config.EnableTLS)
	if config.EnableTLS {
		fmt.Printf("  TLS Cert: %s\n", config.TLSCertFile)
		fmt.Printf("  TLS Key: %s\n", config.TLSKeyFile)
	}
	fmt.Printf("  Read Timeout: %v\n", config.ReadTimeout)
	fmt.Printf("  Write Timeout: %v\n", config.WriteTimeout)
	fmt.Printf("  Idle Timeout: %v\n", config.IdleTimeout)
	fmt.Printf("  Environment: %s\n", config.Environment)
	fmt.Printf("  Debug Mode: %v\n", config.Debug)
	fmt.Printf("  Features: metrics=%v, logging=%v, cors=%v\n", 
		config.EnableMetrics, config.EnableLogging, config.EnableCORS)

	// Create HTTP server with configured timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if config.EnableLogging {
			fmt.Printf("Request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		}

		if config.EnableCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		fmt.Fprintf(w, "Web Server Configuration Demo\n")
		fmt.Fprintf(w, "Environment: %s\n", config.Environment)
		fmt.Fprintf(w, "Debug Mode: %v\n", config.Debug)
		fmt.Fprintf(w, "Server Time: %s\n", time.Now().Format(time.RFC3339))
	})

	if config.EnableMetrics {
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "# Basic metrics endpoint\n")
			fmt.Fprintf(w, "webserver_uptime_seconds %d\n", time.Now().Unix())
			fmt.Fprintf(w, "webserver_environment_info{env=\"%s\"} 1\n", config.Environment)
		})
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	// Start server
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if config.EnableTLS {
		fmt.Printf("Starting HTTPS server on %s\n", address)
		log.Fatal(server.ListenAndServeTLS(config.TLSCertFile, config.TLSKeyFile))
	} else {
		fmt.Printf("Starting HTTP server on %s\n", address)
		log.Fatal(server.ListenAndServe())
	}
}