package osdirs

import (
	"errors"
	"fmt"
	"gearbox/eventbroker/fileperms"
	"github.com/gearboxworks/go-status/only"
	"os"
	"path/filepath"
	"strings"
)

func AddPaths(dir Dir, path ...interface{}) Dir {
	format := strings.Repeat("/%s", len(path))
	s := fmt.Sprintf("%s%s", dir, fmt.Sprintf(format, path...))
	return filepath.FromSlash(s)
}

func AddFilef(dir Dir, format string, fn ...interface{}) File {
	s := fmt.Sprintf("%s/%s", dir, fmt.Sprintf(format, fn...))
	return filepath.FromSlash(s)
}

func DirExists(dir Dir) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

func CreateIfNotExists(dir Dir) (created bool, err error) {
	for range only.Once {
		if DirExists(dir) {
			break
		}

		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("could not create directory '%s'", dir)
			break
		}
		created = true

	}
	return created, err
}

func CheckFileExists(f File) error {

	var err error

	if f == "" {
		err = errors.New("filename is empty")
		return err
	}

	_, err = os.Stat(f)
	if os.IsNotExist(err) {
		//fmt.Printf("Not exists PATH: '%s'\n", f.String())
	}

	return err
}

func FileDelete(f File) (err error) {
	_, err = os.Stat(f)
	if !os.IsNotExist(err) {
		err = os.Remove(f)
	}
	return err
}

func Split(fn string) *Path {

	var pn Path

	d, f := filepath.Split(fn)
	pn.Dir = d
	pn.File = f

	return &pn
}

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
			err = errors.New(fmt.Sprintf("file '%s' not found with error '%v'", f, err))
			break
		}

		if err != nil {
			err = errors.New(fmt.Sprintf("file '%s' access failed with '%v'", f, err))
			break
		}

		if stat == nil {
			err = errors.New(fmt.Sprintf("file '%s' not found", f))
			break
		}

		if stat.IsDir() {
			err = errors.New(fmt.Sprintf("file '%s' is a directory", f))
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
		err = os.Chmod(f, stat.Mode()|fileperms.AllExecute)
		if err != nil {
			break
		}

		// Check again for changes.
		stat, err = FileExists(f)
		if err != nil {
			break
		}

		if (stat.Mode() & fileperms.AllExecute) != fileperms.AllExecute {
			err = errors.New(fmt.Sprintf("file '%s' permissions couldn't be changed", f))
			break
		}
	}

	return stat, err
}
