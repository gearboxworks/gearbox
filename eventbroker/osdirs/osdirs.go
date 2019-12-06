// A simple wrapper around osbridge.OsBridger.
// This makes it much easier to separate the EventBroker code into it's own package later on.
package osdirs

import (
	"gearbox/global"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status/only"
	"github.com/getlantern/errors"
	"path/filepath"
	"sync"
)

type BaseDirs struct {
	LocalDir              Dir
	EventBrokerDir        Dir
	EventBrokerWorkingDir Dir
	EventBrokerLogDir     Dir
	EventBrokerEtcDir     Dir

	osbridge.OsBridger
	mutex sync.RWMutex
}

func (me *BaseDirs) GetUserConfigDir() string {
	return me.OsBridger.GetUserConfigDir()
}

// I don't like this naming convention, but *Dir not being
// a string I can't come up with a cleaner approach.
func (me *BaseDirs) UserConfigDir() Dir {
	return me.GetUserConfigDir()

}

func New(subdir ...string) *BaseDirs {

	var ret BaseDirs

	if len(subdir) == 0 || subdir[0] == "" {
		subdir = []string{DefaultBaseDir}
	}

	b := GetOsBridge(global.Brandname, global.UserDataPath)

	ret.OsBridger = b

	ret.LocalDir = filepath.FromSlash(defaultLocalDir)

	ebDir := ret.AppendToUserConfigDir(subdir[0])
	ret.EventBrokerLogDir = AddPaths(ebDir, defaultLogBaseDir)
	ret.EventBrokerEtcDir = AddPaths(ebDir, defaultEtcBaseDir)
	ret.EventBrokerWorkingDir = ebDir
	ret.EventBrokerDir = ebDir

	return &ret
}

func (me *BaseDirs) AppendToUserConfigDir(p ...interface{}) Dir {
	return AddPaths(me.GetUserConfigDir(), p...)
}

func (me *BaseDirs) AddFileToUserConfigDir(format string, fn ...interface{}) File {
	return AddFilef(me.GetUserConfigDir(), format, fn...)
}

func (me *BaseDirs) EnsureNotNil() (err error) {
	if me == nil {
		err = errors.New("basedirs is nil")
	}
	return err
}

func (me *BaseDirs) CreateIfNotExists() (err error) {

	for range only.Once {
		_, err = CreateIfNotExists(me.EventBrokerDir)
		if err != nil {
			break
		}

		_, err = CreateIfNotExists(me.EventBrokerEtcDir)
		if err != nil {
			break
		}

		_, err = CreateIfNotExists(me.EventBrokerLogDir)
		if err != nil {
			break
		}

		_, err = CreateIfNotExists(me.EventBrokerWorkingDir)
		if err != nil {
			break
		}
	}

	return err
}
