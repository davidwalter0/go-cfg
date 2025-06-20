package types

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEnumValue_Set(t *testing.T) {
	validValues := []string{"debug", "info", "warn", "error", "fatal"}
	
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "valid lowercase",
			input:     "debug",
			wantError: false,
			expected:  "debug",
		},
		{
			name:      "valid uppercase",
			input:     "DEBUG",
			wantError: false,
			expected:  "DEBUG",
		},
		{
			name:      "valid mixed case",
			input:     "WaRn",
			wantError: false,
			expected:  "WaRn",
		},
		{
			name:      "valid with whitespace",
			input:     " info ",
			wantError: false,
			expected:  " info ",
		},
		{
			name:      "invalid value",
			input:     "invalid",
			wantError: true,
		},
		{
			name:      "empty string not allowed",
			input:     "",
			wantError: true,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ev := NewEnumValue(validValues, "", false)
			err := ev.Set(test.input)
			
			if (err != nil) != test.wantError {
				t.Errorf("EnumValue.Set(%s) error = %v, wantError = %v", test.input, err, test.wantError)
				return
			}
			
			if !test.wantError {
				result := ev.GetValue()
				if result != test.expected {
					t.Errorf("EnumValue.Set(%s) result = %s, want %s", test.input, result, test.expected)
				}
			}
		})
	}
}

func TestEnumValue_SetWithEmptyAllowed(t *testing.T) {
	validValues := []string{"option1", "option2", "option3"}
	ev := NewEnumValue(validValues, "", true)
	
	// Test empty string allowed
	err := ev.Set("")
	if err != nil {
		t.Errorf("EnumValue.Set(\"\") should not error when empty is allowed, got: %v", err)
	}
	
	if ev.GetValue() != "" {
		t.Errorf("EnumValue.GetValue() after setting empty should be empty, got: %s", ev.GetValue())
	}
}

func TestEnumValue_IsValid(t *testing.T) {
	validValues := []string{"debug", "info", "warn", "error", "fatal"}
	ev := NewEnumValue(validValues, "info", false)
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"debug", true},
		{"DEBUG", true},
		{"Info", true},
		{"WARN", true},
		{"invalid", false},
		{"", false},
	}
	
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := ev.IsValid(test.input)
			if result != test.expected {
				t.Errorf("EnumValue.IsValid(%s) = %v, want %v", test.input, result, test.expected)
			}
		})
	}
}

func TestEnumValue_GetValidValues(t *testing.T) {
	validValues := []string{"debug", "info", "warn", "error", "fatal"}
	ev := NewEnumValue(validValues, "info", false)
	
	result := ev.GetValidValues()
	if len(result) != len(validValues) {
		t.Errorf("EnumValue.GetValidValues() length = %d, want %d", len(result), len(validValues))
	}
	
	// Verify it's a copy (modifying returned slice shouldn't affect original)
	result[0] = "modified"
	if ev.GetValidValues()[0] == "modified" {
		t.Errorf("EnumValue.GetValidValues() should return a copy, not the original slice")
	}
}

func TestEnumValue_JSONMarshaling(t *testing.T) {
	validValues := []string{"debug", "info", "warn", "error", "fatal"}
	ev := NewEnumValue(validValues, "info", false)
	
	// Test marshaling
	data, err := json.Marshal(ev)
	if err != nil {
		t.Errorf("EnumValue.MarshalJSON() error = %v", err)
		return
	}
	
	expected := `"info"`
	if string(data) != expected {
		t.Errorf("EnumValue.MarshalJSON() = %s, want %s", string(data), expected)
	}
	
	// Test unmarshaling
	var newEv EnumValue
	newEv.validValues = make([]string, len(validValues))
	copy(newEv.validValues, validValues)
	
	err = json.Unmarshal(data, &newEv)
	if err != nil {
		t.Errorf("EnumValue.UnmarshalJSON() error = %v", err)
		return
	}
	
	if newEv.GetValue() != "info" {
		t.Errorf("EnumValue JSON round-trip: got %s, want info", newEv.GetValue())
	}
}

