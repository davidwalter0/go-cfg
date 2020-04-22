package cfg

import (
	"encoding/json"
	"fmt"
	"testing"
)

type S1 struct {
	A string
	X int
}

type S2 struct {
	B string
	Y int
}

func TestSave(t *testing.T) {
	var err error
	var data []byte

	s1 := &S1{A: "A", X: 1}
	s2 := &S2{B: "A", Y: 1}
	fmt.Println(Nest(s1, s2))
	for k, v := range Store {
		fmt.Printf("%-15.15s %p\n", k, v)
	}
	data, err = json.MarshalIndent(Store, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("before store load")
	fmt.Println(string(data))
	Store.Save("config.yaml")
	for k, v := range Store {
		fmt.Printf("%-15.15s %p\n", k, v)
	}
	data, err = json.MarshalIndent(Store, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("before load")
	fmt.Println(string(data))
	Store.Load("config.yaml")
	for k, v := range Store {
		fmt.Printf("%-15.15s %p\n", k, v)
	}
	data, err = json.MarshalIndent(Store, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("after load")
	fmt.Println(string(data))
}
