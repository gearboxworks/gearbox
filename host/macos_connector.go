package host

import (
	"fmt"
)

const macOsSuggestedBasedir = "Sites"

var MacOsConnectorInstance = (*MacOsConnector)(nil)

var _ Connector = MacOsConnectorInstance

type MacOsConnector struct {
	NixConnector
}

func (me *MacOsConnector) GetSuggestedBasedir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		macOsSuggestedBasedir,
	)
}
