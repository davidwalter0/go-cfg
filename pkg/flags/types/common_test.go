package types

import (
	"fmt"
	"os"
	"testing"
)

func TestNamingConventions_ToUpperSnakeCase(t *testing.T) {
	nc := &NamingConventions{}
	
	tests := []struct {
		input    string
		expected string
	}{
		{"APIEndpoint", "API_ENDPOINT"},
		{"DatabaseURL", "DATABASE_URL"},
		{"HTTPSProxy", "HTTPS_PROXY"},
		{"LogLevel", "LOG_LEVEL"},
		{"MaxRetries", "MAX_RETRIES"},
		{"TCPPort", "TCP_PORT"},
		{"XMLParser", "XML_PARSER"},
		{"JSONData", "JSON_DATA"},
		{"HTMLContent", "HTML_CONTENT"},
		{"UserID", "USER_ID"},
		{"ABCTest", "ABC_TEST"},
		{"TestABCTest", "TEST_ABC_TEST"},
	}
	
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := nc.ToUpperSnakeCase(test.input)
			if result != test.expected {
				t.Errorf("ToUpperSnakeCase(%s) = %s, want %s", test.input, result, test.expected)
			}
		})
	}
}

func TestNamingConventions_ToKebabCase(t *testing.T) {
	nc := &NamingConventions{}
	
	tests := []struct {
		input    string
		expected string
	}{
		{"APIEndpoint", "api-endpoint"},
		{"DatabaseURL", "database-url"},
		{"HTTPSProxy", "https-proxy"},
		{"LogLevel", "log-level"},
		{"MaxRetries", "max-retries"},
		{"TCPPort", "tcp-port"},
		{"XMLParser", "xml-parser"},
		{"JSONData", "json-data"},
		{"HTMLContent", "html-content"},
		{"UserID", "user-id"},
		{"ABCTest", "abc-test"},
		{"TestABCTest", "test-abc-test"},
	}
	
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := nc.ToKebabCase(test.input)
			if result != test.expected {
				t.Errorf("ToKebabCase(%s) = %s, want %s", test.input, result, test.expected)
			}
		})
	}
}

func TestGetEnvDefault(t *testing.T) {
	// Set up test environment variables
	testEnvVars := map[string]string{
		"API_ENDPOINT":  "https://test-api.example.com",
		"DATABASE_URL":  "postgres://test:test@localhost/test",
		"LOG_LEVEL":     "debug",
		"MAX_RETRIES":   "5",
		"HTTPS_PROXY":   "https://proxy.example.com:8080",
	}
	
	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}
	
	tests := []struct {
		fieldName string
		expected  string
	}{
		{"APIEndpoint", "https://test-api.example.com"},
		{"DatabaseURL", "postgres://test:test@localhost/test"},
		{"LogLevel", "debug"},
		{"MaxRetries", "5"},
		{"HTTPSProxy", "https://proxy.example.com:8080"},
		{"NonExistentField", ""},
	}
	
	for _, test := range tests {
		t.Run(test.fieldName, func(t *testing.T) {
			result := GetEnvDefault(test.fieldName)
			if result != test.expected {
				t.Errorf("GetEnvDefault(%s) = %s, want %s", test.fieldName, result, test.expected)
			}
		})
	}
}

func TestTypeValidator_ValidateRequired(t *testing.T) {
	tv := &TypeValidator{}
	
	tests := []struct {
		name      string
		value     interface{}
		fieldName string
		wantError bool
	}{
		{"nil value", nil, "TestField", true},
		{"empty string", "", "TestField", true},
		{"non-empty string", "value", "TestField", false},
		{"zero int", 0, "TestField", false},
		{"non-zero int", 42, "TestField", false},
		{"false bool", false, "TestField", false},
		{"true bool", true, "TestField", false},
		{"zero float", 0.0, "TestField", false},
		{"non-zero float", 3.14, "TestField", false},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := tv.ValidateRequired(test.value, test.fieldName)
			if (err != nil) != test.wantError {
				t.Errorf("ValidateRequired(%v, %s) error = %v, wantError = %v", 
					test.value, test.fieldName, err, test.wantError)
			}
		})
	}
}

func TestSafeParser_ParseWithFallback(t *testing.T) {
	sp := &SafeParser{}
	
	fallbackCalled := false
	fallback := func() {
		fallbackCalled = true
	}
	
	// Test successful parsing
	t.Run("successful parsing", func(t *testing.T) {
		fallbackCalled = false
		parser := func(s string) error {
			return nil
		}
		
		err := sp.ParseWithFallback("test", parser, fallback)
		if err != nil {
			t.Errorf("ParseWithFallback should not return error for successful parsing")
		}
		if fallbackCalled {
			t.Errorf("Fallback should not be called for successful parsing")
		}
	})
	
	// Test failed parsing with fallback
	t.Run("failed parsing with fallback", func(t *testing.T) {
		fallbackCalled = false
		parser := func(s string) error {
			return fmt.Errorf("parse error")
		}
		
		err := sp.ParseWithFallback("test", parser, fallback)
		if err == nil {
			t.Errorf("ParseWithFallback should return error for failed parsing")
		}
		if !fallbackCalled {
			t.Errorf("Fallback should be called for failed parsing")
		}
	})
}

// Benchmark tests for performance validation
func BenchmarkToUpperSnakeCase(b *testing.B) {
	nc := &NamingConventions{}
	testInputs := []string{
		"APIEndpoint",
		"DatabaseURL", 
		"HTTPSProxy",
		"VeryLongFieldNameWithManyWords",
		"XMLHTTPRequestHandler",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range testInputs {
			nc.ToUpperSnakeCase(input)
		}
	}
}

func BenchmarkToKebabCase(b *testing.B) {
	nc := &NamingConventions{}
	testInputs := []string{
		"APIEndpoint",
		"DatabaseURL",
		"HTTPSProxy", 
		"VeryLongFieldNameWithManyWords",
		"XMLHTTPRequestHandler",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range testInputs {
			nc.ToKebabCase(input)
		}
	}
}

func BenchmarkGetEnvDefault(b *testing.B) {
	// Set up test environment variable
	os.Setenv("TEST_FIELD", "test-value")
	defer os.Unsetenv("TEST_FIELD")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEnvDefault("TestField")
	}
}