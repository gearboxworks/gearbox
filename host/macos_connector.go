package host

import (
	"fmt"
)

const macOsSuggestedBaseDir = "Sites"

var MacOsConnectorInstance = (*MacOsConnector)(nil)

var _ Connector = MacOsConnectorInstance

type MacOsConnector struct {
	NixConnector
}

func (me *MacOsConnector) GetSuggestedBaseDir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		macOsSuggestedBaseDir,
	)
}
