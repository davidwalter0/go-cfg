package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Level represents logging level
type Level int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production
	DebugLevel Level = iota
	// InfoLevel is the default logging priority
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual human review
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1)
	FatalLevel
)

// String returns the string representation of the logging level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel parses a level string into a Level value
func ParseLevel(s string) (Level, error) {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DebugLevel, nil
	case "INFO":
		return InfoLevel, nil
	case "WARN", "WARNING":
		return WarnLevel, nil
	case "ERROR":
		return ErrorLevel, nil
	case "FATAL":
		return FatalLevel, nil
	default:
		return InfoLevel, fmt.Errorf("unknown level: %s", s)
	}
}

// Logger interface defines the logging contract
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	SetLevel(level Level)
	SetOutput(w io.Writer)
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// DefaultLogger is a simple logger implementation
type DefaultLogger struct {
	level      Level
	logger     *log.Logger
	fields     map[string]interface{}
	timeFormat string
}

// New creates a new logger instance
func New() *DefaultLogger {
	return &DefaultLogger{
		level:      InfoLevel,
		logger:     log.New(os.Stdout, "", 0),
		fields:     make(map[string]interface{}),
		timeFormat: time.RFC3339,
	}
}

// NewWithLevel creates a new logger with specified level
func NewWithLevel(level Level) *DefaultLogger {
	l := New()
	l.level = level
	return l
}

// SetLevel sets the logging level
func (l *DefaultLogger) SetLevel(level Level) {
	l.level = level
}

// SetOutput sets the output destination for the logger
func (l *DefaultLogger) SetOutput(w io.Writer) {
	l.logger = log.New(w, "", 0)
}

// SetTimeFormat sets the time format for log entries
func (l *DefaultLogger) SetTimeFormat(format string) {
	l.timeFormat = format
}

// WithField returns a new logger with an additional field
func (l *DefaultLogger) WithField(key string, value interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &DefaultLogger{
		level:      l.level,
		logger:     l.logger,
		fields:     newFields,
		timeFormat: l.timeFormat,
	}
}

// WithFields returns a new logger with additional fields
func (l *DefaultLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &DefaultLogger{
		level:      l.level,
		logger:     l.logger,
		fields:     newFields,
		timeFormat: l.timeFormat,
	}
}

// formatMessage formats the log message with level, timestamp, and fields
func (l *DefaultLogger) formatMessage(level Level, msg string, args ...interface{}) string {
	timestamp := time.Now().Format(l.timeFormat)
	message := fmt.Sprintf(msg, args...)

	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", level.String()))
	parts = append(parts, timestamp)

	// Add fields if any
	if len(l.fields) > 0 {
		var fieldParts []string
		for k, v := range l.fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, " ")))
	}

	parts = append(parts, message)
	return strings.Join(parts, " ")
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.logger.Println(l.formatMessage(DebugLevel, msg, args...))
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.logger.Println(l.formatMessage(InfoLevel, msg, args...))
	}
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.logger.Println(l.formatMessage(WarnLevel, msg, args...))
	}
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.logger.Println(l.formatMessage(ErrorLevel, msg, args...))
	}
}

// Fatal logs a fatal message and exits the program
func (l *DefaultLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Println(l.formatMessage(FatalLevel, msg, args...))
	os.Exit(1)
}

// Global logger instance
var defaultLogger = New()

// SetGlobalLevel sets the global logger level
func SetGlobalLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetGlobalOutput sets the global logger output
func SetGlobalOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// Global logging functions
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	defaultLogger.Fatal(msg, args...)
}

func WithField(key string, value interface{}) Logger {
	return defaultLogger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) Logger {
	return defaultLogger.WithFields(fields)
}
