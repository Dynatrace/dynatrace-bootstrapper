package version

import "fmt"

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

func Print() {
	fmt.Printf("name: %s, version: %s, commit: %s, build_date: %s", AppName, Version, Commit, BuildDate)
}
