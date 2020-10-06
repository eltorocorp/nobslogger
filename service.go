package nobslogger

import (
	"fmt"
	"io"
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

// InitializeUDP establishes a connection to a specified UDP server (typically
// logstash), starts an internal log message poller, and returns a LogService
// instance through which a more detailed logging context can be established.
// This function will panic if an error occurs while establishing the connection
// to the UDP server.
func InitializeUDP(hostURI string, serviceContext *ServiceContext) LogService {
	cn, err := net.Dial("udp", hostURI)
	if err != nil {
		panic("error occurred while establishing udp connection")
	}
	return InitializeWriter(cn, serviceContext)
}

// InitializeWriter publishes logs via the provider io.Writer.
func InitializeWriter(w io.Writer, serviceContext *ServiceContext) LogService {
	messageChannel := make(chan LogEntry, 2)
	ls := LogService{
		messageChannel: messageChannel,
		LogWriter:      w,
		serviceContext: serviceContext,
	}
	go ls.handleLogs()
	return ls
}

func (ls *LogService) handleLogs() {
	defer func() {
		close(ls.messageChannel)
	}()
	for {
		select {
		case entry := <-ls.messageChannel:
			_, err := ls.LogWriter.Write(entry.Serialize())
			if err != nil {
				_, err = ls.LogWriter.Write(LogEntry{
					ServiceContext: *ls.serviceContext,
					LogContext: LogContext{
						Site:      "log service",
						Operation: "handleLogs",
					},
					LogDetail: LogDetail{
						Level:     LogLevelError,
						Severity:  LogSeverityError,
						Message:   "error occurred while shipping log data",
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
