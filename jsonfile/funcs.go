package jsonfile

import (
	"fmt"
	"gearbox/types"
	"path/filepath"
)

func GetFilepath(basedir types.Dir, path types.Path) types.Filepath {
	return types.Filepath(filepath.FromSlash(fmt.Sprintf("%s/%s/%s",
		basedir,
		path,
		BaseFilename,
	)))
}
