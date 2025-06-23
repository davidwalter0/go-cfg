package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davidwalter0/go-cfg/v3/pkg/discovery"
	"github.com/davidwalter0/go-cfg/v3/pkg/flags"
	"gopkg.in/yaml.v3"
)

// Manager implements ConfigManager interface
type Manager struct {
	options   ConfigOptions
	discovery DiscoveryEngine
	flags     FlagProcessor
	naming    flags.NamingStrategy
}

// NewManager creates a new configuration manager with default options
func NewManager() ConfigManager {
	opts := DefaultOptions()
	return &Manager{
		options:   opts,
		discovery: discovery.NewEngine(),
		flags:     flags.NewProcessor(),
		naming:    flags.NewFlatNestStrategy(),
	}
}

// Load processes configuration from files, environment, and command line
func (m *Manager) Load(config interface{}) error {
	// Stage 1: Load configuration files
	discoveryOpts := discovery.DiscoveryOptions{
		Mode:        discovery.DiscoveryMode(m.options.DiscoveryMode),
		SearchPaths: m.options.SearchPaths,
		EnvOverride: m.options.EnvOverride,
		ProgramName: m.options.ProgramName,
	}

	configFiles, err := m.discovery.FindConfigFiles(discoveryOpts)
	if err != nil {
		return fmt.Errorf("config discovery failed: %w", err)
	}

	if len(configFiles) > 0 {
		if err := m.discovery.MergeConfigs(configFiles, config); err != nil {
			return fmt.Errorf("config loading failed: %w", err)
		}
	}

	// Stage 2: Process command line flags
	flagSet, err := m.flags.GenerateFlags(config)
	if err != nil {
		return fmt.Errorf("flag generation failed: %w", err)
	}

	// Parse command line arguments
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("flag parsing failed: %w", err)
	}

	// Apply flag values to config
	if err := m.flags.ApplyFlagValues(config, flagSet); err != nil {
		return fmt.Errorf("flag application failed: %w", err)
	}

	// Stage 3: Validation (if enabled)
	if m.options.ValidateOnLoad {
		if err := m.validateConfig(config); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// WithDiscovery configures discovery mode
func (m *Manager) WithDiscovery(mode DiscoveryMode) ConfigManager {
	newManager := *m
	newManager.options.DiscoveryMode = mode
	return &newManager
}

// WithFlagStyle configures flag naming style
func (m *Manager) WithFlagStyle(style FlagStyle) ConfigManager {
	newManager := *m
	newManager.options.FlagStyle = style

	// Update naming strategy based on style
	switch style {
	case FlatNest:
		newManager.naming = flags.NewFlatNestStrategy()
	case Nested:
		newManager.naming = flags.NewNestedStrategy()
	case Prefixed:
		newManager.naming = flags.NewPrefixedStrategy()
	}

	return &newManager
}

// WithValidation configures validation
func (m *Manager) WithValidation(enabled bool) ConfigManager {
	newManager := *m
	newManager.options.ValidateOnLoad = enabled
	return &newManager
}

// WithEnvironmentPrefix configures environment variable prefix
func (m *Manager) WithEnvironmentPrefix(prefix string) ConfigManager {
	newManager := *m
	newManager.options.EnvironmentPrefix = prefix
	return &newManager
}

// GenerateSample creates a sample configuration file
func (m *Manager) GenerateSample(format string) ([]byte, error) {
	// This would be implemented to generate sample configs
	// For now, return a placeholder
	sample := map[string]interface{}{
		"service_name": "example-service",
		"version":      "1.0.0",
		"api": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
	}

	switch format {
	case "json":
		return json.MarshalIndent(sample, "", "  ")
	case "yaml":
		return yaml.Marshal(sample)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// GenerateHelp creates help text for the configuration
func (m *Manager) GenerateHelp() string {
	// This would generate comprehensive help text
	// For now, return a placeholder
	return "Configuration help would be generated here"
}

// SaveCurrent saves the current configuration to a file
func (m *Manager) SaveCurrent(config interface{}, filename string) error {
	// Determine format from filename extension
	var data []byte
	var err error

	if filename[len(filename)-5:] == ".yaml" || filename[len(filename)-4:] == ".yml" {
		data, err = yaml.Marshal(config)
	} else {
		data, err = json.MarshalIndent(config, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("marshaling failed: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// validateConfig performs validation on the loaded configuration
func (m *Manager) validateConfig(config interface{}) error {
	// This would implement comprehensive validation
	// For now, just return success
	return nil
}

// Load provides a simple API using default options
func Load(config interface{}) error {
	manager := NewManager()
	return manager.Load(config)
}
