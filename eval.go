/*
go-cfg
package cfg

Modes of configuration

- Unwrap( ptr ... interface{} ) error
  - Flags and Env Vars are unwrapped - no prefix or argument name added
  - Call with one or more structures with uniquely named members
  - One flag + env variable for each entry named
  - When called with duplicate members in one or more structs flags and
  env vars will conflict and error


*/
package cfg

import (
	"log"
	"os"
	"strconv"
)

// decorate var setting, disable prefixes when undecorated
var decorate bool = true

// Args unifies api for recursion
type Args struct {
	Depth    int
	Prefix   string
	Prefixed bool
	UseFlags bool
}

func NewArgs(name string) *Args {
	var prefixed bool
	var prefix string
	var text string
	var ok bool
	if text, ok = LookupEnv(cfgDecorate); ok {
		if v, err := strconv.ParseBool(text); err != nil {
			log.Println(err)
		} else {
			decorate = v
		}
	}
	if decorate {
		prefix, ok = LookupEnv(cfgEnvKeyPrefix)
		if ok && len(prefix) > 0 {
			prefix = prefix + "_"
			prefixed = true
		}
		if len(name) > 0 {
			prefixed = true
			prefix = prefix + name
		}
	}
	return &Args{Depth: 0, Prefixed: prefixed, Prefix: prefix, UseFlags: true}
}

// Undecorate structs with prefix
func Undecorate() bool {
	decorate = false
	os.Setenv(cfgDecorate, "false")
	return decorate
}

// Decorate structs with prefix
func Decorate() bool {
	decorate = true
	os.Setenv(cfgDecorate, "true")
	return decorate
}

func Unprefix() {
	os.Unsetenv(cfgEnvKeyPrefix)
}

// Eval one or more configuration structures
func Eval(ptrs ...interface{}) error {
	args := NewArgs("")
	return Run(args, ptrs...)
}

// Init flags
func Init(ptrs ...interface{}) error {
	args := NewArgs("")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// EvalName one or more configuration structures overriding the name
func EvalName(name string, ptrs ...interface{}) error {
	args := NewArgs(name)
	err := Run(args, ptrs...)
	Freeze()
	return err
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

// SimpleFlags create env vars prefices
func Simple(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	args := NewArgs("")
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
	args := NewArgs("")
	err := Run(args, ptrs...)
	return err
}

// Final freezes calling flag.Parse, no more additions to the
// configuration after Final
func Final() {
	Freeze()
}

// Bare alias of Simple create env vars without prefices
func Bare(ptrs ...interface{}) error {
	Unprefix()
	Undecorate()
	args := NewArgs("")
	err := Run(args, ptrs...)
	return err
}

// Wrap (optionally) with prefix and struct names create env vars without prefices
func Wrap(name string, ptrs ...interface{}) error {
	args := NewArgs(name)
	err := Run(args, ptrs...)
	return err
}

// Add alias of Eval
func Add(ptrs ...interface{}) error {
	args := NewArgs("")
	return Run(args, ptrs...)
}

// Flags alias of Eval
func Flags(ptrs ...interface{}) error {
	args := NewArgs("")
	err := Run(args, ptrs...)
	Freeze()
	return err
}

// Nest objects retaining object hierarchy
func Nest(ptrs ...interface{}) error {
	args := NewArgs("")
	args.Prefixed = false
	err := Run(args, ptrs...)
	return err
}

// NestWrap objects retaining object hierarchy with prefix
func NestWrap(prefix string, ptrs ...interface{}) error {
	args := NewArgs(prefix)
	args.Prefixed = true
	err := Run(args, ptrs...)
	return err
}
