package cfg

import (
	"reflect"
	"testing"
)

func TestRemovePkg(t *testing.T) {
	type args struct {
		pkg string
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

			got1 := RemovePkg(tArgs.pkg)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RemovePkg got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestIsStructPtr(t *testing.T) {
	type args struct {
		config interface{}
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

			got1 := IsStructPtr(tArgs.config)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IsStructPtr got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestIsMap(t *testing.T) {
	type args struct {
		config interface{}
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

			got1 := IsMap(tArgs.config)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IsMap got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestType(t *testing.T) {
	type args struct {
		config interface{}
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

			got1 := Type(tArgs.config)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Type got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		config interface{}
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

			Copy(tArgs.config)

		})
	}
}
