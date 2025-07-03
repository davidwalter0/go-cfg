package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"gopkg.in/yaml.v3"
)

// ConfigFormat is imported from the calling package to avoid circular dependencies
// The actual type will be passed as a parameter

// Engine implements DiscoveryEngine interface
type Engine struct{}

// NewEngine creates a new discovery engine
func NewEngine() *Engine {
	return &Engine{}
}

// DiscoveryOptions provides configuration for file discovery
type DiscoveryOptions struct {
	Mode        DiscoveryMode
	SearchPaths []string
	EnvOverride string
	ProgramName string
}

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

// FindConfigFiles discovers configuration files based on options
func (e *Engine) FindConfigFiles(opts DiscoveryOptions) ([]string, error) {
	var configFiles []string

	// Check environment override first
	if opts.EnvOverride != "" {
		if envFile := os.Getenv(opts.EnvOverride); envFile != "" {
			if _, err := os.Stat(envFile); err == nil {
				configFiles = append(configFiles, envFile)
				// If using First mode with env override, return immediately
				if opts.Mode&First != 0 {
					return configFiles, nil
				}
			}
		}
	}

	// Use provided search paths or generate standard paths
	searchPaths := opts.SearchPaths
	if len(searchPaths) == 0 {
		searchPaths = e.generateStandardPaths(opts.ProgramName)
	}

	// Search through paths
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); err == nil {
			configFiles = append(configFiles, searchPath)
			// If using First mode, return immediately after finding first file
			if opts.Mode&First != 0 {
				return configFiles, nil
			}
		}
	}

	return configFiles, nil
}

// FindConfigFilesWithFormat discovers configuration files with format-aware paths
func (e *Engine) FindConfigFilesWithFormat(opts DiscoveryOptions, format interface{}) ([]string, error) {
	var configFiles []string

	// Check environment override first
	if opts.EnvOverride != "" {
		if envFile := os.Getenv(opts.EnvOverride); envFile != "" {
			if _, err := os.Stat(envFile); err == nil {
				configFiles = append(configFiles, envFile)
				// If using First mode with env override, return immediately
				if opts.Mode&First != 0 {
					return configFiles, nil
				}
			}
		}
	}

	// Use provided search paths or generate format-aware standard paths
	searchPaths := opts.SearchPaths
	if len(searchPaths) == 0 {
		searchPaths = e.generateStandardPathsWithFormat(opts.ProgramName, format)
	}

	// Search through paths
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); err == nil {
			configFiles = append(configFiles, searchPath)
			// If using First mode, return immediately after finding first file
			if opts.Mode&First != 0 {
				return configFiles, nil
			}
		}
	}

	return configFiles, nil
}

// LoadConfigFile loads a single configuration file
func (e *Engine) LoadConfigFile(filename string, config interface{}) error {
	// For indirect mode, check if this is a pointer file
	if isIndirect, actualPath, err := e.checkIndirectConfig(filename); err != nil {
		return err
	} else if isIndirect {
		filename = actualPath
	}

	// Load the configuration file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	// Expand environment variables in the content
	expanded := os.ExpandEnv(string(data))

	// Parse JSON configuration
	if err := json.Unmarshal([]byte(expanded), config); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", filename, err)
	}

	return nil
}

// LoadConfigFileWithFormat loads a configuration file with explicit format
func (e *Engine) LoadConfigFileWithFormat(filename string, config interface{}, format interface{}) error {
	formatStr := fmt.Sprintf("%v", format)
	// For indirect mode, check if this is a pointer file
	if isIndirect, actualPath, err := e.checkIndirectConfig(filename); err != nil {
		return err
	} else if isIndirect {
		filename = actualPath
	}

	// Load the configuration file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	// Expand environment variables in the content
	expanded := os.ExpandEnv(string(data))

	// Parse based on format
	switch formatStr {
	case "yaml":
		if err := yaml.Unmarshal([]byte(expanded), config); err != nil {
			return fmt.Errorf("failed to parse YAML config file %s: %w", filename, err)
		}
	case "json":
		fallthrough
	default:
		if err := json.Unmarshal([]byte(expanded), config); err != nil {
			return fmt.Errorf("failed to parse JSON config file %s: %w", filename, err)
		}
	}

	return nil
}

