//go:build !appengine
// +build !appengine

package cfg // import "github.com/davidwalter0/go-cfg"

import "os"

// LookupEnv provides a compatible interface for environment variable lookup.
// This implementation uses os.LookupEnv directly to ensure consistent behavior
// across all build environments.
var LookupEnv = os.LookupEnv
