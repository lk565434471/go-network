package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go-network/utils"
	"io"
	"os"
)

var defaultTimestampFormatter string = "2006-01-02 15:04:05.000"

type LogSettings struct {
	Settings logrus.Formatter
	LogDir string
	LogFilename string
	LogLevel logrus.Level
	ReportCaller bool
	EnableStdout bool
}

type Logger struct {
	logrus.Logger
	settings LogSettings
}

func NewLogger(settings LogSettings, options ...rotatelogs.Option) *Logger {
	logger := &Logger{}
	logger.SetFormatter(settings.Settings)
	logger.SetLevel(settings.LogLevel)
	logger.SetReportCaller(settings.ReportCaller)

	if settings.EnableStdout && settings.LogFilename != "" {
		writer, err := createStdoutAndFileWriter(settings, options...)

		if err != nil {
			return nil
		}

		logger.SetOutput(writer)
	} else if settings.LogFilename != "" {
		writer, err := createFileWriter(settings, options...)

		if err != nil {
			return nil
		}

		logger.SetOutput(writer)

	} else {
		logger.SetOutput(os.Stdout)
	}

	return logger
}

func createFileWriter(settings LogSettings, options ...rotatelogs.Option) (io.Writer, error) {
	filename := utils.JoinPath(settings.LogDir, settings.LogFilename)

	writer, err := rotatelogs.New(filename, options...)

	if err != nil {
		return nil, err
	}

	return writer, nil
}

func createStdoutAndFileWriter(settings LogSettings,
	options ...rotatelogs.Option) (io.Writer, error) {
	writer, err := createFileWriter(settings, options...)

	if err != nil {
		return nil, err
	}

	writers := []io.Writer{
		os.Stdout,
		writer,
	}

	multiWriter := io.MultiWriter(writers...)

	return multiWriter, nil
}