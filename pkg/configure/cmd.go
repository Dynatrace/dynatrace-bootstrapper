package configure

import (
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/curl"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	inputFolderFlag  = "input-directory"
	configFolderFlag = "config-directory"
)

var (
	inputFolder  string
	configFolder string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&inputFolder, inputFolderFlag, "", "Base path where to look for the configuration files.")

	cmd.PersistentFlags().StringVar(&configFolder, configFolderFlag, "", "Base path where to put the configuration files.")
}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(fs afero.Afero) error {
	if configFolder == "" || inputFolder == "" {
		return nil
	}

	curl.Configure(fs, inputFolder, configFolder)

	return nil
}
