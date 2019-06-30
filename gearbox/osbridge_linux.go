// +build linux

package gearbox

import (
	"gearbox/global"
	"github.com/gearboxworks/go-osbridge"
)

func GetOsBridge(project types.Name, userdata types.RelativePath) *osbridge.OsBridge {
	return osbridge.NewOsBridge(&osbridge.Args{
		ProjectName:  project,
		UserDataPath: userdata,
		AdminPath:    global.NixAdminPath,
		ProjectDir:   global.LinuxSuggestedBasedir,
	})
}
