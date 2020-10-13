package logger

// A LogEntry defines the highest level structured log entry.
type LogEntry struct {
	ServiceContext
	LogContext
	LogDetail
}
