# :cow: nobslogger
NobSlogger. A fast, lightweight no BS structured/leveled logger.

# Installation

**!** *nobSlogger is currently in beta and should be considered unstable*

`go get -u github.com/eltorocorp/nobslogger`

# Quick Start

The basics
```go
logService:= nobslogger.Initialize("logstash.theclouds.com:1234")
logger:= logService.NewContext("production", "grib-app", "foo-service", "instance 12abc")
logger.Fatal("system is borked", "details about the bork")
```

Capturing stdlib.log 
```go
logService:= nobslogger.Initialize("logstash.theclouds.com:1234")
logger:= logService.NewContext("production", "grib-app", "foo-service", "instance 12abc")

log.SetOutput(logService.LogWriter)
log.Println("this is forwarded to logstash along with the other structured logs")

logger.Fatal("system is borked", "details about the bork")
```

Multiple logging contexts
```go
logService:= nobslogger.Initialize("logstash.theclouds.com:1234")

logger1:= logService.NewContext("production", "grib-app", "foo-service", "instance 1")
logger1.Info("Logger 1", "this message was generated from the logger1 context")

logger2:= logService.NewContext("production", "grib-app", "foo-service", "instance 2")
logger2.Info("Logger 2", "this message was generated from the logger2 context")

```

Working across goroutines
```go
logService:= nobslogger.Initialize("logstash.theclouds.com:1234")
logger:= logService.NewContext("production", "grib-app", "foo-service", "instance 12abc")

go func() {
    logger.Debug("This is on one goroutine", "details!")
}()

go func() {
    logger.Debug("This is on another goroutine", "more details!")
}()

go func() {
    newLoggerContext:= logService.NewContext("production", "grib-app", "foo-service", "instance 3")
    newLoggerContext.Warn("Whoa, this is from a different context.", "Crazy details")
}()
```

# Performance

## Accumulated Context

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
