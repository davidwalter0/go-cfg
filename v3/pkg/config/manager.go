package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	programName := m.options.ProgramName
	if programName == "" {
		// Auto-detect from current directory if not set
		if wd, err := os.Getwd(); err == nil {
			programName = filepath.Base(wd)
		}
	}
	
	discoveryOpts := discovery.DiscoveryOptions{
		Mode:        discovery.DiscoveryMode(m.options.DiscoveryMode),
		SearchPaths: m.options.SearchPaths,
		EnvOverride: m.options.EnvOverride,
		ProgramName: programName,
	}

	if m.options.EnableTrace {
		fmt.Fprintf(m.options.TraceWriter, "=== Configuration Discovery Trace ===\n")
		fmt.Fprintf(m.options.TraceWriter, "Discovery mode: %v\n", m.options.DiscoveryMode)
		fmt.Fprintf(m.options.TraceWriter, "Config format: %s\n", m.options.ConfigFormat)
		fmt.Fprintf(m.options.TraceWriter, "Environment override: %s\n", m.options.EnvOverride)
		fmt.Fprintf(m.options.TraceWriter, "Program name: %s\n", m.options.ProgramName)
	}

	var configFiles []string
	var err error
	
	// Use format-aware discovery if a specific format is configured
	if m.options.ConfigFormat != "" {
		configFiles, err = m.discovery.FindConfigFilesWithFormat(discoveryOpts, m.options.ConfigFormat)
	} else {
		configFiles, err = m.discovery.FindConfigFiles(discoveryOpts)
	}
	
	if err != nil {
		return fmt.Errorf("config discovery failed: %w", err)
	}

	if m.options.EnableTrace {
		if len(configFiles) == 0 {
			fmt.Fprintf(m.options.TraceWriter, "No config files found\n")
		} else {
			fmt.Fprintf(m.options.TraceWriter, "Found %d config file(s):\n", len(configFiles))
			for i, file := range configFiles {
				fmt.Fprintf(m.options.TraceWriter, "  %d. %s\n", i+1, file)
			}
		}
	}

	if len(configFiles) > 0 {
		// Use format-aware merging if a specific format is configured
		if m.options.ConfigFormat != "" {
			if err := m.discovery.MergeConfigsWithFormat(configFiles, config, m.options.ConfigFormat); err != nil {
				return fmt.Errorf("config loading failed: %w", err)
			}
		} else {
			if err := m.discovery.MergeConfigs(configFiles, config); err != nil {
				return fmt.Errorf("config loading failed: %w", err)
			}
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

// LoadName loads configuration from a specific file with format detection by extension
func (m *Manager) LoadName(filename string, config interface{}) error {
	if m.options.EnableTrace {
		fmt.Fprintf(m.options.TraceWriter, "Loading config from file: %s\n", filename)
	}
	
	// Detect format from file extension
	format := m.detectFormatFromFilename(filename)
	return m.discovery.LoadConfigFileWithFormat(filename, config, format)
}

// Save saves configuration using default search paths and current format
func (m *Manager) Save(config interface{}) error {
	// Get the first writable path from standard search paths
	programName := m.options.ProgramName
	if programName == "" {
		if len(os.Args) > 0 {
			programName = filepath.Base(os.Args[0])
		} else {
			programName = "config"
		}
	}
	
	// Choose filename based on format
	var filename string
	switch m.options.ConfigFormat {
	case YAML:
		filename = fmt.Sprintf(".%s.yaml", programName)
	default:
		filename = fmt.Sprintf(".%s.json", programName)
	}
	
	if m.options.EnableTrace {
		fmt.Fprintf(m.options.TraceWriter, "Saving config to: %s (format: %s)\n", filename, m.options.ConfigFormat)
	}
	
	return m.discovery.SaveConfigFile(filename, config, m.options.ConfigFormat)
}

// SaveName saves configuration to a specific file with format detection by extension  
func (m *Manager) SaveName(filename string, config interface{}) error {
	if m.options.EnableTrace {
		fmt.Fprintf(m.options.TraceWriter, "Saving config to file: %s\n", filename)
	}
	
	// Detect format from file extension
	format := m.detectFormatFromFilename(filename)
	return m.discovery.SaveConfigFile(filename, config, format)
}

// detectFormatFromFilename determines format from file extension
func (m *Manager) detectFormatFromFilename(filename string) ConfigFormat {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		return YAML
	case ".json":
		return JSON
	default:
		// Fall back to configured format
		return m.options.ConfigFormat
	}
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

// WithConfigFormat configures the configuration file format
func (m *Manager) WithConfigFormat(format ConfigFormat) ConfigManager {
	newManager := *m
	newManager.options.ConfigFormat = format
	return &newManager
}

// WithTrace enables tracing with the specified writer
func (m *Manager) WithTrace(writer io.Writer) ConfigManager {
	newManager := *m
	newManager.options.EnableTrace = true
	if writer == nil {
		newManager.options.TraceWriter = os.Stdout
	} else {
		newManager.options.TraceWriter = writer
	}
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
