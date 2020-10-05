package benchmarks

import (
	"io/ioutil"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
)

func newApexLog() *log.Logger {
	return &log.Logger{
		Handler: json.New(ioutil.Discard),
		Level:   log.DebugLevel,
	}
}

func fakeApexFields() log.Fields {
	return log.Fields{
		field1Name: field1Value,
		field2Name: field2Value,
		field3Name: field3Value,
		field4Name: field4Value,
		field5Name: field5Value,
	}
}
