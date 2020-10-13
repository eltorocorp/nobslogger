package logger_test

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/eltorocorp/nobslogger/mocks/mock_io"
	"github.com/eltorocorp/nobslogger/v2/logger"
	"github.com/golang/mock/gomock"
)

func Test_ServiceInitializeWriterHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)

	loggerService := logger.InitializeWriterWithOptions(writer, logger.ServiceContext{}, logger.LogServiceOptions{})
	logger := loggerService.NewContext("context site", "operation")
	logger.Info("some info")
	loggerService.Finish()
}

func Test_ServiceInitializeWriterPersistentError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().
		Write(gomock.Any()).
		Return(0, fmt.Errorf("test error")).
		Times(2)

	loggerService := logger.InitializeWriterWithOptions(writer, logger.ServiceContext{}, logger.LogServiceOptions{})
	logger := loggerService.NewContext("context site", "operation")
	logger.Info("message")
	loggerService.Finish()
}

// The LogService must support ingestion of logs from LogContexts on different
// goroutines.
func Test_LogServiceSupportsMultipleContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(2)

	serviceContext := logger.ServiceContext{}
	loggerSvc := logger.InitializeWriterWithOptions(writer, serviceContext, logger.LogServiceOptions{
		CancellationDeadline: 1 * time.Second,
	})

	go func() {
		logger := loggerSvc.NewContext("1", "")
		logger.Info("1")
	}()

	go func() {
		logger := loggerSvc.NewContext("2", "")
		logger.Info("2")
	}()

	loggerSvc.Finish()
}

func TestLogServiceEscapesJSON(t *testing.T) {
	// timestamp, level, and severity don't require escaping since they're set
	// internally.
	//
	// environment, system name, service name, and service instance id are all
	// set at the LogService level when the LogServiceContext is ingested.
	//
	// site and operation are set at the LogService level when a new LogContext
	// is created.
	//
	// message and details are set at the LogEntry level when LogEntry.Serialize
	// is called.
	//
	// In some test cases, runes are duplicated to also account for a case where
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
			// 6 standard J methods plus the Write method.
			const numberOfJMethods = 7

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := mock_io.NewMockWriter(ctrl)
			writer.EXPECT().Write(gomock.Any()).
				Times(numberOfJMethods).
				DoAndReturn(
					func(bb []byte) (int, error) {
						msg := string(bb)
						d := json.NewDecoder(strings.NewReader(msg))
						for {
							_, err := d.Token()
							if err != nil && err == io.EOF {
								return len(msg), nil
							}
							if err != nil {
								// Test is failed with a panic to ensure that
								// the logService doesn't try to handle the test
								// error as a normal error. Panicking also
								// simplifies bubbling up the test failure in
								// gomock since we have a lot of goroutines in
								// the mix (both from the test framework and
								// nobslogger internals).
								panic(fmt.Sprintf("Error: %v\nJSON: %s", err, msg))
							}
						}
					},
				)

			logService := logger.InitializeWriterWithOptions(writer,
				logger.ServiceContext{
					Environment:       s,
					ServiceInstanceID: s,
					ServiceName:       s,
					SystemName:        s,
				}, logger.LogServiceOptions{
					CancellationDeadline: 10 * time.Millisecond,
				})
			log := logService.NewContext(s, s)

			jsonLogMethods := []func(string, string){
				log.TraceJ,
				log.InfoJ,
				log.DebugJ,
				log.WarnJ,
				log.ErrorJ,
				log.FatalJ,
			}

			for _, fn := range jsonLogMethods {
				fn(s, s)
			}

			// don't forget to excercise the Write method as well as the other
			// J methods.
			log.Write([]byte(s))

			logService.Finish()
		}
		t.Run("rune:"+s, f)
	}

}

// Here we're just verifying that all of the context methods result in a call
// to the underlaying writer.
func Test_ContextMethodHappyPath(t *testing.T) {
	// 6 standard methods, 6 D methods, 6 J methods, and 1 Write method.
	const numberOfLogMethods = 19

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(numberOfLogMethods)

	loggerService := logger.InitializeWriterWithOptions(writer, logger.ServiceContext{}, logger.LogServiceOptions{})
	logger := loggerService.NewContext("context site", "operation")

	plainLogMethods := []func(string){
		logger.Trace,
		logger.Info,
		logger.Debug,
		logger.Warn,
		logger.Error,
		logger.Fatal,
	}

	for _, logMethod := range plainLogMethods {
		logMethod("test message")
	}

	detailLogMethods := []func(string, string){
		logger.TraceD,
		logger.InfoD,
		logger.DebugD,
		logger.WarnD,
		logger.ErrorD,
		logger.FatalD,
	}

	for _, logMethod := range detailLogMethods {
		logMethod("test message", "test detail")
	}

	jsonLogMethods := []func(string, string){
		logger.TraceJ,
		logger.InfoJ,
		logger.DebugJ,
		logger.WarnJ,
		logger.ErrorJ,
		logger.FatalJ,
	}

	for _, logMethod := range jsonLogMethods {
		logMethod("test message", "test detail")
	}

	// don't forget to excercise the Write method as well as the other
	// J methods.
	logger.Write([]byte("test message"))

	loggerService.Finish()
}
