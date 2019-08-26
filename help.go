package cfg

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var PrefixOverrideText = `
CFG_ENV_VAR_PREFIX_KEY

CFG_ENV_VAR_PREFIX_KEY=ex

The environment variable CFG_ENV_VAR_PREFIX_KEY overrides the prefix
of each environment variable to value. The value will be converted to
upper case like ex - EX

For a struct 

type App struct {
   i int
   s string
}

The environment vars will be
EX_I
EX_S

Without CFG_ENV_VAR_PREFIX_KEY override

APP_I
APP_S

`

var helpText string

// HelpText pre write text for help
func HelpText(text string) {
	helpText = text
}

// Usage add usage text for flags/help processing
var Usage = func() {
	var parts = strings.Split(os.Args[0], "/")
	var l = len(parts)
	var Program = parts[l-1]
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", Program)
	if len(helpText) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", helpText)
	}
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, PrefixOverrideText)
}

func init() {
	flag.Usage = Usage
}
