package logger_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/eltorocorp/nobslogger/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogEntrySerialize(t *testing.T) {
	logEntry := logger.LogEntry{
		ServiceContext: logger.ServiceContext{
			Environment:       "env",
			SystemName:        "sys",
			ServiceName:       "srn",
			ServiceInstanceID: "sid",
		},
		LogContext: logger.LogContext{
			Site:      "sit",
			Operation: "opn",
		},
		LogDetail: logger.LogDetail{
			Level:     logger.LogLevelInfo,
			Severity:  logger.LogSeverityInfo,
			Timestamp: "tms",
			Message:   "msg",
			Details:   "dtl",
		},
	}

	actualResult := string(logEntry.Serialize())
	expectedResult := `{"timestamp":"tms","environment":"env","system_name":"sys","service_name":"srn","service_instance_id":"sid","site":"sit","operation":"opn","level":"300","severity":"info","message":"msg","details":"dtl"}`
	assert.Equal(t, expectedResult, actualResult)
}

func TestLogEntrySerializeEscapesJSON(t *testing.T) {
	// runes in some cases are duplicated to also account for a case where
	// two or more runes must be replaced within a string.
	testCases := []string{
		"\b\b",
		"\f\f",
		"\n\n",
		"\r\r",
		"\t\t",
		`"`,
		`\`}

	for _, s := range testCases {
		f := func(t *testing.T) {
			logEntry := logger.LogEntry{
				ServiceContext: logger.ServiceContext{
					Environment:       s,
					SystemName:        s,
					ServiceName:       s,
					ServiceInstanceID: s,
				},
				LogContext: logger.LogContext{
					Site:      s,
					Operation: s,
				},
				LogDetail: logger.LogDetail{
					// timestamp, level, and severity are set internally
					// and don't need to be escaped.
					Message: s,
					Details: s,
				},
			}

			maybeJSON := logEntry.Serialize()
			d := json.NewDecoder(bytes.NewReader(maybeJSON))
			for {
				_, err := d.Token()
				if err != nil && err == io.EOF {
					return
				}
				if err != nil {
					t.Log(string(maybeJSON))
					t.Error(err)
					return
				}
			}
		}
		t.Run("rune:"+s, f)
	}

}
