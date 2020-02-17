package cfg

import (
	"log"
	"os"
	"strconv"
)

// decorate var setting, disable prefixes when undecorated
var decorate bool = true

// Initialize reading setting from the environment
func Initialize() {
	if text, ok := LookupEnv(cfgDecorate); ok {
		if v, err := strconv.ParseBool(text); err != nil {
			log.Println(err)
		} else {
			decorate = v
		}
	}
}

// Undecorate structs with prefix
func Undecorate() bool {
	decorate = false
	os.Setenv(cfgDecorate, "false")
	log.Printf("decorate %v\n", decorate)
	return decorate
}

// Decorate structs with prefix
func Decorate() bool {
	decorate = true
	log.Printf("decorate %v\n", decorate)
	os.Setenv(cfgDecorate, "true")
	return decorate
}

func Unprefix() {
	os.Unsetenv(cfgEnvKeyPrefix)
}

// // Decorate structs with prefix
// func Decorate() bool {
// 	var err error
// 	log.Printf("decorate %v\n", decorate)
// 	text, ok := LookupEnv(cfgDecorate)
// 	if ok {
// 		decorate, err = strconv.ParseBool(text)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// 	log.Printf("decorate %v\n", decorate)
// 	return decorate
// }

// Eval one or more configuration structures
func Eval(ptrs ...interface{}) error {
	args := NewArgs(false, "", "")
	return Run(args, ptrs...)
}

// Init flags
func Init(ptrs ...interface{}) error {
	args := NewArgs(false, "", "")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// EvalName one or more configuration structures overriding the name
func EvalName(name string, ptrs ...interface{}) error {
	prefix, ok := LookupEnv(cfgEnvKeyPrefix)
	if ok && len(prefix) > 0 && len(name) > 0 {
		prefix = prefix + "_" + name
	}
	if ok && len(name) > 0 {
		prefix = name
	}
	if len(prefix) > 0 {
		os.Setenv(cfgEnvKeyPrefix, prefix+"_")
	}
	args := NewArgs(len(prefix) > 0, prefix, "")
	return Run(args, ptrs...)
}

func Run(args *Args, ptrs ...interface{}) error {
	var err error
	for _, ptr := range ptrs {
		err = Enter(args, ptr)
		if err != nil {
			return err
		}
	}
	return err
}

// SimpleFlags don't prefix with base struct prefix, undecorated.
// type T struct {
//    I int
// }
// var ptr = &T{}
// Default
// Eval(ptr) is decorated/prefixed with the struct name
// type T struct {
// ...
//    I int
// }
// type T struct {
// Eval
// EvalName("D", ptr) is decorated/prefixed with flag "--d-i"
// EvalName("D", ptr) is decorated/prefixed with env var "D_I"
//
// export CFG_KEY_PREFIX=KP
// Eval(p) is decorated/prefixed with flag "--kp-d-i"
// EvalName("D", p) is decorated/prefixed with env var "KP_D_I"
//
// export CFG_KEY_PREFIX=KP
// export CFG_DECORATE=false
// ignore prefix and struct name prefixing
// Eval(p) is decorated/prefixed with flag "--kp-d-i"
// EvalName("D", p) is decorated/prefixed with env var "KP_D_I"

func NewArgs(prefixed bool, prefix, name string) *Args {
	return &Args{Depth: 0, Prefixed: prefixed, Prefix: prefix, Name: name}
}

// SimpleFlags create env vars prefices
func Simple(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	args := NewArgs(false, "", "")
	return Run(args, ptrs...)
}

// SimpleFlags create env vars and flags without prefices
func SimpleFlags(ptrs ...interface{}) error {
	err := Simple(ptrs...)
	Freeze()
	return err
}

// Unwrap alias of Simple create env vars without prefices
func Unwrap(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	Initialize()
	args := NewArgs(false, "", "")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// Bare alias of Simple create env vars without prefices
func Bare(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	Initialize()
	args := NewArgs(false, "", "")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// Wrap (optionally) with prefix and struct names create env vars without prefices
func Wrap(name string, ptrs ...interface{}) error {
	Initialize()
	prefix, ok := LookupEnv(cfgEnvKeyPrefix)
	if ok && len(prefix) > 0 {
		prefix = prefix + "_"
	}

	if len(prefix) > 0 && len(name) > 0 {
		prefix = prefix + name
	}
	args := NewArgs(len(prefix) > 0, prefix, "")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// Add alias of Eval
func Add(ptrs ...interface{}) error {
	Initialize()
	args := NewArgs(false, "", "")
	return Run(args, ptrs...)
}

// Flags alias of Eval
func Flags(ptrs ...interface{}) error {
	Initialize()
	args := NewArgs(false, "", "")
	err := Run(args, ptrs...)
	Freeze()
	return err
}
