package benchmarks

import (
	"io/ioutil"

	"github.com/rs/zerolog"
)

func newZerolog() zerolog.Logger {
	return zerolog.New(ioutil.Discard).With().Timestamp().Logger()
}

func fakeZerologFields(e *zerolog.Event) *zerolog.Event {
	return e.
		Str(field1Name, field1Value).
		Str(field2Name, field2Value).
		Str(field3Name, field3Value).
		Str(field4Name, field4Value).
		Str(field5Name, field5Value).
		Str(field6Name, field6Value).
		Str(field7Name, field7Value)
}

func fakeZerologContext(c zerolog.Context) zerolog.Context {
	return c.
		Str(field1Name, field1Value).
		Str(field2Name, field2Value).
		Str(field3Name, field3Value).
		Str(field4Name, field4Value).
		Str(field5Name, field5Value).
		Str(field6Name, field6Value).
		Str(field7Name, field7Value)
}
