package cfg

// Eval a configuration structure
func Eval(ptr interface{}) error {
	return Enter(0, ptr)
}
