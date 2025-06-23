package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SearchMode defines how configuration files are discovered and loaded
type SearchMode int

const (
	Union    SearchMode = 1 << iota // Merge all found configs (union dominance)
	First                           // Stop on first found config
	Direct                          // Load direct config files
	Indirect                        // Load via autocfg pointer files
)

// ConfigOptions defines configuration discovery options
type ConfigOptions struct {
	SearchMode  SearchMode // How to search for and load configs
	SearchPaths []string   // Custom search paths (overrides defaults)
	EnvOverride string     // Environment variable for config file override
	ProgramName string     // Program name for default path generation (auto-detected if empty)
}

// DefaultConfigOptions returns sensible defaults for configuration discovery
func DefaultConfigOptions() ConfigOptions {
	return ConfigOptions{
		SearchMode:  Union | Direct,     // Merge all found direct config files
		SearchPaths: nil,                // Use standard paths
		EnvOverride: "AUTOCFG_FILENAME", // Standard environment override
		ProgramName: "",                 // Auto-detect from os.Args[0]
	}
}

// AutoCfgCompatible returns options that match autocfg behavior exactly
func AutoCfgCompatible() ConfigOptions {
	return ConfigOptions{
		SearchMode:  Union | Direct, // AutoCfg's default behavior
		SearchPaths: nil,            // Use autocfg standard paths
		EnvOverride: "AUTOCFG_FILENAME",
		ProgramName: "",
	}
}

// ConfigDiscovery handles configuration file discovery and loading
type ConfigDiscovery struct {
	Options ConfigOptions
}

// NewConfigDiscovery creates a configuration discovery instance
func NewConfigDiscovery(opts ConfigOptions) *ConfigDiscovery {
	return &ConfigDiscovery{Options: opts}
}

// DiscoverAndLoad finds and loads configuration files based on search options
func (cd *ConfigDiscovery) DiscoverAndLoad(config interface{}) error {
	configFiles, err := cd.findConfigFiles()
	if err != nil {
		return fmt.Errorf("config discovery failed: %w", err)
	}

	if len(configFiles) == 0 {
		// No config files found - this is not an error
		return nil
	}

	// Load configuration files based on search mode
	if cd.Options.SearchMode&First != 0 {
		// First mode: load only the first found file
		return cd.loadConfigFile(configFiles[0], config)
	} else if cd.Options.SearchMode&Union != 0 {
		// Union mode: load all files with later ones overriding earlier ones
		for _, configFile := range configFiles {
			if err := cd.loadConfigFile(configFile, config); err != nil {
				return fmt.Errorf("failed to load config file %s: %w", configFile, err)
			}
		}
	}

	return nil
}

// findConfigFiles discovers configuration files based on search paths and modes
func (cd *ConfigDiscovery) findConfigFiles() ([]string, error) {
	var configFiles []string

	// Check environment override first
	if cd.Options.EnvOverride != "" {
		if envFile := os.Getenv(cd.Options.EnvOverride); envFile != "" {
			if _, err := os.Stat(envFile); err == nil {
				configFiles = append(configFiles, envFile)
				// If using First mode with env override, return immediately
				if cd.Options.SearchMode&First != 0 {
					return configFiles, nil
				}
			}
		}
	}

	// Get search paths
	searchPaths := cd.getSearchPaths()

	// Search for config files
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); err == nil {
			configFiles = append(configFiles, searchPath)
			// If using First mode, return immediately after finding first file
			if cd.Options.SearchMode&First != 0 {
				return configFiles, nil
			}
		}
	}

	return configFiles, nil
}

// getSearchPaths returns the list of paths to search for configuration files
func (cd *ConfigDiscovery) getSearchPaths() []string {
	if len(cd.Options.SearchPaths) > 0 {
		// Use custom search paths if provided
		return cd.Options.SearchPaths
	}

	// Generate standard AutoCfg-compatible search paths
	programName := cd.getProgramName()

	var paths []string

	// Standard search order (for Union mode)
	// For First mode, this will be reversed during search

	// 1. System-wide configuration
	paths = append(paths, filepath.Join("/etc", programName, "config.json"))

	// 2. User configuration directory (XDG-compliant)
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".config", programName, "config.json"))
	}

	// 3. Current directory - program-specific
	paths = append(paths, fmt.Sprintf(".%s.json", programName))

	// 4. Current directory - generic
	paths = append(paths, "config.json")

	// For First mode, reverse the search order
	if cd.Options.SearchMode&First != 0 {
		// Reverse the slice for First mode (environment override → local → user → system)
		for i, j := 0, len(paths)-1; i < j; i, j = i+1, j-1 {
			paths[i], paths[j] = paths[j], paths[i]
		}
	}

	return paths
}

// getProgramName extracts the program name for configuration file discovery
func (cd *ConfigDiscovery) getProgramName() string {
	if cd.Options.ProgramName != "" {
		return cd.Options.ProgramName
	}

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

// loadConfigFile loads a single configuration file
func (cd *ConfigDiscovery) loadConfigFile(filename string, config interface{}) error {
	// Check if this is an indirect config file (AutoCfg pointer file)
	if cd.Options.SearchMode&Indirect != 0 {
		if isIndirectConfig, actualPath, err := cd.checkIndirectConfig(filename); err != nil {
			return err
		} else if isIndirectConfig {
			filename = actualPath
		}
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

// checkIndirectConfig checks if a file is an AutoCfg indirect configuration file
func (cd *ConfigDiscovery) checkIndirectConfig(filename string) (bool, string, error) {
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

// FlatNestWithDiscovery provides enhanced FlatNest with AutoCfg-style discovery
func FlatNestWithDiscovery(opts ConfigOptions, ptrs ...interface{}) error {
	// Stage 1: Configuration discovery and loading
	discovery := NewConfigDiscovery(opts)

	for _, ptr := range ptrs {
		if err := discovery.DiscoverAndLoad(ptr); err != nil {
			return fmt.Errorf("config discovery failed: %w", err)
		}
	}

	// Stage 2: go-cfg flag processing with clean naming
	if err := FlatNest(ptrs...); err != nil {
		return fmt.Errorf("flag processing failed: %w", err)
	}

	return nil
}
