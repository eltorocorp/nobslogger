package benchmarks

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func newLogrus() *logrus.Logger {
	return &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
}

func fakeLogrusFields() logrus.Fields {
	return logrus.Fields{
		field1Name: field1Value,
		field2Name: field2Value,
		field3Name: field3Value,
		field4Name: field4Value,
		field5Name: field5Value,
	}
}
