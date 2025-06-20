package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/davidwalter0/go-cfg/pkg/flags"
	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

// ManualConfig demonstrates manual flag registration with more control
type ManualConfig struct {
	ServerURL    *types.URLValue
	Timeout      *types.DurationValue
	LogLevel     *types.LogLevelValue
	Environment  *types.EnvironmentValue
	DatabaseHost *types.HostPortValue
	Labels       *types.KeyValueValue
	Retries      int
	MaxSize      int64
	EnableTLS    bool
}

func main() {
	config := &ManualConfig{}
	
	// Create custom flag set
	flagSet := pflag.NewFlagSet("manual-example", pflag.ExitOnError)
	
	// Manual registration with full control
	
	// URL flag with validation
	flags.RegisterURLFlag(flagSet, "ServerURL", &config.ServerURL, 
		"Server endpoint URL (supports HTTP/HTTPS)")
	
	// Duration with multiple format support
	flags.RegisterDurationFlag(flagSet, "Timeout", &config.Timeout, 30*time.Second,
		"Request timeout (e.g., 30s, 1m30s, 90)")
	
	// Enum with predefined values
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	flags.RegisterEnumFlag(flagSet, "LogLevel", &config.LogLevel, validLogLevels, "info",
		"Log level (debug|info|warn|error|fatal)")
	
	// Environment enum
	validEnvs := []string{"development", "staging", "production", "test"}
	flags.RegisterEnumFlag(flagSet, "Environment", &config.Environment, validEnvs, "development",
		"Deployment environment")
	
	// Host:port combination
	config.DatabaseHost = types.NewHostPortValue("localhost", 5432)
	types.RegisterValue(flagSet, config.DatabaseHost, "DatabaseHost",
		"Database host and port (host:port format)")
	
	// Key-value pairs for labels
	flags.RegisterKeyValueFlag(flagSet, "Labels", &config.Labels,
		"Application labels in key=value format")
	
	// Standard types with environment defaults
	retries := 3
	if envRetries := types.GetEnvDefault("Retries"); envRetries != "" {
		fmt.Sscanf(envRetries, "%d", &retries)
	}
	flagSet.IntVar(&config.Retries, "retries", retries, "Number of retry attempts")
	
	maxSize := int64(1024 * 1024) // 1MB
	if envMaxSize := types.GetEnvDefault("MaxSize"); envMaxSize != "" {
		fmt.Sscanf(envMaxSize, "%d", &maxSize)
	}
	flagSet.Int64Var(&config.MaxSize, "max-size", maxSize, "Maximum size in bytes")
	
	enableTLS := false
	if envTLS := types.GetEnvDefault("EnableTLS"); envTLS != "" {
		enableTLS = envTLS == "true"
	}
	flagSet.BoolVar(&config.EnableTLS, "enable-tls", enableTLS, "Enable TLS encryption")
	
	// Custom flag with specialized validation
	portRange := types.NewPortRangeValue(8000, 9000)
	flagSet.Var(portRange, "port-range", "Allowed port range (start-end)")
	
	// Custom CIDR flag
	allowedSubnet := types.NewCIDRValue("10.0.0.0/8")
	flagSet.Var(allowedSubnet, "allowed-subnet", "Allowed subnet in CIDR notation")
	
	// Multiple IPs
	bindIPs := []*types.IPValue{
		types.NewIPValue("127.0.0.1"),
		types.NewIPValue("::1"),
	}
	for i, ip := range bindIPs {
		flagSet.Var(ip, fmt.Sprintf("bind-ip-%d", i+1), 
			fmt.Sprintf("Bind IP address %d", i+1))
	}
	
	// Help flag
	showHelp := flagSet.BoolP("help", "h", false, "Show help message")
	
	// Parse arguments
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}
	
	if *showHelp {
		fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Print(flagSet.FlagUsages())
		os.Exit(0)
	}
	
	// Validation
	if err := validateConfig(config, portRange, allowedSubnet, bindIPs); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	
	// Display configuration
	displayConfig(config, portRange, allowedSubnet, bindIPs)
}

func validateConfig(config *ManualConfig, portRange *types.PortRangeValue, 
	allowedSubnet *types.CIDRValue, bindIPs []*types.IPValue) error {
	
	// Validate required fields
	if config.ServerURL == nil || !config.ServerURL.IsSet() {
		return fmt.Errorf("server-url is required")
	}
	
	// Validate URL scheme
	if !config.ServerURL.IsAbsolute() {
		return fmt.Errorf("server-url must be absolute (include scheme)")
	}
	
	// Environment-specific validation
	if config.Environment != nil && config.Environment.IsProduction() {
		if config.LogLevel != nil && config.LogLevel.IsDebugEnabled() {
			log.Println("Warning: Debug logging enabled in production")
		}
		
		if !config.EnableTLS {
			return fmt.Errorf("TLS must be enabled in production environment")
		}
	}
	
	// Network validation
	if allowedSubnet.GetIPNet() != nil {
		for i, bindIP := range bindIPs {
			if bindIP.IsSet() && !allowedSubnet.Contains(bindIP.GetIP()) {
				return fmt.Errorf("bind-ip-%d (%s) is not in allowed subnet (%s)",
					i+1, bindIP.String(), allowedSubnet.String())
			}
		}
	}
	
	return nil
}

