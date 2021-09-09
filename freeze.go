package cfg

import (
	"os"

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

func Reset() {
	Thaw()
	Store = NewStor()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
}

// Thaw flags
func Thaw() {
	frozen = false
}
