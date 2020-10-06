package nobslogger_test

import (
	"io/ioutil"
	"testing"

	"github.com/eltorocorp/nobslogger"
)

func Test_ServiceInitialize(t *testing.T) {
	loggerService := nobslogger.InitializeWriter(ioutil.Discard, &nobslogger.ServiceContext{})
	log := loggerService.NewContext("context site", "operation")
	log.Info("message")
}
