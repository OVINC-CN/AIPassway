package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat:   time.RFC3339,
			DisableHTMLEscape: true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				// extract only the file name from the full path
				functionPaths := strings.Split(frame.Function, "/")
				if len(functionPaths) > 0 {
					function = functionPaths[len(functionPaths)-1]
				} else {
					function = frame.Function
				}
				// file
				file = fmt.Sprintf("%s:%d", frame.File, frame.Line)
				return function, file
			},
		},
	)
	logger.SetLevel(logrus.InfoLevel)
}

func Logger() *logrus.Logger {
	return logger
}
