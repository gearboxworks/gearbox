package gbMqttClient

import (
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"net/url"
)

type Client struct {
	OsSupport  oss.OsSupporter
	Sts        status.Status

	Server url.URL
}
type Args Client
