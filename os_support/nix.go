// +build !windows

package oss

import (
	"fmt"
	"gearbox/types"
)

const AdminPath = "admin/dist"

type Nix struct {
	Base
}

func (me *Nix) GetUserConfigDir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		UserDataPath,
	))
}

func (me *Nix) GetAdminRootDir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		AdminPath,
	))
}
func (me *Nix) GetCacheDir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		CachePath,
	))
}
