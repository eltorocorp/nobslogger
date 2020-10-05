package nobslogger

// A LogDetail defines low level information for a structured log entry.
type LogDetail struct {
	Level     LogLevel
	Severity  LogSeverity
	Timestamp string
	Message   string
	Details   string
}
