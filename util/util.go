package util

import (
	"gearbox/types"
	"os"
)

func EntryExists(file types.AbsoluteEntry) bool {
	_, err := os.Stat(string(file))
	return !os.IsNotExist(err)
}
