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

// Init flags
func Init(ptr interface{}) error {
	err := Enter(0, ptr)
	Freeze()
	return err
}

// FlagInit flags
func FlagInit() {
	Freeze()
}
