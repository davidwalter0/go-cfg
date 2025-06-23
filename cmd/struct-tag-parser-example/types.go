package main

// DatabaseConfig represents database connection configuration
// These usage tags would be generated from the inline comments above
type DatabaseConfig struct {
	Host     string `usage:"Database server hostname" default:"localhost"`
	Port     int    `usage:"Database server port" default:"5432"`
	Username string `usage:"Database username" default:"postgres"`
	Password string `usage:"Database password (required in production)"`
	Database string `usage:"Database name" default:"myapp"`
	MaxConns int    `usage:"Maximum database connections" default:"25"`
	Timeout  string `usage:"Connection timeout duration" default:"30s"`
	SSLMode  string `usage:"SSL connection mode" default:"prefer"`
	ReadOnly bool   `usage:"Enable read-only database mode"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host         string `usage:"Server bind address" default:"0.0.0.0"`
	Port         int    `usage:"Server listen port" default:"8080"`
	ReadTimeout  string `usage:"HTTP read timeout" default:"30s"`
	WriteTimeout string `usage:"HTTP write timeout" default:"30s"`
	IdleTimeout  string `usage:"HTTP idle timeout" default:"120s"`
	EnableTLS    bool   `usage:"Enable TLS/HTTPS support"`
	CertFile     string `usage:"TLS certificate file path"`
	KeyFile      string `usage:"TLS private key file path"`
}

// MonitoringConfig represents monitoring and observability settings
type MonitoringConfig struct {
	EnableMetrics   bool   `usage:"Enable Prometheus metrics collection" default:"true"`
	MetricsPort     int    `usage:"Metrics server port" default:"9090"`
	EnableTracing   bool   `usage:"Enable distributed tracing"`
	TracingEndpoint string `usage:"Jaeger tracing endpoint URL"`
	LogLevel        string `usage:"Application log level" default:"info"`
	LogFormat       string `usage:"Log output format (json|text)" default:"json"`
}
