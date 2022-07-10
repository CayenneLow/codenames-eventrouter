package logger

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

func Init(logLevel string) {
	logger.SetOutput(os.Stdout)
	if lvl, err := logger.ParseLevel(logLevel); err != nil {
		logger.Fatalf("Unable to parse log level: %s. %s", logLevel, err)
		logger.SetLevel(lvl)
	}
	logger.SetFormatter(&logger.TextFormatter{
		DisableQuote: true,
	})
	logger.Info("Logger Initialized")
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
