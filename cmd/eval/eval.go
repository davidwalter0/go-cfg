package main

import (
	"fmt"
	"os"

	// "github.com/davidwalter0/go-cfg"
	cfg "github.com/davidwalter0/go-cfg"
	yaml "gopkg.in/yaml.v3"
)

// A ...
type A struct {
	A string `json:"a" doc:"member A"`
}

// B ...
type B struct {
	B string `doc:"member B"`
}

// AB ...
type AB struct {
	A `doc:"struct A"`
	B `doc:"struct B"`
}

// ABC ...
type ABC struct {
	// A string
	// B int
	A `doc:"struct A"`
	B `doc:"struct B"`
	C bool
}

// Nest ...
type Nest struct {
	Nested ABC
	A
	B
	C bool
	X string
	Y string
	Z string
}

// TwiceNested ...
type TwiceNested struct {
	Nest1 Nest `doc:"struct Nest1: Nest"`
	Nest2 Nest `doc:"struct Nest2: Nest"`
	A
	B
	C bool
	X string
}

// O2S object to yaml string
func O2S(o interface{}) string {
	text, err := yaml.Marshal(o)
	if err != nil {
		return fmt.Sprintf("O2S: Error: %+v", err)
	}
	return string(text)
}

// main ...
func main() {
	os.Setenv("NEST_A", "Value: NEST_A")
	os.Setenv("NEST_B", "Value: NEST_B")
	os.Setenv("NEST_C", "Value: NEST_C")
	os.Setenv("A_A", "Value: A_A")
	os.Setenv("B_B", "Value: A_B")
	os.Setenv("ABC_A", "Value: ABC_A")
	os.Setenv("ABC_B", "Value: ABC_B")
	os.Setenv("ABC_C", "Value: ABC_C")
	os.Setenv("A", "Value: A")
	os.Setenv("B", "Value: B")

	os.Setenv("NEST_A_A", "Value: NEST_A_A")
	os.Setenv("NEST_B_B", "Value: NEST_B_B")
	os.Setenv("NEST_C", "true")
	os.Setenv("NEST_NESTED_A_A", "Value: NEST_NESTED_A_A")
	os.Setenv("NEST_NESTED_B_B", "Value: NEST_NESTED_B_B")
	os.Setenv("NEST_NESTED_C", "true")
	os.Setenv("NEST_X", "Value: NEST_X")
	os.Setenv("NEST_Y", "Value: NEST_Y")
	os.Setenv("NEST_Z", "Value: NEST_Z")

	os.Setenv("TWICE_NESTED_NEST1_A_A", "Value: TWICE_NESTED_NEST1_A_A")
	os.Setenv("TWICE_NESTED_NEST1_B_B", "Value: TWICE_NESTED_NEST1_B_B")
	os.Setenv("TWICE_NESTED_NEST1_C", "true")
	os.Setenv("TWICE_NESTED_NEST1_NESTED_A_A", "Value: TWICE_NESTED_NEST1_NESTED_A_A")
	os.Setenv("TWICE_NESTED_NEST1_NESTED_B_B", "Value: TWICE_NESTED_NEST1_NESTED_B_B")
	os.Setenv("TWICE_NESTED_NEST1_NESTED_C", "true")

	os.Setenv("TWICE_NESTED_X", "Value: TWICE_NESTED_X")

	os.Setenv("TWICE_NESTED_NEST2_A_A", "Value: TWICE_NESTED_NEST2_A_A")
	os.Setenv("TWICE_NESTED_NEST2_B_B", "Value: TWICE_NESTED_NEST2_B_B")
	os.Setenv("TWICE_NESTED_NEST2_C", "true")
	os.Setenv("TWICE_NESTED_NEST2_NESTED_A_A", "Value: TWICE_NESTED_NEST2_NESTED_A_A")
	os.Setenv("TWICE_NESTED_NEST2_NESTED_B_B", "Value: TWICE_NESTED_NEST2_NESTED_B_B")
	os.Setenv("TWICE_NESTED_NEST2_NESTED_C", "true")
	os.Setenv("TWICE_NESTED_NEST2_X", "Value: TWICE_NESTED_NEST2_X")

	os.Setenv("TWICE_NESTED_X", "Value: NEST_Y")

	Run()
}

// Run ...
func Run() {

	a := &A{A: "a"}
	cfg.Nest(a)
	fmt.Println(O2S(a))

	b := &B{}
	cfg.Nest(b)
	fmt.Println(O2S(b))

	ab := &AB{}
	cfg.Nest(ab)
	fmt.Println(O2S(ab))

	abc := &ABC{}
	cfg.Nest(abc)
	fmt.Println(O2S(abc))

	nest := &Nest{}
	cfg.Nest(nest)
	fmt.Println(O2S(nest))

	twicenested := &TwiceNested{}
	// cfg.NestWrap("2", twicenested)
	cfg.Nest(twicenested)
	fmt.Println(O2S(twicenested))
	cfg.Nest(twicenested)
	fmt.Println(O2S(twicenested))

	cfg.Freeze()
	// flag.Parse()
	// fmt.Println(O2S(twicenested))
	internalRep := cfg.Store
	internalRep.Save("eval.yaml")
	// cfg.Usage()
	mgr := cfg.NewStor()
	mgr.Load("eval.yaml")
	// fmt.Printf(">>%+v\n", mgr)
	data, err := yaml.Marshal(mgr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(data))
	cfg.Usage()
}
