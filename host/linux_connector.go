package host

import (
	"fmt"
)

const linuxSuggestedBaseDir = "projects"

var LinuxConnectorInstance = (*LinuxConnector)(nil)

var _ Connector = LinuxConnectorInstance

type LinuxConnector struct {
	NixConnector
}

func (me *LinuxConnector) GetSuggestedBaseDir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		linuxSuggestedBaseDir,
	)
}
