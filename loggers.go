package errors

import (
	"fmt"
	"log"
)

const (
	LogDebug = iota + 10
	LogInfo
	LogWarning
	LogError
)

type loggerOpts struct {
	verbosity           int
	saveStack           bool
	stackTraceFormatter StackTraceFormatter
	errorFormatter      ErrorFormatter
}

type LoggerOption func(*loggerOpts)

func WithVerbosity(verbosity int) LoggerOption {
	return func(o *loggerOpts) {
		o.verbosity = verbosity
	}
}

func WithStackTraceFormatter(formatter StackTraceFormatter) LoggerOption {
	return func(o *loggerOpts) {
		o.stackTraceFormatter = formatter
	}
}

func WithErrorFormatter(formatter ErrorFormatter) LoggerOption {
	return func(o *loggerOpts) {
		o.errorFormatter = formatter
	}
}

func WithSaveStack(saveStack bool) LoggerOption {
	return func(o *loggerOpts) {
		o.saveStack = saveStack
	}
}

func CustomLogger(opts ...LoggerOption) LoggerFunc {
	options := &loggerOpts{
		verbosity:           0,
		stackTraceFormatter: MultilineStackTraceFormatter,
		errorFormatter:      MultilineFormatter,
		saveStack:           false,
	}
	for _, opt := range opts {
		opt(options)
	}
	return func(e EnhancedError) {
		log.Print(options.errorFormatter(e, options.verbosity, options.stackTraceFormatter))
		if options.saveStack {
			if err := Manager().SaveStack(e); err != nil {
				log.Print(fmt.Errorf("error saving stack: %w", err))
			}
		}
	}
}

func DefaultLogger() LoggerFunc {
	return CustomLogger()
}
