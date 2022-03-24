package cfg_test // import "github.com/davidwalter0/go-cfg"

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/davidwalter0/go-cfg"
	eflag "github.com/davidwalter0/go-flag"
)

func _ignore_() {
	eflag.Parse()
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
	// eflag.Parse()
	// fmt.Println(O2S(abc))
	// nest := &NestStruct{}
	// Eval(nest)
	// fmt.Println(O2S(nest))
	// eflag.Parse()
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
			t.Log(cfg.Caller())
			t.Logf("%d Want: %+v Got: %+v\n   marshal %v\n", i, state.Want, state.Got, err)
			t.Fail()
		}
		if i == 0 && !state.Ok() {
			t.Log(cfg.Caller())
			t.Logf("%d Want: %+v Got: %+v\n", i, state.Want, state.Got)
			t.Fail()
		}
		if i == 1 && state.Ok() {
			t.Log(cfg.Caller())
			t.Logf("%d Want: %+v Got: %+v\n", i, state.Want, state.Got)
			t.Fail()
		}
		if err != nil {
			t.Log(cfg.Caller())
			t.Logf("%d Want: %+v Got: %+v\n   marshal %v\n", i, state.Want, state.Got, err)
			t.Fail()
		}

		if false {
			t.Logf(">> Want: %+v Got: %+v\n   text%s\n", state.Want, state.Got, text)
		}
	}
}

