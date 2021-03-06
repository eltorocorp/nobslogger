package benchmarks

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/internal/ztest"
	"go.uber.org/zap/zapcore"
)

var (
	_messages = fakeMessages(1000)
)

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func getMessage(iter int) string {
	return _messages[iter%1000]
}

func newZapLogger(lvl zapcore.Level) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)
	return zap.New(zapcore.NewCore(
		enc,
		&ztest.Discarder{
			Syncer: ztest.Syncer{},
		},
		lvl,
	))
}

const (
	field1Name  = "environment"
	field1Value = "test"

	field2Name  = "system_name"
	field2Value = "benchmarker"

	field3Name  = "service_name"
	field3Value = "benchmark"

	field4Name  = "service_instance_id"
	field4Value = "00000-11111-22222-33333-44444-55555"

	field5Name = "timestamp"
	// logger tracks this as a unixnano. Doing similar here for an apples to apples comparison.
	field5Value = "1601932749025031000"

	field6Name  = "site"
	field6Value = "benchmark method"

	field7Name  = "operation"
	field7Value = "conducting benchmark"
)

func fakeFields() []zap.Field {
	return []zap.Field{
		zap.String(field1Name, field1Value),
		zap.String(field2Name, field2Value),
		zap.String(field3Name, field3Value),
		zap.String(field4Name, field4Value),
		zap.String(field5Name, field5Value),
		zap.String(field6Name, field6Value),
		zap.String(field7Name, field7Value),
	}
}

func fakeSugarFields() []interface{} {
	return []interface{}{
		field1Name, field1Value,
		field2Name, field2Value,
		field3Name, field3Value,
		field4Name, field4Value,
		field5Name, field5Value,
		field5Name, field5Value,
		field5Name, field5Value,
	}
}
