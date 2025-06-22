package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestEnhancedLogger_OutputFormats(t *testing.T) {
	tests := []struct {
		name   string
		format OutputFormat
	}{
		{"TextFormat", TextFormat},
		{"JSONFormat", JSONFormat},
		{"KVFormat", KVFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewEnhanced(tt.format)
			logger.SetOutput(&buf)
			
			logger.Info("test message")
			
			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for %s format, got empty string", tt.name)
			}
			
			// Basic validation that the message appears in output
			if !strings.Contains(output, "test message") {
				t.Errorf("Expected 'test message' in output, got: %s", output)
			}
		})
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewEnhanced(JSONFormat)
	logger.SetOutput(&buf)
	
	logger.WithField("key", "value").Info("json test message")
	
	output := buf.String()
	
	// Parse JSON to validate structure
	var entry LogEntry
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &entry)
	if err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}
	
	// Validate JSON structure
	if entry.Level != "INFO" {
		t.Errorf("Expected level 'INFO', got '%s'", entry.Level)
	}
	
	if entry.Message != "json test message" {
		t.Errorf("Expected message 'json test message', got '%s'", entry.Message)
	}
	
	if entry.Fields["key"] != "value" {
		t.Errorf("Expected field 'key'='value', got '%v'", entry.Fields["key"])
	}
}

func TestKVFormatter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewEnhanced(KVFormat)
	logger.SetOutput(&buf)
	
	logger.WithFields(map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}).Info("kv test message")
	
	output := buf.String()
	
	// Validate KV format structure
	if !strings.Contains(output, "level=INFO") {
		t.Errorf("Expected 'level=INFO' in KV output, got: %s", output)
	}
	
	if !strings.Contains(output, "user_id=123") {
		t.Errorf("Expected 'user_id=123' in KV output, got: %s", output)
	}
	
	if !strings.Contains(output, "action=login") {
		t.Errorf("Expected 'action=login' in KV output, got: %s", output)
	}
	
	if !strings.Contains(output, `msg="kv test message"`) {
		t.Errorf("Expected quoted message in KV output, got: %s", output)
	}
}

func TestCallerInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewEnhanced(TextFormat)
	logger.SetOutput(&buf)
	logger.SetIncludeCaller(true)
	
	logger.Info("caller test message")
	
	output := buf.String()
	
	// Should contain filename and line number
	if !strings.Contains(output, "enhanced_test.go:") {
		t.Errorf("Expected caller info with filename, got: %s", output)
	}
}

func TestCallerFormats(t *testing.T) {
	tests := []struct {
		name   string
		format CallerFormat
	}{
		{"ShortCaller", ShortCaller},
		{"FullCaller", FullCaller},
		{"FunctionCaller", FunctionCaller},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewEnhanced(TextFormat)
			logger.SetOutput(&buf)
			logger.SetIncludeCaller(true)
			logger.SetCallerFormat(tt.format)
			
			logger.Info("caller format test")
			
			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for %s caller format, got empty string", tt.name)
			}
			
			// All formats should include line number
			if !strings.Contains(output, ":") {
				t.Errorf("Expected line number in caller info, got: %s", output)
			}
		})
	}
}

func TestJSONFormatterWithCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := NewEnhanced(JSONFormat)
	logger.SetOutput(&buf)
	logger.SetIncludeCaller(true)
	
	logger.Info("json caller test")
	
	output := buf.String()
	
	// Parse JSON to validate caller structure
	var entry LogEntry
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &entry)
	if err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}
	
	// Validate caller information is present
	if entry.Caller == nil {
		t.Error("Expected caller information in JSON output")
	} else {
		if entry.Caller.File == "" {
			t.Error("Expected file information in caller")
		}
		if entry.Caller.Line == 0 {
			t.Error("Expected line number in caller")
		}
		if entry.Caller.Function == "" {
			t.Error("Expected function information in caller")
		}
	}
}

func TestEnhancedLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewEnhanced(JSONFormat)
	logger.SetOutput(&buf)
	
	// Test field chaining
	enhanced := logger.WithField("first", "value1").WithField("second", "value2")
	enhanced.Info("chained fields test")
	
	output := buf.String()
	
	var entry LogEntry
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &entry)
	if err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}
	
	if entry.Fields["first"] != "value1" {
		t.Errorf("Expected field 'first'='value1', got '%v'", entry.Fields["first"])
	}
	
	if entry.Fields["second"] != "value2" {
		t.Errorf("Expected field 'second'='value2', got '%v'", entry.Fields["second"])
	}
}

func TestEnhancedLoggerImmutability(t *testing.T) {
	logger := NewEnhanced(TextFormat)
	originalFieldCount := len(logger.fields)
	
	// Create new logger with fields
	newLogger := logger.WithField("test", "value")
	
	// Original logger should not be modified
	if len(logger.fields) != originalFieldCount {
		t.Errorf("Original logger was modified, expected %d fields, got %d", originalFieldCount, len(logger.fields))
	}
	
	// New logger should have the additional field
	newEnhanced := newLogger.(*EnhancedLogger)
	if len(newEnhanced.fields) != originalFieldCount+1 {
		t.Errorf("New logger should have %d fields, got %d", originalFieldCount+1, len(newEnhanced.fields))
	}
}

func TestEnhancedLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := NewEnhanced(JSONFormat)
	logger.SetOutput(&buf)
	logger.SetLevel(WarnLevel)
	
	// Debug and Info should not be logged
	logger.Debug("debug message")
	logger.Info("info message")
	
	if buf.Len() > 0 {
		t.Errorf("Expected no output for DEBUG/INFO when level is WARN, got: %s", buf.String())
	}
	
	// Warn should be logged
	logger.Warn("warn message")
	output := buf.String()
	
	if !strings.Contains(output, "warn message") {
		t.Errorf("Expected warn message in output, got: %s", output)
	}
}

func TestGetCaller(t *testing.T) {
	caller := getCaller(1) // Get caller of this function
	
	if caller == nil {
		t.Fatal("getCaller returned nil")
	}
	
	if caller.File == "" {
		t.Error("Expected file information in caller")
	}
	
	if caller.Line == 0 {
		t.Error("Expected line number in caller")
	}
	
	if caller.Function == "" {
		t.Error("Expected function information in caller")
	}
	
	// Should contain this test file name
	if !strings.Contains(caller.File, "enhanced_test.go") {
		t.Errorf("Expected caller file to contain 'enhanced_test.go', got: %s", caller.File)
	}
}

func TestFormatCaller(t *testing.T) {
	caller := &CallerInfo{
		File:     "/path/to/file.go",
		Line:     42,
		Function: "package.function",
	}
	
	tests := []struct {
		name     string
		format   CallerFormat
		expected string
	}{
		{"ShortCaller", ShortCaller, "file.go:42"},
		{"FullCaller", FullCaller, "/path/to/file.go:42"},
		{"FunctionCaller", FunctionCaller, "package.function:42"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCaller(caller, tt.format)
			if result != tt.expected {
				t.Errorf("formatCaller(%s) = %s, want %s", tt.name, result, tt.expected)
			}
		})
	}
}

func BenchmarkEnhancedLogger_Text(b *testing.B) {
	var buf bytes.Buffer
	logger := NewEnhanced(TextFormat)
	logger.SetOutput(&buf)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message %d", i)
	}
}

func BenchmarkEnhancedLogger_JSON(b *testing.B) {
	var buf bytes.Buffer
	logger := NewEnhanced(JSONFormat)
	logger.SetOutput(&buf)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message %d", i)
	}
}

func BenchmarkEnhancedLogger_KV(b *testing.B) {
	var buf bytes.Buffer
	logger := NewEnhanced(KVFormat)
	logger.SetOutput(&buf)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message %d", i)
	}
}

func BenchmarkEnhancedLogger_WithCaller(b *testing.B) {
	var buf bytes.Buffer
	logger := NewEnhanced(TextFormat)
	logger.SetOutput(&buf)
	logger.SetIncludeCaller(true)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message %d", i)
	}
}