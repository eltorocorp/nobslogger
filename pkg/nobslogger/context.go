package nobslogger

// LogContext defines high level information for a structured log entry.
// Information in LogContext is applicable to multiple log calls, and describe
// the general environment in which a series of related log calls will be made.
type LogContext struct {
	logService *LogService

	// Site specifies a general location in a codebase from which a group of
	// log messages may emit.
	Site string

	// Operation specifies the a general operation being conducted within this
	// context. All log messages within this context have a direct relationship
	// to this operation.
	Operation string
}

// Trace logs the most granular information about system state.
func (l *LogContext) Trace(message string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelTrace,
		Severity: LogSeverityTrace,
		Message:  message,
	})
}

// TraceD logs the most granular about system state along with extra
// detail.
func (l *LogContext) TraceD(message, details string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelTrace,
		Severity: LogSeverityTrace,
		Message:  message,
		Details:  details,
	})
}

// Debug logs fairly graunlar information about system state.
func (l *LogContext) Debug(message string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelDebug,
		Severity: LogSeverityDebug,
		Message:  message,
	})
}

// DebugD logs relatively detailed information about system state along with
// extra detail.
func (l *LogContext) DebugD(message, details string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelDebug,
		Severity: LogSeverityDebug,
		Message:  message,
		Details:  details,
	})
}

// Info logs general informational messages useful for describing system state.
func (l *LogContext) Info(message string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelInfo,
		Severity: LogSeverityInfo,
		Message:  message,
	})
}

// InfoD logs general informational messages useful for describing system state
// along with extra detail.
func (l *LogContext) InfoD(message, details string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelInfo,
		Severity: LogSeverityInfo,
		Message:  message,
		Details:  details,
	})
}

// Warn logs information about potentially harmful situations of interest.
func (l *LogContext) Warn(message string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelWarn,
		Severity: LogSeverityWarn,
		Message:  message,
	})
}

// WarnD logs information about potentially harmful situations of interest along
// with extra detail.
func (l *LogContext) WarnD(message, details string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelWarn,
		Severity: LogSeverityWarn,
		Message:  message,
		Details:  details,
	})
}

// Error logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running.
func (l *LogContext) Error(message string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelError,
		Severity: LogSeverityError,
		Message:  message,
	})
}

// ErrorD logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running along
// with extra detail.
func (l *LogContext) ErrorD(message, details string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelError,
		Severity: LogSeverityError,
		Message:  message,
		Details:  details,
	})
}

// Fatal logs the most severe events. Fatal events are likely to have caused
// a service to terminate.
func (l *LogContext) Fatal(message string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelFatal,
		Severity: LogSeverityFatal,
		Message:  message,
	})
}

// FatalD logs the most severe events. Fatal events are likely to have caused
// a service to terminate.
func (l *LogContext) FatalD(message, details string) {
	l.logService.submitAsync(*l.logService.serviceContext, *l, LogDetail{
		Level:    LogLevelFatal,
		Severity: LogSeverityFatal,
		Message:  message,
		Details:  details,
	})
}
