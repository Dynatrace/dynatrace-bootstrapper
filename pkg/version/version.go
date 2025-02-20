package version

import (
	"github.com/go-logr/logr"
)

var (

	// AppName contains the name of the application
	AppName = "dynatrace-bootsrapper"

	// Version contains the version of the Bootstrapper. Assigned externally.
	Version = "snapshot"

	// Commit indicates the Git commit hash the binary was build from. Assigned externally.
	Commit = ""

	// BuildDate is the date when the binary was build. Assigned externally.
	BuildDate = ""
)

func Print(log logr.Logger) {
	log.Info("version info", "name", AppName, "version", Version, "commit", Commit, "build_date", BuildDate)
}
