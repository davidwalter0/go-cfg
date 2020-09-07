package cfg

import (
	"reflect"
	"testing"
)

func TestColor(t *testing.T) {
	type args struct {
		colorString string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want func(...interface{}) string
	}{
		{
			name: "red",
			args: func(t *testing.T) args {
				return args{colorString: Red("text")}
			},
			want: func(...interface{}) string {
				return Color("\033[1;31m%s\033[0m")("text")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			got := tArgs.colorString
			// fmt.Println("Color got", got, "want:", tt.want())
			if !reflect.DeepEqual(got, tt.want()) {
				t.Log("Color got", got, "want:", tt.want())
				t.Errorf("Color got = %v, want: %v", got, tt.want())
			}
		})
	}
}
