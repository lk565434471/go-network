package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"time"
)

func init() {
	path := "G:\\Go\\go-network\\test.log"
	writer, err := rotatelogs.New(
		path + ".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(2)*time.Second),
	)

	if err != nil {
		return
	}

	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(writer)
}

func main() {

	for  {
		logrus.Debug("Hello, world.")
		time.Sleep(5 * time.Second)
	}
}
