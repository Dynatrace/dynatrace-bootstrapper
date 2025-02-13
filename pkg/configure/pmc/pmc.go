package pmc

import (
	"github.com/spf13/afero"
)

const (

)

func Configure(fs afero.Afero, inputDir, configDir, targetDir string) error {
	// TODO: Ipml + decide if we merge or overwrite what is in the targetDir
	return nil
}