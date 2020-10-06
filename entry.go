package nobslogger

// A LogEntry defines the highest level structured log entry.
type LogEntry struct {
	ServiceContext
	LogContext
	LogDetail
}

const (
	braceOpenToken         = "{"
	braceCloseToken        = "}"
	fieldOpenToken         = ":\""
	fieldCloseToken        = "\","
	timestampToken         = "\"timestamp\""
	environmentToken       = "\"environment\""
	systemNameToken        = "\"system_name\""
	serviceNameToken       = "\"service_name\""
	serviceInstanceIDToken = "\"service_instance_id\""
	levelToken             = "\"level\""
	severityToken          = "\"severity\""
	siteToken              = "\"site\""
	operationToken         = "\"operation\""
	messageToken           = "\"message\""
	detailsToken           = "\"details\""
)

// Serialize marshals the LogEntry into a JSON format.
// This method constructs the JSON response manually just for the sake of being
// no bullshit and really fast. This is less cute than using higher abstractions
// but is also ~140 times faster than `json.MarshalIndent`, ~30 times faster
// than `json.Marhsal`, and ~20 times faster than `fmt.Sprintf`.
func (le LogEntry) Serialize() []byte {
	return []byte(braceOpenToken +
		timestampToken + fieldOpenToken + le.Timestamp + fieldCloseToken +
		environmentToken + fieldOpenToken + le.Environment + fieldCloseToken +
		systemNameToken + fieldOpenToken + le.SystemName + fieldCloseToken +
		serviceNameToken + fieldOpenToken + le.ServiceName + fieldCloseToken +
		serviceInstanceIDToken + fieldOpenToken + le.ServiceInstanceID + fieldCloseToken +
		levelToken + fieldOpenToken + string(le.Level) + fieldCloseToken +
		severityToken + fieldOpenToken + string(le.Severity) + fieldCloseToken +
		siteToken + fieldOpenToken + le.Site + fieldCloseToken +
		operationToken + fieldOpenToken + le.Operation + fieldCloseToken +
		messageToken + fieldOpenToken + le.Message + fieldCloseToken +
		detailsToken + fieldOpenToken + le.Details + fieldCloseToken +
		braceCloseToken)
}
