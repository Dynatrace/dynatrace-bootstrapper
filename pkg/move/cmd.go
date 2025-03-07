package move

import (
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	SourceFolderFlag = "source"
	TargetFolderFlag = "target"
	WorkFolderFlag   = "work"
	TechnologyFlag   = "technology"
)

var (
	workFolder string
	technology string

	// SourceFolder holds the value defined in the --source flag, only should be used if AddFlags was also used and part of a cobra.Command
	SourceFolder string
	// TargetFolder holds the value defined in the --target flag, only should be used if AddFlags was also used and part of a cobra.Command
	TargetFolder string
)

func AddFlags(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&SourceFolder, SourceFolderFlag, "", "Base path where to copy the codemodule FROM.")
	_ = cmd.MarkPersistentFlagRequired(SourceFolderFlag)

	cmd.PersistentFlags().StringVar(&TargetFolder, TargetFolderFlag, "", "Base path where to copy the codemodule TO.")
	_ = cmd.MarkPersistentFlagRequired(TargetFolderFlag)

	cmd.PersistentFlags().StringVar(&workFolder, WorkFolderFlag, "", "(Optional) Base path for a tmp folder, this is where the command will do its work, to make sure the operations are atomic. It must be on the same disk as the target folder.")

	cmd.PersistentFlags().StringVar(&technology, TechnologyFlag, "", "(Optional) Comma-separated list of technologies to filter files.")

}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(log logr.Logger, fs afero.Afero, from, to string) error {
	copy := simpleCopy

	if technology != "" {
		copy = copyByTechnology
	}

	if workFolder != "" {
		copy = atomic(workFolder, copy)
	}

	return copy(log, fs, from, to)
}
