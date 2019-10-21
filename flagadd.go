package cfg

import (
	"fmt"
	"strings"

	"github.com/davidwalter0/go-flag"
)

var allFlagNames = make(map[string]bool)
var announceDuplicates bool

// AddFlag if not created and warn if it is a duplicate
func (field *Field) AddFlag() {
	if _, ok := allFlagNames[field.FlagName]; !ok {
		var usage string
		if len(field.Doc) > 0 {
			usage = "usage: " + field.Doc
		}
		flag.MakeVar(field.FieldPtr, field.FlagName, field.Default, usage+fmt.Sprintf(" Env %-32s : (%v)", field.KeyName, field.StructField.Type), field.Value)
	} else {
		if !announceDuplicates {
			fmt.Printf("Duplicate flag(s)/env vars found\n")
			fmt.Println(strings.ToUpper(fmt.Sprintf("%-20s %-20s", "flag", "env vars")))
			fmt.Println("-----------------------------------------")
			announceDuplicates = true
		}
		fmt.Printf("%-20s %-20s\n", field.FlagName, field.KeyName)
	}
}
