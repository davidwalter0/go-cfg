package legacy

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// Legacy types for backward compatibility
// These maintain the same interface as the original go-cfg library

// Field represents a configuration field with metadata (legacy compatibility)
type Field struct {
	StructField reflect.StructField
	FieldPtr    interface{}
	UseFlags    bool
	Depth       int
	Name        string
	Prefix      string
	KeyName     string
	FlagName    string
	Type        string
	Error       error
}

// SetField applies field configuration (legacy compatibility)
func (f *Field) SetField() {
	// In v3, field processing is handled internally
	// This is a no-op for compatibility
}

// Stor represents configuration storage (legacy compatibility)
type Stor map[string]interface{}

// Bytes returns JSON representation (legacy compatibility)
func (s Stor) Bytes() []byte {
	// Simple JSON marshaling for compatibility
	data, _ := json.Marshal(s)
	return data
}

// Store is the global configuration store (legacy compatibility)
var Store = make(map[string]interface{})

// NewStor creates a new storage instance (legacy compatibility)
func NewStor() map[string]interface{} {
	return make(map[string]interface{})
}

// Arg represents processing arguments (legacy compatibility)
type Arg struct {
	Depth      int
	Prefix     string
	Prefixed   bool
	Unprefixed bool
	Reprefix   bool
	UseFlags   bool
}

// NewArg creates processing arguments (legacy compatibility)
func NewArg(name string) *Arg {
	return &Arg{
		Depth:    0,
		Prefixed: name != "",
		Prefix:   name,
		UseFlags: true,
	}
}

// Legacy interfaces for extensibility

// DeepCopier interface for custom copying (legacy compatibility)
type DeepCopier interface {
	DeepCopy() interface{}
}

// Loader interface for custom loading (legacy compatibility)
type Loader interface {
	Load() error
}

// Storer interface for custom storage (legacy compatibility)
type Storer interface {
	Store() error
}

// Adder interface for custom field addition (legacy compatibility)
type Adder interface {
	Add(field *Field) error
}

// Legacy error types
var (
	ErrInvalidArgStructPointerRequired = errors.New("argument must be a struct pointer")
	ErrInvalidArgPointerRequired       = errors.New("argument must be a pointer")
)

// Legacy utility functions

// RemovePkg from a type name (legacy compatibility)
func RemovePkg(pkg string) string {
	regx := regexp.MustCompile(`\*?[^.]*\.`)
	text := regx.ReplaceAll([]byte(pkg), []byte{})
	return string(text)
}

// IsStructPtr checks if interface is a struct pointer (legacy compatibility)
func IsStructPtr(config interface{}) bool {
	return reflect.TypeOf(config).Kind() == reflect.Ptr &&
		reflect.ValueOf(config).Elem().Kind() == reflect.Struct
}

// IsMap checks if interface is a map (legacy compatibility)
func IsMap(config interface{}) bool {
	return reflect.TypeOf(config).Kind() == reflect.Map
}

// Type returns type information (legacy compatibility)
func Type(config interface{}) string {
	if reflect.TypeOf(config).Kind() == reflect.Ptr {
		return fmt.Sprintf("(%s) %s", reflect.TypeOf(config).Kind(), reflect.ValueOf(config).Elem().Kind())
	}
	return fmt.Sprintf("%s", reflect.ValueOf(config).Elem().Kind())
}

// CheckArgs validates configuration arguments (legacy compatibility)
func CheckArgs(ptrs ...interface{}) error {
	for _, ptr := range ptrs {
		if !IsStructPtr(ptr) {
			return ErrInvalidArgStructPointerRequired
		}
	}
	return nil
}