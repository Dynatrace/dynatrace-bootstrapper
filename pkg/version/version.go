// Package version provides versioning information for the dynatrace-bootstrapper application.
package version

import (
	"runtime/debug"

	"github.com/go-logr/logr"
)

var (
	// AppName contains the name of the application.
	AppName = "dynatrace-bootstrapper"

	// Version contains the version of the Bootstrapper. Assigned externally.
	Version = ""

	// Commit indicates the Git commit hash the binary was build from. Assigned externally.
	Commit = ""

	// BuildDate is the date when the binary was build. Assigned externally.
	BuildDate = ""

	// ModuleSum is the module checksum of the main module. Assigned externally.
	ModuleSum = ""
)

func init() {
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if Version == "" {
		Version = i.Main.Version
	}

	ModuleSum = i.Main.Sum
}

// Print logs the version information using the provided logger.
func Print(log logr.Logger) {
	keyValues := []any{"name", AppName, "version", Version}

	if ModuleSum != "" {
		keyValues = append(keyValues, "module-sum", ModuleSum)
	}

	if Commit != "" {
		keyValues = append(keyValues, "commit", Commit)
	}

	if BuildDate != "" {
		keyValues = append(keyValues, "build_date", BuildDate)
	}

	log.Info("version info", keyValues...)
}
