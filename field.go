package cfg

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/davidwalter0/go-flag"
)

// FieldPtr for the struct field field
type FieldPtr interface{}

// Field holds the parsed struct tag information
type Field struct {
	reflect.StructField
	FieldPtr
	StructName string // The name of the current owning structure
	// AppName    string // The name of the current App
	Prefix   string // env variable application prefix: PREFIX_....name...
	Name     string // from var name if tag name is present replace tag name with tag
	KeyName  string // ENV variable name prefix Prf + "_" + name CamelCase -> PRF_CAMEL_CASE
	Default  string // default from tag, empty string for default
	EnvText  string // environment text, empty string for default
	Short    string // short flag name
	Doc      string // description
	FlagName string // Hyphenated flag name CamelCase -> camel-case
	Value    string // if env use, else if default tag use, else use type's default
	Omit     bool   // obey json:"...,omitempty" or json:"...,omit" or json:"-"
	Required bool   // set to force field to have a value
	Depth    int    // struct nesting depth
	Ignore   bool   // don't store or load the corresponding Attribute
	Error    error
}

// Get a tag from the struct tags
func (field *Field) Get(name string) string {
	text := field.StructField.Tag.Get(name)
	if len(text) > 0 {
		return text
	}
	return ""
}

func NewField(i, depth int, ptr interface{}, attr reflect.StructField, prefix, structName string) *Field {
	switch reflect.ValueOf(ptr).Elem().Kind() {
	case reflect.Struct:
		elem := reflect.ValueOf(ptr).Elem()
		eType := elem.Type()
		name := eType.Name()
		Parse(depth+1, prefix+"_"+name, ptr)
		return nil
	}

	var field = &Field{
		StructField: attr,
		FieldPtr:    ptr,
		StructName:  structName,
		Prefix:      prefix,
		Depth:       depth,
	}
	field.Parse(prefix)
	if field.Error != nil {
		log.Println(field.Error)
		return nil
	}
	byte, err := json.Marshal(field)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(byte))
	return field
}

// // 	m := reflect.ValueOf(ptr).Elem()
// switch m.Kind() {
// case reflect.Struct:
// 	e := reflect.ValueOf(ptr).Elem()
// 	eType := e.Type()
// 	structName := eType.Name()

// 	prefix := p.Prefix
// 	for i := 0; i < eType.NumField(); i++ {
// 		structField := elementType.Field(i)
// 		ptr := element.Field(i).Addr().Interface()
// 		field := p.Field(i, depth, e, eType.Field(i), structName, prefix)
// 	}
// }

// Parse the struct tags
func (field *Field) Parse(prefix string) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	field.Prefix = prefix
	field.SetDefault()
	field.SetName()
	field.SetIgnore()
	field.SetDoc()
	field.SetShort()
	field.SetOmit()
	field.SetRequired()
	field.SetKeyName()
	field.SetFlagName()
	field.SetValueFromEnv()
	var usage string
	if len(field.Doc) > 0 {
		usage = "usage: " + field.Doc
	}
	flag.MakeVar(field.FieldPtr, field.FlagName, field.Default,
		usage+fmt.Sprintf(" Env %-32s : (%v)",
			field.KeyName, field.StructField.Type), field.Value)
}

// SetOmit read tag omit option and set when enabled, via ,omit
// ,omitempty or '-' the hyphen option
func (field *Field) SetOmit() {
	json := field.Get("json")
	field.Omit = json == "-" || strings.Index(json, ",omitempty") >= 0 || strings.Index(json, ",omit") >= 0
}

// SetDefault read tag default option and save the text
func (field *Field) SetDefault() {
	field.Default = field.Get("default")
	field.Value = field.Default
}

// SetIgnore read tag ignore option and save the text
func (field *Field) SetIgnore() {
	text := field.Get("ignore")
	if v, err := strconv.ParseBool(text); err == nil {
		field.Ignore = v
	}
}

// SetDoc read tag doc option and save the text
func (field *Field) SetDoc() {
	field.Doc = field.Get("doc")
}

// SetShort read tag short option and save the text
func (field *Field) SetShort() {
	field.Short = field.Get("short")
}

// SetPrefix read tag prefix option and save the text
func (field *Field) SetPrefix() {
	field.Prefix = field.Get("prefix")
}

// SetName read tag name option and save the text
func (field *Field) SetName() {
	name := field.Get("json")
	i := strings.Index(name, ",")
	if i > 0 {
		name = name[0:i]
	}
	field.Name = name
}

// SetIgnore read tag required option and save the text
func (field *Field) SetRequired() {
	text := field.Get("required")
	if v, err := strconv.ParseBool(text); err == nil {
		field.Omit = v
	}
}

