package logger

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/kpango/fastime"
)

func init() {
	// If changing the output format, be sure to also update the serializer, as
	// it is expecting a 32 byte value.
	fastime.SetFormat(time.RFC3339Nano)
}

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
// The Write method always inspects the inbound message and escapes any JSON
// characters to avoid unintentionally mangling the expected log entry.
func (l LogContext) Write(message []byte) (int, error) {
	l.Trace(escape(string(message)))
	return len(message), nil
}

func (l LogContext) submit(sc *ServiceContext, lc *LogContext, ld LogDetail) {
	atomic.AddUint32(&l.logService.waiters, 1)
	for {
		if atomic.CompareAndSwapInt32(&l.logService.locked, 0, 1) {
			break
		}
		// Need to give other goroutines a chance to execute. Verify
		// benchmarks if altering this call.
		runtime.Gosched()
	}
	offset := serialize(l.buffer, sc, lc, ld)
	l.logService.writeEntry(lc.buffer[0:offset])
	atomic.SwapInt32(&l.logService.locked, 0)
	atomic.AddUint32(&l.logService.waiters, ^uint32(0))
}

func serialize(buffer []byte, sc *ServiceContext, lc *LogContext, ld LogDetail) (offset int) {
	// Avoiding a loop-construct saves a few cycles.
	// Since we're being opinionated and know ahead of time how many fields
	// we're processing, we can just explicitly construct the outbound message
	// token by token.
	// Verify with benchmarks if altering this section.
	offset += copy(buffer[offset:offset+len(braceOpenToken)], braceOpenToken)
	offset += copy(buffer[offset:offset+len(timestampToken)], timestampToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	// init function sets format to rfc3339nano, which is always 35 bytes long.
	offset += copy(buffer[offset:offset+32], fastime.FormattedNow())
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(environmentToken)], environmentToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(sc.Environment)], sc.Environment)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(systemNameToken)], systemNameToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(sc.SystemName)], sc.SystemName)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(serviceNameToken)], serviceNameToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(sc.ServiceName)], sc.ServiceName)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(serviceInstanceIDToken)], serviceInstanceIDToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(sc.ServiceInstanceID)], sc.ServiceInstanceID)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(siteToken)], siteToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(lc.Site)], lc.Site)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(operationToken)], operationToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(lc.Operation)], lc.Operation)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(levelToken)], levelToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(string(ld.Level))], string(ld.Level))
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(severityToken)], severityToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(string(ld.Severity))], string(ld.Severity))
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(messageToken)], messageToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(ld.Message)], ld.Message)
	offset += copy(buffer[offset:offset+len(fieldCloseToken)], fieldCloseToken)
	offset += copy(buffer[offset:offset+len(detailsToken)], detailsToken)
	offset += copy(buffer[offset:offset+len(fieldOpenToken)], fieldOpenToken)
	offset += copy(buffer[offset:offset+len(ld.Details)], ld.Details)
	offset += copy(buffer[offset:offset+len(finalFieldCloseToken)], finalFieldCloseToken)
	offset += copy(buffer[offset:offset+len(braceCloseToken)], braceCloseToken)
	return
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
