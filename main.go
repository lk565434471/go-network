package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go-network/application"
	"time"
)

func main() {
	app := application.GetApp()()

	if !app.Init(application.AppSettings{
		AppLogSettings: application.AppLogSettings{
			LogFormatter: &logrus.TextFormatter{},
			LoggingLevel: logrus.TraceLevel,
			LogFilename: "test.log.%Y-%m-%d",
			ReportCaller: true,
			EnableStdout: true,
			Options: []rotatelogs.Option{
				rotatelogs.WithMaxAge(time.Hour * 24 * 180),
				rotatelogs.WithRotationTime(time.Hour * 24),
			},
		},
	}) {
		return
	}

	app.Debug("Hello, world.")

	app.Run()
}
