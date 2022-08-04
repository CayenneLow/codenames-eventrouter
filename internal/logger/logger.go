package logger

import (
	logger "github.com/sirupsen/logrus"
	"os"
)

func Init(logLevel string) {
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logger.JSONFormatter{
		TimestampFormat:   "2006-01-02T15:04:05Z07:00",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       true,
	})
	lvl, err := logger.ParseLevel(logLevel)
	if err != nil {
		logger.Fatalf("Unable to parse log level: %s. %s", logLevel, err)
	}
	logger.SetLevel(lvl)
	logger.Infof("Logger Initialized with log level:%s", logger.GetLevel())
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
