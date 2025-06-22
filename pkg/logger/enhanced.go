package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// OutputFormat represents different output formats
type OutputFormat int

const (
	// TextFormat outputs plain text logs (default)
	TextFormat OutputFormat = iota
	// JSONFormat outputs structured JSON logs
	JSONFormat
	// KVFormat outputs key-value structured logs
	KVFormat
)

// CallerInfo represents source code location information
type CallerInfo struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// LogEntry represents a complete log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    *CallerInfo            `json:"caller,omitempty"`
}

// Formatter interface for different output formats
type Formatter interface {
	Format(entry *LogEntry) ([]byte, error)
}

// TextFormatter formats logs as plain text
type TextFormatter struct {
	TimestampFormat string
	IncludeCaller   bool
	CallerFormat    CallerFormat
}

// JSONFormatter formats logs as JSON
type JSONFormatter struct {
	IncludeCaller bool
}

// KVFormatter formats logs as key-value pairs
type KVFormatter struct {
	TimestampFormat string
	IncludeCaller   bool
	CallerFormat    CallerFormat
}

// CallerFormat specifies how caller information is displayed
type CallerFormat int

const (
	// ShortCaller shows only filename:line (default)
	ShortCaller CallerFormat = iota
	// FullCaller shows full path:line
	FullCaller
	// FunctionCaller shows package.function:line
	FunctionCaller
)

// Format implements Formatter for TextFormatter
func (f *TextFormatter) Format(entry *LogEntry) ([]byte, error) {
	var parts []string
	
	// Add level
	parts = append(parts, fmt.Sprintf("[%s]", entry.Level))
	
	// Add timestamp
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}
	parts = append(parts, entry.Timestamp.Format(timestampFormat))
	
	// Add caller info
	if f.IncludeCaller && entry.Caller != nil {
		callerInfo := formatCaller(entry.Caller, f.CallerFormat)
		parts = append(parts, callerInfo)
	}
	
	// Add fields
	if len(entry.Fields) > 0 {
		var fieldParts []string
		for k, v := range entry.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, " ")))
	}
	
	// Add message
	parts = append(parts, entry.Message)
	
	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

// Format implements Formatter for JSONFormatter
func (f *JSONFormatter) Format(entry *LogEntry) ([]byte, error) {
	if !f.IncludeCaller {
		entry.Caller = nil
	}
	
	data, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	
	// Add newline for readability
	result := append(data, '\n')
	return result, nil
}

// Format implements Formatter for KVFormatter
func (f *KVFormatter) Format(entry *LogEntry) ([]byte, error) {
	var parts []string
	
	// Add basic fields
	parts = append(parts, fmt.Sprintf("level=%s", entry.Level))
	
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}
	parts = append(parts, fmt.Sprintf("time=%s", entry.Timestamp.Format(timestampFormat)))
	
	// Add caller info
	if f.IncludeCaller && entry.Caller != nil {
		callerInfo := formatCaller(entry.Caller, f.CallerFormat)
		parts = append(parts, fmt.Sprintf("caller=%s", callerInfo))
	}
	
	// Add custom fields
	for k, v := range entry.Fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	
	// Add message
	parts = append(parts, fmt.Sprintf("msg=%q", entry.Message))
	
	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

// formatCaller formats caller information based on the specified format
func formatCaller(caller *CallerInfo, format CallerFormat) string {
	switch format {
	case FullCaller:
		return fmt.Sprintf("%s:%d", caller.File, caller.Line)
	case FunctionCaller:
		return fmt.Sprintf("%s:%d", caller.Function, caller.Line)
	case ShortCaller:
		fallthrough
	default:
		return fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line)
	}
}

// getCaller retrieves caller information from the runtime stack
func getCaller(skip int) *CallerInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}
	
	fn := runtime.FuncForPC(pc)
	var funcName string
	if fn != nil {
		funcName = fn.Name()
	}
	
	return &CallerInfo{
		File:     file,
		Line:     line,
		Function: funcName,
	}
}

// EnhancedLogger extends DefaultLogger with advanced formatting capabilities
type EnhancedLogger struct {
	*DefaultLogger
	formatter     Formatter
	includeCaller bool
	callerSkip    int
}

