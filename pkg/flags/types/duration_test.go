package types

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestDurationValue_Set(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantError  bool
		wantResult time.Duration
	}{
		{
			name:       "standard Go duration",
			input:      "1h30m",
			wantError:  false,
			wantResult: 90 * time.Minute,
		},
		{
			name:       "seconds only",
			input:      "45s",
			wantError:  false,
			wantResult: 45 * time.Second,
		},
		{
			name:       "integer seconds",
			input:      "90",
			wantError:  false,
			wantResult: 90 * time.Second,
		},
		{
			name:       "float seconds",
			input:      "90.5",
			wantError:  false,
			wantResult: time.Duration(90.5 * float64(time.Second)),
		},
		{
			name:       "minutes",
			input:      "5m",
			wantError:  false,
			wantResult: 5 * time.Minute,
		},
		{
			name:       "hours",
			input:      "2h",
			wantError:  false,
			wantResult: 2 * time.Hour,
		},
		{
			name:       "complex duration",
			input:      "2h30m45s",
			wantError:  false,
			wantResult: 2*time.Hour + 30*time.Minute + 45*time.Second,
		},
		{
			name:       "empty string",
			input:      "",
			wantError:  false,
			wantResult: 0,
		},
		{
			name:      "invalid format",
			input:     "not-a-duration",
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dv := NewDurationValue(0)
			err := dv.Set(test.input)

			if (err != nil) != test.wantError {
				t.Errorf("DurationValue.Set(%s) error = %v, wantError = %v", test.input, err, test.wantError)
				return
			}

			if !test.wantError {
				result := dv.GetDuration()
				if result != test.wantResult {
					t.Errorf("DurationValue.Set(%s) result = %v, want %v", test.input, result, test.wantResult)
				}
			}
		})
	}
}

func TestDurationValue_String(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30s",
		},
		{
			name:     "5 minutes",
			duration: 5 * time.Minute,
			expected: "5m0s",
		},
		{
			name:     "1 hour 30 minutes",
			duration: 90 * time.Minute,
			expected: "1h30m0s",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dv := NewDurationValue(test.duration)
			result := dv.String()
			if result != test.expected {
				t.Errorf("DurationValue.String() = %s, want %s", result, test.expected)
			}
		})
	}
}

func TestDurationValue_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30",
		},
		{
			name:     "90 seconds (1.5 minutes)",
			duration: 90 * time.Second,
			expected: "90",
		},
		{
			name:     "3600 seconds (1 hour)",
			duration: time.Hour,
			expected: "3600",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dv := NewDurationValue(test.duration)

			// Test marshaling
			data, err := json.Marshal(dv)
			if err != nil {
				t.Errorf("DurationValue.MarshalJSON() error = %v", err)
				return
			}

			if string(data) != test.expected {
				t.Errorf("DurationValue.MarshalJSON() = %s, want %s", string(data), test.expected)
			}

			// Test unmarshaling
			var newDv DurationValue
			err = json.Unmarshal(data, &newDv)
			if err != nil {
				t.Errorf("DurationValue.UnmarshalJSON() error = %v", err)
				return
			}

			if newDv.GetDuration() != test.duration {
				t.Errorf("DurationValue round-trip: got %v, want %v", newDv.GetDuration(), test.duration)
			}
		})
	}
}

func TestDurationValue_YAMLMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30s\n",
		},
		{
			name:     "90 seconds",
			duration: 90 * time.Second,
			expected: "1m30s\n",
		},
		{
			name:     "1 hour",
			duration: time.Hour,
			expected: "1h0m0s\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dv := NewDurationValue(test.duration)

			// Test marshaling
			data, err := yaml.Marshal(dv)
			if err != nil {
				t.Errorf("DurationValue.MarshalYAML() error = %v", err)
				return
			}

			if string(data) != test.expected {
				t.Errorf("DurationValue.MarshalYAML() = %s, want %s", string(data), test.expected)
			}

			// Test unmarshaling
			var newDv DurationValue
			err = yaml.Unmarshal(data, &newDv)
			if err != nil {
				t.Errorf("DurationValue.UnmarshalYAML() error = %v", err)
				return
			}

			if newDv.GetDuration() != test.duration {
				t.Errorf("DurationValue YAML round-trip: got %v, want %v", newDv.GetDuration(), test.duration)
			}
		})
	}
}

