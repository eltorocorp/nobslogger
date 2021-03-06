module go.uber.org/zap/benchmarks

go 1.15

// replace go.uber.org/zap => ../
replace github.com/eltorocorp/nobslogger/v2 => ../

require (
	github.com/apex/log v1.1.1
	github.com/eltorocorp/nobslogger/v2 v0.0.0-00010101000000-000000000000
	github.com/go-kit/kit v0.9.0
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/rs/zerolog v1.16.0
	github.com/sirupsen/logrus v1.4.2
	go.uber.org/zap v1.16.0
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20180818164646-67afb5ed74ec
)
