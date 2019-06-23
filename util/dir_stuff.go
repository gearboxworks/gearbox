package util

import (
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"
)

type (
	Status           = status.Status
	AbsoluteDir      = types.Dir
	AbsoluteEntry    = types.FileSystemEntry
	AbsoluteFilepath = types.Filepath
)

func DirExists(dir types.Dir) bool {
	return EntryExists(types.FileSystemEntry(dir))
}
func MaybeMakeDir(dir types.Dir, perms os.FileMode) (err error) {
	if !DirExists(dir) {
		err = os.MkdirAll(string(dir), perms)
	}
	return err
}
func FileDir(file types.Filepath) types.Dir {
	return types.Dir(filepath.Dir(string(file)))
}
func ParentDir(file types.Dir) types.Dir {
	return types.Dir(filepath.Dir(string(file)))
}

//func GetExecutableDir() string {
//	return util.FileDir(GetExecutableFilepath())
//}
//func GetProjectDir() string {
//	return util.FileDir(GetExecutableDir())
//}

func ExtractRelativePath(fulldir types.Filepath, basedir types.Dir) (path types.Path) {
	if strings.HasPrefix(string(fulldir), string(basedir)) {
		path = types.Path(string([]byte(fulldir)[len(basedir):]))
	} else {
		path = types.Path(fulldir)
	}
	return path
}

func MaybeExpandFilepath(fp AbsoluteFilepath) (nd AbsoluteFilepath, sts Status) {
	var e AbsoluteEntry
	e, sts = MaybeExpandEntry(AbsoluteEntry(fp))
	return AbsoluteFilepath(e), sts
}

func MaybeExpandDir(dir AbsoluteDir) (nd AbsoluteDir, sts Status) {
	var e AbsoluteEntry
	e, sts = MaybeExpandEntry(AbsoluteEntry(dir))
	return AbsoluteDir(e), sts
}

func MaybeExpandEntry(entry types.FileSystemEntry) (ne AbsoluteEntry, sts Status) {
	for range only.Once {
		ne = entry
		if !strings.HasPrefix(string(entry), "~") {
			break
		}
		newentry, err := homedir.Expand(string(entry))
		if err != nil {
			sts = status.Wrap(err).SetMessage("could not expand entry '%s': %s",
				entry,
				err.Error(),
			)
			break
		}
		ne = types.FileSystemEntry(newentry)
	}
	if is.Success(sts) {
		sts = status.Success("directory expanded from '%s' to '%s'",
			entry,
			ne,
		)
	}
	return ne, sts

}