func TestEnumValue_YAMLMarshaling(t *testing.T) {
	validValues := []string{"debug", "info", "warn", "error", "fatal"}
	ev := NewEnumValue(validValues, "warn", false)
	
	// Test marshaling
	data, err := yaml.Marshal(ev)
	if err != nil {
		t.Errorf("EnumValue.MarshalYAML() error = %v", err)
		return
	}
	
	expected := "warn\n"
	if string(data) != expected {
		t.Errorf("EnumValue.MarshalYAML() = %s, want %s", string(data), expected)
	}
	
	// Test unmarshaling
	var newEv EnumValue
	newEv.validValues = make([]string, len(validValues))
	copy(newEv.validValues, validValues)
	
	err = yaml.Unmarshal(data, &newEv)
	if err != nil {
		t.Errorf("EnumValue.UnmarshalYAML() error = %v", err)
		return
	}
	
	if newEv.GetValue() != "warn" {
		t.Errorf("EnumValue YAML round-trip: got %s, want warn", newEv.GetValue())
	}
}

func TestLogLevelValue_Methods(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		debugEnabled bool
		infoEnabled  bool
		warnEnabled  bool
	}{
		{
			name:         "debug level",
			level:        "debug",
			debugEnabled: true,
			infoEnabled:  true,
			warnEnabled:  true,
		},
		{
			name:         "info level",
			level:        "info",
			debugEnabled: false,
			infoEnabled:  true,
			warnEnabled:  true,
		},
		{
			name:         "warn level",
			level:        "warn",
			debugEnabled: false,
			infoEnabled:  false,
			warnEnabled:  true,
		},
		{
			name:         "error level",
			level:        "error",
			debugEnabled: false,
			infoEnabled:  false,
			warnEnabled:  false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			llv := NewLogLevelValue(test.level)
			
			if llv.IsDebugEnabled() != test.debugEnabled {
				t.Errorf("LogLevelValue(%s).IsDebugEnabled() = %v, want %v", 
					test.level, llv.IsDebugEnabled(), test.debugEnabled)
			}
			
			if llv.IsInfoEnabled() != test.infoEnabled {
				t.Errorf("LogLevelValue(%s).IsInfoEnabled() = %v, want %v", 
					test.level, llv.IsInfoEnabled(), test.infoEnabled)
			}
			
			if llv.IsWarnEnabled() != test.warnEnabled {
				t.Errorf("LogLevelValue(%s).IsWarnEnabled() = %v, want %v", 
					test.level, llv.IsWarnEnabled(), test.warnEnabled)
			}
			
			if llv.GetLevel() != test.level {
				t.Errorf("LogLevelValue(%s).GetLevel() = %s, want %s", 
					test.level, llv.GetLevel(), test.level)
			}
		})
	}
}

