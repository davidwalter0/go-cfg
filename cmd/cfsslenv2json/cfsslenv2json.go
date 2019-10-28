/*

Example command line use of package cfg to configure
using environment variables or flags

*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/davidwalter0/go-cfg"
)

/*
If you need to override the environment variable prefix at runtime it
can be done setting the APP_OVERRIDE_PREFIX environment variable


export APP_OVERRIDE_PREFIX=APP
export APP_CN="example.com"
export APP_NAME_C="US"
export APP_NAME_L="San Francisco"
export APP_NAME_O="Widgets Inc"
export APP_NAME_OU="WWW"
export APP_NAME_ST="California"
export APP_NAMES='{"C":"US","L":"San Francisco","O":"Widgets Inc","OU":"WWW","ST":"California"}'
export APP_HOSTS="1.example.com,2.example.com"
go run csffse2json.go ;

*/

type Name struct {
	C  string `default:"US"`
	L  string `default:"San Francisco"`
	O  string `default:"Internet Widgets, Inc."`
	OU string `default:"WWW"`
	ST string `default:"California"`
}

type Key struct {
	Algo string `default:"rsa"  json:"algo"`
	Size int    `default:"4096" json:"size"`
}

type CFSSL struct {
	CN    string
	Names []Name   `json:"names" help:"array of cert names struct"`
	Key   Key      `json:"key"`
	Hosts []string `json:"hosts" help:"comma separated list of hosts"`
}

type App struct {
	CN    string `default:"example.com"`
	Name  Name
	Key   Key
	Hosts []string `name:"host-comma-list" help:"comma separated list of hosts"`
	Outer struct {
		Inner struct {
			S string `json:"string"`
			I int
		}
	}
}

func main() {
	var app App

	cfg.Init(&app)

	var appCfssl CFSSL = CFSSL{
		CN:    app.CN,
		Names: []Name{app.Name},
		Key:   app.Key,
		Hosts: app.Hosts,
	}
	var jsonText []byte
	jsonText, _ = json.MarshalIndent(&appCfssl, "", "  ")
	fmt.Printf("\n%v\n", string(jsonText))
	// ioutil.WriteFile("configured.json", jsonText, 0777)
}
