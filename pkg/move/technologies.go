package move

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
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

func CopyByTechnologyWrapper(technology string) CopyFunc {
	return func(log logr.Logger, from, to string) error {
		return CopyByTechnology(log, from, to, technology)
	}
}

func CopyByTechnology(log logr.Logger, from string, to string, technology string) error {
	log.Info("starting to copy (filtered)", "from", from, "to", to, "technology", technology)

	filteredPaths, err := filterFilesByTechnology(log, from, strings.Split(technology, ","))
	if err != nil {
		return err
	}

	err = copyByList(log, from, to, filteredPaths)
	if err != nil {
		return err
	}

	log.Info("successfully copied (filtered)", "from", from, "to", to, "technology", technology)

	return nil
}

func copyByList(log logr.Logger, from string, to string, paths []string) error {
	oldUmask := unix.Umask(noPermissionsMask)
	defer unix.Umask(oldUmask)

	fromStat, err := os.Stat(from)
	if err != nil {
		log.Error(err, "error checking stat mode from source folder")

		return err
	}

	err = os.MkdirAll(to, fromStat.Mode())
	if err != nil {
		log.Error(err, "error creating target folder")

		return err
	}

	for _, path := range paths {
		splitPath := strings.Split(path, string(filepath.Separator))
		walkedPath := ""

		for _, subPath := range splitPath {
			walkedPath = filepath.Join(walkedPath, subPath)
			sourcePath := filepath.Join(from, walkedPath)
			targetPath := filepath.Join(to, walkedPath)

			sourceStat, err := os.Stat(sourcePath)
			if err != nil {
				log.Error(err, "failed checking stat mode from source", "path", sourcePath)

				return err
			}

			if sourceStat.IsDir() {
				err := os.Mkdir(targetPath, sourceStat.Mode())
				if err != nil && !os.IsExist(err) {
					log.Error(err, "failed to create new dir", "path", targetPath)

					return err
				}

				log.V(1).Info("created new dir", "from", sourcePath, "to", targetPath, "mode", sourceStat.Mode())

				continue
			}

			log.V(1).Info("copying file", "from", sourcePath, "to", targetPath, "mode", sourceStat.Mode())

			err = fsutils.CopyFile(sourcePath, targetPath)
			if err != nil {
				log.Error(err, "error copying file")

				return err
			}
		}
	}

	return nil
}

func filterFilesByTechnology(log logr.Logger, source string, technologies []string) ([]string, error) {
	manifestPath := filepath.Join(source, "manifest.json")

	manifestFile, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open manifest.json")
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestFile, &manifest); err != nil {
		return nil, errors.WithMessage(err, "failed to parse manifest.json")
	}

	var paths []string

	for _, tech := range technologies {
		tech := strings.TrimSpace(tech)
		techData, exists := manifest.Technologies[tech]

		if !exists {
			log.Info("technology not found", "tech", tech)

			continue
		}

		for arch, files := range techData {
			log.V(1).Info("collecting files for technology", "tech", tech, "arch", arch)

			for _, file := range files {
				paths = append(paths, file.Path)
			}
		}
	}

	return paths, nil
}
