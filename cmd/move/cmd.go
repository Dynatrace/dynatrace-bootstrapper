package move

import (
	impl "github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	WorkFolderFlag = "work"
	TechnologyFlag = "technology"
)

var (
	workFolder string
	technology string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&workFolder, WorkFolderFlag, "", "(Optional) Base path for a tmp folder, this is where the command will do its work, to make sure the operations are atomic. It must be on the same disk as the target folder.")

	cmd.PersistentFlags().StringVar(&technology, TechnologyFlag, "", "(Optional) Comma-separated list of technologies to filter files.")

}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(log logr.Logger, fs afero.Afero, from, to string) error {
	copyFunc := impl.SimpleCopy

	if technology != "" {
		copyFunc = impl.CopyByTechnologyWrapper(technology)
	}

	if workFolder != "" {
		copyFunc = impl.Atomic(workFolder, copyFunc)
	}

	err := copyFunc(log, fs, from, to)
	if err != nil {
		return err
	}

	return impl.CreateCurrentSymlink(log, fs.Fs, to)
}
