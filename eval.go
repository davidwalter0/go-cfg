package cfg

import (
	"log"
	"os"
	"strconv"
)

// decorate var setting, disable prefixes when undecorated
var decorate bool = true

// Eval a configuration structure
func Eval(ptr interface{}) error {
	return Enter(0, ptr)
}

func init() {
	var err error
	text, ok := LookupEnv(cfgDecorate)
	if ok {
		decorate, err = strconv.ParseBool(text)
		if err != nil {
			log.Println(err)
		}
	}
}

// EvalName a configuration structure
func EvalName(name string, ptr interface{}) error {
	os.Setenv(cfgEnvKeyPrefix, name)
	return Enter(0, ptr)
}

// Decorate structs with prefix
func Decorate(enable bool) {
	decorate = enable
}
