package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
