package gearbox

import (
	"gearbox/util"
	"log"
	"os"
	"path/filepath"
)

var rootDir string
var execDir string

func init() {
	executable := os.Args[0]
	file, err := filepath.Abs(executable)
	if err != nil {
		log.Fatal(err)
	}
	execDir = util.FileDir(file)
	rootDir = util.ParentDir(execDir)
}
