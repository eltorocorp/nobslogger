# nobslogger
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

## Without fields
|                     | ns/op | % to nobslogger | allocs/op |
|---------------------|-------|-----------------|-----------|
| **nobslogger**      | 3352  | 0%              | 3         |
| zerolog.Check       | 4303  | 28%             | 0         |
| zerolog             | 4324  | 29%             | 0         |
| stdlib.Println      | 5549  | 66%             | 2         |
| kit/log             | 5597  | 67%             | 11        |
| Zap                 | 5959  | 78%             | 1         |
| Zap.Check           | 6225  | 86%             | 1         |
| Zap.Sugar           | 6979  | 108%            | 3         |
| apex/log            | 32501 | 870%            | 7         |
| log15               | 34636 | 933%            | 23        |
| logrus              | 40382 | 1105%           | 26        |

## Accumulated Context
|                     | ns/op | % to nobslogger | allocs/op |
|---------------------|-------|-----------------|-----------|
| **nobslogger**      | 3531  | 0%              | 3         |
| zerolog.Check       | 4708  | 33%             | 0         |
| zerolog             | 4744  | 34%             | 0         |
| Zap                 | 6412  | 82%             | 1         |
| Zap.Check           | 6564  | 86%             | 1         |
| Zap.Sugar           | 7638  | 116%            | 3         |
| kit/log             | 15979 | 353%            | 23        |
| apex/log            | 48602 | 1276%           | 20        |
| logrus              | 61984 | 1655%           | 33        |
| log15               | 75207 | 2030%           | 35        |

## Adding Fields
|                     | ns/op | % to nobslogger | allocs/op |
|---------------------|-------|-----------------|-----------|
| **nobslogger**      | 3628  | 0%              | 3         |
| zerolog             | 7458  | 106%            | 0         |
| zerolog.Check       | 7458  | 106%            | 0         |
| Zap.Check           | 11625 | 220%            | 2         |
| Zap                 | 11931 | 229%            | 2         |
| Zap.Sugar           | 12736 | 251%            | 2         |
| kit/log             | 12841 | 254%            | 18        |
| apex/log            | 50982 | 1305%           | 24        |
| log15               | 64141 | 1668%           | 33        |
| logrus              | 70013 | 1830%           | 37        |