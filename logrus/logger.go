package logger

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
)

var defaultTimestampFormatter string = "2006-01-02 15:04:05.000"

type LogSettings struct {
	Settings logrus.Formatter
	LogFilename string
	LogFilenameSuffix string
	LogDir string
	Writer io.Writer
	LogLevel logrus.Level
	ReportCaller bool
}

type Logger struct {
	logrus.Logger
	settings LogSettings
}

func NewLogger(settings LogSettings) *Logger {
	logger := &Logger{}

	if isJsonFormatter(settings.Settings) {
		newSettings, err := buildJsonFormatter(settings)

		if err != nil {
			return nil
		}

		logger.SetFormatter(newSettings)

	} else if isTextFormatter(settings.Settings) {
		newSettings, err := buildTextFormatter(settings)

		if err != nil {
			return nil
		}

		logger.SetFormatter(newSettings)
	}

	var writer io.Writer
	writer, err := buildMultiWriter(settings)

	if err != nil {
		writer = os.Stdout
	}

	logger.SetLevel(settings.LogLevel)
	logger.SetOutput(writer)
	logger.SetReportCaller(settings.ReportCaller)

	return logger
}

func buildMultiWriter(settings LogSettings) (io.Writer, error) {
	if settings.LogFilename == "" {
		return nil, errors.New("")
	}

	filenameSuffix := settings.LogFilenameSuffix
	filename := settings.LogDir + settings.LogFilename + filenameSuffix + ".log"

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.ModePerm)

	if err != nil {
		return nil, err
	}

	writers := []io.Writer{
		f, os.Stdout,
	}

	return io.MultiWriter(writers...), nil
}

func isJsonFormatter(settings logrus.Formatter) bool {
	_, ok := settings.(*logrus.JSONFormatter)

	return ok
}

func buildJsonFormatter(settings LogSettings) (*logrus.JSONFormatter, error) {
	newSettings, ok := settings.Settings.(*logrus.JSONFormatter)

	if !ok {
		return nil, errors.New("")
	}

	if !newSettings.DisableTimestamp && len(newSettings.TimestampFormat) == 0 {
		newSettings.TimestampFormat = defaultTimestampFormatter
	}

	return newSettings, nil
}

func isTextFormatter(settings logrus.Formatter) bool {
	_, ok := settings.(*logrus.TextFormatter)

	return ok
}

func buildTextFormatter(settings LogSettings) (*logrus.TextFormatter, error) {
	newSettings, ok := settings.Settings.(*logrus.TextFormatter)

	if !ok {
		return nil, errors.New("")
	}

	if !newSettings.DisableTimestamp && len(newSettings.TimestampFormat) == 0 {
		newSettings.TimestampFormat = defaultTimestampFormatter
	}

	return newSettings, nil
}

func NewDefaultJsonFormatter() LogSettings {
	return LogSettings{
		Settings: &logrus.JSONFormatter{
				TimestampFormat: defaultTimestampFormatter,
		},
		LogLevel: logrus.TraceLevel,
		Writer: os.Stdout,
	}
}

func NewDefaultJsonLogger() *Logger {
	return NewLogger(NewDefaultJsonFormatter())
}

func NewDefaultTextFormatter() LogSettings {
	return LogSettings{
		Settings: &logrus.TextFormatter{
			TimestampFormat: defaultTimestampFormatter,
			DisableSorting: true,
		},
		LogLevel: logrus.TraceLevel,
		Writer: os.Stdout,
	}
}

func NewDefaultTextLogger() *Logger {
	return NewLogger(NewDefaultTextFormatter())
}