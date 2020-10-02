package nobslogger

// A LogDetail defines low level information for a structured log entry.
type LogDetail struct {
	Level   LogLevel
	Message string
	Details string
}
