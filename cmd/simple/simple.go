// Run with
// go run simple.go

package main

import (
	"encoding/json"
	"fmt"

	"github.com/davidwalter0/go-cfg"
)

type App struct {
	I      int `default:"-1"`
	Nested struct {
		Y float64
	}
}

func main() {
	var myapp App

	if err := cfg.Init(&myapp); err != nil { // parse tags, environment, flags
		fmt.Errorf("%v", err)
	}
	fmt.Printf("%v %T\n", myapp, myapp)
	jsonText, _ := json.MarshalIndent(&myapp, "", "  ")
	fmt.Printf("\n%v\n", string(jsonText))
}
