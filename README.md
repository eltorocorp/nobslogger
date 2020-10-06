# :cow: NobSlogger
NobSlogger. A fast, lightweight, no-BS, static-structured/leveled logger.

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/eltorocorp/nobslogger)
[![Go Report Card](https://goreportcard.com/badge/github.com/eltorocorp/nobslogger)](https://goreportcard.com/report/github.com/eltorocorp/nobslogger)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/eltorocorp/nobslogger/master/LICENSE)
[![Coverage](http://gocover.io/_badge/github.com/rs/zerolog)](http://gocover.io/github.com/eltorocorp/nobslogger)

> NobSlogger is currently pre-release, and its API should be considered unstable.

No BS:
 - NobSlogger doesn't try to bend to everybody's idea of what should be logged and how it should be structured.
 - It has a staticly structured log format that is focused on use for microservice activity logs.
 - It is focused on providing structured information that helps identify what is happening, where it is happening, and in what order it is happening, so that issues can be identified quickly.
 - If you want more flexibility on which fields you want to log or how they are formatted, use Zap or Zerolog.
 - If you just want something that is minimal, fast, and doesn't deal with any BS, then use this.

# Installation

`go get -u github.com/eltorocorp/nobslogger`

# Performance

NobSlogger is very opinionated. And it is fast as a result\*.
| Package                               | Time         |  Time % to nobSlogger | Objects Allocated |
|---------------------------------------|--------------|-----------------------|-------------------|
| :cow: eltorocorp/nobslogger.Info-4    | 3456 ns/op   | +0%                   | 3 allocs/op       |
| :cow: eltorocorp/nobslogger.InfoD-4   | 3665 ns/op   | +6%                   | 3 allocs/op       |
| rs/zerolog.Check-4                    | 4710 ns/op   | +36%                  | 0 allocs/op       |
| rs/zerolog-4                          | 4729 ns/op   | +37%                  | 0 allocs/op       |
| Zap-4                                 | 6454 ns/op   | +87%                  | 1 allocs/op       |
| Zap.Check-4                           | 6542 ns/op   | +89%                  | 1 allocs/op       |
| Zap.Sugar-4                           | 7545 ns/op   | +118%                 | 3 allocs/op       |
| go-kit/kit/log-4                      | 20891 ns/op  | +504%                 | 25 allocs/op      |
| apex/log-4                            | 71477 ns/op  | +1968%                | 27 allocs/op      |
| inconshreveable/log15-4               | 81934 ns/op  | +2271%                | 35 allocs/op      |
| sirupsen/logrus-4                     | 84882 ns/op  | +2356%                | 39 allocs/op      |

> *\*NobSlogger's benchmarks are based on the accumulated context benchmark suite used by the Zap and Zerolog loggers.
Since NobSlogger is so much more opinionated than Zap and Zerolog, the other benchmarks they often use to compare eachother don't apply well to NobSlogger. However, the accumulated context benchmark suite is a fair representation of what is likely the most apples to apples use case between all three systems.*

# Log Structure
As mentioned in the No BS section above, NobSlogger gets its performance by being very opinionated about what constitutes a log entry. It does not try to be all things to all people, but does succeed at doing what will work for most scenarios really well.

Logs are structred at three levels. The Service, Context, and Entry (described in more detail below).
Each structural level presents progressively more detail about the context within which each log entry occurs. This is designed to help identify what, where, and when things are going on with a system without getting bogged down in too much BS.

## Log Service Level
*Values applied to all logs within this instance across all goroutines and contexts.*
- Environment: *The deployment stage this service is active within. i.e. "stage", "dev", "prod", etc.*
- System Name: *The name of the broader system within which this service participates.*
- Service Name: *The name of this service in particular.*
- Service Instance ID: *An ID that uniquely designates this service instance independent from other parallel instances within the current environment and system.*

## Log Context Level
*Values applied in a more narrow context. Typically within a package, or other "sub-module" within a service.*

- Site: *The general region to which logs for the current context apply.*
- Operation *The general operation being performed within this region of the system.*

## Log Entry Level
*Values that are specific to a discrete log entry.*

- Timestamp: *A 64 bit UTC unixnano timestamp. This value is numeric because that is faster than providing a formatted value.*
- Message: *A concrete message describing system state.*
- Detail: *Additional information in support of the message.*
- Severity: *A value describing the nature of the log message. One of trace, debug, info, warn, error, or fatal.*
- Level: *A numeric value with respect to the log severity. One of 100, 200, 300, 400, 500, 600.*

# Examples

The basics
```go
loggerSvc := nobslogger.InitializeUDP("logstash.theclouds.com:1234", &nobslogger.ServiceContext{
    Environment:       "dev",
    ServiceInstanceID: "123456789",
    SystemName:        "grib-app",
    ServiceName:       "foo-service",
})
logger := loggerSvc.NewContext("entrypoint", "initializing service")
logger.Info("starting up")
```

Hook into stdlib/log 
```go
loggerSvc := nobslogger.InitializeUDP("logstash.theclouds.com:1234", &nobslogger.ServiceContext{
    Environment:       "dev",
    ServiceInstanceID: "123456789",
    SystemName:        "grib-app",
    ServiceName:       "foo-service",
})
logger := loggerSvc.NewContext("entrypoint", "initializing service")
logger.Info("starting up")

log.SetOutput(logService.LogWriter)
log.Println("this is forwarded to logstash along with the other structured logs")

logger.FatalD("system is borked", "details about the bork")
```

Multiple logging contexts
```go
loggerSvc := nobslogger.InitializeUDP("logstash.theclouds.com:1234", &nobslogger.ServiceContext{
    Environment:       "dev",
    ServiceInstanceID: "123456789",
    SystemName:        "grib-app",
    ServiceName:       "foo-service",
})

logger1:= logService.NewContext("log context 1", "doing work within context 1")
logger1.InfoD("Logger 1", "this message was generated from the logger1 context")

logger2:= logService.NewContext("log context 2", "doing work within context 2")
logger2.InfoD("Logger 2", "this message was generated from the logger2 context")
```

Working across goroutines
```go
loggerSvc := nobslogger.InitializeUDP("logstash.theclouds.com:1234", &nobslogger.ServiceContext{
    Environment:       "dev",
    ServiceInstanceID: "123456789",
    SystemName:        "grib-app",
    ServiceName:       "foo-service",
})

logger := loggerSvc.NewContext("first logger", "this is one of two log contexts we'll establish")

go func() {
    logger.DebugD("This is on one goroutine", "details!")
}()

go func() {
    logger.DebugD("This is on another goroutine", "more details!")
}()

go func() {
    newLoggerContext:= logService.NewContext("second logger", "this is the second of two log contexts.")
    newLoggerContext.WarnD("Whoa, this is from a different context.", "Crazy details")
}()
```

