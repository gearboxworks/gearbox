package gearbox

import (
	"log"
	"os"
	"path/filepath"
)

var rootDir string

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	rootDir = filepath.Dir(dir)
}
