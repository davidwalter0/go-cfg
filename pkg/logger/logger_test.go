package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		hasError bool
	}{
		{"DEBUG", DebugLevel, false},
		{"debug", DebugLevel, false},
		{"INFO", InfoLevel, false},
		{"info", InfoLevel, false},
		{"WARN", WarnLevel, false},
		{"WARNING", WarnLevel, false},
		{"warn", WarnLevel, false},
		{"ERROR", ErrorLevel, false},
		{"error", ErrorLevel, false},
		{"FATAL", FatalLevel, false},
		{"fatal", FatalLevel, false},
		{"INVALID", InfoLevel, true},
		{"", InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseLevel(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("ParseLevel(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseLevel(%s) unexpected error: %v", tt.input, err)
				}
				if got != tt.expected {
					t.Errorf("ParseLevel(%s) = %v, want %v", tt.input, got, tt.expected)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Fatal("New() returned nil")
	}
	if logger.level != InfoLevel {
		t.Errorf("New() default level = %v, want %v", logger.level, InfoLevel)
	}
	if logger.fields == nil {
		t.Error("New() fields map is nil")
	}
	if logger.logger == nil {
		t.Error("New() internal logger is nil")
	}
}

func TestNewWithLevel(t *testing.T) {
	logger := NewWithLevel(DebugLevel)
	if logger == nil {
		t.Fatal("NewWithLevel() returned nil")
	}
	if logger.level != DebugLevel {
		t.Errorf("NewWithLevel(DebugLevel) level = %v, want %v", logger.level, DebugLevel)
	}
}

func TestSetLevel(t *testing.T) {
	logger := New()
	logger.SetLevel(ErrorLevel)
	if logger.level != ErrorLevel {
		t.Errorf("SetLevel(ErrorLevel) level = %v, want %v", logger.level, ErrorLevel)
	}
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)
	logger.SetLevel(WarnLevel)

	// Test that debug and info messages are not logged when level is WARN
	logger.Debug("debug message")
	logger.Info("info message")

	if buf.Len() > 0 {
		t.Errorf("Expected no output for DEBUG/INFO when level is WARN, got: %s", buf.String())
	}

	// Test that warn, error messages are logged
	logger.Warn("warn message")
	output := buf.String()
	if !strings.Contains(output, "WARN") || !strings.Contains(output, "warn message") {
		t.Errorf("Expected WARN message in output, got: %s", output)
	}

	buf.Reset()
	logger.Error("error message")
	output = buf.String()
	if !strings.Contains(output, "ERROR") || !strings.Contains(output, "error message") {
		t.Errorf("Expected ERROR message in output, got: %s", output)
	}
}

func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)

	fieldLogger := logger.WithField("key", "value")
	fieldLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected field 'key=value' in output, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message' in output, got: %s", output)
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)

	fields := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}
	fieldLogger := logger.WithFields(fields)
	fieldLogger.Info("user action")

	output := buf.String()
	if !strings.Contains(output, "user_id=123") {
		t.Errorf("Expected field 'user_id=123' in output, got: %s", output)
	}
	if !strings.Contains(output, "action=login") {
		t.Errorf("Expected field 'action=login' in output, got: %s", output)
	}
	if !strings.Contains(output, "user action") {
		t.Errorf("Expected 'user action' in output, got: %s", output)
	}
}

func TestFormatMessage(t *testing.T) {
	logger := New()
	logger.SetTimeFormat("2006-01-02T15:04:05Z")

	// Test basic message formatting
	msg := logger.formatMessage(InfoLevel, "test message")
	if !strings.Contains(msg, "[INFO]") {
		t.Errorf("Expected '[INFO]' in formatted message, got: %s", msg)
	}
	if !strings.Contains(msg, "test message") {
		t.Errorf("Expected 'test message' in formatted message, got: %s", msg)
	}

	// Test message formatting with arguments
	msg = logger.formatMessage(ErrorLevel, "error: %v", "file not found")
	if !strings.Contains(msg, "[ERROR]") {
		t.Errorf("Expected '[ERROR]' in formatted message, got: %s", msg)
	}
	if !strings.Contains(msg, "error: file not found") {
		t.Errorf("Expected 'error: file not found' in formatted message, got: %s", msg)
	}
}

func TestGlobalFunctions(t *testing.T) {
	var buf bytes.Buffer
	SetGlobalOutput(&buf)
	SetGlobalLevel(InfoLevel)

	Info("global info message")
	output := buf.String()
	if !strings.Contains(output, "INFO") || !strings.Contains(output, "global info message") {
		t.Errorf("Expected global INFO message in output, got: %s", output)
	}

	buf.Reset()
	globalFieldLogger := WithField("global", "test")
	globalFieldLogger.Info("global field message")
	output = buf.String()
	if !strings.Contains(output, "global=test") {
		t.Errorf("Expected global field in output, got: %s", output)
	}
}

func TestChainedFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)

	// Test chaining WithField calls
	chainedLogger := logger.WithField("first", "value1").WithField("second", "value2")
	chainedLogger.Info("chained message")

	output := buf.String()
	if !strings.Contains(output, "first=value1") {
		t.Errorf("Expected field 'first=value1' in output, got: %s", output)
	}
	if !strings.Contains(output, "second=value2") {
		t.Errorf("Expected field 'second=value2' in output, got: %s", output)
	}
}

func TestLoggerImmutability(t *testing.T) {
	logger := New()
	originalFields := len(logger.fields)

	// Create a new logger with fields
	newLogger := logger.WithField("test", "value")

	// Original logger should not be modified
	if len(logger.fields) != originalFields {
		t.Errorf("Original logger was modified, expected %d fields, got %d", originalFields, len(logger.fields))
	}

	// New logger should have the additional field
	newLoggerImpl := newLogger.(*DefaultLogger)
	if len(newLoggerImpl.fields) != originalFields+1 {
		t.Errorf("New logger should have %d fields, got %d", originalFields+1, len(newLoggerImpl.fields))
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message %d", i)
	}
}

func BenchmarkLoggerWithField(b *testing.B) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithField("iteration", i).Info("benchmark message")
	}
}

func BenchmarkLoggerWithFields(b *testing.B) {
	var buf bytes.Buffer
	logger := New()
	logger.SetOutput(&buf)

	fields := map[string]interface{}{
		"component": "benchmark",
		"version":   "1.0.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(fields).Info("benchmark message %d", i)
	}
}
