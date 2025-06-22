// Package logger provides a structured logging interface for the go-cfg library.
//
// The logger package implements a leveled logging system with support for
// structured fields and configurable output. It's designed to be simple,
// fast, and compatible with the autocfg configuration system.
//
// Basic Usage:
//
//	import "github.com/davidwalter0/go-cfg/pkg/logger"
//
//	// Use global logger functions
//	logger.Info("Application started")
//	logger.Error("Failed to connect: %v", err)
//
//	// Create a new logger instance
//	log := logger.New()
//	log.SetLevel(logger.DebugLevel)
//	log.Debug("Debug message")
//
//	// Use structured logging with fields
//	log.WithField("user_id", 123).Info("User logged in")
//	log.WithFields(map[string]interface{}{
//		"method": "POST",
//		"path":   "/api/users",
//		"status": 201,
//	}).Info("HTTP request completed")
//
// Log Levels:
//
// The logger supports five log levels in order of increasing severity:
//   - DEBUG: Detailed information for diagnosing problems
//   - INFO: General information about program execution
//   - WARN: Warning messages for potentially harmful situations
//   - ERROR: Error events that allow the application to continue
//   - FATAL: Critical errors that cause the application to exit
//
// Configuration Integration:
//
// The logger is designed to integrate seamlessly with go-cfg's configuration
// system. You can configure the log level through environment variables,
// command-line flags, or configuration files:
//
//	type Config struct {
//		LogLevel string `default:"info" json:"log_level"`
//	}
//
//	// Parse log level from configuration
//	level, err := logger.ParseLevel(config.LogLevel)
//	if err != nil {
//		log.Fatal("Invalid log level: %v", err)
//	}
//	logger.SetGlobalLevel(level)
//
// Thread Safety:
//
// The logger is safe for concurrent use by multiple goroutines.
// Each logger instance maintains its own state and can be used
// independently across different parts of your application.
package logger