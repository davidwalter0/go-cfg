package cfg

import "log"

// Eval a configuration structure
func Eval(ptr interface{}) {
	parser, err := NewParser(ptr)
	if err != nil {
		log.Fatal(err)
	}
	parser.Eval(0)
}