func TestNewArg(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 *cfg.Arg
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := cfg.NewArg(tArgs.name)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewArg got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestUndecorate(t *testing.T) {
	tests := []struct {
		name string

		want1 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := cfg.Undecorate()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Undecorate got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestDecorate(t *testing.T) {
	tests := []struct {
		name string

		want1 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := cfg.Decorate()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Decorate got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestUnprefix(t *testing.T) {
	tests := []struct {
		name string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Unprefix()

		})
	}
}

func TestEval(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Eval(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Eval error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Init(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Init error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		args *cfg.Arg
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Run(tArgs.args, tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Run error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestSimple(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Simple(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Simple error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestSimpleFlags(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.SimpleFlags(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("SimpleFlags error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Unwrap(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Unwrap error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestFinal(t *testing.T) {
	tests := []struct {
		name string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Final()

		})
	}
}

func TestBare(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Bare(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Bare error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		name string
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Wrap(tArgs.name, tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Wrap error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Add(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Add error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestFlags(t *testing.T) {
	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := cfg.Flags(tArgs.ptrs...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Flags error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestNest(t *testing.T) {
	type args struct {
		prefix string
		ptrs   []interface{}
		AName  string
		AValue string
		BName  string
		BValue string

		set func()
		clr func()
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
		AMatch     bool
		BMatch     bool
		// set flags values
		Flags []string
	}{
		{
			name: "use defaults",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "A_A",
					AValue: "env-value-a",
					BName:  "B_B",
					BValue: "env-value-b",
				}
				rc.set = func() {
					set(rc.AName, rc.AValue)
					set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
		},
		{
			name: "use flag",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "A_A",
					AValue: "flag-value-a",
					BName:  "B_B",
					BValue: "env-value-b",
				}
				rc.set = func() {
					// mismatch value from env
					set(rc.AName, "X"+rc.AValue)
					set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
			Flags:  []string{"--a-a=flag-value-a"},
		},
		{
			name: "set both values with flags",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "A_A",
					AValue: "flag-value-a",
					BName:  "B_B",
					BValue: "flag-value-b",
				}
				rc.set = func() {
					// set(rc.AName, rc.AValue)
					// set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
			Flags: []string{
				"--a-a=flag-value-a",
				"--b-b=flag-value-b",
			},
		},
		{
			name: "use env and flag value mismatch",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "A_A",
					AValue: "env-value-a",
					BName:  "B_B",
					BValue: "value-b",
				}
				rc.set = func() {
					// set(rc.AName, rc.AValue)
					set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				return rc
			},
			AMatch: false,
			BMatch: false,
			Flags:  []string{"--b-b=flag-value-b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer flagReset()()
			tArgs := tt.args(t)
			// fmt.Println(os.Args)
			os.Args = append(os.Args, tt.Flags...)
			// fmt.Println(os.Args)
			tArgs.set()
			defer tArgs.clr()
			err := cfg.Nest(tArgs.ptrs...)
			eflag.Parse()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Nest error = %v, wantErr: %t", err, tt.wantErr)
			}
			if (tArgs.AValue == tArgs.ptrs[0].(*A).A) != tt.AMatch {
				t.Fatalf("Nest expected %v %v %v", tArgs.AValue, Match(tt.AMatch), tArgs.ptrs[0].(*A).A)
			}

			if (tArgs.BValue == tArgs.ptrs[1].(*B).B) != tt.BMatch {
				t.Fatalf("Nest expected %v %v %v", tArgs.BValue, Match(tt.BMatch), tArgs.ptrs[1].(*B).B)
			}
		})
	}
}

func TestNestWrap(t *testing.T) {

	type args struct {
		prefix string
		ptrs   []interface{}
		AName  string
		AValue string
		BName  string
		BValue string

		set func()
		see func()
		clr func()
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
		AMatch     bool
		BMatch     bool
		// set flags values
		Flags []string
	}{
		{
			name: "use defaults",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "PREFIX_A_A",
					AValue: "value-a",
					BName:  "PREFIX_B_B",
					BValue: "value-b",
				}
				rc.set = func() {
					// set(rc.AName, rc.AValue)
					// set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				rc.see = func() {
					see(rc.AName)
					see(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
		},
		{
			name: "use env settings",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "PREFIX_A_A",
					AValue: "env-value-a",
					BName:  "PREFIX_B_B",
					BValue: "env-value-b",
				}
				rc.set = func() {
					set(rc.AName, rc.AValue)
					set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				rc.see = func() {
					see(rc.AName)
					see(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
		},
		{
			name: "use flag for a",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "PREFIX_A_A",
					AValue: "flag-value-a",
					BName:  "PREFIX_B_B",
					BValue: "env-value-b",
				}
				rc.set = func() {
					// reject env for flag
					set(rc.AName, "X"+rc.AValue)
					set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				rc.see = func() {
					see(rc.AName)
					see(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
			Flags:  []string{"--prefix-a-a=flag-value-a"},
		},
		{
			name: "set both values with flags",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "PREFIX_A_A",
					AValue: "flag-value-a",
					BName:  "PREFIX_B_B",
					BValue: "flag-value-b",
				}
				rc.set = func() {
					// set(rc.AName, rc.AValue)
					// set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				rc.see = func() {
					see(rc.AName)
					see(rc.BName)
				}
				return rc
			},
			AMatch: true,
			BMatch: true,
			Flags: []string{
				"--prefix-a-a=flag-value-a",
				"--prefix-b-b=flag-value-b",
			},
		},
		{
			name: "use env and flag value mismatch",
			args: func(t *testing.T) args {
				rc := args{
					prefix: "Prefix",
					ptrs:   []interface{}{&A{}, &B{}},
					AName:  "PREFIX_A_A",
					AValue: "env-value-a",
					BName:  "PREFIX_B_B",
					BValue: "value-b",
				}
				rc.set = func() {
					// set(rc.AName, rc.AValue)
					set(rc.BName, rc.BValue)
				}
				rc.clr = func() {
					clr(rc.AName)
					clr(rc.BName)
				}
				rc.see = func() {
					see(rc.AName)
					see(rc.BName)
				}
				return rc
			},
			AMatch: false,
			BMatch: false,
			Flags:  []string{"--prefix-b-b=flag-value-b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer flagReset()()
			tArgs := tt.args(t)
			// fmt.Println(os.Args)
			os.Args = append(os.Args, tt.Flags...)
			// fmt.Println(os.Args)
			tArgs.set()
			defer tArgs.clr()
			err := cfg.NestWrap(tArgs.prefix, tArgs.ptrs...)
			eflag.Parse()
			// tArgs.see()
			// cfg.Usage()
			if (err != nil) != tt.wantErr {
				t.Fatalf("NestWrap error = %v, wantErr: %t", err, tt.wantErr)
			}
			if (tArgs.AValue == tArgs.ptrs[0].(*A).A) != tt.AMatch {
				t.Fatalf("NestWrap expected %v %v %v", tArgs.AValue, Match(tt.AMatch), tArgs.ptrs[0].(*A).A)
			}

			if (tArgs.BValue == tArgs.ptrs[1].(*B).B) != tt.BMatch {
				t.Fatalf("NestWrap expected %v %v %v", tArgs.BValue, Match(tt.BMatch), tArgs.ptrs[1].(*B).B)
			}
		})
	}
}

func TestCheckArgs(t *testing.T) {
	type A struct {
	}
	type B struct {
	}

	type args struct {
		ptrs   []interface{}
		expect error
	}
	tests := []struct {
		name  string
		args  func(t *testing.T) args
		isNil bool
	}{
		{
			name: "ptrs ok",
			args: func(t *testing.T) args {
				return args{
					ptrs:   []interface{}{&A{}, &B{}},
					expect: nil,
				}
			},
			isNil: true,
		},
		{
			name: "ptrs aren't pointers",
			args: func(t *testing.T) args {
				return args{
					ptrs:   []interface{}{A{}, &B{}},
					expect: fmt.Errorf(""),
				}
			},
			isNil: false,
		},
		{
			name: "ptrs are ** pointers to pointers",
			args: func(t *testing.T) args {
				a, b := &A{}, B{}
				return args{
					ptrs:   []interface{}{&a, &b},
					expect: fmt.Errorf(""),
				}
			},
			isNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tArgs := tt.args(t)
			err = cfg.CheckArgs(tArgs.ptrs...)
			// Using unwrap might enable a constant test response and still
			// fulfill the end user output formatting

			// if err != tArgs.expect {
			// 	t.Errorf("want %v got: %v", tArgs.expect, err)
			// }
			if (err == nil) != tt.isNil {
				t.Errorf("want <nil> error got: %v", err)
			}
		})
	}
}

func TestCaller(t *testing.T) {
	tests := []struct {
		name string

		want1 string
	}{
		{
			name:  "Check if Caller() returns parent and path",
			want1: "go-cfg_test.TestCaller.func",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.Caller()
			regx := regexp.MustCompile("\\.go:[0-9]*:([^/]*/)*go-cfg_test\\.TestCaller\\.func.*")
			if len(regx.Find([]byte(got))) <= 0 {
				t.Errorf("Caller got1 = %v, want1: %v", got, tt.want1)
			}
		})
	}
}
