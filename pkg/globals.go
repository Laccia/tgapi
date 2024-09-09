package pkg

import "time"

const timeFormat = time.Stamp

// LogFormat custom type for logger format.
type LogFormat int64

const (
	// LogFormatDefault default PRANA format
	LogFormatDefault LogFormat = iota
	// LogFormatCEFGPN GPN's configs
	LogFormatCEFGPN
)

type LogFile string

// LogDestination custom type for logger destination.
type LogDestination int64

const (
	// LogDestinationUnknown ..
	LogDestinationUnknown LogDestination = iota
	// LogDestinationConsoleOut ..
	LogDestinationConsoleOut
	// LogDestinationConsoleErr ..
	LogDestinationConsoleErr
)

// LogLevel custom type for logger level
type LogLevel string

const (
	// LogLevelTrace is the value used for the trace level field.
	LogLevelTrace LogLevel = "trace"
	// LogLevelDebug is the value used for the debug level field.
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo is the value used for the info level field.
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn is the value used for the warn level field.
	LogLevelWarn LogLevel = "warn"
	// LogLevelError is the value used for the error level field.
	LogLevelError LogLevel = "error"
	// LogLevelFatal is the value used for the fatal level field.
	LogLevelFatal LogLevel = "fatal"
	// LogLevelPanic is the value used for the panic level field.
	LogLevelPanic LogLevel = "panic"
)
