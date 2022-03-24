//go:build !appengine
// +build !appengine

package eflag

import "syscall"

var LookupEnv = syscall.Getenv
