package integration

import (
	"github.com/davidwalter0/go-cfg"
)

// LayeredConfigureWithHelp provides layered integration with enhanced help text
func LayeredConfigureWithHelp(config interface{}, helpText string) error {
	// Set custom help text before processing
	if helpText != "" {
		cfg.HelpText(helpText)
	}

	return LayeredConfigure(config)
}

// ConfigureWithDefaults provides a simplified API with sensible defaults
func ConfigureWithDefaults(config interface{}) error {
	defaultHelp := generateDefaultHelp()
	return LayeredConfigureWithHelp(config, defaultHelp)
}

// Configure provides the simplest API for layered integration
func Configure(config interface{}) error {
	return LayeredConfigure(config)
}

func generateDefaultHelp() string {
	return `
Layered Configuration Loading

CONFIGURATION SOURCES (precedence low to high):
1. Configuration files (discovered automatically)
2. Environment variables (UPPER_SNAKE_CASE)
3. Command-line flags (kebab-case with clean nesting)

CONFIGURATION FILE DISCOVERY:
- $AUTOCFG_FILENAME (environment override)
- .{program-name}.json (current directory)
- ~/.config/{program-name}/config.json
- /etc/{program-name}/config.json

FLAGS:
All nested configuration generates clean, intuitive flags:
• Top-level fields: --field-name
• Nested fields: --parent-field-name (e.g., --api-host, --database-port)

ENVIRONMENT VARIABLES:
All flags can be set via environment variables using UPPER_SNAKE_CASE:
• --service-name → SERVICE_NAME
• --api-host → API_HOST
• --database-port → DATABASE_PORT

Use --help for detailed flag information.
`
}
