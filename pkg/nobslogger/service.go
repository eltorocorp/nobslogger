package nobslogger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

// GlobalContext defines structural log elements that are applied to every
// log entry from this log service instance.
type GlobalContext struct {
	// Environment, i.e. dev, stage, prod, etc.
	Environment       string

	// The name of the system at large which this service acts within.
	SystemName        string

	// The name of this particular service.
	ServiceName       string

	// An ID that defines this service instance uniquely from other instances
	// of the same service within this system and environment.
	ServiceInstanceID string
}

// Initialize establishes a connection to a specified UDP server (typically
// logstash), starts an internal log message poller, and returns a LogService
// instance through which a more detailed logging context can be established.
// This function will panic if an error occurs while establishing the connection
// to the UDP server.
// Special case: If hostURI is supplied as an empty string, the logger will
// run, but all log messages are sent to a "null writer" (see `ioutil.Discard`).
func Initialize(hostURI string, globalContext GlobalContext) *LogService {
	var conn io.Writer
	var err error
	if hostURI == "" {
		conn = ioutil.Discard
	} else {
		conn, err = net.Dial("udp", hostURI)
		if err != nil {
			panic(fmt.Errorf("error occured while establishing udp connection: %v", err))
		}
	}
	messageChannel := make(chan *LogEntry)
	go func() {
		defer func() {
			log.Println("halting log message poller")
			close(messageChannel)
		}()
	LongPoll:
		for {
			select {
			case entry := <-messageChannel:
				bytes := entry.Serialize()
				_, err = conn.Write(bytes)
				if err != nil {
					fmt.Printf("error occured while transmitting log packet: %v", err)
					break LongPoll
				}
			}
		}
	}()
	return &LogService{
		messageChannel: messageChannel,
		LogWriter:      conn,
		globalContext:  &globalContext,
	}
}

// LogService provides access to an upstream UDP log server (such as LogStash).
type LogService struct {
	messageChannel chan *LogEntry

	// LogWriter is an io.Writer that is exposed to allow the standard library's
	// logger to also transmit logs a la `log.SetOutput(logService.LogWriter)`.
	LogWriter io.Writer

	globalContext *GlobalContext
}

// NewContext provides high level structured information used to decorate
// log messages, and exposes methods for writing at various log levels.
func (ls *LogService) NewContext(site, operation string) *LogContext {
	return &LogContext{
		logService: ls,
		Site: site,
		Operation: operation
	}
}

func (ls *LogService) submitAsync(logContext *LogContext, logDetail *LogDetail) {
	logDetail.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	ls.messageChannel <- &LogEntry{
		LogContext: *logContext,
		LogDetail:  *logDetail,
	}
}
