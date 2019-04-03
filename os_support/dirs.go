package oss

import (
	"fmt"
	"gearbox/types"
	"runtime"
)

type OsSupporter interface {
	GetAdminRootDir() types.AbsoluteDir
	GetUserConfigDir() types.AbsoluteDir
	GetSuggestedBasedir() types.AbsoluteDir
	GetUserHomeDir() types.AbsoluteDir
	GetCacheDir() types.AbsoluteDir
}

func Get() OsSupporter {
	switch runtime.GOOS {
	case "darwin":
	case "windows":
	case "linux":
	default:
		msg := "Sadly, Gearbox does not currently run on '%s.'\nIf you would like to offer us support to change that please contact us via https://gearbox.works.\n"
		fmt.Printf(msg, runtime.GOOS)
	}
	return &OsSupport{}
}
