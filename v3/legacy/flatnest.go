package legacy

import (
	"github.com/davidwalter0/go-cfg/v3/pkg/config"
)

// FlatNest provides the modern flat-nesting configuration API
// This is the enhanced function that combines the best of legacy APIs
// with modern v3 capabilities.
func FlatNest(ptrs ...interface{}) error {
	manager := config.NewManager().
		WithDiscovery(config.Union | config.Direct).
		WithFlagStyle(config.FlatNest).
		WithValidation(true)

	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}

	frozen = true
	return nil
}

// FlatNestWithDiscovery provides enhanced configuration with AutoCfg-style discovery
func FlatNestWithDiscovery(mode config.DiscoveryMode, ptrs ...interface{}) error {
	manager := config.NewManager().
		WithDiscovery(mode).
		WithFlagStyle(config.FlatNest).
		WithValidation(true)

	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}

	frozen = true
	return nil
}

// FlatNestWithOptions provides the most flexible configuration API
func FlatNestWithOptions(opts config.ConfigOptions, ptrs ...interface{}) error {
	manager := config.NewManager().
		WithDiscovery(opts.DiscoveryMode).
		WithFlagStyle(opts.FlagStyle).
		WithValidation(opts.ValidateOnLoad).
		WithEnvironmentPrefix(opts.EnvironmentPrefix)

	for _, ptr := range ptrs {
		if err := manager.Load(ptr); err != nil {
			return err
		}
	}

	frozen = true
	return nil
}

// Quick provides a simple one-line configuration API using modern defaults
func Quick(ptr interface{}) error {
	return config.Load(ptr)
}