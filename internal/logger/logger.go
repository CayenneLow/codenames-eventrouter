package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func Init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})
	logger.Info("Logger Initialized")
}

func Debug(args ...interface{}) {
	logger.Out = os.Stdout
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
