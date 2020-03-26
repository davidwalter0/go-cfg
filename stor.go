package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// Stor configuration representation, restorable object from, saveable
// to persistence.
type Stor map[string]interface{}

type Loader interface {
	Load(filename string)
}

type Storer interface {
	Store(filename string)
}

type Adder interface {
	AddStor(name string, o interface{})
}

// NewStor object to from persistence
func NewStor() Stor {
	return Stor{}
}

// // Read object from io.Reader
// func (stor Stor) Read(r io.Reader) (int, error) {
// 	var data = []byte{}
// 	n, err := r.Read(data)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return n, err
// }

// // Write object from io.Writer
// func (stor Stor) Write(w io.Writer) (int, error) {
// 	data, err := yaml.Marshal(stor)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return w.Write(data)
// }

// AddStor object to from persistence
func (stor Stor) AddStor(name string, o interface{}) {
	stor[name] = o
}

// Load object from persistence
func (stor *Stor) Load(filename string) error {
	var err error
	var data []byte

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	switch strings.ToLower(path.Ext(filename)) {
	case ".json":
		err = json.Unmarshal(data, stor)
		if err != nil {
			return err
		}
	case ".yaml":
		err = yaml.Unmarshal(data, stor)
		if err != nil {
			return err
		}
	default:
		err = yaml.Unmarshal(data, stor)
		if err != nil {
			return err
		}
	}
	return err
}

// Stor object to persistence
func (stor Stor) Stor(filename string) error {
	return stor.Save(filename)
}

// Save object to persistence
func (stor Stor) Save(filename string) error {
	var err error
	var data []byte
	switch strings.ToLower(path.Ext(filename)) {
	case ".json":
		data, err = json.Marshal(stor)
		if err != nil {
			return err
		}
	case ".yaml":
		data, err = yaml.Marshal(stor)
		if err != nil {
			return err
		}
	default:
		data, err = yaml.Marshal(stor)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// String stor interfaceable
func (stor Stor) Bytes() []byte {
	text, err := json.Marshal(stor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stor to bytes: %s", err.Error())
		return nil
	}
	return text
}

// String stor interfaceable
func (stor Stor) String() string {
	return string(stor.Bytes())
}
