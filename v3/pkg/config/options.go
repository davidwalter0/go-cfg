package config

// DiscoveryMode defines how configuration files are discovered and merged
type DiscoveryMode int

const (
	Union    DiscoveryMode = 1 << iota // Merge all found configs
	First                              // Stop on first found
	Direct                             // Load config files directly
	Indirect                           // Support AutoCfg pointer files
)

// FlagStyle defines how flags are generated and named
type FlagStyle int

const (
	FlatNest FlagStyle = iota // Clean flat naming: --api-host
	Nested                    // Hierarchical naming: --api.host
	Prefixed                  // Prefixed naming: --config-api-host
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

	// Flag processing settings
	FlagStyle         FlagStyle
	EnvironmentPrefix string
	RequiredFields    bool

	// Output and generation settings
	HelpTemplate   string
	SampleFormat   string
	ValidateOnLoad bool
}

// DefaultOptions provides sensible defaults for most use cases
func DefaultOptions() ConfigOptions {
	return ConfigOptions{
		DiscoveryMode:     Union | Direct,
		FlagStyle:         FlatNest,
		EnvOverride:       "AUTOCFG_FILENAME",
		RequiredFields:    true,
		ValidateOnLoad:    true,
		SampleFormat:      "json",
		EnvironmentPrefix: "",
	}
}
