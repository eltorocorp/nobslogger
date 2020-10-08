package logger_test

import (
	"fmt"
	"regexp"
	"time"

	"github.com/eltorocorp/nobslogger/logger"
)

type fakeWriter struct{}

// Write just replaces the timestamp internally assigned by the LogService
// with a constant value so the tests remain deterministic.
func (fakeWriter) Write(message []byte) (int, error) {
	re := regexp.MustCompile(`\d{19}`)
	msg := re.ReplaceAllString(string(message), "1234567890123456789")
	fmt.Println(msg)
	return len(msg), nil
}

func ExampleLogService_InitializeWriter() {
	// Establish a ServiceContext.
	// This records the highest level information about the system being logged.
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}

	// Initialize the LogService.
	loggerSvc := logger.InitializeWriter(new(fakeWriter), serviceContext)

	// Get a new logger (LogContext) from the LogService.
	logger := loggerSvc.NewContext("ExampleInitializeWriter_ServiceContext", "running example")

	// Log something.
	logger.Info("Here is some info")

	// Calling Cancel signals to the LogService to begin flushing the internal
	// log queue.
	loggerSvc.Cancel()

	// Wait always blocks while the LogService is active, and will only unblock
	// after the Cancel method has been called and has finished flusing the
	// log message queue.
	loggerSvc.Wait()

	// Output: {"timestamp":"1234567890123456789","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"ExampleInitializeWriter_ServiceContext","operation":"running example","level":"300","severity":"info","message":"Here is some info","details":""}
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
	loggerSvc := logger.InitializeWriter(new(fakeWriter), serviceContext)

	go func() {
		logger := loggerSvc.NewContext("goroutine 1", "running example")
		logger.Info("Here is some info from goroutine 1")
	}()

	go func() {
		logger := loggerSvc.NewContext("goroutine 2", "running example")
		logger.Info("Here is some info from goroutine 2")
	}()

	loggerSvc.Cancel()
	loggerSvc.Wait()
}

// LogService supports having a one (or more) log contexts span multiple
// goroutines.
func ExampleLogService_contextAcrossGoroutines() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}
	loggerSvc := logger.InitializeWriter(new(fakeWriter), serviceContext)
	logger := loggerSvc.NewContext("single context", "used across multiple goroutines")

	logger.Info("Log from goroutine 1")

	go func() {
		logger.Info("Log from goroutine 2")
	}()

	go func() {
		logger.Info("Log from goroutine 3")
	}()

	loggerSvc.Cancel()
	loggerSvc.Wait()
}

// LogService supports cancellation from a separate goroutine from where the
// service was originally initialized.
func ExampleLogService_cancelFromSeparateGoroutine() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}
	loggerSvc := logger.InitializeWriter(new(fakeWriter), serviceContext)
	logger := loggerSvc.NewContext("single context", "used across multiple goroutines")

	go func() {
		for {
			logger.Warn("infinite loop")
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		logger.Info("Cancelling")
		loggerSvc.Cancel()
	}()

	loggerSvc.Wait()
}

// LogService contexts support a vartiety of log methods, including, but not
// limmitted to those shown in this example.
func ExampleLogService_variousContextMethods() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}
	loggerSvc := logger.InitializeWriter(new(fakeWriter), serviceContext)
	logger := loggerSvc.NewContext("goroutine 1", "running example")

	logger.Info("An info-level message.")
	logger.InfoD("An info-level message.", "With more details!")

	logger.Debug("A debug-level message")
	logger.DebugD("A debug-level message.", "With extra details!")

	loggerSvc.Cancel()
	loggerSvc.Wait()
}
