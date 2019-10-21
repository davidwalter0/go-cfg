package cfg

// Eval a configuration structure
func Eval(ptr interface{}) error {
	parser, err := NewParser(ptr)
	if err == nil {
		parser.Eval(0)
	}
	return err
}

// Eval a configuration structure
func NamedEval(name string, ptr interface{}) error {
	parser, err := NewParserPrefixed(name, ptr)
	if err == nil {
		parser.Eval(0)
	}
	return err
}
