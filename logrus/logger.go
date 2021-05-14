package logger

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
)

var defaultTimestampFormatter string = "2006-01-02 15:04:05.000"

type LogSettings struct {
	Settings logrus.Formatter
	LogFileName string
	LogFileNameSuffix string
	LogPath string
	Writer io.Writer
	LogLevel logrus.Level
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

	logger.SetLevel(settings.LogLevel)
	logger.SetOutput(settings.Writer)

	return logger
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
	return LogSettings{}
}

func NewDefaultTextFormatter() LogSettings {
	return LogSettings{}
}