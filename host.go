package gearbox

import (
	"gearbox/host"
)

var Host host.Connector

func init() {
	Host = host.GetConnector()
}
