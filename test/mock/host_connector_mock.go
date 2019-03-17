package mock

import (
	"fmt"
	"gearbox/host"
	"os"
	"path/filepath"
	"testing"
)

const testSuggestedBasedir = "sites"

var HostConnectorInstance = (*HostConnector)(nil)

var _ host.Connector = HostConnectorInstance

type HostConnector struct {
	host.BaseConnector
	T                 *testing.T
	UserHomePath      string
	SuggestedBasePath string
	UserConfigPath    string
	AdminRootPath     string
	CachePath         string
}

func NewHostConnector(t *testing.T) host.Connector {
	return &HostConnector{
		T: t,
	}
}

func (me *HostConnector) GetUserHomeDir() string {
	if me.UserHomePath == "" {
		me.UserHomePath = "user-home"
	}
	dir, err := os.Getwd()
	if err != nil {
		me.T.Error(fmt.Sprintf("failed to get current working directory: %s", err.Error()))
	}
	dir = filepath.Dir(dir)
	return fmt.Sprintf("%s/%s", dir, me.UserHomePath)
}

func (me *HostConnector) GetSuggestedBasedir() string {
	if me.SuggestedBasePath == "" {
		me.SuggestedBasePath = testSuggestedBasedir
	}
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		me.SuggestedBasePath,
	)
}

func (me *HostConnector) GetUserConfigDir() string {
	if me.UserConfigPath == "" {
		me.UserConfigPath = host.UserDataPath
	}
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		me.UserConfigPath,
	)
}

func (me *HostConnector) GetAdminRootDir() string {
	if me.AdminRootPath == "" {
		me.AdminRootPath = me.GetAdminPath()
	}
	return fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		me.AdminRootPath,
	)
}

func (me *HostConnector) GetCacheDir() string {
	if me.CachePath == "" {
		me.CachePath = host.CachePath
	}
	return fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		me.CachePath,
	)
}
