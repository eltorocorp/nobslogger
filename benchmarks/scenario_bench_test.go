package benchmarks

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/eltorocorp/nobslogger/pkg/nobslogger"
	"go.uber.org/zap"
)

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Log(getMessage(0), getMessage(1))
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("stdlib.Println", func(b *testing.B) {
		logger := log.New(ioutil.Discard, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog.Check", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if e := logger.Info(); e.Enabled() {
					e.Msg(getMessage(0))
				}
			}
		})
	})
	b.Run("eltorocorp/nobslogger.Info", func(b *testing.B) {
		logService := nobslogger.Initialize("", nobslogger.ServiceContext{
			Environment:       "",
			SystemName:        "",
			ServiceName:       "",
			ServiceInstanceID: "",
		})
		logger := logService.NewContext("", "")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("eltorocorp/nobslogger.InfoD", func(b *testing.B) {
		logService := nobslogger.Initialize("", nobslogger.ServiceContext{
			Environment:       "",
			SystemName:        "",
			ServiceName:       "",
			ServiceInstanceID: "",
		})
		logger := logService.NewContext("", "")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoD(getMessage(0), getMessage(1))
			}
		})
	})
}

func BenchmarkAccumulatedContext(b *testing.B) {
	b.Logf("Logging with some accumulated context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog().WithFields(fakeApexFields())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog(fakeSugarFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Log(getMessage(0), getMessage(1))
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15().New(fakeSugarFields())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus().WithFields(fakeLogrusFields())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := fakeZerologContext(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog.Check", func(b *testing.B) {
		logger := fakeZerologContext(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if e := logger.Info(); e.Enabled() {
					e.Msg(getMessage(0))
				}
			}
		})
	})
	b.Run("eltorocorp/nobslogger.Info", func(b *testing.B) {
		logService := nobslogger.Initialize("", nobslogger.ServiceContext{
			Environment:       "",
			SystemName:        "",
			ServiceName:       "",
			ServiceInstanceID: "",
		})
		logger := logService.NewContext("", "")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("eltorocorp/nobslogger.InfoD", func(b *testing.B) {
		logService := nobslogger.Initialize("", nobslogger.ServiceContext{
			Environment:       "",
			SystemName:        "",
			ServiceName:       "",
			ServiceInstanceID: "",
		})
		logger := logService.NewContext("", "")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoD(getMessage(0), getMessage(1))
			}
		})
	})
}

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("Logging with additional context at each log site.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeFields()...)
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write(fakeFields()...)
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeApexFields()).Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Log(fakeSugarFields()...)
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fakeZerologFields(logger.Info()).Msg(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog.Check", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if e := logger.Info(); e.Enabled() {
					fakeZerologFields(e).Msg(getMessage(0))
				}
			}
		})
	})
	b.Run("eltorocorp/nobslogger.Info", func(b *testing.B) {
		logService := nobslogger.Initialize("", nobslogger.ServiceContext{
			Environment:       field1Value,
			SystemName:        field2Value,
			ServiceName:       field3Value,
			ServiceInstanceID: field4Value,
		})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			logger := logService.NewContext(field6Value, field7Value)
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("eltorocorp/nobslogger.InfoD", func(b *testing.B) {
		logService := nobslogger.Initialize("", nobslogger.ServiceContext{
			Environment:       field1Value,
			SystemName:        field2Value,
			ServiceName:       field3Value,
			ServiceInstanceID: field4Value,
		})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			logger := logService.NewContext(field6Value, field7Value)
			for pb.Next() {
				logger.InfoD(getMessage(0), getMessage(1))
			}
		})
	})
}
