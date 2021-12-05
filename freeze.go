package cfg

import (
	"github.com/davidwalter0/go-flag"
)

var frozen bool

// Freeze flags
func Freeze() {
	if !frozen {
		flag.Parse()
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
	flag.CommandLine = flag.NewFlagSet(name, ErrorHandlerModel)
}

// ErrorHandlerModel enables reconfiguring flag.ErrorHandling for the
// flag handlers
var ErrorHandlerModel = flag.ContinueOnError

func Reset(name string) {
	Thaw()
	Store = NewStor()
	flag.CommandLine = flag.NewFlagSet(name, flag.PanicOnError)
}

// Thaw flags
func Thaw() {
	frozen = false
}
