package oss

import (
	"gearbox/types"
	"github.com/mitchellh/go-homedir"
)

const (
	UserDataPath = ".gearbox"
	CachePath    = "cache"
)

type Base struct{}

func (me *Base) GetUserHomeDir() types.AbsoluteDir {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return types.AbsoluteDir(homeDir)
}
