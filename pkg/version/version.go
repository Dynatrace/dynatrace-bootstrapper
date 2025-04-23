package version

import (
	"runtime/debug"

	"github.com/go-logr/logr"
)

var (

	// AppName contains the name of the application
	AppName = "dynatrace-bootsrapper"

	// Version contains the version of the Bootstrapper. Assigned externally.
	Version = ""

	// Commit indicates the Git commit hash the binary was build from. Assigned externally.
	Commit = ""

	// BuildDate is the date when the binary was build. Assigned externally.
	BuildDate = ""
)

func Print(log logr.Logger) {
	keyValues := []any{"name", AppName,}

	i, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if Version == "" {
		Version = i.Main.Version
	}

	keyValues = append(keyValues, "version", Version)

	if i.Main.Sum != "" {
		keyValues = append(keyValues, "module-sum", i.Main.Sum)
	}

	if Commit != "" {
		keyValues = append(keyValues, "commit", Commit)
	}

	if BuildDate != "" {
		keyValues = append(keyValues, "build_date", BuildDate)
	}


	log.Info("version info", keyValues...)
}
