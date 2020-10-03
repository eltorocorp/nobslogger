package nobslogger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

// Initialize establishes a connection to a specified UDP server (typically
// logstash), starts an internal log message poller, and returns a LogService
// instance through which a more detailed logging context can be established.
// This function will panic if an error occurs while establishing the connection
// to the UDP server.
// Special case: If hostURI is supplied as an empty string, the logger will
// run, but all log messages are sent to a "null writer" (see `ioutil.Discard`).
func Initialize(hostURI string) *LogService {
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
	}
}

// LogService provides access to an upstream UDP log server (such as LogStash).
type LogService struct {
	messageChannel chan *LogEntry

	// LogWriter is an io.Writer that is exposed to allow the standard library's
	// logger to also transmit logs a la `log.SetOutput(logService.LogWriter)`.
	LogWriter io.Writer
}

// NewContext provides high level structured information used to decorate
// log messages, and exposes methods for writing at various log levels.
func (ls *LogService) NewContext(environment, systemName, serviceName, serviceInstanceID string) *LogContext {
	return &LogContext{
		logService:        ls,
		Environment:       environment,
		SystemName:        systemName,
		ServiceName:       serviceName,
		ServiceInstanceID: serviceInstanceID,
	}
}

func (ls *LogService) submitAsync(logContext *LogContext, logDetail *LogDetail) {
	ls.messageChannel <- &LogEntry{
		LogContext: *logContext,
		LogDetail:  *logDetail,
	}
}
