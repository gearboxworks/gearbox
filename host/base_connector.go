package host

import (
	"github.com/mitchellh/go-homedir"
	"runtime"
)

const UserDataPath = ".gearbox"

const CachePath = "cache"

type BaseConnector struct{}

func (me *BaseConnector) GetAdminPath() string {
	if runtime.GOOS == "windows" {
		return "admin\\dist"
	}
	return "admin/dist"
}

func (me *BaseConnector) GetUserHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
