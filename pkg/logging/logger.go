package logging

import (
	"github.com/rs/zerolog"
	"os"
)

// LoggerOption is a function that can be passed to NewLogger
// to customize logger options. It should accept zerolog Logger,
// modify it and return modified instance of zerolog Logger
type LoggerOption func(logger zerolog.Logger) zerolog.Logger

// WithLogLevel accepts zerolog.Level and returns LoggerOption that
// set zerolog logging level to passed level
func WithLogLevel(level zerolog.Level) LoggerOption {
	return func(logger zerolog.Logger) zerolog.Logger {
		return logger.Level(level)
	}
}

// NewLogger creates a new zerolog logger that writes to os.Stdout
// via zerolog.ConsoleWriter and has logging level set to zerolog.TraceLevel.
// It also adds timestamps to logger entries.
//
// It accepts a slice of LoggerOption and applies it to built logger.
func NewLogger(opts ...LoggerOption) zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	for _, opt := range opts {
		logger = opt(logger)
	}

	return logger
}