// NewEnhanced creates a new enhanced logger with specified formatter
func NewEnhanced(format OutputFormat) *EnhancedLogger {
	base := New()
	enhanced := &EnhancedLogger{
		DefaultLogger: base,
		includeCaller: false,
		callerSkip:    3, // Default skip level for caller detection
	}
	
	enhanced.SetOutputFormat(format)
	return enhanced
}

// SetOutputFormat configures the output format
func (l *EnhancedLogger) SetOutputFormat(format OutputFormat) {
	switch format {
	case JSONFormat:
		l.formatter = &JSONFormatter{IncludeCaller: l.includeCaller}
	case KVFormat:
		l.formatter = &KVFormatter{
			TimestampFormat: l.timeFormat,
			IncludeCaller:   l.includeCaller,
			CallerFormat:    ShortCaller,
		}
	case TextFormat:
		fallthrough
	default:
		l.formatter = &TextFormatter{
			TimestampFormat: l.timeFormat,
			IncludeCaller:   l.includeCaller,
			CallerFormat:    ShortCaller,
		}
	}
}

// SetIncludeCaller enables or disables caller information in logs
func (l *EnhancedLogger) SetIncludeCaller(include bool) {
	l.includeCaller = include
	// Update the current formatter
	switch f := l.formatter.(type) {
	case *TextFormatter:
		f.IncludeCaller = include
	case *JSONFormatter:
		f.IncludeCaller = include
	case *KVFormatter:
		f.IncludeCaller = include
	}
}

// SetCallerFormat sets the format for caller information
func (l *EnhancedLogger) SetCallerFormat(format CallerFormat) {
	switch f := l.formatter.(type) {
	case *TextFormatter:
		f.CallerFormat = format
	case *KVFormatter:
		f.CallerFormat = format
	}
}

// logWithCaller logs a message with caller information
func (l *EnhancedLogger) logWithCaller(level Level, msg string, args ...interface{}) {
	if l.level > level {
		return
	}
	
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   fmt.Sprintf(msg, args...),
		Fields:    make(map[string]interface{}),
	}
	
	// Copy fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	
	// Add caller information if enabled
	if l.includeCaller {
		entry.Caller = getCaller(l.callerSkip)
	}
	
	// Format and output
	formatted, err := l.formatter.Format(entry)
	if err != nil {
		// Fallback to basic formatting if formatter fails
		fallback := fmt.Sprintf("[%s] %s %s\n", level.String(), time.Now().Format(time.RFC3339), entry.Message)
		l.logger.Print(fallback)
		return
	}
	
	l.logger.Print(string(formatted))
}

// Override logging methods to use enhanced formatting
func (l *EnhancedLogger) Debug(msg string, args ...interface{}) {
	l.logWithCaller(DebugLevel, msg, args...)
}

func (l *EnhancedLogger) Info(msg string, args ...interface{}) {
	l.logWithCaller(InfoLevel, msg, args...)
}

func (l *EnhancedLogger) Warn(msg string, args ...interface{}) {
	l.logWithCaller(WarnLevel, msg, args...)
}

func (l *EnhancedLogger) Error(msg string, args ...interface{}) {
	l.logWithCaller(ErrorLevel, msg, args...)
}

func (l *EnhancedLogger) Fatal(msg string, args ...interface{}) {
	l.logWithCaller(FatalLevel, msg, args...)
	l.DefaultLogger.Fatal("") // Trigger exit
}

// WithField returns a new enhanced logger with an additional field
func (l *EnhancedLogger) WithField(key string, value interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	
	enhanced := &EnhancedLogger{
		DefaultLogger: &DefaultLogger{
			level:      l.level,
			logger:     l.logger,
			fields:     newFields,
			timeFormat: l.timeFormat,
		},
		formatter:     l.formatter,
		includeCaller: l.includeCaller,
		callerSkip:    l.callerSkip,
	}
	
	return enhanced
}

// WithFields returns a new enhanced logger with additional fields
func (l *EnhancedLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	
	enhanced := &EnhancedLogger{
		DefaultLogger: &DefaultLogger{
			level:      l.level,
			logger:     l.logger,
			fields:     newFields,
			timeFormat: l.timeFormat,
		},
		formatter:     l.formatter,
		includeCaller: l.includeCaller,
		callerSkip:    l.callerSkip,
	}
	
	return enhanced
}

// SetOutput sets the output destination and preserves enhanced functionality
func (l *EnhancedLogger) SetOutput(w io.Writer) {
	l.DefaultLogger.SetOutput(w)
}