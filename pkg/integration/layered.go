package integration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidwalter0/go-cfg"
)

// LayeredConfigure implements the layered integration approach
// Stage 1: AutoCfg-style file discovery and loading
// Stage 2: go-cfg handles flag processing with clean naming
func LayeredConfigure(config interface{}) error {
	// Stage 1: File-based configuration discovery
	if err := loadConfigFiles(config); err != nil {
		log.Printf("File loading failed (continuing): %v", err)
	}

	// Stage 2: go-cfg flag processing with FlatNest
	if err := cfg.FlatNest(config); err != nil {
		return fmt.Errorf("flag processing failed: %w", err)
	}

	return nil
}

// loadConfigFiles implements AutoCfg-style file discovery and loading
func loadConfigFiles(config interface{}) error {
	// Get program name for config discovery
	programName := getProgramName()

	// AutoCfg search order:
	// 1. $AUTOCFG_FILENAME environment variable override
	// 2. .{program}.json (current directory)
	// 3. ~/.config/{program}/config.json
	// 4. /etc/{program}/config.json

	configFile := os.Getenv("AUTOCFG_FILENAME")
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			log.Printf("Loading configuration from AUTOCFG_FILENAME: %s", configFile)
			return loadJSONConfig(configFile, config)
		}
		log.Printf("AUTOCFG_FILENAME specified but file not found: %s", configFile)
	}

	// Search standard locations
	candidates := []string{
		fmt.Sprintf(".%s.json", programName),
		"config.json",
		".config.json",
	}

	// Add user config directory if available
	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(homeDir, ".config", programName, "config.json"),
		)
	}

	// Add system config directory
	candidates = append(candidates,
		filepath.Join("/etc", programName, "config.json"),
	)

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			log.Printf("Loading configuration from: %s", candidate)
			return loadJSONConfig(candidate, config)
		}
	}

	log.Printf("No configuration files found, using defaults and environment/flags")
	return nil // No config file found, continue with defaults
}

// loadJSONConfig loads configuration from a JSON file
func loadJSONConfig(filename string, config interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	// Expand environment variables in the JSON content
	expanded := os.ExpandEnv(string(data))

	if err := json.Unmarshal([]byte(expanded), config); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", filename, err)
	}

	return nil
}

// getProgramName extracts the program name for config file discovery
func getProgramName() string {
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
