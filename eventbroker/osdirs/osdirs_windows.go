// +build windows

package osdirs

import (
	"github.com/gearboxworks/go-osbridge"
)

func GetOsBridge(project Name, userdata Dir) *osbridge.OsBridge {
	return osbridge.NewOsBridge(&osbridge.Args{
		ProjectName:  osbridge.Name(project),
		UserDataPath: osbridge.Path(userdata),
		AdminPath:    string(WindowsAdminPath),
		ProjectDir:   string(WindowsProjectBaseDir),
	})
}
