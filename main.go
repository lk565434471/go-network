package main

import (
	logger2 "go-network/logrus"
)

func main() {

	logger := logger2.NewStdoutLogger()
	logger.Info("aaaa")
	logger.Trace("aaaa")
}
