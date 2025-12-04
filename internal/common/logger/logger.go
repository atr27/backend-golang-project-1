package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog
type Logger struct {
	logger zerolog.Logger
}

// Config holds logger configuration
type Config struct {
	Level  string
	Format string
}

// New creates a new logger instance
func New(config Config) *Logger {
	// Set log level
	level := parseLevel(config.Level)
	zerolog.SetGlobalLevel(level)

	// Configure output
	var output io.Writer
	if config.Format == "json" {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Create logger with context
	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{logger: logger}
}

// parseLevel parses log level string
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}

// ErrorWithErr logs an error with error object
func (l *Logger) ErrorWithErr(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Msgf(format, v...)
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &Logger{logger: newLogger}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := l.logger.With().Fields(fields).Logger()
	return &Logger{logger: newLogger}
}

// WithError adds an error to the logger
func (l *Logger) WithError(err error) *Logger {
	newLogger := l.logger.With().Err(err).Logger()
	return &Logger{logger: newLogger}
}

// Global logger functions for convenience
var defaultLogger *Logger

// Init initializes the default logger
func Init(config Config) {
	defaultLogger = New(config)
	log.Logger = defaultLogger.logger
}

// Debug logs a debug message using default logger
func Debug(msg string) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg)
	}
}

// Debugf logs a formatted debug message using default logger
func Debugf(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debugf(format, v...)
	}
}

// Info logs an info message using default logger
func Info(msg string) {
	if defaultLogger != nil {
		defaultLogger.Info(msg)
	}
}

// Infof logs a formatted info message using default logger
func Infof(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Infof(format, v...)
	}
}

// Warn logs a warning message using default logger
func Warn(msg string) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg)
	}
}

// Warnf logs a formatted warning message using default logger
func Warnf(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warnf(format, v...)
	}
}

// Error logs an error message using default logger
func Error(msg string) {
	if defaultLogger != nil {
		defaultLogger.Error(msg)
	}
}

// Errorf logs a formatted error message using default logger
func Errorf(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Errorf(format, v...)
	}
}

// ErrorWithErr logs an error with error object using default logger
func ErrorWithErr(err error, msg string) {
	if defaultLogger != nil {
		defaultLogger.ErrorWithErr(err, msg)
	}
}

// Fatal logs a fatal message and exits using default logger
func Fatal(msg string) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg)
	}
}

// Fatalf logs a formatted fatal message and exits using default logger
func Fatalf(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatalf(format, v...)
	}
}

// WithField adds a field to the default logger
func WithField(key string, value interface{}) *Logger {
	if defaultLogger != nil {
		return defaultLogger.WithField(key, value)
	}
	return nil
}

// WithFields adds multiple fields to the default logger
func WithFields(fields map[string]interface{}) *Logger {
	if defaultLogger != nil {
		return defaultLogger.WithFields(fields)
	}
	return nil
}

// WithError adds an error to the default logger
func WithError(err error) *Logger {
	if defaultLogger != nil {
		return defaultLogger.WithError(err)
	}
	return nil
}
