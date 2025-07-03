package config

import (
	"io"
	"strings"
)

// DiscoveryMode defines how configuration files are discovered and merged
type DiscoveryMode int

const (
	Union    DiscoveryMode = 1 << iota // Merge all found configs
	First                              // Stop on first found
	Direct                             // Load config files directly
	Indirect                           // Support AutoCfg pointer files
)

// String returns a human-readable representation of the DiscoveryMode
func (d DiscoveryMode) String() string {
	var modes []string
	
	if d&Union != 0 {
		modes = append(modes, "Union")
	}
	if d&First != 0 {
		modes = append(modes, "First")
	}
	if d&Direct != 0 {
		modes = append(modes, "Direct")
	}
	if d&Indirect != 0 {
		modes = append(modes, "Indirect")
	}
	
	if len(modes) == 0 {
		return "None"
	}
	
	return strings.Join(modes, " | ")
}

// FlagStyle defines how flags are generated and named
type FlagStyle int

const (
	FlatNest FlagStyle = iota // Clean flat naming: --api-host
	Nested                    // Hierarchical naming: --api.host
	Prefixed                  // Prefixed naming: --config-api-host
)

// ConfigFormat defines supported configuration file formats
type ConfigFormat string

const (
	JSON ConfigFormat = "json"
	YAML ConfigFormat = "yaml"
)

// DiscoveryOptions provides configuration for file discovery
type DiscoveryOptions struct {
	Mode        DiscoveryMode
	SearchPaths []string
	EnvOverride string
	ProgramName string
}

// ConfigOptions provides comprehensive configuration control
type ConfigOptions struct {
	// Discovery settings
	DiscoveryMode DiscoveryMode
	SearchPaths   []string
	EnvOverride   string
	ProgramName   string

	// Format settings
	ConfigFormat ConfigFormat
	
	// Flag processing settings
	FlagStyle         FlagStyle
	EnvironmentPrefix string
	RequiredFields    bool

	// Output and generation settings
	HelpTemplate   string
	SampleFormat   string
	ValidateOnLoad bool
	
	// Tracing settings
	TraceWriter io.Writer
	EnableTrace bool
}

// DefaultOptions provides sensible defaults for most use cases
func DefaultOptions() ConfigOptions {
	return ConfigOptions{
		DiscoveryMode:     Union | Direct,
		ConfigFormat:      JSON,
		FlagStyle:         FlatNest,
		EnvOverride:       "AUTOCFG_FILENAME",
		RequiredFields:    true,
		ValidateOnLoad:    true,
		SampleFormat:      "json",
		EnvironmentPrefix: "",
		EnableTrace:       false,
		TraceWriter:       nil,
	}
}
