package nobslogger_test

import (
	"testing"

	"github.com/eltorocorp/nobslogger/pkg/nobslogger"
)

func Test_ServiceInitialize(t *testing.T) {
	loggerService := nobslogger.Initialize("", &nobslogger.ServiceContext{})
	log := loggerService.NewContext("context site", "operation")
	log.Info("message")
}
