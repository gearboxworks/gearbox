package ospaths

import (
	"errors"
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"os"
)

func FileExists(f string) (os.FileInfo, error) {

	var err error
	var stat os.FileInfo

	for range only.Once {
		if f == "" {
			err = errors.New("file is nil")
			break
		}

		stat, err = os.Stat(f)
		if os.IsNotExist(err) {
			err = errors.New(fmt.Sprintf("file %s not found with error '%v'", f, err))
			break
		}

		if err != nil {
			err = errors.New(fmt.Sprintf("file %s access failed with '%v'", f, err))
			break
		}

		if stat == nil {
			err = errors.New(fmt.Sprintf("file %s not found", f))
			break
		}

		if stat.IsDir() {
			err = errors.New(fmt.Sprintf("file %s is a directory", f))
			break
		}
	}

	return stat, err
}

func FileSetExecutePerms(f string) (os.FileInfo, error) {

	var err error
	var stat os.FileInfo

	for range only.Once {
		// Read in current permissions so we can add +x.
		stat, err = FileExists(f)
		if err != nil {
			break
		}

		// Set +x permissions.
		err = os.Chmod(f, stat.Mode()|OS_ALL_X)
		if err != nil {
			break
		}

		// Check again for changes.
		stat, err = FileExists(f)
		if err != nil {
			break
		}
		if (stat.Mode() & OS_ALL_X) != OS_ALL_X {
			err = errors.New(fmt.Sprintf("file %s permissions couldn't be changed", f))
			break
		}
	}

	return stat, err
}
