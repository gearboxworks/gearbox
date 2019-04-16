package mock

import (
	"fmt"
	"gearbox/os_support"
	"gearbox/test/user-home"
	"gearbox/types"
	"strings"
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
	if strings.HasPrefix(string(testconst.UserHomeDir), "ERROR:") {
		me.T.Error(fmt.Sprintf("failed to get current working directory: %s", testconst.UserHomeDir))
	}
	return testconst.UserHomeDir
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
	dir := types.AbsoluteDir(fmt.Sprintf("%s/%s",
		testconst.UserHomeDir,
		me.UserConfigPath,
	))
	return dir
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
