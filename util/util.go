package util

import (
	"os"
)

func EntryExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
