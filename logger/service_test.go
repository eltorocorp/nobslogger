package logger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/eltorocorp/nobslogger/logger"
	"github.com/eltorocorp/nobslogger/mocks/mock_io"
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

// The LogService must support ingestion of logs from LogContexts on different
// goroutines.
func Test_LogServiceSupportsMultipleContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(2)

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

// The LogService must support cancellation from a separate goroutine.
func Test_LogServiceSupportsCancellationFromSeparateGoroutine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(3)

	serviceContext := logger.ServiceContext{}
	loggerSvc := logger.InitializeWriter(writer, serviceContext)

	logger := loggerSvc.NewContext("n/a", "n/a")

	go func() {
		logger.Warn("extremely long operation")
		time.Sleep(24 * 365.25 * time.Hour)
	}()

	go func() {
		logger.Info("logging")
	}()

	go func() {
		logger.Info("cancelling")
		loggerSvc.Cancel()
	}()

	loggerSvc.Wait()
}
