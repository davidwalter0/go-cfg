// Package flags provides a modern replacement for the vendor go-flag implementation
// This compatibility layer allows existing code to use the new implementation
// with minimal changes.
package flags

import (
	"fmt"
	"os"
	"sync"

	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

// Global flag configuration instance for compatibility
var globalFlagConfig *FlagConfig
var globalMutex sync.RWMutex

// GetGlobalFlagConfig returns the global flag configuration instance
func GetGlobalFlagConfig() *FlagConfig {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return globalFlagConfig
}

// SetGlobalFlagConfig sets the global flag configuration instance
func SetGlobalFlagConfig(fc *FlagConfig) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalFlagConfig = fc
}

// Flags is a compatibility wrapper for the old cfg.Flags() function
// It registers the given struct with the global flag configuration and parses command line arguments
func Flags(config interface{}) error {
	flagConfig := NewFlagConfig()

	err := flagConfig.RegisterCustomTypes(config)
	if err != nil {
		return err
	}

	// Store globally for other compatibility functions
	SetGlobalFlagConfig(flagConfig)

	// Parse with os.Args (excluding program name)
	args := os.Args[1:]
	return flagConfig.Parse(args)
}

// Usage returns the usage string for all registered flags
// This is a compatibility wrapper for the old cfg.Usage() function
func Usage() string {
	fc := GetGlobalFlagConfig()
	if fc == nil {
		return "No flags registered"
	}
	return fc.Usage()
}

// IsSet returns whether a flag has been explicitly set
// This is a compatibility wrapper for the old cfg.IsSet() function
func IsSet(name string) bool {
	fc := GetGlobalFlagConfig()
	if fc == nil {
		return false
	}

	// Convert to kebab-case as our system uses kebab-case flag names
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(name)

	flag := fc.FlagSet.Lookup(flagName)
	if flag == nil {
		return false
	}

	return flag.Changed
}

// Ok is an alias for IsSet for compatibility
// This is a compatibility wrapper for the old cfg.Ok() function
func Ok(name string) bool {
	return IsSet(name)
}

// Required returns whether a flag is marked as required
// This is a compatibility wrapper for the old cfg.Required() function
func Required(name string) bool {
	fc := GetGlobalFlagConfig()
	if fc == nil {
		return false
	}

	// Convert to kebab-case as our system uses kebab-case flag names
	nc := &types.NamingConventions{}
	flagName := nc.ToKebabCase(name)

	return fc.Required[flagName]
}

// NestWrap is a compatibility wrapper for the old cfg.NestWrap() function
// This function registers multiple structs with nested prefixes
func NestWrap(prefix string, configs ...interface{}) error {
	flagConfig := NewFlagConfig()

	// Register each config
	for _, config := range configs {
		err := flagConfig.RegisterCustomTypes(config)
		if err != nil {
			return err
		}
	}

	// Store globally for other compatibility functions
	SetGlobalFlagConfig(flagConfig)

	return nil
}

// HelpText sets help text for the application
// This is a compatibility wrapper function
func HelpText(text string) {
	// For compatibility - in the new system, help text can be handled
	// through the usage function or custom help formatting
}

// AddValidator adds a custom validator for a field
// This extends the compatibility layer with validation support
func AddValidator(fieldName string, validator func(interface{}) error) {
	fc := GetGlobalFlagConfig()
	if fc != nil {
		fc.AddValidator(fieldName, validator)
	}
}

// SetRequired marks a field as required
// This extends the compatibility layer with required field support
func SetRequired(fieldName string) {
	fc := GetGlobalFlagConfig()
	if fc != nil {
		fc.SetRequired(fieldName)
	}
}

// Parse parses the command line arguments using the global flag configuration
// This provides an explicit parsing function for the compatibility layer
func Parse(args []string) error {
	fc := GetGlobalFlagConfig()
	if fc == nil {
		return fmt.Errorf("no flag configuration available - call Flags() first")
	}

	return fc.Parse(args)
}
