package nobslogger_test

import (
	"fmt"
	"regexp"

	"github.com/eltorocorp/nobslogger"
)

type fakeWriter struct{}

func (fakeWriter) Write(message []byte) (int, error) {
	re := regexp.MustCompile(`\d{19}`)
	msg := re.ReplaceAllString(string(message), "1234567890123456789")
	fmt.Println(msg)
	return len(msg), nil
}

func ExampleInitializeWriter_ServiceContext() {
	serviceContext := nobslogger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}
	loggerSvc := nobslogger.InitializeWriter(new(fakeWriter), serviceContext)
	logger := loggerSvc.NewContext("ExampleInitializeWriter_ServiceContext", "running example")
	logger.Info("Here is some info")
	loggerSvc.Cancel()
	loggerSvc.Wait()
	// Output: {"timestamp":"1234567890123456789","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"ExampleInitializeWriter_ServiceContext","operation":"running example","level":"300","severity":"info","message":"Here is some info","details":""}
}
