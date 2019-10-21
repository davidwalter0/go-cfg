package cfg

import (
	"flag"
	"log"
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
	parser, err := NewParser(ptr)
	if err == nil {
		parser.Eval(0)
	}
	return err
}

// FlagInit flags
func FlagInit() {
	Freeze()
}

func NamedInit(name string, ptr interface{}) error {
	log.Printf("NamedInit prefix: %s / ptr: %+v\n", name, ptr)
	parser, err := NewParserPrefixed(name, ptr)
	log.Println("NamedInit", err)
	if err == nil {
		parser.Eval(0)
	}
	Freeze()
	return err
}
