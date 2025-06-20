package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davidwalter0/go-cfg/pkg/flags/generators"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: generator_example <type>")
		fmt.Println("Types: loglevel, environment, string-slice")
		os.Exit(1)
	}
	
	typeToGenerate := os.Args[1]
	
	var generated string
	var err error
	
	switch typeToGenerate {
	case "loglevel":
		generated, err = generators.GenerateLogLevelType()
	case "environment":
		generated, err = generators.GenerateEnvironmentType()
	case "string-slice":
		generated, err = generators.GenerateStringSliceType()
	default:
		log.Fatalf("Unknown type: %s", typeToGenerate)
	}
	
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}
	
	fmt.Println(generated)
}

// Example usage:
//   go run generator_example.go loglevel > loglevel_generated.go
//   go run generator_example.go environment > environment_generated.go