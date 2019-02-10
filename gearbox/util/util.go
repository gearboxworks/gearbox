package util

import (
	"os"
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
