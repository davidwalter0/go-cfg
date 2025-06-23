// This file handles command-line flag registration and duplicate detection for go-cfg.
// It provides the integration layer between Field processing and the pflag system,
// ensuring that each configuration field gets properly registered as a command-line flag.
//
// Key Features:
//   - Duplicate flag detection and warning system
//   - Automatic flag type inference from Field metadata
//   - Enhanced help text generation combining field documentation with environment variable info
//   - Conditional flag registration to prevent conflicts
//
// The AddFlag method is called during Field.SetField() processing to register each
// configuration field as a command-line flag. The system tracks all registered flag
// names to detect conflicts and provides detailed warnings when duplicates are found.
//
// Integration Flow:
//  1. Field.SetField() calls field.AddFlag() if UseFlags is true
//  2. AddFlag() checks for duplicate flag names using allFlagNames map
//  3. If not duplicate: calls MakeVar() to register with pflag system
//  4. If duplicate: displays formatted warning table with conflicting names
//
// The duplicate detection system helps users identify configuration conflicts when
// multiple structs contain fields that would generate the same flag names.
package cfg

import (
	"fmt"
	"strings"
)

// allFlagNames tracks all registered flag names to detect duplicates.
// This map is populated by AddFlag() and used to prevent flag registration
// conflicts when multiple configuration fields would generate the same flag name.
//
// The map key is the flag name (e.g., "server-port") and the value is always true.
// This provides O(1) lookup for duplicate detection during flag registration.
var allFlagNames = make(map[string]bool)

// announceDuplicates tracks whether the duplicate warning header has been displayed.
// This ensures the warning table header is only printed once, even when multiple
// duplicate flags are detected during a single configuration processing run.
var announceDuplicates bool

// AddFlag registers a command-line flag for this field if not already registered.
// This method provides the core integration between go-cfg's field processing and
// pflag's flag registration system, with built-in duplicate detection and warning.
//
// The method creates flags with enhanced help text that includes both the field
// documentation and environment variable information, providing users with complete
// configuration guidance in a single help display.
//
// Flag Registration Process:
//  1. Check if flag name already exists in allFlagNames map
//  2. If new: register flag via MakeVar() with enhanced usage text
//  3. If duplicate: display formatted warning table showing conflicts
//
// Usage Text Format:
//
//	The generated usage text combines multiple information sources:
//	- "usage: {field.Doc}" - Field documentation from struct tags
//	- "Env {ENV_VAR_NAME} : ({field_name}) ({type})" - Environment variable info
//
// Duplicate Handling:
//
//	When duplicates are detected, the method prints a formatted table showing:
//	- Header: "pflag" and "env vars" columns
//	- Separator line for visual clarity
//	- Each duplicate: flag name and corresponding environment variable
//
// The duplicate warning helps users identify configuration conflicts and adjust
// their struct field names or json tags to avoid flag name collisions.
//
// Example Usage Text Output:
//
//	--server-port int
//	      usage: HTTP server port Env SERVER_PORT : (port) (int) (default 8080)
func (field *Field) AddFlag() {
	if _, ok := allFlagNames[field.FlagName]; !ok {
		var usage string
		if len(field.Doc) > 0 {
			usage = "usage: " + field.Doc
		}
		isset := len(field.Default) > 0 || len(field.EnvText) > 0
		MakeVar(field.FieldPtr, field.FlagName, field.Default, usage+fmt.Sprintf(" Env %-32s : (%s) (%v)", field.KeyName, field.Name, field.Type), field.Value, field.Required, isset)
	} else {
		if !announceDuplicates {
			fmt.Printf("Duplicate flag(s)/env vars found\n")
			fmt.Println(strings.ToUpper(fmt.Sprintf("%-20s %-20s", "pflag", "env vars")))
			fmt.Println("-----------------------------------------")
			announceDuplicates = true
		}
		fmt.Printf("%-20s %-20s\n", field.FlagName, field.KeyName)
	}
}
