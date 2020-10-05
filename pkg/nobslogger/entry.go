package nobslogger

// A LogEntry defines the highest level structured log entry.
type LogEntry struct {
	LogContext
	LogDetail
}

// Serialize marshals the LogEntry into a JSON format.
// This method constructs the JSON response manually just for the sake of being
// no bullshit and really fast. This is less cute than using higher abstractions
// but is also ~140 times faster than `json.MarshalIndent`, ~30 times faster
// than `json.Marhsal`, and ~20 times faster than `fmt.Sprintf`.
func (le *LogEntry) Serialize() []byte {
	return []byte("{" +
		"\"timestamp\":\"" + le.Timestamp + "\"," +
		"\"environment\":\"" + le.logService.serviceContext.Environment + "\"," +
		"\"system_name\":\"" + le.logService.serviceContext.SystemName + "\"," +
		"\"service_name\":\"" + le.logService.serviceContext.ServiceName + "\"," +
		"\"service_instance_id\":\"" + le.logService.serviceContext.ServiceInstanceID + "\"," +
		"\"level\":\"" + string(le.Level) + "\"," +
		"\"severity\":\"" + string(le.Severity) + "\"," +
		"\"site\":\"" + le.Site + "\"," +
		"\"operation\":\"" + le.Operation + "\"," +
		"\"message\":\"" + le.Message + "\"," +
		"\"details\":\"" + le.Details + "\"" +
		"}")
}
