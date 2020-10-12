package logger

import (
	"time"

	"github.com/kpango/fastime"
)

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

func init() {
	// If changing the output format, be sure to also update the serializer, as
	// it is expecting a 32 byte value.
	fastime.SetFormat(time.RFC3339Nano)
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

func escape(s string) string {
	for i := 0; i < len(s); {
		switch s[i] {
		case '\b':
			s = s[:i] + "\\b" + s[i+1:]
			i += 2
		case '\f':
			s = s[:i] + "\\f" + s[i+1:]
			i += 2
		case '\n':
			s = s[:i] + "\\n" + s[i+1:]
			i += 2
		case '\r':
			s = s[:i] + "\\r" + s[i+1:]
			i += 2
		case '\t':
			s = s[:i] + "\\t" + s[i+1:]
			i += 2
		case '"':
			s = s[:i] + `\"` + s[i+1:]
			i += 2
		case '\\':
			s = s[:i] + `\\` + s[i+1:]
			i += 2
		default:
			i++
		}
	}
	return s
}
