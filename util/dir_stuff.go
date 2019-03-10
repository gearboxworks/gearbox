package util

import (
	"os"
)

func DirExists(dir string) bool {
	return EntryExists(dir)
}
func MaybeMakeDir(dir string, perms os.FileMode) (err error) {
	if !DirExists(dir) {
		err = os.MkdirAll(dir, perms)
	}
	return err
}
//func GetExecutableDir() string {
//	return filepath.Dir(GetExecutableFilepath())
//}
//func GetProjectDir() string {
//	return filepath.Dir(GetExecutableDir())
//}
