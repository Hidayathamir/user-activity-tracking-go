package x

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func SetupLogger(logLevel string) {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logger.SetReportCaller(true)
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.JSONFormatter{})

	Logger = logger
}

func SetLogger(logger *logrus.Logger) {
	if logger == nil {
		logger = logrus.New()
	}
	Logger = logger
}
