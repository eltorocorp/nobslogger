package logger

import (
	"runtime"
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
	l.TraceD(message, "")
}

// TraceD logs the most granular information about system state along with extra
// detail.
func (l *LogContext) TraceD(message, details string) {
	l.submit(l.logService.serviceContext, l, LogDetail{
		Level:    LogLevelTrace,
		Severity: LogSeverityTrace,
		Message:  message,
		Details:  details,
	})
}

// TraceJ logs the most granular information about system state along with extra
// detail, while also escaping all reserved JSON characters.
func (l *LogContext) TraceJ(message, details string) {
	l.TraceD(escape(message), escape(details))
}

// Debug logs fairly graunlar information about system state.
func (l *LogContext) Debug(message string) {
	l.DebugD(message, "")
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

// DebugJ logs relatively detailed information about system state along with
// extra detail, while also escaping all reserved JSON characters.
func (l *LogContext) DebugJ(message, details string) {
	l.DebugD(escape(message), escape(details))
}

// Info logs general informational messages useful for describing system state.
func (l *LogContext) Info(message string) {
	l.InfoD(message, "")
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

// InfoJ logs general informational messages useful for describing system state
// along with extra detail, while also escaping all reserved JSON characters.
func (l *LogContext) InfoJ(message, details string) {
	l.InfoD(escape(message), escape(details))
}

// Warn logs information about potentially harmful situations of interest.
func (l *LogContext) Warn(message string) {
	l.WarnD(message, "")
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

// WarnJ logs information about potentially harmful situations of interest along
// with extra detail, while also escaping all reserved JSON characters.
func (l *LogContext) WarnJ(message, details string) {
	l.WarnD(escape(message), escape(details))
}

// Error logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running.
func (l *LogContext) Error(message string) {
	l.ErrorD(message, "")
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

// ErrorJ logs events of considerable importance that will prevent normal program
// execution, but might still allow the application to continue running along
// with extra detail, while also escaping all reserved JSON characters.
func (l *LogContext) ErrorJ(message, details string) {
	l.ErrorD(escape(message), escape(details))
}

// Fatal logs the most severe events. Fatal events are likely to have caused
// a service to terminate.
func (l *LogContext) Fatal(message string) {
	l.FatalD(message, "")
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

// FatalJ logs the most severe events. Fatal events are likely to have caused
// a service to terminate, while also escaping all reserved JSON characters.
func (l *LogContext) FatalJ(message, details string) {
	l.FatalD(escape(message), escape(details))
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
