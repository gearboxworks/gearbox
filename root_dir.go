package gearbox

import (
	"gearbox/types"
	"gearbox/util"
	"log"
	"os"
	"path/filepath"
)

var rootDir string

func init() {
	executable := types.AbsoluteFilepath(os.Args[0])
	file, err := filepath.Abs(string(executable))
	if err != nil {
		log.Fatal(err)
	}
	rootDir = string(util.FileDir(types.AbsoluteFilepath(file)))
}
