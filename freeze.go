package cfg

import (
	"os"
	
	"github.com/spf13/pflag"
)

var frozen bool

// Freeze flags
func Freeze() {
	if !frozen {
		globalFlagSet.Parse(os.Args[1:])
		frozen = true
	}
}

// FlagInit flags
func FlagInit() {
	Freeze()
}

// Reset from frozen and enable re-evaluation with ErrorHandlerModel
func Reset(name string) {
	Thaw()
	Store = NewStor()
	globalFlagSet = pflag.NewFlagSet(name, ErrorHandlerModel)
	CommandLine = globalFlagSet
}

// ErrorHandlerModel enables reconfiguring flag.ErrorHandling for the
// flag handlers
var ErrorHandlerModel = pflag.ContinueOnError

// Thaw flags
func Thaw() {
	frozen = false
}