func TestHumanDurationValue_Set(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantError  bool
		wantResult time.Duration
	}{
		{
			name:       "standard Go duration",
			input:      "1h30m",
			wantError:  false,
			wantResult: 90 * time.Minute,
		},
		{
			name:       "week format",
			input:      "1w",
			wantError:  false,
			wantResult: 7 * 24 * time.Hour,
		},
		{
			name:       "day format",
			input:      "1d",
			wantError:  false,
			wantResult: 24 * time.Hour,
		},
		{
			name:       "human readable hours",
			input:      "2h",
			wantError:  false,
			wantResult: 2 * time.Hour,
		},
		{
			name:       "human readable minutes",
			input:      "30m",
			wantError:  false,
			wantResult: 30 * time.Minute,
		},
		{
			name:       "human readable seconds",
			input:      "45s",
			wantError:  false,
			wantResult: 45 * time.Second,
		},
		{
			name:       "mixed human format",
			input:      "1h30m",
			wantError:  false,
			wantResult: 90 * time.Minute,
		},
		{
			name:       "empty string",
			input:      "",
			wantError:  false,
			wantResult: 0,
		},
		{
			name:      "invalid format",
			input:     "not-a-duration",
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hdv := NewHumanDurationValue(0)
			err := hdv.Set(test.input)

			if (err != nil) != test.wantError {
				t.Errorf("HumanDurationValue.Set(%s) error = %v, wantError = %v", test.input, err, test.wantError)
				return
			}

			if !test.wantError {
				result := hdv.GetDuration()
				if result != test.wantResult {
					t.Errorf("HumanDurationValue.Set(%s) result = %v, want %v", test.input, result, test.wantResult)
				}
			}
		})
	}
}

func TestHumanDurationValue_String(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30s",
		},
		{
			name:     "5 minutes",
			duration: 5 * time.Minute,
			expected: "5m",
		},
		{
			name:     "1 hour 30 minutes",
			duration: 90 * time.Minute,
			expected: "1h30m",
		},
		{
			name:     "1 day",
			duration: 24 * time.Hour,
			expected: "1d",
		},
		{
			name:     "1 week",
			duration: 7 * 24 * time.Hour,
			expected: "7d",
		},
		{
			name:     "complex duration",
			duration: 2*24*time.Hour + 3*time.Hour + 30*time.Minute + 45*time.Second,
			expected: "2d3h30m45s",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hdv := NewHumanDurationValue(test.duration)
			result := hdv.String()
			if result != test.expected {
				t.Errorf("HumanDurationValue.String() = %s, want %s", result, test.expected)
			}
		})
	}
}

func TestDurationValue_Type(t *testing.T) {
	dv := NewDurationValue(0)
	if dv.Type() != "duration" {
		t.Errorf("DurationValue.Type() = %s, want duration", dv.Type())
	}
}

func TestHumanDurationValue_Type(t *testing.T) {
	hdv := NewHumanDurationValue(0)
	if hdv.Type() != "humanDuration" {
		t.Errorf("HumanDurationValue.Type() = %s, want humanDuration", hdv.Type())
	}
}

// Benchmark tests
func BenchmarkDurationValue_Set(b *testing.B) {
	dv := NewDurationValue(0)
	testInputs := []string{"1h30m", "90s", "45", "2h30m45s"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range testInputs {
			_ = dv.Set(input)
		}
	}
}

func BenchmarkHumanDurationValue_Set(b *testing.B) {
	hdv := NewHumanDurationValue(0)
	testInputs := []string{"1 hour", "30 minutes", "1 day", "1 week"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range testInputs {
			_ = hdv.Set(input)
		}
	}
}

func BenchmarkHumanDurationValue_String(b *testing.B) {
	hdv := NewHumanDurationValue(2*24*time.Hour + 3*time.Hour + 30*time.Minute + 45*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hdv.String()
	}
}
