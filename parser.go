package cfg

import (
	"encoding/json"
	"log"
	"reflect"
	//	"strconv"
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

// Args unifies api for recursion
type Args struct {
	Depth    int
	Prefix   string
	Prefixed bool
	Name     string
}

// Enter recursively processes object configurations
func Enter(args *Args, ptr interface{}) error {
	var err error
	elem := reflect.ValueOf(ptr).Elem()
	etype := elem.Type()
	name := etype.Name()
	kind := reflect.TypeOf(ptr).Kind()
	if kind != reflect.Ptr {
		log.Printf("%s %+v\n", name, ErrInvalidArgPointerRequired)
		return ErrInvalidArgPointerRequired
	}

	Store[name] = ptr

	// prefix, found := LookupEnv(cfgEnvKeyPrefix)
	// log.Println("name", name)
	// log.Println("cfgEnvKeyPrefix", cfgEnvKeyPrefix)
	// log.Println("prefix", prefix)

	// if !found && debug {
	// 	log.Println("The configuration parser will run without a prefix override")
	// }

	if args.Prefixed {
		name = args.Prefix
	}

	if !decorate {
		name = ""
	}

	// log.Println("name", name)
	// log.Println("cfgEnvKeyPrefix", cfgEnvKeyPrefix)
	// log.Println("prefix", prefix)
	// log.Println("decorate", decorate)
	err = ParseStruct(args, ptr, name, emptyStructField)
	if err != nil {
		log.Printf("error parsing %s %+v\n", etype.Name(), err)
	}

	return err
}

// ParseStruct recursively processes object configurations
func ParseStruct(args *Args, ptr interface{}, prefix string, structField reflect.StructField) error {
	depth := args.Depth
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
			var name string
			ptr := elem.Field(i).Addr().Interface()
			if args.Prefixed && len(args.Prefix) != 0 {
				name = Capitalize(args.Prefix) + "-" + Capitalize(indirect.Type().Field(i).Name)
			} else {
				name = Capitalize(indirect.Type().Field(i).Name)
			}
			args.Depth++
			err = ParseStruct(args, ptr, name, etype.Field(i))
			if err != nil {
				panic(err)
				break
			}
			args.Depth--
		}
	default:
		elem := reflect.ValueOf(ptr).Elem()
		etype := elem.Type()
		var field = &Field{
			StructField: structField,
			FieldPtr:    ptr,
			Depth:       depth,
			Name:        structField.Name,
			Prefix:      prefix,
			KeyName:     prefix,
			FlagName:    prefix,
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

/*
// EnterName recursively processes object configurations with prefix
func EnterName(depth int, prefix string, ptrs ...interface{}) error {
	var err error
	for _, ptr := range ptrs {
		kind := reflect.TypeOf(ptr).Kind()
		elem := reflect.ValueOf(ptr).Elem()
		etype := elem.Type()
		name := etype.Name()
		if kind != reflect.Ptr {
			log.Printf("%s %+v\n", name, ErrInvalidArgPointerRequired)
			continue
		}
		Store[name] = ptr

		prefix, found := LookupEnv(cfgEnvKeyPrefix)
		log.Println("name", name)
		log.Println("cfgEnvKeyPrefix", cfgEnvKeyPrefix)
		log.Println("prefix", prefix)

		// if !found && debug {
		// 	log.Println("The configuration parser will run without a prefix override")
		// }
		if len(prefix) > 0 {
			name = prefix + name
		}

		if found {
			name = prefix + name
		}

		if !decorate {
			name = ""
		}

		log.Println("name", name)
		log.Println("cfgEnvKeyPrefix", cfgEnvKeyPrefix)
		log.Println("prefix", prefix)
		log.Println("decorate", decorate)

		err = ParseStruct(0, ptr, name, emptyStructField)
		if err != nil {
			log.Printf("error parsing %s %+v\n", etype.Name(), err)
		}
	}
	return err
}
*/