func displayConfig(config *ManualConfig, portRange *types.PortRangeValue,
	allowedSubnet *types.CIDRValue, bindIPs []*types.IPValue) {
	
	fmt.Println("=== Manual Registration Example ===")
	
	// URL configuration
	if config.ServerURL != nil && config.ServerURL.IsSet() {
		url := config.ServerURL.GetURL()
		fmt.Printf("Server URL: %s\n", url.String())
		fmt.Printf("  Scheme: %s\n", url.Scheme)
		fmt.Printf("  Host: %s\n", url.Host)
		fmt.Printf("  Path: %s\n", url.Path)
	}
	
	// Timing configuration
	if config.Timeout != nil {
		fmt.Printf("Timeout: %v\n", config.Timeout.GetDuration())
	}
	
	// Logging configuration
	if config.LogLevel != nil {
		level := config.LogLevel.GetLevel()
		fmt.Printf("Log Level: %s\n", level)
		fmt.Printf("  Debug enabled: %t\n", config.LogLevel.IsDebugEnabled())
		fmt.Printf("  Info enabled: %t\n", config.LogLevel.IsInfoEnabled())
	}
	
	// Environment configuration
	if config.Environment != nil {
		env := config.Environment.GetEnvironment()
		fmt.Printf("Environment: %s\n", env)
		fmt.Printf("  Production: %t\n", config.Environment.IsProduction())
		fmt.Printf("  Development: %t\n", config.Environment.IsDevelopment())
	}
	
	// Database configuration
	if config.DatabaseHost != nil {
		fmt.Printf("Database: %s (host: %s, port: %d)\n",
			config.DatabaseHost.String(),
			config.DatabaseHost.GetHost(),
			config.DatabaseHost.GetPort())
	}
	
	// Labels
	if config.Labels != nil && config.Labels.Len() > 0 {
		fmt.Printf("Labels (%d):\n", config.Labels.Len())
		for _, key := range config.Labels.Keys() {
			if value, exists := config.Labels.GetValue(key); exists {
				fmt.Printf("  %s = %s\n", key, value)
			}
		}
	}
	
	// Standard configurations
	fmt.Printf("Retries: %d\n", config.Retries)
	fmt.Printf("Max Size: %d bytes (%.2f MB)\n", 
		config.MaxSize, float64(config.MaxSize)/(1024*1024))
	fmt.Printf("TLS Enabled: %t\n", config.EnableTLS)
	
	// Network configuration
	if portRange.Size() > 0 {
		fmt.Printf("Port Range: %s (%d ports)\n", 
			portRange.String(), portRange.Size())
	}
	
	if allowedSubnet.GetIPNet() != nil {
		fmt.Printf("Allowed Subnet: %s\n", allowedSubnet.String())
	}
	
	fmt.Println("Bind IPs:")
	for i, bindIP := range bindIPs {
		if bindIP.IsSet() {
			fmt.Printf("  IP %d: %s (IPv4: %t, IPv6: %t)\n", 
				i+1, bindIP.String(), bindIP.IsIPv4(), bindIP.IsIPv6())
		}
	}
	
	// Environment variable suggestions
	fmt.Println("\n=== Environment Variable Examples ===")
	fmt.Println("export SERVER_URL='https://api.production.com'")
	fmt.Println("export TIMEOUT='2m'")
	fmt.Println("export LOG_LEVEL='warn'")
	fmt.Println("export ENVIRONMENT='production'")
	fmt.Println("export DATABASE_HOST='db.example.com:5432'")
	fmt.Println("export LABELS='service=api,version=v2.1.0,team=backend'")
	fmt.Println("export RETRIES=5")
	fmt.Println("export MAX_SIZE=10485760")  // 10MB
	fmt.Println("export ENABLE_TLS=true")
}

// Example usage:
//
// Full configuration:
//   ./manual_registration \
//     --server-url https://api.example.com/v1 \
//     --timeout 90s \
//     --log-level debug \
//     --environment development \
//     --database-host localhost:5432 \
//     --labels env=dev,service=test-api \
//     --retries 5 \
//     --max-size 5242880 \
//     --enable-tls \
//     --port-range 8000-9000 \
//     --allowed-subnet 10.0.0.0/16 \
//     --bind-ip-1 10.0.1.100 \
//     --bind-ip-2 10.0.1.101
//
// Production configuration:
//   export ENVIRONMENT=production
//   export SERVER_URL=https://prod-api.example.com
//   export ENABLE_TLS=true
//   export LOG_LEVEL=warn
//   ./manual_registration --database-host prod-db:5432