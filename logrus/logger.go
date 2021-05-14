package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Logger struct {
	logrus.Logger
	defaultTimestampFormat string
	outputMode int
}

func (logger *Logger) Init(formatter logrus.Formatter) {
	logger.SetFormatter(formatter)
}

func NewLogger(output io.Writer, level logrus.Level) *Logger {
	logger := &Logger{}

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	logger.SetOutput(output)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetReportCaller(true)

	return logger
}

func NewStdoutLogger() *Logger {
	return nil
}

func NewFileLogger(path string) *Logger {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		return nil
	}

	return NewLogger(f)
}