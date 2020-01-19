package main

import (
	"fmt"
	"os"
)

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}

type X struct {
}

type Y struct {
}

func (x *X) Marshal(o interface{}) (serial []byte, err error) {
	fmt.Print(" Marshal ")
	return []byte("Marshal"), nil
}

func (x *X) Unmarshal(serial []byte, o interface{}) (err error) {
	fmt.Print(" Unmarshal ")
	return nil
}

func main() {
	var i interface{} = &X{}
	if l, r := i.(Unmarshaler); true {
		fmt.Printf("%T ", l)
		if r {
			i.(Unmarshaler).Unmarshal([]byte{}, i)
		}
		fmt.Println(l, r)

	}
	if l, r := i.(Marshaler); true {
		fmt.Printf("%T ", l)
		if r {
			i.(Marshaler).Marshal(i)
		}
		fmt.Println(l, r)
	}
	i = &Y{}
	if l, r := i.(Unmarshaler); true {
		fmt.Printf("%T ", l)
		if r {
			i.(Unmarshaler).Unmarshal([]byte{}, i)
		}
		fmt.Println(l, r)
	}
	if l, r := i.(Marshaler); true {
		fmt.Printf("%T ", l)
		if r {
			i.(Marshaler).Marshal(i)
		}
		fmt.Println(l, r)
	}
	os.Exit(0)
}
