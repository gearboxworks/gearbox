package host

import (
	"github.com/mitchellh/go-homedir"
)

const UserDataPath = ".gearbox"

type BaseConnector struct{}

func (me *BaseConnector) GetUserHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
