// This file implements help text generation and usage formatting for go-cfg.
// It provides customizable help output that integrates with pflag's usage system
// to display comprehensive flag documentation along with go-cfg specific features.
//
// Features:
//   - Automatic program name detection from os.Args[0]
//   - Custom help text injection via HelpText()
//   - Standard pflag flag help formatting
//   - Environment variable prefix documentation
//   - Consistent help output formatting
//
// The help system is automatically configured during package initialization
// and integrates with pflag's Usage function interface. Custom help text
// appears between the usage line and flag definitions.
//
// Integration:
//
//	The help system automatically sets pflag's Usage function during initialization,
//	ensuring that help output is consistent whether triggered by --help flags or
//	directly by calling the usage function.
package cfg

import (
	"fmt"
	"os"
	"strings"
)

// PrefixOverrideText documents the CFG_KEY_PREFIX environment variable behavior.
// This text explains how users can override the default prefixing behavior to
// customize environment variable and flag naming across their application.
//
// The documentation covers:
//   - How CFG_KEY_PREFIX affects environment variable names
//   - How CFG_KEY_PREFIX affects flag names
//   - Examples with and without the override
//   - Empty string behavior for minimal prefixing
//
// This text can be included in help output or documentation to explain
// the prefix override functionality to end users.
var PrefixOverrideText = `
The environment variable CFG_KEY_PREFIX overrides the prefix
of each environment variable to value. The value will be converted to
upper case like ex - EX

For a struct 

type App struct {
   i int
   s string
}

When CFG_KEY_PREFIX=ex, the environment vars will be
EX_I
EX_S

and the flags will be

--ex-i --ex-s

When CFG_KEY_PREFIX="", the environment vars will be
I
S

and the flags will be

--i --s

Without CFG_KEY_PREFIX override

APP_I
APP_S

and the flags will be
--app-i --app-s
`

// helpText stores custom help text that will be displayed in usage output.
// This variable is set by the HelpText() function and inserted into the
// standard help output between the usage line and flag definitions.
var helpText string

// HelpText sets custom text to be displayed in help output.
// This function allows applications to provide context, examples, or additional
// information that appears before the flag definitions in help text.
//
// Parameters:
//
//	text - Custom help text to display (supports multiple lines)
//
// Usage:
//
//	cfg.HelpText("MyApp v1.0\nAdvanced configuration management tool\n\nExamples:\n  myapp --port 8080 --debug")
//	cfg.FlatNest(&config)  // Help will include custom text
//
// The custom text appears in the following help format:
//
//	Usage of myapp:
//
//	{custom help text}
//
//	{flag definitions}
//
// Multiple calls to HelpText will replace the previous custom text.
func HelpText(text string) {
	helpText = text
}

// UsageFunc generates formatted help output for the application.
// This function is automatically assigned to pflag's Usage variable during
// package initialization, ensuring consistent help formatting.
//
// Output Format:
//
//	Usage of {program}:
//
//	{custom help text if set via HelpText()}
//
//	{flag definitions from pflag.PrintDefaults()}
//
// Program Name Detection:
//
//	The function automatically extracts the program name from os.Args[0],
//	handling full paths by taking only the final component (basename).
//
// Custom Text Integration:
//
//	If custom help text was set via HelpText(), it appears between the
//	usage line and flag definitions, providing application-specific context.
//
// Flag Documentation:
//
//	The function delegates to pflag's PrintDefaults() for consistent flag
//	formatting that includes names, types, defaults, and usage descriptions.
//
// This function writes to os.Stderr following standard Unix conventions
// for help and error output.
var UsageFunc = func() {
	var parts = strings.Split(os.Args[0], "/")
	var l = len(parts)
	var Program = parts[l-1]
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", Program)
	if len(helpText) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", helpText)
	}

	globalFlagSet.PrintDefaults()
	// fmt.Fprintf(os.Stderr, PrefixOverrideText)
}

// init performs package initialization to configure the help system.
// This function executes during package initialization to set up
// the integration between go-cfg and pflag's usage system.
func init() {
	Usage = UsageFunc
}
