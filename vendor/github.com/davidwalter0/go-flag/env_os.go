//go:build appengine
// +build appengine

package eflag

import "os"

var LookupEnv = os.LookupEnv
