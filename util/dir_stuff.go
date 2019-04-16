package util

import (
	"gearbox/types"
	"os"
	"path/filepath"
	"strings"
)

func DirExists(dir types.AbsoluteDir) bool {
	return EntryExists(types.AbsoluteEntry(dir))
}
func MaybeMakeDir(dir types.AbsoluteDir, perms os.FileMode) (err error) {
	if !DirExists(dir) {
		err = os.MkdirAll(string(dir), perms)
	}
	return err
}
func FileDir(file types.AbsoluteFilepath) types.AbsoluteDir {
	return types.AbsoluteDir(filepath.Dir(string(file)))
}

func ParentDir(file types.AbsoluteDir) types.AbsoluteDir {
	return types.AbsoluteDir(filepath.Dir(string(file)))
}

//func GetExecutableDir() string {
//	return util.FileDir(GetExecutableFilepath())
//}
//func GetProjectDir() string {
//	return util.FileDir(GetExecutableDir())
//}

func ExtractRelativePath(fulldir types.AbsoluteFilepath, basedir types.AbsoluteDir) (path types.RelativePath) {
	if strings.HasPrefix(string(fulldir), string(basedir)) {
		path = types.RelativePath(string([]byte(fulldir)[len(basedir):]))
	} else {
		path = types.RelativePath(fulldir)
	}
	return path
}
