// +build !appengine

package cfgflag

import "syscall"

var LookupEnv = syscall.Getenv
