package cfg

import (
	eflag "github.com/davidwalter0/go-flag"
)

var frozen bool

// Freeze flags
func Freeze() {
	if !frozen {
		eflag.Parse()
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
	eflag.CommandLine = eflag.NewFlagSet(name, ErrorHandlerModel)
}

// ErrorHandlerModel enables reconfiguring eflag.ErrorHandling for the
// flag handlers
var ErrorHandlerModel = eflag.ContinueOnError

// Thaw flags
func Thaw() {
	frozen = false
}
