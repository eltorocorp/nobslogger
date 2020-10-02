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
func (lp *LogEntry) Serialize() []byte {
	return []byte("{" +
		"\"environment\":\"" + lp.Environment + "\"," +
		"\"system_name\":\"" + lp.SystemName + "\"," +
		"\"service_name\":\"" + lp.ServiceName + "\"," +
		"\"service_instance_id\":\"" + lp.ServiceInstanceID + "\"," +
		"\"log_level\":\"" + string(lp.Level) + "\"," +
		"\"message\":\"" + lp.Message + "\"," +
		"\"details\":\"" + lp.Details + "\"" +
		"}")
}
