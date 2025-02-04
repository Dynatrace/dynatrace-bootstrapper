package move

import (
	"encoding/json"
	"os"
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

func copyByTechnology(fs afero.Afero) error {
	logrus.Infof("Starting to copy (filtered) from %s to %s", sourceFolder, targetFolder)

	filteredPaths, err := filterFilesByTechnology(fs, sourceFolder, strings.Split(technology, ","))
	if err != nil {
		return err
	}

	for _, path := range filteredPaths {
		targetFile := filepath.Join(targetFolder, strings.Split(path, sourceFolder)[1])

		err = fs.MkdirAll(filepath.Dir(targetFile), os.ModePerm) // check source dir's Stat.Mode
		if err != nil {
			logrus.Errorf("Error checking folder: %v", err)

			return err
		}

		err = copyFile(fs, path, targetFile)

		if err != nil {
			logrus.Errorf("Error moving folder: %v", err)

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
			logrus.Warnf("technology %s not found", tech)
			continue
		}

		for _, files := range techData {
			logrus.Infof("processing technology %s", tech)

			for _, file := range files {
				paths = append(paths, filepath.Join(source, file.Path))
			}
		}
	}

	return paths, nil
}
