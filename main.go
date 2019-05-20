package main

import (
	"github.com/sirupsen/logrus"
)

// log holds our main logger instance
var log *logrus.Logger

func main() {

	log = logrus.New()
	initCmd()
	executeCmd()

	log.Debug("done")
}
