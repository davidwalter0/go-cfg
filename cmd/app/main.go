package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-flag"
)

type Key string
type Value float64

// X flag
var X = flag.String("FLAG", "STRING_VALUE", "FLAG USAGE...", false, false)

func main() {

	{
		var myapp App

		if err := cfg.Nest(&myapp); err != nil {
			log.Fatalf("%v\n", err)

		}
		log.Printf("%v %T\n", myapp, myapp)
		jsonText, _ := json.MarshalIndent(&myapp, "", "  ")
		log.Printf("\n%v\n", string(jsonText))
		cfg.Usage()
		/*
			flag.Usage()
			// Error can't call parse again
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("\n\n\n\n")
					fmt.Printf("***Error*** %v\n", err)
					fmt.Printf("\n\n\n\n")
					// sti.Parse()
				}
			}()
			sti.Parse()
		*/
	}

	os.Exit(0)
}
