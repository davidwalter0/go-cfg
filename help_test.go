package cfg

import (
	"testing"
)

func TestHelpText(t *testing.T) {
	type args struct {
		text string
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

			HelpText(tArgs.text)

		})
	}
}
