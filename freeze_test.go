package cfg

import (
	"testing"
)

func TestFreeze(t *testing.T) {
	tests := []struct {
		name string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Freeze()

		})
	}
}

func TestFlagInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FlagInit()

		})
	}
}