// SetValueFromEnv uses the value from the environment for this
// structure tag replacing the default tag value
func (field *Field) SetValueFromEnv() {
	if len(field.KeyName) > 0 {
		field.EnvText, _ = LookupEnv(field.KeyName)
	}
	if len(field.EnvText) != 0 {
		field.Value = field.EnvText
	}
}

// SetKeyName read tag keyword option and save the text
func (field *Field) SetKeyName() {
	field.KeyName = strings.Replace(field.KeyName, "-", "_", -1)
	field.UnderScoreCamelCaseWords()
}

// SetKeyName read tag keyword option and save the text
func (field *Field) SetFlagName() {
	if len(field.Prefix) > 0 {
		field.FlagName = field.Prefix + "-" + Capitalize(field.Name)
	}
	field.FlagName = strings.Replace(field.FlagName, "_", "-", -1)
	field.HyphenateCamelCaseWords()
}

var regExpr = regexp.MustCompile("([^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")

// UnderScoreCamelCaseWords split on CamelCase words CAMEL_CASE, if
// field.Prefix is set, prepend to the text for environment variables
// split on camel case regular expression
func (field *Field) UnderScoreCamelCaseWords() {
	words := regExpr.FindAllStringSubmatch(field.Name, -1)
	if len(words) > 0 {
		var nameParts []string
		for _, words := range words {
			nameParts = append(nameParts, words[0])
		}
		field.KeyName = strings.Join(nameParts, "_")
		if len(field.Prefix) > 0 {
			field.KeyName = strings.ToUpper(field.Prefix + "_" + field.KeyName)
		} else {
			field.KeyName = strings.ToUpper(field.KeyName)
		}
		field.KeyName = strings.Replace(field.KeyName, "-", "_", -1)
	}
}

// HyphenateCamelCaseWords for flags CamelCase to camel-case
// hyphenated, split on camel case regular expression
func (field *Field) HyphenateCamelCaseWords() {
	words := regExpr.FindAllStringSubmatch(field.FlagName, -1)
	if len(words) > 0 {
		var name []string
		for _, words := range words {
			name = append(name, strings.ToLower(words[0]))
		}
		field.FlagName = strings.Join(name, "-")
	}
	for n := strings.Index(field.FlagName, "--"); n > 0; n = strings.Index(field.FlagName, "--") {
		field.FlagName = strings.Replace(field.FlagName, "--", "-", -1)
	}
}

// Capitalize text
func Capitalize(text string) string {
	var capitalized = text
	if len(text) > 0 {
		capitalized = strings.ToUpper(text[0:1]) + text[1:]
	}
	return capitalized
}

/*
func (field *MemberType) Parse(prefix string, fld reflect.StructField, ptr interface{}, depth int) (err error) {
	field.TagInit(fld.Tag)

	if len(field.Name) == 0 {
		field.Name = fld.Name
	}

	field.EnvVarPrefix = prefix
	field.FlagName = field.Name
	field.UnderScoreCamelCaseWords()
	field.HyphenateCamelCaseWords()
	field.EnvInit()
	// Env variable names historically didn't allow hyphenation
	field.KeyName = strings.Replace(field.KeyName, "-", "_", -1)

	switch fld.Type.Kind() {
	case reflect.Struct:
		depth = depth + 1
		element := reflect.ValueOf(ptr).Elem()
		elementType := element.Type()
		if len(prefix) > 0 {
			prefix = prefix + "_" + Capitalize(field.Name)
		} else {
			prefix = Capitalize(field.Name)
		}
		for i := 0; i < elementType.NumField(); i++ {
			fld := elementType.Field(i)
			ptr := element.Field(i).Addr().Interface()
			if err = field.Parse(prefix, fld, ptr, depth); err != nil {
				return
			}
		}
	default:
		if !field.Ignore {
			var usage string
			if len(field.Usage) > 0 {
				usage = "usage: " + field.Usage
			}
			if _, ok := allFlagNames[field.FlagName]; !ok {
				v := reflect.ValueOf(ptr).Elem()
				typeOf := v.Type()
				fmt.Printf("Type %T v %v\n", typeOf, v)
				if field.Required && len(field.Value) == 0 {
					return fmt.Errorf("Required field: Name [%s] FlagName [%s] Env var name [%s] has no default and is unset", field.Name, field.FlagName, field.KeyName)
				}
				flag.MakeVar(ptr, field.FlagName, field.Default,
					usage+fmt.Sprintf(" Env %-32s : (%v)",
						field.KeyName, fld.Type), field.Value)
				allFlagNames[field.FlagName] = true
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
	}
	return
}
*/
