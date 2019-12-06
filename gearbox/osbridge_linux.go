// +build linux

package gearbox

import (
	"gearbox/global"
	"gearbox/types"
	"github.com/gearboxworks/go-osbridge"
)

func GetOsBridge(project types.Name, userdata types.Path) *osbridge.OsBridge {
	return osbridge.NewOsBridge(&osbridge.Args{
		ProjectName:  project,
		UserDataPath: userdata,
		AdminPath:    global.NixAdminPath,
		ProjectDir:   global.LinuxSuggestedBasedir,
	})
}
