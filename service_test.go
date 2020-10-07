package nobslogger_test

import (
	"testing"

	"github.com/eltorocorp/nobslogger"
	"github.com/eltorocorp/nobslogger/mocks/mock_io"
	"github.com/golang/mock/gomock"
)

func Test_ServiceInitializeWriterHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// expectedMsg := `"timestamp":"1602046435294762000","environment":"","system_name":"","service_name":"","service_instance_id":"","level":"300","severity":"info","site":"context site","operation":"operation","message":"some info","details":"",}`

	writer := mock_io.NewMockWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)

	loggerService := nobslogger.InitializeWriter(writer, &nobslogger.ServiceContext{}, nobslogger.LogServiceOptions{})
	logger := loggerService.NewContext("context site", "operation")
	logger.Info("some info")
	loggerService.Cancel()
	loggerService.Wait()
}

// func Test_ServiceInitializeWriterPersistentError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	writer := mock_io.NewMockWriter(ctrl)
// 	writer.EXPECT().
// 		Write(gomock.Any()).
// 		Return(0, fmt.Errorf("test error")).
// 		Times(2)

// 	loggerService := nobslogger.InitializeWriter(writer, &nobslogger.ServiceContext{})
// 	logger := loggerService.NewContext("context site", "operation")
// 	logger.Info("message")
// 	loggerService.Cancel()
// 	loggerService.Wait()
// }

// Tests
// At least one happy path test confirming output format
// A singlur error condition
// a persistent error condition

// Examples
// Recreate examples from readme (just with a fake UDP client)
// Use of cancel and wait methods
