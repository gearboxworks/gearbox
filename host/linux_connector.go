package host

import (
	"fmt"
)

const linuxSuggestedProjectsPath = "projects"

var LinuxConnectorInstance = (*LinuxConnector)(nil)

var _ Connector = LinuxConnectorInstance

type LinuxConnector struct {
	NixConnector
}

func (me *LinuxConnector) GetSuggestedProjectRoot() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		linuxSuggestedProjectsPath,
	)
}
