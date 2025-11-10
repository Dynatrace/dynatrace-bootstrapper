package preload

import (
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	LibAgentProcPath = "agent/lib64/liboneagentproc.so"
	ConfigPath       = "oneagent/ld.so.preload"
)

func Configure(log logr.Logger, configDir, installPath string) error {
	log.Info("configuring ld.so.preload", "config-directory", configDir, "install-path", installPath)

	return fsutils.CreateFile(filepath.Join(configDir, ConfigPath), filepath.Join(installPath, LibAgentProcPath))
}
