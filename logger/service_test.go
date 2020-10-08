package logger_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/eltorocorp/nobslogger/logger"
	"github.com/eltorocorp/nobslogger/mocks/mock_io"
	"github.com/golang/mock/gomock"
)

func Test_ServiceInitializeWriterHappyPath(t *testing.T) {
	// Important: Avoid making assertions about the expected input value for
	// writer.Write in this test. gomock would evaluate this on a separate
	// goroutine, and if the assertion fails on a separate goroutine, the test
	// will either hang indefinitely or timeout without a clear failure mode.
	// Ask me how I know.
	// Instead, test the LogEntry.Serialize method directly.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)

	loggerService := logger.InitializeWriterWithOptions(writer, logger.ServiceContext{}, logger.LogServiceOptions{})
	logger := loggerService.NewContext("context site", "operation")
	logger.Info("some info")
	loggerService.Cancel()
	loggerService.Wait()
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
	loggerService.Cancel()
	loggerService.Wait()
}

// The LogService should support ingestion of logs from LogContexts on different
// goroutines.
func Test_LogServiceSupportsMultipleContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := func(bb []byte) (int, error) {
		re := regexp.MustCompile(`\d{19}`)
		msg := re.ReplaceAllString(string(bb), "1234567890123456789")
		return len(msg), nil
	}
	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).DoAndReturn(w).Times(2)

	serviceContext := logger.ServiceContext{}
	loggerSvc := logger.InitializeWriter(writer, serviceContext)

	go func() {
		logger := loggerSvc.NewContext("1", "")
		logger.Info("1")
	}()

	go func() {
		logger := loggerSvc.NewContext("2", "")
		logger.Info("2")
	}()

	loggerSvc.Cancel()
	loggerSvc.Wait()
}
