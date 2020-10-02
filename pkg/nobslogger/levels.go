package nobslogger

// LogLevel defines the general severity of a log entry.
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
