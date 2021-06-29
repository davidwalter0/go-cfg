package cfg

import (
	"reflect"
	"testing"
)

func TestField_AddFlag(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *Field
		inspect func(r *Field, t *testing.T) //inspects receiver after test run

	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			receiver.AddFlag()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestIsSet(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := IsSet(tArgs.name)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IsSet got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestOk(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := Ok(tArgs.name)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Ok got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := Required(tArgs.name)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Required got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
