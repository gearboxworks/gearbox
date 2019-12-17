// +build windows

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
		AdminPath:    global.WindowsAdminPath,
		ProjectDir:   global.WindowsSuggestedBasedir,
	})
}