// SaveConfigFile saves a configuration file in the specified format
func (e *Engine) SaveConfigFile(filename string, config interface{}, format interface{}) error {
	formatStr := fmt.Sprintf("%v", format)
	var data []byte
	var err error

	// Marshal based on format
	switch formatStr {
	case "yaml":
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}
	case "json":
		fallthrough
	default:
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %w", err)
		}
	}

	// Write the file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", filename, err)
	}

	return nil
}

// MergeConfigs loads and merges multiple configuration files
func (e *Engine) MergeConfigs(files []string, config interface{}) error {
	for _, file := range files {
		if err := e.LoadConfigFile(file, config); err != nil {
			// Log error but continue with other files in Union mode
			fmt.Printf("Warning: Failed to load config file %s: %v\n", file, err)
		}
	}
	return nil
}

// MergeConfigsWithFormat loads and merges multiple configuration files with explicit format
func (e *Engine) MergeConfigsWithFormat(files []string, config interface{}, format interface{}) error {
	for _, file := range files {
		if err := e.LoadConfigFileWithFormat(file, config, format); err != nil {
			// Log error but continue with other files in Union mode
			fmt.Printf("Warning: Failed to load config file %s: %v\n", file, err)
		}
	}
	return nil
}

// generateStandardPaths creates AutoCfg-compatible search paths
func (e *Engine) generateStandardPaths(programName string) []string {
	if programName == "" {
		programName = e.getProgramName()
	}

	var paths []string

	// Standard search order
	// 1. System-wide configuration
	paths = append(paths, filepath.Join("/etc", programName, "config.json"))

	// 2. User configuration directory (XDG-compliant)
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".config", programName, "config.json"))
	}

	// 3. Local configuration (hidden file)
	paths = append(paths, fmt.Sprintf(".%s.json", programName))

	// 4. Current directory
	paths = append(paths, "config.json")

	return paths
}

// generateStandardPathsWithFormat creates format-aware search paths
func (e *Engine) generateStandardPathsWithFormat(programName string, format interface{}) []string {
	if programName == "" {
		programName = e.getProgramName()
	}
	
	formatStr := fmt.Sprintf("%v", format)
	var ext string
	switch formatStr {
	case "yaml":
		ext = "yaml"
	default:
		ext = "json"
	}

	var paths []string

	// Standard search order
	// 1. System-wide configuration
	paths = append(paths, filepath.Join("/etc", programName, fmt.Sprintf("config.%s", ext)))

	// 2. User configuration directory (XDG-compliant)
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".config", programName, fmt.Sprintf("config.%s", ext)))
	}

	// 3. Local configuration (hidden file)
	paths = append(paths, fmt.Sprintf(".%s.%s", programName, ext))

	// 4. Current directory
	paths = append(paths, fmt.Sprintf("config.%s", ext))

	return paths
}

// getProgramName extracts the program name for configuration file discovery
func (e *Engine) getProgramName() string {
	// Auto-detect from os.Args[0]
	if len(os.Args) > 0 {
		programPath := os.Args[0]
		programName := filepath.Base(programPath)

		// Remove common extensions
		for _, ext := range []string{".exe", ".bin"} {
			if strings.HasSuffix(programName, ext) {
				programName = strings.TrimSuffix(programName, ext)
				break
			}
		}

		return programName
	}

	return "app" // fallback
}

// checkIndirectConfig checks if a file is an AutoCfg indirect configuration file
func (e *Engine) checkIndirectConfig(filename string) (bool, string, error) {
	// Read the file to check if it's an indirect config
	data, err := os.ReadFile(filename)
	if err != nil {
		return false, "", err
	}

	// Try to parse as indirect config
	var indirectConfig struct {
		Path string            `json:"path"`
		Env  map[string]string `json:"env"`
	}

	if err := json.Unmarshal(data, &indirectConfig); err != nil {
		// Not a valid indirect config, treat as direct
		return false, filename, nil
	}

	// Check if it has the required "path" field
	if indirectConfig.Path == "" {
		// Not an indirect config, treat as direct
		return false, filename, nil
	}

	// Set environment variables if specified
	for key, value := range indirectConfig.Env {
		os.Setenv(key, value)
	}

	// Expand environment variables in the path
	actualPath := os.ExpandEnv(indirectConfig.Path)

	// Expand tilde to home directory
	if strings.HasPrefix(actualPath, "~/") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			actualPath = filepath.Join(homeDir, actualPath[2:])
		}
	}

	return true, actualPath, nil
}
