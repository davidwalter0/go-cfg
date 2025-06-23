package flags

import (
	"testing"
)

func TestFlatNestStrategy_FlagName(t *testing.T) {
	strategy := NewFlatNestStrategy()

	tests := []struct {
		name      string
		fieldPath []string
		expected  string
	}{
		{
			name:      "single field",
			fieldPath: []string{"ServiceName"},
			expected:  "service-name",
		},
		{
			name:      "nested field",
			fieldPath: []string{"API", "Host"},
			expected:  "api-host",
		},
		{
			name:      "deeply nested",
			fieldPath: []string{"Database", "ConnectionPool", "MaxConnections"},
			expected:  "database-connection-pool-max-connections",
		},
		{
			name:      "already kebab case",
			fieldPath: []string{"api-host"},
			expected:  "api-host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.FlagName(tt.fieldPath)
			if result != tt.expected {
				t.Errorf("FlagName(%v) = %s, expected %s", tt.fieldPath, result, tt.expected)
			}
		})
	}
}

func TestFlatNestStrategy_EnvVarName(t *testing.T) {
	strategy := NewFlatNestStrategy()

	tests := []struct {
		name      string
		fieldPath []string
		expected  string
	}{
		{
			name:      "single field",
			fieldPath: []string{"ServiceName"},
			expected:  "SERVICE_NAME",
		},
		{
			name:      "nested field",
			fieldPath: []string{"API", "Host"},
			expected:  "API_HOST",
		},
		{
			name:      "deeply nested",
			fieldPath: []string{"Database", "ConnectionPool", "MaxConnections"},
			expected:  "DATABASE_CONNECTION_POOL_MAX_CONNECTIONS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.EnvVarName(tt.fieldPath)
			if result != tt.expected {
				t.Errorf("EnvVarName(%v) = %s, expected %s", tt.fieldPath, result, tt.expected)
			}
		})
	}
}

func TestNestedStrategy_FlagName(t *testing.T) {
	strategy := NewNestedStrategy()

	tests := []struct {
		name      string
		fieldPath []string
		expected  string
	}{
		{
			name:      "single field",
			fieldPath: []string{"ServiceName"},
			expected:  "service-name",
		},
		{
			name:      "nested field",
			fieldPath: []string{"API", "Host"},
			expected:  "api.host",
		},
		{
			name:      "deeply nested",
			fieldPath: []string{"Database", "ConnectionPool", "MaxConnections"},
			expected:  "database.connection-pool.max-connections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.FlagName(tt.fieldPath)
			if result != tt.expected {
				t.Errorf("FlagName(%v) = %s, expected %s", tt.fieldPath, result, tt.expected)
			}
		})
	}
}

func TestPrefixedStrategy_FlagName(t *testing.T) {
	strategy := NewPrefixedStrategy()

	tests := []struct {
		name      string
		fieldPath []string
		expected  string
	}{
		{
			name:      "single field",
			fieldPath: []string{"ServiceName"},
			expected:  "config-service-name",
		},
		{
			name:      "nested field",
			fieldPath: []string{"API", "Host"},
			expected:  "config-api-host",
		},
		{
			name:      "deeply nested",
			fieldPath: []string{"Database", "ConnectionPool", "MaxConnections"},
			expected:  "config-database-connection-pool-max-connections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.FlagName(tt.fieldPath)
			if result != tt.expected {
				t.Errorf("FlagName(%v) = %s, expected %s", tt.fieldPath, result, tt.expected)
			}
		})
	}
}

func TestPrefixedStrategy_EnvVarName(t *testing.T) {
	strategy := NewPrefixedStrategy()

	tests := []struct {
		name      string
		fieldPath []string
		expected  string
	}{
		{
			name:      "single field",
			fieldPath: []string{"ServiceName"},
			expected:  "CONFIG_SERVICE_NAME",
		},
		{
			name:      "nested field",
			fieldPath: []string{"API", "Host"},
			expected:  "CONFIG_API_HOST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.EnvVarName(tt.fieldPath)
			if result != tt.expected {
				t.Errorf("EnvVarName(%v) = %s, expected %s", tt.fieldPath, result, tt.expected)
			}
		})
	}
}

func TestCamelToKebab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ServiceName", "service-name"},
		{"API", "api"},
		{"HTTPPort", "http-port"},
		{"MaxConnectionPoolSize", "max-connection-pool-size"},
		{"alreadyLowerCase", "already-lower-case"},
		{"A", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := camelToKebab(tt.input)
			if result != tt.expected {
				t.Errorf("camelToKebab(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCamelToUpperSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ServiceName", "SERVICE_NAME"},
		{"API", "API"},
		{"HTTPPort", "HTTP_PORT"},
		{"MaxConnectionPoolSize", "MAX_CONNECTION_POOL_SIZE"},
		{"alreadyLowerCase", "ALREADY_LOWER_CASE"},
		{"A", "A"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := camelToUpperSnake(tt.input)
			if result != tt.expected {
				t.Errorf("camelToUpperSnake(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAllStrategies_HelpText(t *testing.T) {
	field := FieldInfo{
		Name:      "TestField",
		HelpText:  "This is a test field",
		Required:  true,
		FieldPath: []string{"Test", "Field"},
	}

	strategies := []NamingStrategy{
		NewFlatNestStrategy(),
		NewNestedStrategy(),
		NewPrefixedStrategy(),
	}

	for _, strategy := range strategies {
		helpText := strategy.HelpText(field)
		if helpText != "This is a test field" {
			t.Errorf("Expected help text 'This is a test field', got '%s'", helpText)
		}
	}

	// Test with empty help text
	emptyField := FieldInfo{
		Name:      "TestField",
		HelpText:  "",
		FieldPath: []string{"Test", "Field"},
	}

	for _, strategy := range strategies {
		helpText := strategy.HelpText(emptyField)
		if helpText != "TestField" {
			t.Errorf("Expected help text 'TestField', got '%s'", helpText)
		}
	}
}