func TestEnvironmentValue_Methods(t *testing.T) {
	tests := []struct {
		name          string
		env           string
		isProduction  bool
		isDevelopment bool
		isStaging     bool
		isTesting     bool
	}{
		{
			name:          "production",
			env:           "production",
			isProduction:  true,
			isDevelopment: false,
			isStaging:     false,
			isTesting:     false,
		},
		{
			name:          "prod",
			env:           "prod",
			isProduction:  true,
			isDevelopment: false,
			isStaging:     false,
			isTesting:     false,
		},
		{
			name:          "development",
			env:           "development",
			isProduction:  false,
			isDevelopment: true,
			isStaging:     false,
			isTesting:     false,
		},
		{
			name:          "dev",
			env:           "dev",
			isProduction:  false,
			isDevelopment: true,
			isStaging:     false,
			isTesting:     false,
		},
		{
			name:          "staging",
			env:           "staging",
			isProduction:  false,
			isDevelopment: false,
			isStaging:     true,
			isTesting:     false,
		},
		{
			name:          "test",
			env:           "test",
			isProduction:  false,
			isDevelopment: false,
			isStaging:     false,
			isTesting:     true,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			env := NewEnvironmentValue(test.env)
			
			if env.IsProduction() != test.isProduction {
				t.Errorf("EnvironmentValue(%s).IsProduction() = %v, want %v", 
					test.env, env.IsProduction(), test.isProduction)
			}
			
			if env.IsDevelopment() != test.isDevelopment {
				t.Errorf("EnvironmentValue(%s).IsDevelopment() = %v, want %v", 
					test.env, env.IsDevelopment(), test.isDevelopment)
			}
			
			if env.IsStaging() != test.isStaging {
				t.Errorf("EnvironmentValue(%s).IsStaging() = %v, want %v", 
					test.env, env.IsStaging(), test.isStaging)
			}
			
			if env.IsTesting() != test.isTesting {
				t.Errorf("EnvironmentValue(%s).IsTesting() = %v, want %v", 
					test.env, env.IsTesting(), test.isTesting)
			}
			
			if env.GetEnvironment() != test.env {
				t.Errorf("EnvironmentValue(%s).GetEnvironment() = %s, want %s", 
					test.env, env.GetEnvironment(), test.env)
			}
		})
	}
}

func TestMultiEnumValue_Set(t *testing.T) {
	validValues := []string{"read", "write", "execute", "delete"}
	
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  []string
	}{
		{
			name:      "single value",
			input:     "read",
			wantError: false,
			expected:  []string{"read"},
		},
		{
			name:      "multiple values",
			input:     "read,write,execute",
			wantError: false,
			expected:  []string{"read", "write", "execute"},
		},
		{
			name:      "values with whitespace",
			input:     "read, write , execute",
			wantError: false,
			expected:  []string{"read", "write", "execute"},
		},
		{
			name:      "duplicate values",
			input:     "read,write,read",
			wantError: false,
			expected:  []string{"read", "write"}, // Duplicates should be removed
		},
		{
			name:      "invalid value",
			input:     "read,invalid,write",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: false,
			expected:  []string{}, // Should clear the values
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mev := NewMultiEnumValue(validValues, nil, true)
			err := mev.Set(test.input)
			
			if (err != nil) != test.wantError {
				t.Errorf("MultiEnumValue.Set(%s) error = %v, wantError = %v", test.input, err, test.wantError)
				return
			}
			
			if !test.wantError {
				result := mev.GetValues()
				if len(result) != len(test.expected) {
					t.Errorf("MultiEnumValue.Set(%s) length = %d, want %d", 
						test.input, len(result), len(test.expected))
					return
				}
				
				for i, expected := range test.expected {
					if result[i] != expected {
						t.Errorf("MultiEnumValue.Set(%s) result[%d] = %s, want %s", 
							test.input, i, result[i], expected)
					}
				}
			}
		})
	}
}

func TestMultiEnumValue_Contains(t *testing.T) {
	validValues := []string{"read", "write", "execute"}
	mev := NewMultiEnumValue(validValues, nil, false)
	// Set the values we want to test
	mev.Set("read,write")
	
	tests := []struct {
		value    string
		expected bool
	}{
		{"read", true},
		{"write", true},
		{"execute", false},
		{"READ", true}, // Case insensitive
		{"invalid", false},
	}
	
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			result := mev.Contains(test.value)
			if result != test.expected {
				t.Errorf("MultiEnumValue.Contains(%s) = %v, want %v", test.value, result, test.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkEnumValue_Set(b *testing.B) {
	validValues := []string{"debug", "info", "warn", "error", "fatal"}
	ev := NewEnumValue(validValues, "info", false)
	testValues := []string{"debug", "info", "warn", "error", "fatal"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, value := range testValues {
			ev.Set(value)
		}
	}
}

func BenchmarkLogLevelValue_IsDebugEnabled(b *testing.B) {
	llv := NewLogLevelValue("debug")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		llv.IsDebugEnabled()
	}
}