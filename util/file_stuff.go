package util

import (
	"fmt"
	"gearbox/only"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func FileExists(file string) bool {
	return EntryExists(file)
}
func GetExecutableFilepath() string {
	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
//func GetProjectDir() string {
//	return filepath.Dir(GetExecutableDir())
//}
func ErrorIsFileDoesNotExist(err error) bool {
	pe, ok := err.(*os.PathError)
	return ok && pe.Op == "open" && pe.Err == syscall.ENOENT
}

func ReadBytes(filepath string) (b []byte, err error) {
	for range only.Once {
		var err error
		b, err = ioutil.ReadFile(filepath)
		if err != nil && ErrorIsFileDoesNotExist(err) {
			err = nil
		}
		if err != nil {
			err = AddHelpToError(
				fmt.Errorf("cannot read from '%s' file", filepath),
				fmt.Sprintf("confirm file '%s' is readable", filepath),
			)
			break
		}
	}
	return b, err
}
