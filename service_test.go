package nobslogger_test

import (
	"fmt"
	"testing"

	"github.com/eltorocorp/nobslogger"
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

	loggerService := nobslogger.InitializeWriter(writer, &nobslogger.ServiceContext{}, nobslogger.LogServiceOptions{})
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

	loggerService := nobslogger.InitializeWriter(writer, &nobslogger.ServiceContext{}, nobslogger.LogServiceOptions{})
	logger := loggerService.NewContext("context site", "operation")
	logger.Info("message")
	loggerService.Cancel()
	loggerService.Wait()
}

// Tests
// At least one happy path test confirming output format
// A singlur error condition
// a persistent error condition

// Examples
// Recreate examples from readme (just with a fake UDP client)
// Use of cancel and wait methods
