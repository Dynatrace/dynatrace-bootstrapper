// Package preload provides configuration for OneAgent preload files.
package preload

import (
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

// LibAgentProcPath is the path to the OneAgent process library.
// ConfigPath is the path to the preload configuration file.
const (
	LibAgentProcPath = "agent/lib64/liboneagentproc.so"
	ConfigPath       = "oneagent/ld.so.preload"
)

// Configure sets up the preload configuration file for OneAgent.
func Configure(log logr.Logger, fs afero.Afero, configDir, installPath string) error {
	log.Info("configuring ld.so.preload", "config-directory", configDir, "install-path", installPath)

	return fsutils.CreateFile(fs, filepath.Join(configDir, ConfigPath), filepath.Join(installPath, LibAgentProcPath))
}
