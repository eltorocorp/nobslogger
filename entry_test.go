package nobslogger_test

import (
	"testing"

	"github.com/eltorocorp/nobslogger"
	"github.com/stretchr/testify/assert"
)

func TestLogEntrySerialize(t *testing.T) {
	logEntry := nobslogger.LogEntry{
		ServiceContext: nobslogger.ServiceContext{
			Environment:       "env",
			SystemName:        "sys",
			ServiceName:       "srn",
			ServiceInstanceID: "sid",
		},
		LogContext: nobslogger.LogContext{
			Site:      "sit",
			Operation: "opn",
		},
		LogDetail: nobslogger.LogDetail{
			Level:     nobslogger.LogLevelInfo,
			Severity:  nobslogger.LogSeverityInfo,
			Timestamp: "tms",
			Message:   "msg",
			Details:   "dtl",
		},
	}

	actualResult := string(logEntry.Serialize())
	expectedResult := `{"timestamp":"tms","environment":"env","system_name":"sys","service_name":"srn","service_instance_id":"sid","site":"sit","operation":"opn","level":"300","severity":"info","message":"msg","details":"dtl"}`
	assert.Equal(t, expectedResult, actualResult)
}
