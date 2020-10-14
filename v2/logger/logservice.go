package logger

import (
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// 64kib is the initial allocation for the message buffer, as this matches the
// theorhtical (max) MTU for UDP transmissions.
const initialMsgBufferAllocation = 64 * 1024

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
	CancellationDeadline time.Duration
}

func defaultLogServiceOptions() LogServiceOptions {
	return LogServiceOptions{
		CancellationDeadline: 30 * time.Second,
	}
}

// LogService provides access to a writer such as that for a file system or an
// upstream UDP endpoint.
type LogService struct {
	locked         int32
	waiters        uint32
	errMsgBuffer   []byte
	serviceContext *ServiceContext
	options        LogServiceOptions
	logWriter      io.Writer
}

// InitializeUDP establishes a connection to a specified UDP server (such as
// logstash),  and returns a LogService instance through which more detailed
// logging contexts can be spawned (see NewContext)
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

// InitializeWriter establishes a logging service that transmits logs to the
// provided io.Writer.
func InitializeWriter(writer io.Writer, serviceContext ServiceContext) LogService {
	return InitializeWriterWithOptions(writer, serviceContext, defaultLogServiceOptions())
}

// InitializeWriterWithOptions is the same as InitializeWriter, but with custom
// LogServiceOptions supplied. See InitializeWriter.
func InitializeWriterWithOptions(w io.Writer, serviceContext ServiceContext, options LogServiceOptions) LogService {
	serviceContext.Environment = escape(serviceContext.Environment)
	serviceContext.ServiceInstanceID = escape(serviceContext.ServiceInstanceID)
	serviceContext.ServiceName = escape(serviceContext.ServiceName)
	serviceContext.SystemName = escape(serviceContext.SystemName)

	ls := LogService{
		locked:         0,
		waiters:        0,
		errMsgBuffer:   make([]byte, initialMsgBufferAllocation),
		serviceContext: &serviceContext,
		options:        options,
		logWriter:      w,
	}

	return ls
}

func (ls *LogService) writeEntry(msg []byte) {
	_, err := ls.logWriter.Write(msg)
	if err != nil {
		// We dump the original message to stdErr and try to transmit the error
		// notification back to the writer. If error transmission fails, we
		// write the original error transmission as well as the newest error
		// to stdErr and continue on.
		//
		// An error message context is constructed manually here and submitted
		// directly to `logWriter.Write` (as opposed to instantiating via
		// `NewContext` and using the `LogContext.submit` method). This allows
		// the error message to retain priority in the message pipeline.
		stdErr := log.New(os.Stderr, "", 0)
		offset := serialize(ls.errMsgBuffer, ls.serviceContext,
			&LogContext{
				Site:      "log service",
				Operation: "handleLogs",
			},
			LogDetail{
				Level:     LogLevelError,
				Severity:  LogSeverityError,
				Message:   "error occurred while shipping log data",
				Details:   err.Error(),
				Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			},
		)
		stdErr.Println(string(msg))
		_, err = ls.logWriter.Write(ls.errMsgBuffer[0:offset])
		if err != nil {
			stdErr.Println(string(ls.errMsgBuffer[0:offset]))
			stdErr.Println(err.Error())
		}
	}
}

// NewContext provides high level structured information used to decorate
// log messages, and exposes methods for writing at various log levels.
func (ls *LogService) NewContext(site, operation string) LogContext {
	return LogContext{
		logService: ls,
		buffer:     make([]byte, initialMsgBufferAllocation),
		Site:       escape(site),
		Operation:  escape(operation),
	}
}

// Finish sets a deadline for any concurrent LogContexts to finish sending any
// remaining messages; then blocks until that deadline has expired.
// Finish will reset its internal deadline if any messages are received during
// the waiting period. Finish will only unblock once the full deadline duration
// has elapsed with no inbound log activity.
//
// Finish makes a good faith attempt to flush all inbound messages during the
// waiting period. However, it remains the host system's responsibility to wind
// down all components that might be broadcasting logs to LogContexts before
// calling Finish.
//
// If the host system continues to send log messages to the log service while
// Finishing, the log service will either a) never exit because it keeps
// receiving messages or b) exit before all messages have been processed.
func (ls *LogService) Finish() {
	deadline := time.Now().Add(ls.options.CancellationDeadline)
	for {
		runtime.Gosched()
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
