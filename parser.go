package cfg

import (
	"encoding/json"
	"log"
	"reflect"
)

const (
	cfgEnvKeyPrefix = "CFG_KEY_PREFIX"
	cfgDecorate     = "CFG_DECORATE"
)

/*
// Args to parse
type Args []string

// Key struct name
type Key string

// KV map of key(struct element) + tags
type KV map[Key]*Field

// Parser recursively parses a struct
type Parser struct {
	ptr interface{}
	kv  KV
	mgr Mgr
}

// KV returns a map of object definitions
func (parser *Parser) KV() KV {
	return parser.kv
}

// NewParser return an initialized parser using args
func NewParser(ptr interface{}) (*Parser, error) {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return nil, ErrInvalidArgPointerRequired
	}
	return &Parser{
		ptr: ptr,
		kv:  make(KV),
		mgr: NewCache(),
	}, nil
}
*/
var emptyStructField = reflect.StructField{}

// Store persistable representation
var Store = NewStor()

// Enter recursively processes object configurations
func Enter(depth int, ptr interface{}) error {
	kind := reflect.TypeOf(ptr).Kind()
	if kind != reflect.Ptr {
		panic(ErrInvalidArgPointerRequired)
	}
	elem := reflect.ValueOf(ptr).Elem()
	etype := elem.Type()
	name := etype.Name()

	Store[name] = ptr

	prefix, found := LookupEnv(cfgEnvKeyPrefix)
	if !found && debug {
		log.Println("The configuration parser will run without a prefix override")
	}
	if found {
		name = prefix
	}
	if !decorate {
		name = ""
	}
	return ParseStruct(0, ptr, name, emptyStructField)
}

// ParseStruct recursively processes object configurations
func ParseStruct(depth int, ptr interface{}, parseName string, structField reflect.StructField) error {
	var err error
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		panic(ErrInvalidArgPointerRequired)
	}
	kind := reflect.ValueOf(ptr).Elem().Kind()
	switch kind {
	case reflect.Struct:
		elem := reflect.ValueOf(ptr).Elem()
		etype := elem.Type()
		indirect := reflect.Indirect(reflect.ValueOf(ptr))
		for i := 0; i < etype.NumField(); i++ {
			ptr := elem.Field(i).Addr().Interface()
			if len(parseName) != 0 {
				parseName = Capitalize(parseName) + "-" + Capitalize(indirect.Type().Field(i).Name)
			} else {
				parseName = Capitalize(indirect.Type().Field(i).Name)
			}
			err = ParseStruct(depth+1, ptr, parseName, etype.Field(i))
			if err != nil {
				panic(err)
				break
			}
		}
	default:
		elem := reflect.ValueOf(ptr).Elem()
		etype := elem.Type()
		var field = &Field{
			StructField: structField,
			FieldPtr:    ptr,
			Depth:       depth,
			Name:        structField.Name,
			Prefix:      parseName,
			KeyName:     parseName,
			FlagName:    parseName,
			Type:        etype.Name(),
		}

		field.SetField()
		if field.Error != nil {
			log.Println(field.Error)
			panic(field.Error)
			return nil
		}

		if debug {
			byte, err := json.Marshal(field)
			if err != nil {
				log.Println(err)
				panic(err)
			}
			log.Println(">> ", string(byte))
		}

		return err
	}
	return err
}
