package move

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Manifest struct {
	Technologies TechEntries `json:"technologies"`
	Version      string      `json:"version"`
}

type TechEntries map[string]ArchEntries
type ArchEntries map[string][]FileEntry

type FileEntry struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	MD5     string `json:"md5"`
}

var _ copyFunc = copyByTechnology

func copyByTechnology(fs afero.Afero, from string, to string) error {
	logrus.Infof("Starting to copy (filtered) from %s to %s", from, to)

	filteredPaths, err := filterFilesByTechnology(fs, from, strings.Split(technology, ","))
	if err != nil {
		return err
	}

	for _, path := range filteredPaths {
		targetFile := filepath.Join(to, strings.Split(path, from)[1])

		sourceStatMode, err := fs.Stat(from)
		if err != nil {
			logrus.Errorf("Error checking stat mode from source folder: %v", err)

			return err
		}

		err = fs.MkdirAll(filepath.Dir(targetFile), sourceStatMode.Mode())
		if err != nil {
			logrus.Errorf("Error creating target folder: %v", err)

			return err
		}

		logrus.Infof("Copying file %s to %s", path, targetFile)

		err = copyFile(fs, path, targetFile)
		if err != nil {
			logrus.Errorf("Error copying file: %v", err)

			return err
		}
	}

	return nil
}

func filterFilesByTechnology(fs afero.Afero, source string, technologies []string) ([]string, error) {
	manifestPath := filepath.Join(source, "manifest.json")

	manifestFile, err := fs.ReadFile(manifestPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open manifest.json")
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestFile, &manifest); err != nil {
		return nil, errors.WithMessage(err, "failed to parse manifest.json")
	}

	var paths []string

	for _, tech := range technologies {
		techData, exists := manifest.Technologies[tech]
		if !exists {
			logrus.Warnf("Technology %s not found", tech)
			continue
		}

		for arch, files := range techData {
			logrus.Infof("Collecting files for technology %s for arch %s", tech, arch)

			for _, file := range files {
				paths = append(paths, filepath.Join(source, file.Path))
			}
		}
	}

	return paths, nil
}
