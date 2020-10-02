package nobslogger_test

import (
	"testing"

	"github.com/eltorocorp/nobslogger/pkg/nobslogger"
)

func BenchmarkEntrySerialize(b *testing.B) {
	entry := nobslogger.LogEntry{
		LogContext: nobslogger.LogContext{
			Environment:       "0123456789",
			SystemName:        "0123456789",
			ServiceInstanceID: "0123456789",
			ServiceName:       "0123456789",
		},
		LogDetail: nobslogger.LogDetail{
			Level:   nobslogger.LogLevelTrace,
			Message: "0123456789",
			Details: "0123456789",
		},
	}
	for n := 0; n < b.N; n++ {
		entry.Serialize()
	}
}
