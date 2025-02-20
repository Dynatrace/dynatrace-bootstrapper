package preload

import (
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	libAgentProcPath = "agent/lib64/liboneagentproc.so"
	configPath       = "oneagent/ld.so.preload"
)

func Configure(log logr.Logger, fs afero.Afero, configDir, installPath string) error {
	log.Info("Configuring ld.so.preload", "config-directory", configDir, "install-path", installPath)

	return fsutils.CreateFile(fs, filepath.Join(configDir, configPath), filepath.Join(installPath, libAgentProcPath))
}
