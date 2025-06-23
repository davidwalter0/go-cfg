package config // import "github.com/davidwalter0/go-cfg/v3/pkg/config"


FUNCTIONS

func Load(config interface{}) error
    Load provides a simple API using default options


TYPES

type ConfigManager interface {
	// Core configuration loading
	Load(config interface{}) error

	// Builder pattern configuration
	WithDiscovery(mode DiscoveryMode) ConfigManager
	WithFlagStyle(style FlagStyle) ConfigManager
	WithValidation(enabled bool) ConfigManager
	WithEnvironmentPrefix(prefix string) ConfigManager

	// Advanced features
	GenerateSample(format string) ([]byte, error)
	GenerateHelp() string
	SaveCurrent(config interface{}, filename string) error
}
    ConfigManager provides unified configuration management

func NewManager() ConfigManager
    NewManager creates a new configuration manager with default options

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
    ConfigOptions provides comprehensive configuration control

func DefaultOptions() ConfigOptions
    DefaultOptions provides sensible defaults for most use cases

type DiscoveryEngine interface {
	FindConfigFiles(opts discovery.DiscoveryOptions) ([]string, error)
	LoadConfigFile(filename string, config interface{}) error
	MergeConfigs(files []string, config interface{}) error
}
    DiscoveryEngine handles configuration file discovery and loading

type DiscoveryMode int
    DiscoveryMode defines how configuration files are discovered and merged

const (
	Union    DiscoveryMode = 1 << iota // Merge all found configs
	First                              // Stop on first found
	Direct                             // Load config files directly
	Indirect                           // Support AutoCfg pointer files
)
type DiscoveryOptions struct {
	Mode        DiscoveryMode
	SearchPaths []string
	EnvOverride string
	ProgramName string
}
    DiscoveryOptions provides configuration for file discovery

type FlagProcessor interface {
	ProcessStruct(config interface{}, flagSet *pflag.FlagSet) error
	GenerateFlags(config interface{}) (*pflag.FlagSet, error)
	ApplyFlagValues(config interface{}, flagSet *pflag.FlagSet) error
}
    FlagProcessor handles command line flag processing

type FlagStyle int
    FlagStyle defines how flags are generated and named

const (
	FlatNest FlagStyle = iota // Clean flat naming: --api-host
	Nested                    // Hierarchical naming: --api.host
	Prefixed                  // Prefixed naming: --config-api-host
)
type Manager struct {
	// Has unexported fields.
}
    Manager implements ConfigManager interface

func (m *Manager) GenerateHelp() string
    GenerateHelp creates help text for the configuration

func (m *Manager) GenerateSample(format string) ([]byte, error)
    GenerateSample creates a sample configuration file

func (m *Manager) Load(config interface{}) error
    Load processes configuration from files, environment, and command line

func (m *Manager) SaveCurrent(config interface{}, filename string) error
    SaveCurrent saves the current configuration to a file

func (m *Manager) WithDiscovery(mode DiscoveryMode) ConfigManager
    WithDiscovery configures discovery mode

func (m *Manager) WithEnvironmentPrefix(prefix string) ConfigManager
    WithEnvironmentPrefix configures environment variable prefix

func (m *Manager) WithFlagStyle(style FlagStyle) ConfigManager
    WithFlagStyle configures flag naming style

func (m *Manager) WithValidation(enabled bool) ConfigManager
    WithValidation configures validation

