package cfg

import (
	"log"
	"reflect"
)

const (
	cfgEnvPrefixKey = "CFG_PREFIX_KEY"
)

// Args to parse
type Args []string

// Key struct name
type Key string

// KV map of key(struct element) + tags
type KV map[Key]*Field

// parser recursively parses a struct
type Parser struct {
	ptr    interface{}
	prefix string
	kv     KV
}

// KV returns a map of object definitions
func (parser *Parser) KV() KV {
	return parser.kv
}

// Parse recursively processes object configurations
func (parser *Parser) Parse() {
	Parse(0, parser.prefix, parser.ptr)
}

// NewParser return an initialized parser
func NewParser(ptr interface{}) (*Parser, error) {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return nil, ErrInvalidArgPointerRequired
	}

	prefix, found := LookupEnv(cfgEnvPrefixKey)
	if !found {
		log.Println("parser will run without a prefix override")
	}
	return &Parser{
		prefix: prefix,
		ptr:    ptr,
		kv:     make(KV),
	}, nil
}

// Eval recursively processes object configurations
func (parser *Parser) Eval(depth int) {
	ptr := parser.ptr
	prefix := parser.prefix
	switch reflect.ValueOf(ptr).Elem().Kind() {
	case reflect.Struct:
		elem := reflect.ValueOf(ptr).Elem()
		eType := elem.Type()
		name := eType.Name()
		for i := 0; i < eType.NumField(); i++ {
			attr := eType.Field(i)
			ptr := elem.Field(i).Addr().Interface()
			NewField(i, depth, ptr, attr, prefix, name)
		}
	}
}

// Parse recursively processes object configurations
func Parse(depth int, prefix string, ptr interface{}) {
	switch reflect.ValueOf(ptr).Elem().Kind() {
	case reflect.Struct:
		elem := reflect.ValueOf(ptr).Elem()
		eType := elem.Type()
		name := eType.Name()
		for i := 0; i < eType.NumField(); i++ {
			attr := eType.Field(i)
			ptr := elem.Field(i).Addr().Interface()
			NewField(i, depth, ptr, attr, prefix, name)
		}
	}
}

// NewParserPrefixed return an initialized parser using args
func NewParserPrefixed(prefix string, ptr interface{}) (*Parser, error) {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return nil, ErrInvalidArgPointerRequired
	}

	return &Parser{
		prefix: prefix,
		ptr:    ptr,
		kv:     make(KV),
	}, nil
}
