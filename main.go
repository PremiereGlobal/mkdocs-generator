package main

import (
	"github.com/sirupsen/logrus"
)

// log holds our main logger instance
var log *logrus.Logger

var version string

func main() {

	if version == "" {
		version = "local"
	}

	log = logrus.New()
	initCmd()
	executeCmd()

	log.Debug("done")
}
