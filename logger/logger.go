package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/tsaikd/KDGoLib/logrusutil"
)

var logger = logrusutil.NewConsoleLogger()

// Logger return logger instance
func Logger() *logrus.Logger {
	return logger
}
