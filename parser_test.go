package cfg

import (
	"reflect"
	"testing"
)

func TestSFWrap_GetName(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *SFWrap
		inspect func(r *SFWrap, t *testing.T) //inspects receiver after test run

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.GetName()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SFWrap.GetName got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestEnter(t *testing.T) {
	type args struct {
		args *Arg
		ptr  interface{}
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

			err := Enter(tArgs.args, tArgs.ptr)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Enter error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestParseStruct(t *testing.T) {
	type args struct {
		args        *Arg
		ptr         interface{}
		prefix      string
		structField reflect.StructField
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

			err := ParseStruct(tArgs.args, tArgs.ptr, tArgs.prefix, tArgs.structField)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseStruct error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
