package nobslogger

// LogContext defines high level information for a structured log entry.
// Information in LogContext is applicable to multiple log calls, and describe
// the general environment in which a serious of related log calls will be made.
type LogContext struct {
	logService        *LogService
	Environment       string
	SystemName        string
	ServiceName       string
	ServiceInstanceID string
}

// Trace logs the most detailed infomration about system state.
func (l *LogContext) Trace(message, details string) {
	l.logService.submitAsync(l, &LogDetail{
		Level:   LogLevelTrace,
		Message: message,
		Details: details,
	})
}

// Debug logs relatively detailed information about system state.
func (l *LogContext) Debug(message, details string) {
	l.logService.submitAsync(l, &LogDetail{
		Level:   LogLevelTrace,
		Message: message,
		Details: details,
	})
}

// Info logs general informational messages useful for describing system state.
func (l *LogContext) Info(message, details string) {
	l.logService.submitAsync(l, &LogDetail{
		Level:   LogLevelInfo,
		Message: message,
		Details: details,
	})
}

// Warn logs information about potentially harmful situations of interest.
func (l *LogContext) Warn(message, details string) {
	l.logService.submitAsync(l, &LogDetail{
		Level:   LogLevelWarn,
		Message: message,
		Details: details,
	})
}

// Error logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running.
func (l *LogContext) Error(message, details string) {
	l.logService.submitAsync(l, &LogDetail{
		Level:   LogLevelError,
		Message: message,
		Details: details,
	})
}

// Fatal logs the most severe events. Fatal events are likely to have caused
// a service to terminate.
func (l *LogContext) Fatal(message, details string) {
	l.logService.submitAsync(l, &LogDetail{
		Level:   LogLevelFatal,
		Message: message,
		Details: details,
	})
}