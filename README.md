# :cow: NobSlogger 
NobSlogger. A fast, opinionated, lightweight, no-BS, static-structured/leveled logger.

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/eltorocorp/nobslogger/v2/logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/eltorocorp/nobslogger/v2)](https://goreportcard.com/report/github.com/eltorocorp/nobslogger/v2)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/eltorocorp/nobslogger/master/LICENSE)
[![Coverage](http://gocover.io/_badge/github.com/eltorocorp/nobslogger/v2/logger)](http://gocover.io/github.com/eltorocorp/nobslogger/logger/v2)

No BS:
 - NobSlogger is opinionated. 
 - It has a staticly structured log format that is focused on use for microservice activity logs.
 - It is focused on providing structured information that helps identify what is happening, where it is happening, and in what order it is happening, so that issues can be identified quickly.
 - If you want more flexibility on which fields you want to log or how they are formatted, use Zap or Zerolog.
 - If you just want something that is minimal, fast, and doesn't deal with any BS, then use this.

# Installation

`go get -u github.com/eltorocorp/nobslogger/v2`

# Performance

NobSlogger is very opinionated. And it is fast as a result\*.
|Package|Time|Time %|Allocations|
|-------|----|------|-----------|
|:cow: eltorocorp/logger.Info-4 |309  ns/op|    0%|0  allocs/op|
|:cow: eltorocorp/logger.InfoD-4         |324  ns/op|    5%|0  allocs/op|
|rs/zerolog.Check-4       |355  ns/op|   15%|0  allocs/op|
|rs/zerolog-4             |357  ns/op|   16%|0  allocs/op|
|Zap-4         |476  ns/op|   54%|0  allocs/op|
|Zap.Check-4   |505  ns/op|   63%|0  allocs/op|
|Zap.Sugar-4   |971  ns/op|  214%|2  allocs/op|
|go-kit/kit/log-4         |3189  ns/op|  932%|24  allocs/op|
|apex/log-4    |8199  ns/op| 2553%|25  allocs/op|
|sirupsen/logrus-4        |8556  ns/op| 2669%|37  allocs/op|
|inconshreveable/log15-4  |9620  ns/op| 3013%|31  allocs/op|

> *\*NobSlogger's benchmarks are based on the accumulated context benchmark suite used by the Zap and Zerolog loggers. Additional benchmarks (single field, and adding fields) are available in the benchmarks directory. 

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

- Timestamp: *RFC3339Nano (2006-01-02T15:04:05.999999999Z07:00)*
- Message: *A concrete message describing system state.*
- Detail: *Additional information in support of the message.*
- Severity: *A value describing the nature of the log message. One of trace, debug, info, warn, error, or fatal.*
- Level: *A numeric value with respect to the log severity. One of 100, 200, 300, 400, 500, 600.*

# Examples

- view examples in the docs [here](https://pkg.go.dev/github.com/eltorocorp/nobslogger/v2/logger#pkg-examples)
- or view the same examples in code [here](v2/logger/examples_test.go)