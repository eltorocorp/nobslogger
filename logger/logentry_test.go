package logger_test

// func TestLogEntrySerialize(t *testing.T) {
// 	logEntry := logger.LogEntry{
// 		ServiceContext: logger.ServiceContext{
// 			Environment:       "env",
// 			SystemName:        "sys",
// 			ServiceName:       "srn",
// 			ServiceInstanceID: "sid",
// 		},
// 		LogContext: logger.LogContext{
// 			Site:      "sit",
// 			Operation: "opn",
// 		},
// 		LogDetail: logger.LogDetail{
// 			Level:     logger.LogLevelInfo,
// 			Severity:  logger.LogSeverityInfo,
// 			Timestamp: "tms",
// 			Message:   "msg",
// 			Details:   "dtl",
// 		},
// 	}

// 	actualResult := string(logEntry.Serialize())
// 	expectedResult := `{"timestamp":"tms","environment":"env","system_name":"sys","service_name":"srn","service_instance_id":"sid","site":"sit","operation":"opn","level":"300","severity":"info","msg":"msg","details":"dtl"}`
// 	assert.Equal(t, expectedResult, actualResult)
// }
