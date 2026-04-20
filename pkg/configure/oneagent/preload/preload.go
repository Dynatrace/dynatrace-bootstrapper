package preload

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	LibAgentProcPath = "agent/lib64/liboneagentproc.so"
	ConfigPath       = "oneagent/ld.so.preload"
)

func Configure(log logr.Logger, configDir, installPath string) error {
	log.Info("configuring ld.so.preload", "config-directory", configDir, "install-path", installPath)

	if err := validateInstallPath(installPath); err != nil {
		return err
	}

	return fsutils.CreateFile(filepath.Join(configDir, ConfigPath), filepath.Join(installPath, LibAgentProcPath))
}

func validateInstallPath(installPath string) error {
	if !filepath.IsAbs(installPath) {
		return fmt.Errorf("install path must be absolute, got: %s", installPath)
	}

	if strings.ContainsAny(installPath, "\n\r\t\x00 ,:") {
		return errors.New("install path must be a single path with no separators or whitespace")
	}

	if cleaned := filepath.Clean(installPath); cleaned != installPath {
		return fmt.Errorf("install path must be a clean path, got: %s (expected %s)", installPath, cleaned)
	}

	return nil
}
