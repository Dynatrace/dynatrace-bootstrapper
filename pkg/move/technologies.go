package move

import (
	"encoding/json"
	"path/filepath"
	"strings"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
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

func copyByTechnology(log logr.Logger, fs afero.Afero, from string, to string) error {
	log.Info("starting to copy (filtered)", "from", from, "to", to)

	filteredPaths, err := filterFilesByTechnology(log, fs, from, strings.Split(technology, ","))
	if err != nil {
		return err
	}

	for _, sourceFilePath := range filteredPaths {
		targetFilePath := filepath.Join(to, strings.Split(sourceFilePath, from)[1])

		sourceStatMode, err := fs.Stat(from)
		if err != nil {
			log.Error(err, "error checking stat mode from source folder")

			return err
		}

		err = fs.MkdirAll(filepath.Dir(targetFilePath), sourceStatMode.Mode())
		if err != nil {
			log.Error(err, "error creating target folder")

			return err
		}

		log.V(1).Info("copying file %s to %s", "from", sourceFilePath, "to", targetFilePath)

		err = fsutils.CopyFile(fs, sourceFilePath, targetFilePath)
		if err != nil {
			log.Error(err, "error copying file")

			return err
		}
	}

	return nil
}

func filterFilesByTechnology(log logr.Logger, fs afero.Afero, source string, technologies []string) ([]string, error) {
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
			log.Info("technology not found", "tech", tech)
			continue
		}

		for arch, files := range techData {
			log.V(1).Info("collecting files for technology", "tech", tech, "arch", arch)

			for _, file := range files {
				paths = append(paths, filepath.Join(source, file.Path))
			}
		}
	}

	return paths, nil
}
