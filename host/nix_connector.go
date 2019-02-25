package host

import (
	"fmt"
)

type NixConnector struct {
	BaseConnector
}

func (me *NixConnector) GetUserConfigDir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		UserDataPath,
	)
}

func (me *NixConnector) GetAdminRootDir() string {
	return fmt.Sprintf("%s/admin/dist",
		me.GetUserConfigDir(),
	)
}
