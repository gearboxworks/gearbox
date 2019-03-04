package util

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func EntryExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
func FileExists(file string) bool {
	return EntryExists(file)
}
func DirExists(dir string) bool {
	return EntryExists(dir)
}
func GetExecutableFilepath() string {
	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
func GetExecutableDir() string {
	return filepath.Dir(GetExecutableFilepath())
}
func GetProjectDir() string {
	return filepath.Dir(GetExecutableDir())
}
func ErrorIsFileDoesNotExist(err error) bool {
	pe, ok := err.(*os.PathError)
	return ok && pe.Op == "open" && pe.Err == syscall.ENOENT
}
func MaybeMakeDir(dir string, perms os.FileMode) (err error) {
	if !DirExists(dir) {
		err = os.MkdirAll(dir, perms)
	}
	return err
}
