package main

// This package only exists to assist with performance tuning, particularly to
// assist in identifying heap allocations outside the context of a benchmark
// run. (benchmark runs significantly muddy the allocation stack trace dump, so
// it's useful to inspect a binary in isolation)
//
// Example use:
//   go build main.go
//   GODEBUG=allocfreetrace=1 ./main 2> allocfreetrace.txt
//	 vim allocfreetrace.txt

import "github.com/eltorocorp/nobslogger/pkg/nobslogger"

func main() {
	loggerSvc := nobslogger.Initialize("", &nobslogger.ServiceContext{
		Environment:       "dev",
		ServiceInstanceID: "123456789",
		ServiceName:       "allocation_catcher",
		SystemName:        "allocation_catcher",
	})
	logger := loggerSvc.NewContext("main", "finding allocations")
	logger.InfoD("a message", "some details")
}
