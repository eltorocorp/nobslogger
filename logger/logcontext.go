package logger

import (
	"sync/atomic"
)

// LogContext defines high level information for a structured log entry.
// Information in LogContext is applicable to multiple log calls, and describe
// the general environment in which a series of related log calls will be made.
type LogContext struct {
	logService *LogService
	buffer     []byte

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

func (l LogContext) submit(sc *ServiceContext, lc *LogContext, ld LogDetail) {
	// ld.Timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	atomic.AddUint32(&l.logService.waiters, 1)
	for {
		if !atomic.CompareAndSwapInt32(&l.logService.locked, 0, 1) {
			continue
		}

		offset := 0

		copy(lc.buffer[offset:offset+len(braceOpenToken)], braceOpenToken)
		offset += len(braceOpenToken)
		copy(lc.buffer[offset:offset+len(timestampToken)], timestampToken)
		offset += len(timestampToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(ld.Timestamp)], ld.Timestamp)
		offset += len(ld.Timestamp)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(environmentToken)], environmentToken)
		offset += len(environmentToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(sc.Environment)], sc.Environment)
		offset += len(sc.Environment)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(systemNameToken)], systemNameToken)
		offset += len(systemNameToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(sc.SystemName)], sc.SystemName)
		offset += len(sc.SystemName)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(serviceNameToken)], serviceNameToken)
		offset += len(serviceNameToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(sc.ServiceName)], sc.ServiceName)
		offset += len(sc.ServiceName)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(serviceInstanceIDToken)], serviceInstanceIDToken)
		offset += len(serviceInstanceIDToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(sc.ServiceInstanceID)], sc.ServiceInstanceID)
		offset += len(sc.ServiceInstanceID)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(siteToken)], siteToken)
		offset += len(siteToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(lc.Site)], lc.Site)
		offset += len(lc.Site)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(operationToken)], operationToken)
		offset += len(operationToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(lc.Operation)], lc.Operation)
		offset += len(lc.Operation)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(levelToken)], levelToken)
		offset += len(levelToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(string(ld.Level))], string(ld.Level))
		offset += len(string(ld.Level))
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(severityToken)], severityToken)
		offset += len(severityToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(string(ld.Severity))], string(ld.Severity))
		offset += len(string(ld.Severity))
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(messageToken)], messageToken)
		offset += len(messageToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(ld.Message)], ld.Message)
		offset += len(ld.Message)
		copy(lc.buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
		offset += len(fieldCloseToken)
		copy(lc.buffer[offset:offset+len(detailsToken)], detailsToken)
		offset += len(detailsToken)
		copy(lc.buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
		offset += len(fieldOpenToken)
		copy(lc.buffer[offset:offset+len(ld.Details)], ld.Details)
		offset += len(ld.Details)
		copy(lc.buffer[offset:offset+len(finalFieldCloseToken)], finalFieldCloseToken)
		offset += len(finalFieldCloseToken)
		copy(lc.buffer[offset:offset+len(braceCloseToken)], braceCloseToken)
		offset += len(braceCloseToken)

		l.logService.writeEntry(lc.buffer[0:offset])
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