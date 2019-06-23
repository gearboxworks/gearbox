// +build darwin

package gearbox

import (
	"gearbox/types"
	"github.com/gearboxworks/go-osbridge"
)

func GetOsBridge(project types.Name, userdata types.Path) *osbridge.OsBridge {
	return osbridge.NewOsBridge(&osbridge.Args{
		ProjectName:  project,
		UserDataPath: userdata,
		AdminPath:    NixAdminPath,
		ProjectDir:   MacOsSuggestedBasedir,
	})
}
