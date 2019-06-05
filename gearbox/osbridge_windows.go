// +build windows

package gearbox

import "github.com/gearboxworks/go-osbridge"

func GetOsBridge(project types.Name, userdata types.RelativePath) *osbridge.OsBridge {
	return osbridge.NewOsBridge(&osbridge.Args{
		ProjectName:  project,
		UserDataPath: userdata,
		AdminPath:    WindowsAdminPath,
		ProjectDir:   WindowsSuggestedBasedir,
	})
}
