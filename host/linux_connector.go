package host

import (
	"fmt"
)

const linuxSuggestedBasedir = "projects"

var LinuxConnectorInstance = (*LinuxConnector)(nil)

var _ Connector = LinuxConnectorInstance

type LinuxConnector struct {
	NixConnector
}

func (me *LinuxConnector) GetSuggestedBasedir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		linuxSuggestedBasedir,
	)
}
