// +build windows

package oss

import (
	"fmt"
	"gearbox/types"
)

const (
	AdminPath        types.RelativePath = "admin\\dist"
	SuggestedBasedir types.AbsoluteDir  = "Parent Sites"
)

var NilOsSupport = (*OsSupport)(nil)

var _ OsSupporter = NilOsSupport

type OsSupport struct {
	Base
}

func (me *OsSupport) GetSuggestedBasedir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		SuggestedBasedir,
	))
}

func (me *OsSupport) GetUserConfigDir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		UserDataPath,
	))
}
func (me *OsSupport) GetAdminRootDir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s\\%s",
		me.GetUserConfigDir(),
		AdminPath,
	))
}
func (me *OsSupport) GetCacheDir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s\\%s",
		me.GetUserConfigDir(),
		CachePath,
	))
}
