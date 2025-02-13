package move

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	workFolderFlag = "work"
	technologyFlag = "technology"
)

var (
	workFolder string
	technology string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&workFolder, workFolderFlag, "", "(Optional) Base path for a tmp folder, this is where the command will do its work, to make sure the operations are atomic. It must be on the same disk as the target folder.")

	cmd.PersistentFlags().StringVar(&technology, technologyFlag, "", "(Optional) Comma-separated list of technologies to filter files.")

}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(fs afero.Afero, from, to string) error {
	copy := simpleCopy

	if technology != "" {
		copy = copyByTechnology
	}

	if workFolder != "" {
		copy = atomic(workFolder, copy)
	}

	return copy(fs, from, to)
}
