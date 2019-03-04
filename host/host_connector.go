package host

import (
	"fmt"
	"runtime"
)

type Connector interface {
	GetAdminPath() string
	GetAdminRootDir() string
	GetUserConfigDir() string
	GetSuggestedBaseDir() string
	GetUserHomeDir() string
	GetCacheDir() string
}

func GetConnector() Connector {
	var host Connector
	switch runtime.GOOS {
	case "darwin":
		host = &MacOsConnector{}
	case "windows":
		host = &WinConnector{}
	case "linux":
		host = &LinuxConnector{}
	default:
		msg := "Sadly, Gearbox does not currently run on '%s.'\nIf you would like to offer us support to change that please contact us via https://gearbox.works.\n"
		fmt.Printf(msg, runtime.GOOS)
	}
	return host
}
