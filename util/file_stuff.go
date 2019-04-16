package util

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func FileExists(file types.AbsoluteFilepath) bool {
	return EntryExists(types.AbsoluteEntry(file))
}
func GetExecutableFilepath() types.AbsoluteFilepath {
	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	return types.AbsoluteFilepath(fp)
}

//func GetProjectDir() string {
//	return util.FileDir(GetExecutableDir())
//}
func ErrorIsFileDoesNotExist(err error) bool {
	pe, ok := err.(*os.PathError)
	return ok && pe.Op == "open" && pe.Err == syscall.ENOENT
}

func ReadBytes(filepath types.AbsoluteFilepath) (b []byte, sts status.Status) {
	for range only.Once {
		var err error
		b, err = ioutil.ReadFile(string(filepath))
		if err != nil && ErrorIsFileDoesNotExist(err) {
			sts = status.Success("read %d bytes from '%s'",
				len(b),
				filepath,
			)
		}
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("cannot read from '%s' file", filepath),
				Help:    fmt.Sprintf("confirm file '%s' is readable", filepath),
			})
			break
		}
	}
	return b, sts
}
