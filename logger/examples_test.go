package logger_test

import (
	"fmt"
	"log"
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

func ExampleInitializeWriter() {
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

	// Output: {"timestamp":"1234567890123456789","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"ExampleInitializeWriter_ServiceContext","operation":"running example","level":"300","severity":"info","msg":"Here is some info","details":""}
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

// LogService supports having one (or more) log contexts span multiple
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
// limited to those shown in this example.
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

// LogContexts also support the io.Writer interface, so they can be used to
// hook into any external writer, such as through the use of `log.SetOutput`.
func ExampleLogContext_Write() {
	serviceContext := logger.ServiceContext{
		Environment:       "test",
		SystemName:        "examples",
		ServiceName:       "example runner",
		ServiceInstanceID: "1",
	}
	loggerSvc := logger.InitializeWriter(new(fakeWriter), serviceContext)
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

	loggerSvc.Cancel()
	loggerSvc.Wait()

	// Output: {"timestamp":"1234567890123456789","environment":"test","system_name":"examples","service_name":"example runner","service_instance_id":"1","site":"ExampleLogContext","operation":"Write","level":"100","severity":"trace","msg":"Hello from the standard library logger!\n","details":""}
}
