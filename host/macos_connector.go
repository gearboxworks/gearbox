package host

import (
	"fmt"
)

const macOsSuggestedProjectsPath = "Sites"

var MacOsConnectorInstance = (*MacOsConnector)(nil)

var _ Connector = MacOsConnectorInstance

type MacOsConnector struct {
	NixConnector
}

func (me *MacOsConnector) GetSuggestedProjectRoot() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		macOsSuggestedProjectsPath,
	)
}
