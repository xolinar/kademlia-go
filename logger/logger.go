package logger

import (
	"github.com/xolinar/kademlia-go/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILogger defines a flexible logging interface for different implementations.
type ILogger interface {
	Debug(msg string, args ...any) // Logs a message at DEBUG level
	Info(msg string, args ...any)  // Logs a message at INFO level
	Warn(msg string, args ...any)  // Logs a message at WARN level
	Error(msg string, args ...any) // Logs a message at ERROR level
}

// Logger is an implementation of ILogger using the Zap logging library.
type Logger struct {
	log *zap.SugaredLogger
}

// NewLogger creates a new Logger with the specified zap.Config.
// If config is nil, the default configuration is used.
func NewLogger(cfg *config.LoggingConfig, customConfig *zap.Config) (*Logger, error) {
	// If custom zap.Config is provided, use it. Otherwise, use the default config.
	var zapConfig *zap.Config
	if customConfig != nil {
		zapConfig = customConfig
	} else {
		zapConfig = defaultZapConfig(cfg)
	}

	// Build the logger
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{log: logger.Sugar()}, nil
}

// defaultZapConfig provides a default zap.Config with JSON encoding.
// It accepts a LoggingConfig for customization or uses defaults if nil.
func defaultZapConfig(cfg *config.LoggingConfig) *zap.Config {
	// Set default log level and output
	level := "info"
	out := "stdout"

	// Override defaults if LoggingConfig is provided
	if cfg != nil {
		if cfg.LogLevel != "" {
			level = cfg.LogLevel
		}
		if cfg.LogOutput != "" {
			out = cfg.LogOutput
		}
	}

	// Map string log level to zapcore.Level
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
		OutputPaths:      []string{out},      // Output to the specified log output
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
