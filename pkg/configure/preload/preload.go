package preload

import (
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	libAgentProcPath = "var/lib/dynatrace/oneagent/agent/lib64/liboneagentproc.so" // TODO: make configurable?
	configPath       = "oneagent/ld.so.preload"
)

func Configure(fs afero.Afero, configDir string) error {
	logrus.Infof("Configuring ld.so.preload, config-directory: %s", configDir)

	return fsutils.CreateFile(fs, filepath.Join(configDir, configPath), libAgentProcPath)
}
