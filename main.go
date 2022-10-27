package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/plunder-app/shack/cmd"
)

// Version is populated from the Makefile and is tied to the release TAG
var Version string

// Build is the last GIT commit
var Build string

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	cmd.Release.Version = Version
	cmd.Release.Build = Build
	cmd.Execute()
}
