package network

import (
	"time"
)

const (
	// DefaultEntityId = "eventbroker-zeroconf"
	DefaultWaitTime   = time.Millisecond * 1000
	DefaultDomain     = "local"
	DefaultRetries    = 12
	DefaultRetryDelay = time.Second * 8

	//SelectRandomPort = "0"
)

const (
	Package               = "network"
	InterfaceTypeZeroConf = "*" + Package + ".ZeroConf"
	InterfaceTypeService  = "*" + Package + ".Service"
	//InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
)
