package nobslogger

// A LogDetail defines low level information for a structured log entry.
type LogDetail struct {
	Level     LogLevel
	Timestamp string
	Message   string
	Details   string
}
