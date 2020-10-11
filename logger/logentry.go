package logger

// A LogEntry defines the highest level structured log entry.
type LogEntry struct {
	ServiceContext
	LogContext
	LogDetail
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
