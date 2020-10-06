package nobslogger

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

// ServiceContext defines structural log elements that are applied to every
// log entry from this log service instance.
type ServiceContext struct {
	// Environment, i.e. dev, stage, prod, etc.
	Environment string

	// The name of the system at large which this service acts within.
	SystemName string

	// The name of this particular service.
	ServiceName string

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
func Initialize(hostURI string, serviceContext *ServiceContext) LogService {
	var conn io.Writer
	if hostURI == "" {
		conn = ioutil.Discard
	} else {
		cn, err := net.Dial("udp", hostURI)
		if err != nil {
			panic("error occured while establishing udp connection")
		}
		conn = cn
	}
	messageChannel := make(chan LogEntry, 2)
	ls := LogService{
		messageChannel: messageChannel,
		LogWriter:      conn,
		serviceContext: serviceContext,
	}
	go ls.handleLogs(conn, messageChannel)
	return ls
}

func (ls *LogService) handleLogs(conn io.Writer, messageChannel chan LogEntry) {
	defer func() {
		close(messageChannel)
	}()
	for {
		select {
		case entry := <-messageChannel:
			_, err := conn.Write(entry.Serialize())
			if err != nil {
				// since this error occured while trying to ship a log, we
				// immediately try to transmit a notification of the situation
				// to the host (rather than queueing this response via the
				// message channel)
				_, err = conn.Write(LogEntry{
					ServiceContext: *ls.serviceContext,
					LogContext: LogContext{
						Site:      "log service",
						Operation: "handleLogs",
					},
					LogDetail: LogDetail{
						Level:     LogLevelError,
						Severity:  LogSeverityError,
						Message:   "error occured while shipping log data",
						Details:   err.Error(),
						Timestamp: strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
					},
				}.Serialize())
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

// LogService provides access to an upstream UDP log server (such as LogStash).
type LogService struct {
	messageChannel chan LogEntry

	// LogWriter is an io.Writer that is exposed to allow the standard library's
	// logger to also transmit logs a la `log.SetOutput(logService.LogWriter)`.
	LogWriter io.Writer

	serviceContext *ServiceContext
}

// NewContext provides high level structured information used to decorate
// log messages, and exposes methods for writing at various log levels.
func (ls *LogService) NewContext(site, operation string) LogContext {
	return LogContext{
		logService: ls,
		Site:       site,
		Operation:  operation,
	}
}

func (ls *LogService) submitAsync(sc ServiceContext, lc LogContext, ld LogDetail) {
	// timestamping with unixnano is faster than converting to a std format.
	ld.Timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	ls.messageChannel <- LogEntry{
		ServiceContext: sc,
		LogContext:     lc,
		LogDetail:      ld,
	}
}
