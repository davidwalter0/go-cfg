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
		Y float64 `required:"true"`
	}
}

func main() {
	var myapp App

	if err := cfg.Flags(&myapp); err != nil { // parse tags, environment, flags
		fmt.Errorf("%v", err)
	}
	fmt.Printf("export Y=v to get ok %v %v\n", cfg.Ok("y"), myapp.Nested.Y)
	fmt.Printf("IsSet %v required %v ok %v app %+v\n", cfg.IsSet("i"), cfg.Required("i"), cfg.Ok("i"), myapp)
	fmt.Printf("IsSet %v required %v ok %v app %+v\n", cfg.IsSet("y"), cfg.Required("y"), cfg.Ok("y"), myapp)
	fmt.Printf("type %T\n", myapp)
	fmt.Printf("value %+v\n", myapp)
	if !cfg.Ok("y") {
		panic(fmt.Sprintf("! ok => export Y=v to get ok %v %v\n", cfg.Ok("y"), myapp.Nested.Y))
	}
	jsonText, _ := json.MarshalIndent(&myapp, "", "  ")
	fmt.Printf("\n%v\n", string(jsonText))
}
