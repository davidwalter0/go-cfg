package cfg_test

import (
	"fmt"
	"os"

	eflag "github.com/davidwalter0/go-flag"
)

func set(name, value string) {
	os.Setenv(name, value)
}

func clr(name string) {
	os.Unsetenv(name)
}

func see(name string) {
	fmt.Println(name, os.Getenv(name))
}

type CfgValueTester interface {
	Value() string
	Set()
	Clr()
}

type CfgVT CfgValueTester

// A ...
type A struct {
	A string `json:"a" default:"value-a"`
}

func (a *A) Value() string {
	return a.A
}

// B ...
type B struct {
	B string `json:"b" default:"value-b"`
}

func (b *B) Value() string {
	return b.B
}

func resetFlags() {
	eflag.CommandLine = eflag.NewFlagSet(os.Args[0], eflag.ExitOnError)
}

func isTestFlag(s string) (rc bool) {
	var flg = "-test"
	if len(s) > len(flg) && s[:(len(flg))] == flg {
		rc = true
	}
	return
}

var origArgs = os.Args
var filteredArgs = filterTestFlags()

func flagReset() func() {
	resetFlags()
	os.Args = filteredArgs
	return func() {
		os.Args = origArgs
	}
}

func filterTestFlags() (args []string) {
	for _, arg := range os.Args {
		if !isTestFlag(arg) {
			args = append(args, arg)
		}
	}
	return
}

// AB ...
type AB struct {
	A
	B
}

// ABC ...
type ABC struct {
	// A string
	// B int
	A
	B
	C bool
}

// ABC ...
type NestStruct struct {
	Nested ABC
	A
	B
	C bool
	X string
	Y string
	Z string
}

// Match string for test text
func Match(b bool) string {
	if b {
		return "eq"
	}
	return "ne"
}
