package host

import (
	"fmt"
)

const winSuggestedBasedir = "Gearbox Sites"

var WinConnectorInstance = (*WinConnector)(nil)

var _ Connector = WinConnectorInstance

type WinConnector struct {
	BaseConnector
}

func (me *WinConnector) GetSuggestedBasedir() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		winSuggestedBasedir,
	)
}

func (me *WinConnector) GetUserConfigDir() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		UserDataPath,
	)
}
func (me *WinConnector) GetAdminRootDir() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserConfigDir(),
		me.GetAdminPath(),
	)
}
func (me *WinConnector) GetCacheDir() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserConfigDir(),
		CachePath,
	)
}
