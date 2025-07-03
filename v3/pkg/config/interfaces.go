package config

import (
	"io"
	
	"github.com/davidwalter0/go-cfg/v3/pkg/discovery"
	"github.com/spf13/pflag"
)

// ConfigManager provides unified configuration management
type ConfigManager interface {
	// Core configuration loading
	Load(config interface{}) error
	LoadName(filename string, config interface{}) error
	
	// Core configuration saving
	Save(config interface{}) error
	SaveName(filename string, config interface{}) error

	// Builder pattern configuration
	WithDiscovery(mode DiscoveryMode) ConfigManager
	WithFlagStyle(style FlagStyle) ConfigManager
	WithValidation(enabled bool) ConfigManager
	WithEnvironmentPrefix(prefix string) ConfigManager
	WithConfigFormat(format ConfigFormat) ConfigManager
	WithTrace(writer io.Writer) ConfigManager

	// Advanced features
	GenerateSample(format string) ([]byte, error)
	GenerateHelp() string
	SaveCurrent(config interface{}, filename string) error
}

// DiscoveryEngine handles configuration file discovery and loading
type DiscoveryEngine interface {
	FindConfigFiles(opts discovery.DiscoveryOptions) ([]string, error)
	FindConfigFilesWithFormat(opts discovery.DiscoveryOptions, format interface{}) ([]string, error)
	LoadConfigFile(filename string, config interface{}) error
	LoadConfigFileWithFormat(filename string, config interface{}, format interface{}) error
	MergeConfigs(files []string, config interface{}) error
	MergeConfigsWithFormat(files []string, config interface{}, format interface{}) error
	SaveConfigFile(filename string, config interface{}, format interface{}) error
}

// FlagProcessor handles command line flag processing
type FlagProcessor interface {
	ProcessStruct(config interface{}, flagSet *pflag.FlagSet) error
	GenerateFlags(config interface{}) (*pflag.FlagSet, error)
	ApplyFlagValues(config interface{}, flagSet *pflag.FlagSet) error
}

// We use the flags package's NamingStrategy interface directly
// to avoid duplication and interface mismatch
