package jsonfile

import (
	"fmt"
	"gearbox/types"
	"path/filepath"
)

func GetFilepath(basedir types.AbsoluteDir, path types.RelativePath) types.AbsoluteFilepath {
	return types.AbsoluteFilepath(filepath.FromSlash(fmt.Sprintf("%s/%s/%s",
		basedir,
		path,
		BaseFilename,
	)))
}
