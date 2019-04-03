package mock

import (
	"fmt"
	"gearbox/os_support"
	"gearbox/types"
	"gearbox/util"
	"os"
	"testing"
)

const (
	SuggestedBasedir = "sites"
	AdminPath        = "admin/dist"
)

var NilOsSupport = (*OsSupport)(nil)

var _ oss.OsSupporter = NilOsSupport

type OsSupport struct {
	oss.Base
	T                 *testing.T
	UserHomePath      string
	SuggestedBasePath string
	UserConfigPath    string
	AdminRootPath     string
	CachePath         string
}

func NewOsSupport(t *testing.T) oss.OsSupporter {
	return &OsSupport{
		T: t,
	}
}

func (me *OsSupport) GetUserHomeDir() types.AbsoluteDir {
	if me.UserHomePath == "" {
		me.UserHomePath = "user-home"
	}
	dir, err := os.Getwd()
	if err != nil {
		me.T.Error(fmt.Sprintf("failed to get current working directory: %s", err.Error()))
	}
	absdir := util.ParentDir(types.AbsoluteDir(dir))
	return types.AbsoluteDir(fmt.Sprintf("%s/%s", absdir, me.UserHomePath))
}

func (me *OsSupport) GetSuggestedBasedir() types.AbsoluteDir {
	if me.SuggestedBasePath == "" {
		me.SuggestedBasePath = SuggestedBasedir
	}
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		me.SuggestedBasePath,
	))
}

func (me *OsSupport) GetUserConfigDir() types.AbsoluteDir {
	if me.UserConfigPath == "" {
		me.UserConfigPath = oss.UserDataPath
	}
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		me.UserConfigPath,
	))
}

func (me *OsSupport) GetAdminRootDir() types.AbsoluteDir {
	if me.AdminRootPath == "" {
		me.AdminRootPath = AdminPath
	}
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		me.AdminRootPath,
	))
}

func (me *OsSupport) GetCacheDir() types.AbsoluteDir {
	if me.CachePath == "" {
		me.CachePath = oss.CachePath
	}
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		me.CachePath,
	))
}
