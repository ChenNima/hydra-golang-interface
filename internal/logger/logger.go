package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger
var once sync.Once

func newLogger() logrus.FieldLogger {
	newLogger := logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	newLogger.SetFormatter(formatter)
	return newLogger
}

// GetLogger returns logger singleton
func GetLogger() logrus.FieldLogger {
	once.Do(func() {
		logger = newLogger()
	})
	return logger
}
