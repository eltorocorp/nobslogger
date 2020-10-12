package logger

import (
	"io"
	"net"
	"sync/atomic"
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

// LogServiceOptions exposes configuration settings for LogService behavior.
type LogServiceOptions struct {
	// MaxFlushAttempts is deprecated.
	MaxFlushAttempts int

	// TimeBetweenFlushAttempts is deprecated.
	TimeBetweenFlushAttempts time.Duration

	CancellationDeadline time.Duration
}

func defaultLogServiceOptions() LogServiceOptions {
	return LogServiceOptions{
		MaxFlushAttempts:         10,
		TimeBetweenFlushAttempts: 10 * time.Millisecond,
	}
}

// LogService provides access to a writer such as that for a file system or an
// upstream UDP endpoint.
type LogService struct {
	locked         int32
	waiters        uint32
	serviceContext *ServiceContext
	options        LogServiceOptions
	logWriter      io.Writer
}

// InitializeUDP establishes a connection to a specified UDP server (such as
// logstash), starts an internal log message poller, and returns a LogService
// instance through which more detailed logging contexts can be spawned (see
// NewContext)
//
// This function will panic if an error occurs while establishing the connection
// to the UDP server.
//
// NobSlogger does not make any attempts at UDP MTU discovery, and will not
// prohibit the host system from attempting to send log messages that exceed
// the network's UDP MTU limit. If this limit is exceeded, one of two things
// may occur:
//
// 1) The LogService may return an error while attempting to transmit the
// message. In this case, the LogService will try to log (via UDP) that it
// has received an error while shipping log data. This message will have a
// severity of "Error". If the error log transmission fails, the LogService will
// post the resulting error message to StdErr, and continue on.
//
// 2) If an outbound UDP packet is split or lost downstream, the LogService may
// not have any awareness that the it was lost. In this case the destination
// system might receive a partial log message. As such, it is recommended that
// the destination service be running a json codec that is able to identify
// and flag if/when an inbound message is incomplete.
//
// hostURI: Must be a fully qualified URI including port.
func InitializeUDP(hostURI string, serviceContext ServiceContext) LogService {
	return InitializeUDPWithOptions(hostURI, serviceContext, defaultLogServiceOptions())
}

// InitializeUDPWithOptions is the samme as InitializeUDP, but with custom
// LogServiceOptions supplied. See InitializeUDP.
func InitializeUDPWithOptions(hostURI string, serviceContext ServiceContext, options LogServiceOptions) LogService {
	cn, err := net.Dial("udp", hostURI)
	if err != nil {
		panic("error occurred while establishing udp connection")
	}
	return InitializeWriterWithOptions(cn, serviceContext, options)
}

// InitializeWriter publishes logs via the provided io.Writer.
//
// InitializeWriter initiates a long-poll operation that transmits log messages
// to the specified writer any time a log message is available to write.
func InitializeWriter(writer io.Writer, serviceContext ServiceContext) LogService {
	return InitializeWriterWithOptions(writer, serviceContext, defaultLogServiceOptions())
}

// InitializeWriterWithOptions is the same as InitializeWriter, but with custom
// LogServiceOptions supplied. See InitializeWriter.
func InitializeWriterWithOptions(w io.Writer, serviceContext ServiceContext, options LogServiceOptions) LogService {
	// serviceContext.Environment = escape(serviceContext.Environment)
	// serviceContext.ServiceInstanceID = escape(serviceContext.ServiceInstanceID)
	// serviceContext.ServiceName = escape(serviceContext.ServiceName)
	// serviceContext.SystemName = escape(serviceContext.SystemName)

	ls := LogService{
		locked:         0,
		waiters:        0,
		serviceContext: &serviceContext,
		options:        options,
		logWriter:      w,
	}

	return ls
}

func (ls *LogService) writeEntry(msg []byte) {
	_, err := ls.logWriter.Write(msg)
	if err != nil {
		// stdErr := log.New(os.Stderr, "", 0)
		// errLogEntry := LogEntry{
		// 	ServiceContext: *ls.serviceContext,
		// 	LogContext: LogContext{
		// 		Site:      "log service",
		// 		Operation: "handleLogs",
		// 	},
		// 	LogDetail: LogDetail{
		// 		Level:     LogLevelError,
		// 		Severity:  LogSeverityError,
		// 		Message:   "error occurred while shipping log data",
		// 		Details:   err.Error(),
		// 		Timestamp: strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		// 	},
		// }.Serialize()
		// stdErr.Println(ls.messageBuffer)
		// _, err = ls.logWriter.Write(errLogEntry)
		// if err != nil {
		// 	stdErr.Println(string(errLogEntry))
		// 	stdErr.Println(err.Error())
		// }
	}
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

// Finish sets a deadline for any concurrent LogContexts to finish sending any
// remaining messages.
func (ls *LogService) Finish() {
	deadline := time.Now().Add(ls.options.CancellationDeadline)
	for {
		if time.Now().Before(deadline) {
			continue
		}
		if atomic.LoadUint32(&ls.waiters) > 0 {
			deadline = time.Now().Add(ls.options.CancellationDeadline)
			continue
		}
		return
	}
}
