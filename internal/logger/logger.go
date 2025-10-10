package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat:   time.RFC3339,
			DisableHTMLEscape: true,
		},
	)
	logger.SetLevel(logrus.InfoLevel)
}

func Logger() *logrus.Logger {
	return logger
}
