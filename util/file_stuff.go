package util

import (
	"fmt"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func FileExists(file types.Filepath) bool {
	return EntryExists(types.FileSystemEntry(file))
}
func GetExecutableFilepath() types.Filepath {
	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	return types.Filepath(fp)
}

//func GetProjectDir() string {
//	return util.FileDir(GetExecutableDir())
//}
func ErrorIsFileDoesNotExist(err error) bool {
	pe, ok := err.(*os.PathError)
	return ok && pe.Op == "open" && pe.Err == syscall.ENOENT
}

func ReadBytes(filepath types.Filepath) (b []byte, sts status.Status) {
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
