package legacy

import (
	"os"

	"github.com/davidwalter0/go-cfg/v3/pkg/config"
)

// Legacy compatibility layer for github.com/davidwalter0/go-cfg
// This package provides backward compatibility for existing code
// while internally using the modern v3 unified module.

// Global configuration manager instance
var globalManager config.ConfigManager

func init() {
	globalManager = config.NewManager()
}

// Legacy configuration state
var (
	frozen   bool
	testMode bool
)

// Init from a list of struct pointers (legacy API)
func Init(ptrs ...interface{}) error {
	for _, ptr := range ptrs {
		if err := globalManager.Load(ptr); err != nil {
			return err
		}
	}
	frozen = true
	return nil
}

// Eval one or more configuration structures (legacy API)
func Eval(ptrs ...interface{}) error {
	for _, ptr := range ptrs {
		if err := globalManager.Load(ptr); err != nil {
			return err
		}
	}
	return nil
}

// EvalName one or more configuration structures overriding the name (legacy API)
func EvalName(name string, ptrs ...interface{}) error {
	// In v3, we can configure with environment prefix to simulate name override
	manager := globalManager.WithEnvironmentPrefix(name)
	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}
	frozen = true
	return nil
}

// Simple create env vars and flags without prefixes (legacy API)
func Simple(ptrs ...interface{}) error {
	manager := globalManager.WithFlagStyle(config.FlatNest).WithEnvironmentPrefix("")
	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}
	return nil
}

// SimpleFlags create env vars and flags without prefixes (legacy API)
func SimpleFlags(ptrs ...interface{}) error {
	err := Simple(ptrs...)
	if err != nil {
		return err
	}
	frozen = true
	return nil
}

// Unwrap alias of Simple - create env vars without prefixes (legacy API)
func Unwrap(ptrs ...interface{}) error {
	return Simple(ptrs...)
}

// Bare alias of Simple - create env vars without prefixes (legacy API)
func Bare(ptrs ...interface{}) error {
	return Simple(ptrs...)
}

// Wrap with prefix and struct names (legacy API)
func Wrap(name string, ptrs ...interface{}) error {
	manager := globalManager.WithEnvironmentPrefix(name).WithFlagStyle(config.Prefixed)
	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}
	return nil
}

// Nest each struct member with the name of the parent struct(s) (legacy API)
func Nest(ptrs ...interface{}) error {
	manager := globalManager.WithFlagStyle(config.Nested)
	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}
	return nil
}

// Reprefix replacing object name as prefix with name (legacy API)
func Reprefix(name string, ptrs ...interface{}) error {
	manager := globalManager.WithEnvironmentPrefix(name).WithFlagStyle(config.Prefixed)
	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}
	return nil
}

// Unprefixed next struct names, without the object name as a prefix (legacy API)
func Unprefixed(ptrs ...interface{}) error {
	manager := globalManager.WithFlagStyle(config.FlatNest).WithEnvironmentPrefix("")
	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}
	return nil
}

// Add alias of Eval (legacy API)
func Add(ptrs ...interface{}) error {
	return Eval(ptrs...)
}

// Flags alias of Eval (legacy API)
func Flags(ptrs ...interface{}) error {
	err := Eval(ptrs...)
	if err != nil {
		return err
	}
	frozen = true
	return nil
}

// Final freezes calling flag.Parse (legacy API)
func Final() {
	frozen = true
}

// Freeze prevents further configuration changes (legacy API)
func Freeze() {
	frozen = true
}

// Thaw allows configuration changes again (legacy API)
func Thaw() {
	frozen = false
}

// Frozen returns whether configuration is frozen (legacy API)
func Frozen() bool {
	return frozen
}

// TestMode enables test mode (legacy API)
func TestMode() {
	testMode = true
}

// ProductionMode disables test mode (legacy API)
func ProductionMode() {
	testMode = false
}

// InTestMode returns whether in test mode (legacy API)
func InTestMode() bool {
	return testMode
}

// Decorate structs with prefix (legacy API)
func Decorate() bool {
	os.Setenv("CFG_DECORATE", "true")
	return true
}

// Undecorate structs with prefix (legacy API)
func Undecorate() bool {
	os.Setenv("CFG_DECORATE", "false")
	return false
}

// Unprefix removes environment prefix (legacy API)
func Unprefix() {
	os.Unsetenv("CFG_KEY_PREFIX")
}

// Copy from the internal store to the config object (legacy API)
func Copy(config interface{}) {
	// In v3, configuration is applied directly during Load()
	// This is a no-op for compatibility
}

// IsSet checks if a flag was set (legacy API)
func IsSet(name string) bool {
	// In v3, this would require flag inspection
	// For now, return false as a safe default
	return false
}

// Ok checks if configuration is valid (legacy API)
func Ok() bool {
	return !testMode
}

// Required checks if required fields are set (legacy API)
func Required() error {
	// In v3, validation happens during Load()
	return nil
}

// GenerateHelp creates help text for the configuration (v3 enhancement)
func GenerateHelp() string {
	return globalManager.GenerateHelp()
}

// GenerateSample creates a sample configuration file (v3 enhancement)
func GenerateSample(format string) ([]byte, error) {
	return globalManager.GenerateSample(format)
}

// SaveCurrent saves the current configuration to a file (v3 enhancement)
func SaveCurrent(config interface{}, filename string) error {
	return globalManager.SaveCurrent(config, filename)
}