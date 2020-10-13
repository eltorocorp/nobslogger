package logger_test

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/eltorocorp/nobslogger/v2/logger"
)

type fakeWriter struct{}

// Write just replaces the timestamp internally assigned by the LogService
// with a constant value so the tests remain deterministic.
func (fakeWriter) Write(message []byte) (int, error) {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.(\d{5}|\d{6})-\d{2}:\d{2}`)
	msg := re.ReplaceAllString(string(message), "2009-01-20T12:05:00.000000-04:00")
	fmt.Println(msg)
	return len(msg), nil
}

func ExampleInitializeWriter() {
	// Establish a ServiceContext.
	// This records the highest level information about the system being logged.
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}

	// Setup options for the service. In this case, we're instructing the service
	// to wait at least one second for any remaining log entries to flush
	// before exiting.
	serviceOptions := logger.LogServiceOptions{
		CancellationDeadline: 10 * time.Millisecond,
	}

	// Initialize the LogService.
	loggerSvc := logger.InitializeWriterWithOptions(new(fakeWriter), serviceContext, serviceOptions)

	// Get a new logger (LogContext) from the LogService.
	logger := loggerSvc.NewContext("ExampleInitializeWriter_ServiceContext", "running example")

	// Log something.
	logger.Info("Here is some info")

	loggerSvc.Finish()

	// Output: {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"ExampleInitializeWriter_ServiceContext","operation":"running example","level":"300","severity":"info","msg":"Here is some info","details":""}
}

// LogService supports having multiple logging contexts that may be initialized
// from separate goroutines.
func ExampleLogService_multipleContexts() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}

	serviceOptions := logger.LogServiceOptions{
		CancellationDeadline: 10 * time.Millisecond,
	}

	loggerSvc := logger.InitializeWriterWithOptions(new(fakeWriter), serviceContext, serviceOptions)

	go func() {
		logger := loggerSvc.NewContext("goroutine 1", "running example")
		logger.Info("Here is some info from goroutine 1")
	}()

	go func() {
		logger := loggerSvc.NewContext("goroutine 2", "running example")
		logger.Info("Here is some info from goroutine 2")
	}()

	loggerSvc.Finish()

	// Unordered Output:
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"goroutine 1","operation":"running example","level":"300","severity":"info","msg":"Here is some info from goroutine 1","details":""}
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"goroutine 2","operation":"running example","level":"300","severity":"info","msg":"Here is some info from goroutine 2","details":""}
}

// LogService supports having one (or more) log contexts span multiple
// goroutines.
func ExampleLogService_contextAcrossGoroutines() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}

	serviceOptions := logger.LogServiceOptions{
		CancellationDeadline: 10 * time.Millisecond,
	}

	loggerSvc := logger.InitializeWriterWithOptions(new(fakeWriter), serviceContext, serviceOptions)

	logger := loggerSvc.NewContext("single context", "used across multiple goroutines")

	logger.Info("Log from goroutine 1")

	go func() {
		logger.Info("Log from goroutine 2")
	}()

	go func() {
		logger.Info("Log from goroutine 3")
	}()

	loggerSvc.Finish()

	// Unordered Output:
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"single context","operation":"used across multiple goroutines","level":"300","severity":"info","msg":"Log from goroutine 1","details":""}
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"single context","operation":"used across multiple goroutines","level":"300","severity":"info","msg":"Log from goroutine 2","details":""}
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"single context","operation":"used across multiple goroutines","level":"300","severity":"info","msg":"Log from goroutine 3","details":""}
}

// LogService contexts support a vartiety of log methods, including, but not
// limited to those shown in this example.
func ExampleLogService_variousContextMethods() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}

	serviceOptions := logger.LogServiceOptions{
		CancellationDeadline: 10 * time.Millisecond,
	}

	loggerSvc := logger.InitializeWriterWithOptions(new(fakeWriter), serviceContext, serviceOptions)

	logger := loggerSvc.NewContext("goroutine 1", "running example")

	logger.Info("An info-level message.")
	logger.InfoD("An info-level message.", "With more details!")

	logger.Debug("A debug-level message")
	logger.DebugD("A debug-level message.", "With extra details!")

	loggerSvc.Finish()

	// Unordered Output:
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"goroutine 1","operation":"running example","level":"300","severity":"info","msg":"An info-level message.","details":""}
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"goroutine 1","operation":"running example","level":"300","severity":"info","msg":"An info-level message.","details":"With more details!"}
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"goroutine 1","operation":"running example","level":"200","severity":"debug","msg":"A debug-level message","details":""}
	// {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"goroutine 1","operation":"running example","level":"200","severity":"debug","msg":"A debug-level message.","details":"With extra details!"}

}

// LogContexts also support the io.Writer interface, so they can be used to
// hook into any external writer, such as through the use of `log.SetOutput`.
func ExampleLogContext_Write() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}

	serviceOptions := logger.LogServiceOptions{
		CancellationDeadline: 10 * time.Millisecond,
	}

	loggerSvc := logger.InitializeWriterWithOptions(new(fakeWriter), serviceContext, serviceOptions)

	logger := loggerSvc.NewContext("ExampleLogContext", "Write")

	// Here we simulate hooking the LogContext into the an existing std/logger.
	// In this example we create a new logger via `log.New`, the this could also
	// be done using `log.SetOutput`.
	stdlibLogger := log.New(logger, "", 0)

	// When we call Println, the current LogContext will log the message at the
	// Trace level.
	// Note that the expected output will include an escaped newline character.
	// This is added by the Println function, and is properly escaped by
	// nobslogger to prevent mangling the JSON output.
	stdlibLogger.Println("Hello from the standard library logger!")

	loggerSvc.Finish()

	// Output: {"timestamp":"2009-01-20T12:05:00.000000-04:00","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"ExampleLogContext","operation":"Write","level":"100","severity":"trace","msg":"Hello from the standard library logger!\n","details":""}
}
