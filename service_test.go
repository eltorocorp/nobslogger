package nobslogger_test

import (
	"fmt"
	"testing"

	"github.com/eltorocorp/nobslogger"
)

type fmtWriter struct{}

func (*fmtWriter) Write(bb []byte) (int, error) {
	return 0, fmt.Errorf("test error")
}

func Test_ServiceInitialize(t *testing.T) {
	loggerService := nobslogger.InitializeWriter(new(fmtWriter), &nobslogger.ServiceContext{})
	logger := loggerService.NewContext("context site", "operation")
	logger.Info("message")
	loggerService.Cancel()
	loggerService.Done()
}

// Tests
// At least one happy path test confirming output format
// A singlur error condition
// a persistent error condition

// Examples
// Recreate examples from readme (just with a fake UDP client)
// Use of cancel and wait methods
