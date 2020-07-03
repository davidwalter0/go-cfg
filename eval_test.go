package cfg // import "github.com/davidwalter0/go-cfg"

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/davidwalter0/go-flag"
)

// A ...
type A struct {
	A string `json:"a"`
}

// B ...
type B struct {
	B string
}

// AB ...
type AB struct {
	A
	B
}

// ABC ...
type ABC struct {
	// A string
	// B int
	A
	B
	C bool
}

// ABC ...
type NestStruct struct {
	Nested ABC
	A
	B
	C bool
	X string
	Y string
	Z string
}

func _ignore_() {
	flag.Parse()
}

// O2S
func O2S(o interface{}) string {
	text, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return fmt.Sprintf("O2S: Error: %+v", err)
	}
	return string(text)
}

// test_setup ...
func test_setup() {
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
}

/*
// TestParseA
func TestParseA(t *testing.T) {
	// fmt.Printf("%+v\n", Eval(&A{}))
	a := &A{A: "a"}
	// Eval(a)
	// fmt.Println(O2S(a))
	Init(a)
	fmt.Println(O2S(a))

	///////// cfg = Init(o)
	///////// fmt.Println(O2S(cfg))
	// }
	// // TestParseB
	// func TestParseB(t *testing.T) {
	// fmt.Printf("%+v\n", Eval(&B{}))
	b := &B{}
	Eval(b)
	fmt.Println(O2S(b))
	///////// cfg = Init(o)
	///////// fmt.Println(O2S(cfg))
	// }

// TestParseABC
func TestParseABC(t *testing.T) {
	fmt.Printf("%+v\n", Eval(&ABC{}))
	ab := &AB{}
	Eval(ab)
	fmt.Println(O2S(ab))

	// abc := &ABC{}
	// Eval(abc)
	// fmt.Println(O2S(abc))
	// Init(abc)
	// fmt.Println(O2S(abc))
	// flag.Parse()
	// fmt.Println(O2S(abc))
	// nest := &NestStruct{}
	// Eval(nest)
	// fmt.Println(O2S(nest))
	// flag.Parse()
	// fmt.Println(O2S(nest))
	// Usage()
}

*/

// Decl declaration
type Decl struct {
	S string `json:"s"`
	I int    `json:"i"`
}

// Object arbitrary test object
type Object struct {
	Text string
	Decl Decl
}

type State struct {
	Want Object
	Got  Object
}

var state = State{
	Want: Object{
		`{s:"text",i:0}`,
		Decl{
			S: "text",
			I: 0,
		},
	},
	Got: Object{
		`{s:"text",i:0}`,
		Decl{
			S: "text",
			I: 0,
		},
	},
}

var states = []State{
	State{
		Want: Object{
			Text: `{s:"text",i:0}`,
			Decl: Decl{
				S: "text",
				I: 0,
			},
		},
		Got: Object{
			Text: `{s:"text",i:0}`,
			Decl: Decl{
				S: "text",
				I: 0,
			},
		},
	},
	State{
		Want: Object{
			Text: `{s:"text",i:0}`,
			Decl: Decl{
				S: "text",
				I: 0,
			},
		},
		Got: Object{
			Text: `{ s: "text", i: 1 }`,
			Decl: Decl{
				S: "text",
				I: 1,
			},
		},
	},
}

// Compare two objects
func (lhs *Object) Compare(rhs *Object) bool {
	return lhs.Text == rhs.Text &&
		lhs.Decl.S == rhs.Decl.S &&
		lhs.Decl.I == rhs.Decl.I
}

// Ok when want == got
func (s *State) Ok() bool {
	return s.Want.Compare(&s.Got)
}

func TestEvalName(t *testing.T) {
	test_setup()
	if !states[0].Ok() {
		t.Fatal()
	}
	if states[1].Ok() {
		t.Fatal()
	}
	for i, state := range states {
		// t.Logf("Want: %+v Got: %+v\n", state.Want, state.Got)
		text, err := json.Marshal(state.Want.Decl)
		if err != nil {
			t.Log(Caller())
			t.Logf("%d Want: %+v Got: %+v\n   marshal %v\n", i, state.Want, state.Got, err)
			t.Fail()
		}
		if i == 0 && !state.Ok() {
			t.Log(Caller())
			t.Logf("%d Want: %+v Got: %+v\n", i, state.Want, state.Got)
			t.Fail()
		}
		if i == 1 && state.Ok() {
			t.Log(Caller())
			t.Logf("%d Want: %+v Got: %+v\n", i, state.Want, state.Got)
			t.Fail()
		}
		if err != nil {
			t.Log(Caller())
			t.Logf("%d Want: %+v Got: %+v\n   marshal %v\n", i, state.Want, state.Got, err)
			t.Fail()
		}

		if false {
			t.Logf(">> Want: %+v Got: %+v\n   text%s\n", state.Want, state.Got, text)
		}
	}
}
