package cfg

import (
	"reflect"
	"testing"
)

func TestDebug(t *testing.T) {
	type args struct {
		args []bool
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

			Debug(tArgs.args...)

		})
	}
}

func TestField_Get(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *Field
		inspect func(r *Field, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1 := receiver.Get(tArgs.name)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Field.Get got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestField_SetField(t *testing.T) {
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
			receiver.SetField()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_String(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *Field
		inspect func(r *Field, t *testing.T) //inspects receiver after test run

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
				t.Errorf("Field.String got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestField_SetOmit(t *testing.T) {
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
			receiver.SetOmit()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetDefault(t *testing.T) {
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
			receiver.SetDefault()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetIgnore(t *testing.T) {
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
			receiver.SetIgnore()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetDoc(t *testing.T) {
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
			receiver.SetDoc()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetShort(t *testing.T) {
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
			receiver.SetShort()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetName(t *testing.T) {
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
			receiver.SetName()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetRequired(t *testing.T) {
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
			receiver.SetRequired()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetValueFromEnv(t *testing.T) {
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
			receiver.SetValueFromEnv()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetKeyName(t *testing.T) {
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
			receiver.SetKeyName()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetFlagName(t *testing.T) {
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
			receiver.SetFlagName()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestField_SetType(t *testing.T) {
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
			receiver.SetType()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}
