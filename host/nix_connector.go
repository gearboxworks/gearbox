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
	return fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		me.GetAdminPath(),
	)
}
func (me *NixConnector) GetCacheDir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserConfigDir(),
		CachePath,
	)
}
