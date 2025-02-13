package move

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	sourceFolderFlag       = "source"
	targetFolderFlag       = "target"
	workFolderFlag         = "work"
	technologyFlag         = "technology"
	configDirFlag          = "config-directory"
	attributeFlag          = "attribute"
	attributeContainerFlag = "attribute-container"
)

var (
	sourceFolder        string
	targetFolder        string
	workFolder          string
	technology          string
	configDirectory     string
	attribute           []string
	attributeContainers []string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&sourceFolder, sourceFolderFlag, "", "Base path where to copy the codemodule FROM.")
	_ = cmd.MarkPersistentFlagRequired(sourceFolderFlag)

	cmd.PersistentFlags().StringVar(&targetFolder, targetFolderFlag, "", "Base path where to copy the codemodule TO.")
	_ = cmd.MarkPersistentFlagRequired(targetFolderFlag)

	cmd.PersistentFlags().StringVar(&workFolder, workFolderFlag, "", "(Optional) Base path for a tmp folder, this is where the command will do its work, to make sure the operations are atomic. It must be on the same disk as the target folder.")

	cmd.PersistentFlags().StringVar(&technology, technologyFlag, "", "(Optional) Comma-separated list of technologies to filter files.")

	cmd.PersistentFlags().StringVar(&configDirectory, configDirFlag, "", "(Optional) Path where enrichment/configuration files will be written.")

	cmd.PersistentFlags().StringArrayVar(&attribute, attributeFlag, []string{}, "(Optional) Pod-specific attributes in key=value format.")

	cmd.PersistentFlags().StringArrayVar(&attributeContainers, attributeContainerFlag, []string{}, "(Optional) Container-specific attributes in JSON format.")

}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(fs afero.Afero) error {
	if configDirectory != "" {
		err := writeConfigFiles(fs)
		if err != nil {
			return err
		}
	}

	copy := simpleCopy

	if technology != "" {
		copy = copyByTechnology
	}

	if workFolder != "" {
		copy = atomic(workFolder, copy)
	}

	return copy(fs, sourceFolder, targetFolder)
}
