package host

import (
	"fmt"
)

const winSuggestedBaseDir = "Gearbox Sites"

var WinConnectorInstance = (*WinConnector)(nil)

var _ Connector = WinConnectorInstance

type WinConnector struct {
	BaseConnector
}

func (me *WinConnector) GetSuggestedBaseDir() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		winSuggestedBaseDir,
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
