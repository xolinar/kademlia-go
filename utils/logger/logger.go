package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines a flexible logging interface for different implementations.
type ILogger interface {
	Debug(msg string, args ...any) // Logs a message at DEBUG level
	Info(msg string, args ...any)  // Logs a message at INFO level
	Warn(msg string, args ...any)  // Logs a message at WARN level
	Error(msg string, args ...any) // Logs a message at ERROR level
}

// Logger is an implementation of Logger using the Zap logging library.
type Logger struct {
	log *zap.SugaredLogger
}

// NewZapLogger creates a new ZapLogger with the specified logging level and optional config.
// If config is nil, the default configuration is used.
func NewLogger(level string, config *zap.Config) (*Logger, error) {
	// Use the provided config or default if nil
	if config == nil {
		config = defaultZapConfig(level)
	}

	// Build the logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{log: logger.Sugar()}, nil
}

// defaultZapConfig provides a default zap.Config with JSON encoding.
func defaultZapConfig(level string) *zap.Config {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel // Default to INFO level
	}

	return &zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,  // Use production settings
		Encoding:    "json", // JSON format for structured logs
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",                           // Timestamp key
			LevelKey:       "level",                        // Log level key
			MessageKey:     "msg",                          // Message key
			CallerKey:      "caller",                       // Caller information key
			StacktraceKey:  "stacktrace",                   // Stacktrace key for errors
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // Lowercase level names
			EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 time format
			EncodeCaller:   zapcore.ShortCallerEncoder,     // Short caller format
			EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
			LineEnding:     zapcore.DefaultLineEnding,      // Default line ending
		},
		OutputPaths:      []string{"stdout"}, // Output to standard output
		ErrorOutputPaths: []string{"stderr"}, // Errors to standard error
	}
}

// Debug logs a message at the DEBUG level.
func (l *Logger) Debug(msg string, args ...any) {
	l.log.Debugf(msg, args...)
}

// Info logs a message at the INFO level.
func (l *Logger) Info(msg string, args ...any) {
	l.log.Infof(msg, args...)
}

// Warn logs a message at the WARN level.
func (l *Logger) Warn(msg string, args ...any) {
	l.log.Warnf(msg, args...)
}

// Error logs a message at the ERROR level.
func (l *Logger) Error(msg string, args ...any) {
	l.log.Errorf(msg, args...)
}
