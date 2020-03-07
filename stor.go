package cfg

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type Loader interface {
	Load(filename string)
}

type Storer interface {
	Store(filename string)
}

type Adder interface {
	AddStor(name string, o interface{})
}

// Stor configuration representation, restorable object from, saveable
// to persistence.
type Stor map[string]interface{}

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
func (stor Stor) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, stor)
	return err
}

// Stor object to persistence
func (stor Stor) Stor(filename string) error {
	return stor.Save(filename)
}

// Save object to persistence
func (stor Stor) Save(filename string) error {
	data, err := yaml.Marshal(stor)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	return err
}
