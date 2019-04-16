package gearbox

import (
	"gearbox/types"
	"gearbox/util"
	"log"
	"os"
	"path/filepath"
)

var rootDir string
var execDir string

func init() {
	executable := types.AbsoluteFilepath(os.Args[0])
	file, err := filepath.Abs(string(executable))
	if err != nil {
		log.Fatal(err)
	}
	execDir = string(util.FileDir(types.AbsoluteFilepath(file)))
	rootDir = string(util.ParentDir(types.AbsoluteDir(execDir)))
}
