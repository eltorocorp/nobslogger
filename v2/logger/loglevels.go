package logger

// LogLevel numerically defines the general severity of a log entry.
type LogLevel string

// LogLevel constants.
// LogLevels are backed by strings rather than ints to avoid conversions to
// string when serializing for UDP (see `LogEntry.Serialize()`)
const (
	LogLevelTrace LogLevel = "100"
	LogLevelDebug LogLevel = "200"
	LogLevelInfo  LogLevel = "300"
	LogLevelWarn  LogLevel = "400"
	LogLevelError LogLevel = "500"
	LogLevelFatal LogLevel = "600"
)

// LogSeverity defines the general severity of a log entry.
type LogSeverity string

// LogSeverity constants
const (
	LogSeverityTrace LogSeverity = "trace"
	LogSeverityDebug LogSeverity = "debug"
	LogSeverityInfo  LogSeverity = "info"
	LogSeverityWarn  LogSeverity = "warn"
	LogSeverityError LogSeverity = "error"
	LogSeverityFatal LogSeverity = "fatal"
)
