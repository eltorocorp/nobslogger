package logger

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	// Run benchmarks to verify performance gains before altering messageChannel
	// buffer size.
	messageChannelBufferSize = 10
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
	// MaxFlushAttempts is the maximum number of times that the LogService will
	// try to flush its queue before finalizing a cancellation request.
	MaxFlushAttempts int

	// TimeBetweenFlushAttempts is the amount of time the LogService will wait
	// for new inbound messages to arrive in its queue before deciding that the
	// queue is empty and finalizing the a cancellation request.
	TimeBetweenFlushAttempts time.Duration
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
	messageChannel chan LogEntry
	cancelChannel  chan struct{}
	doneChannel    chan struct{}
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
// logger does not make any attempts at UDP MTU discovery, and will not
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
	messageChannel := make(chan LogEntry, messageChannelBufferSize)
	cancelChannel := make(chan struct{}, 1)
	doneChannel := make(chan struct{}, 1)

	ls := LogService{
		messageChannel: messageChannel,
		cancelChannel:  cancelChannel,
		doneChannel:    doneChannel,
		serviceContext: &serviceContext,
		options:        options,
		logWriter:      w,
	}

	go func() {
		ls.initiateMessagePoll()
		ls.flushPendingMessages()
		ls.doneChannel <- struct{}{}
	}()

	return ls
}

func (ls *LogService) initiateMessagePoll() {
	for {
		select {
		case <-ls.cancelChannel:
			return
		case entry := <-ls.messageChannel:
			ls.writeEntry(entry)
		}
	}
}

func (ls *LogService) flushPendingMessages() {
	remainingFlushAttempts := ls.options.MaxFlushAttempts
	for {
		select {
		case entry := <-ls.messageChannel:
			remainingFlushAttempts = ls.options.MaxFlushAttempts
			ls.writeEntry(entry)
		default:
			if remainingFlushAttempts == 0 || int64(ls.options.TimeBetweenFlushAttempts) == 0 {
				return
			}
			time.Sleep(ls.options.TimeBetweenFlushAttempts)
			remainingFlushAttempts--
		}
	}
}

func (ls *LogService) writeEntry(entry LogEntry) {
	entryBytes := entry.Serialize()
	_, err := ls.logWriter.Write(entryBytes)
	if err != nil {
		stdErr := log.New(os.Stderr, "", 0)
		errLogEntry := LogEntry{
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
		}.Serialize()
		stdErr.Println(string(entryBytes))
		_, err = ls.logWriter.Write(errLogEntry)
		if err != nil {
			stdErr.Println(string(errLogEntry))
			stdErr.Println(err.Error())
		}
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

func (ls *LogService) submitAsync(sc ServiceContext, lc LogContext, ld LogDetail) {
	// timestamping with unixnano is faster than converting to a std format.
	ld.Timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	ls.messageChannel <- LogEntry{
		ServiceContext: sc,
		LogContext:     lc,
		LogDetail:      ld,
	}
}

// Cancel notifies the LogService that the host system is attempting to
// wind down gracefully. When Cancel is called, LogService will begin
// flushing any backlogged messages that remain within its message queue.
//
// Calling cancel more than once has no additional effect.
//
// Note that it is the host system's responsibility to gracefully
// wind down operations. The host system must call Cancel AFTER the host
// system is quiet (and presumably no longer initiating new log messages via any
// spawned LogContexts). If any LogContext's continue to send log messages to
// the LogService after Cancel is called, the LogService will either never halt
// or may halt before all messages are processed. Note that the Wait method
// will always block unless Cancel is called, and will continue to block until
// the cancellation process (as described above) is finalized.
func (ls *LogService) Cancel() {
	ls.cancelChannel <- struct{}{}
}

// Wait blocks until Cancel is called and all logs in LogService's internal
// queue have been flushed.
func (ls *LogService) Wait() {
	<-ls.doneChannel
}
