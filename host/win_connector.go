package host

import (
	"fmt"
)

const winSuggestedProjectsPath = "Gearbox Sites"

var WinConnectorInstance = (*WinConnector)(nil)

var _ Connector = WinConnectorInstance

type WinConnector struct {
	BaseConnector
}

func (me *WinConnector) GetSuggestedProjectRoot() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		winSuggestedProjectsPath,
	)
}

func (me *WinConnector) GetUserConfigDir() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		UserDataPath,
	)
}
func (me *WinConnector) GetAdminRootDir() string {
	return fmt.Sprintf("%s\\admin\\dist",
		me.GetUserConfigDir(),
	)
}
