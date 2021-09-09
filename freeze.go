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

func Reset() {
	Thaw()
	Store = NewStor()
	flag.CommandLine.Reset()
}

// Thaw flags
func Thaw() {
	frozen = false
}
