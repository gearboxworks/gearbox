package util

import (
	"log"
	"os"
	"path/filepath"
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
