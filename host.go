package gearbox

import (
	"fmt"
	"gearbox/host"
	"runtime"
)

var Host host.Connector

func init() {
	switch runtime.GOOS {
	case "darwin":
		Host = &host.MacOsConnector{}
	case "windows":
		Host = &host.WinConnector{}
	case "linux":
		Host = &host.LinuxConnector{}
	default:
		msg := "Sadly, Gearbox does not currently run on '%s.'\nIf you would like to offer us support to change that please contact us via https://gearbox.works.\n"
		fmt.Printf(msg, runtime.GOOS)
	}

}
