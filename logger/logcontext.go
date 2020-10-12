package logger

import (
	"sync/atomic"
)

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
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelTrace,
		Severity: LogSeverityTrace,
		Message:  message,
	})
}

// TraceD logs the most granular about system state along with extra
// detail.
func (l *LogContext) TraceD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelTrace,
		Severity: LogSeverityTrace,
		Message:  message,
		Details:  details,
	})
}

// Debug logs fairly graunlar information about system state.
func (l *LogContext) Debug(message string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelDebug,
		Severity: LogSeverityDebug,
		Message:  message,
	})
}

// DebugD logs relatively detailed information about system state along with
// extra detail.
func (l *LogContext) DebugD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelDebug,
		Severity: LogSeverityDebug,
		Message:  message,
		Details:  details,
	})
}

// Info logs general informational messages useful for describing system state.
func (l *LogContext) Info(message string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelInfo,
		Severity: LogSeverityInfo,
		Message:  message,
	})
}

// InfoD logs general informational messages useful for describing system state
// along with extra detail.
func (l *LogContext) InfoD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelInfo,
		Severity: LogSeverityInfo,
		Message:  message,
		Details:  details,
	})
}

// Warn logs information about potentially harmful situations of interest.
func (l *LogContext) Warn(message string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelWarn,
		Severity: LogSeverityWarn,
		Message:  message,
	})
}

// WarnD logs information about potentially harmful situations of interest along
// with extra detail.
func (l *LogContext) WarnD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelWarn,
		Severity: LogSeverityWarn,
		Message:  message,
		Details:  details,
	})
}

// Error logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running.
func (l *LogContext) Error(message string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelError,
		Severity: LogSeverityError,
		Message:  message,
	})
}

// ErrorD logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running along
// with extra detail.
func (l *LogContext) ErrorD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelError,
		Severity: LogSeverityError,
		Message:  message,
		Details:  details,
	})
}

// Fatal logs the most severe events. Fatal events are likely to have caused
// a service to terminate.
func (l *LogContext) Fatal(message string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelFatal,
		Severity: LogSeverityFatal,
		Message:  message,
	})
}

// FatalD logs the most severe events. Fatal events are likely to have caused
// a service to terminate.
func (l *LogContext) FatalD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelFatal,
		Severity: LogSeverityFatal,
		Message:  message,
		Details:  details,
	})
}

// Write enables this context to be used as an io.Writer.
// Messages sent via the Write method are interpretted as Trace level events.
// The content of messages is not parsed, and is merely forwarded to the
// LogContext.Trace method as a blob.
func (l LogContext) Write(message []byte) (int, error) {
	l.Trace(string(message))
	return len(message), nil
}

var bb []byte = make([]byte, 60000)

func (l LogContext) submit(sc *ServiceContext, lc *LogContext, ld LogDetail) {
	// ld.Timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	atomic.AddUint32(&l.logService.waiters, 1)
	for {
		if !atomic.CompareAndSwapInt32(&l.logService.locked, 0, 1) {
			continue
		}
		// l.logService.writeEntry([]byte(braceOpenToken +
		// 	timestampToken + fieldOpenToken + ld.Timestamp + fieldCloseToken +
		// 	environmentToken + fieldOpenToken + sc.Environment + fieldCloseToken +
		// 	systemNameToken + fieldOpenToken + sc.SystemName + fieldCloseToken +
		// 	serviceNameToken + fieldOpenToken + sc.ServiceName + fieldCloseToken +
		// 	serviceInstanceIDToken + fieldOpenToken + sc.ServiceInstanceID + fieldCloseToken +
		// 	siteToken + fieldOpenToken + lc.Site + fieldCloseToken +
		// 	operationToken + fieldOpenToken + lc.Operation + fieldCloseToken +
		// 	levelToken + fieldOpenToken + string(ld.Level) + fieldCloseToken +
		// 	severityToken + fieldOpenToken + string(ld.Severity) + fieldCloseToken +
		// 	messageToken + fieldOpenToken + ld.Message + fieldCloseToken +
		// 	detailsToken + fieldOpenToken + ld.Details + finalFieldCloseToken +
		// 	braceCloseToken))
		offset := 0

		copy(bb[offset:offset+len(sc.Environment)], sc.Environment)
		offset += len(sc.Environment)

		l.logService.writeEntry(bb[0:offset])
		atomic.SwapInt32(&l.logService.locked, 0)
		atomic.AddUint32(&l.logService.waiters, ^uint32(0))
		break
	}
}

const (
	braceOpenToken         = "{"
	braceCloseToken        = "}"
	fieldOpenToken         = ":\""
	fieldCloseToken        = "\","
	finalFieldCloseToken   = "\""
	timestampToken         = "\"timestamp\""
	environmentToken       = "\"environment\""
	systemNameToken        = "\"system_name\""
	serviceNameToken       = "\"service_name\""
	serviceInstanceIDToken = "\"service_instance_id\""
	levelToken             = "\"level\""
	severityToken          = "\"severity\""
	siteToken              = "\"site\""
	operationToken         = "\"operation\""
	messageToken           = "\"msg\""
	detailsToken           = "\"details\""
)

// // Serialize marshals the LogEntry into a JSON format.
// // This method constructs the JSON response manually just for the sake of being
// // no bullshit and really fast. This is less cute than using higher abstractions
// // but is also ~140 times faster than `json.MarshalIndent`, ~30 times faster
// // than `json.Marhsal`, and ~20 times faster than `fmt.Sprintf`.
// func (le LogEntry) Serialize() []byte {
// 	// JSON escapement: See TestLogServiceEscapesJSON
// 	return []byte(braceOpenToken +
// 		timestampToken + fieldOpenToken + le.Timestamp + fieldCloseToken +
// 		environmentToken + fieldOpenToken + le.Environment + fieldCloseToken +
// 		systemNameToken + fieldOpenToken + le.SystemName + fieldCloseToken +
// 		serviceNameToken + fieldOpenToken + le.ServiceName + fieldCloseToken +
// 		serviceInstanceIDToken + fieldOpenToken + le.ServiceInstanceID + fieldCloseToken +
// 		siteToken + fieldOpenToken + le.Site + fieldCloseToken +
// 		operationToken + fieldOpenToken + le.Operation + fieldCloseToken +
// 		levelToken + fieldOpenToken + string(le.Level) + fieldCloseToken +
// 		severityToken + fieldOpenToken + string(le.Severity) + fieldCloseToken +
// 		messageToken + fieldOpenToken + escape(le.Message) + fieldCloseToken +
// 		detailsToken + fieldOpenToken + escape(le.Details) + finalFieldCloseToken +
// 		braceCloseToken)
// }
