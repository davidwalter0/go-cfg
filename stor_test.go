package cfg

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func TestStor_DeepCopyInto(t *testing.T) {
	type args struct {
		dst       interface{}
		depthArgs []int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) Stor
		inspect func(r Stor, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			receiver.DeepCopyInto(tArgs.dst, tArgs.depthArgs...)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestDeepCopyInto(t *testing.T) {
	type args struct {
		src       Stor
		dst       map[string]interface{}
		depthArgs []int
	}
	tests := []struct {
		name string
		args func(t *testing.T) args
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			DeepCopyInto(tArgs.src, tArgs.dst, tArgs.depthArgs...)

		})
	}
}

func TestNewStor(t *testing.T) {
	tests := []struct {
		name string

		want1 Stor
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := NewStor()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewStor got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestStor_AddStor(t *testing.T) {
	type args struct {
		name string
		o    interface{}
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) Stor
		inspect func(r Stor, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			receiver.AddStor(tArgs.name, tArgs.o)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestStor_Load(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *Stor
		inspect func(r *Stor, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			err := receiver.Load(tArgs.filename)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("Stor.Load error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCopyOut(t *testing.T) {
	type args struct {
		in  interface{}
		out interface{}
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

			err := CopyOut(tArgs.in, tArgs.out)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CopyOut error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestDump(t *testing.T) {
	type args struct {
		o interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := Dump(tArgs.o)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Dump got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestStor_Stor(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) Stor
		inspect func(r Stor, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			err := receiver.Stor(tArgs.filename)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("Stor.Stor error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestStor_Save(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) Stor
		inspect func(r Stor, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			err := receiver.Save(tArgs.filename)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("Stor.Save error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestStor_Bytes(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) Stor
		inspect func(r Stor, t *testing.T) //inspects receiver after test run

		want1 []byte
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.Bytes()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Stor.Bytes got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestStor_String(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) Stor
		inspect func(r Stor, t *testing.T) //inspects receiver after test run

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.String()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Stor.String got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
