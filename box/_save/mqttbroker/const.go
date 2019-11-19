package mqttBroker

import (
	"gearbox/eventbroker/messages"
	//	oss "gearbox/os_support"
	"time"
)

const (
	DefaultEntityId = "eventbroker-mqttbroker"
	defaultWaitTime = time.Millisecond * 2000
	defaultDomain   = "local"
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 3
	DefaultServer = "tcp://127.0.0.1:1883"
)

type mqttBroker struct {
	EntityId        messages.MessageAddress

	restartAttempts int
	waitTime        time.Duration
	domain          string
	osSupport       osbridge.OsBridger
}
type Args mqttBroker


const mqttBrokerJson = `
{
	"Name": "com.gearbox.mqtt",
	"DisplayName": "Gearbox [UNFSD]",
	"Description": "UNFSD is required to support Gearbox",

	"Url": "tcp://{{.Host}}:{{.Port}}",

	"WorkingDirectory": "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/unfsd",
	"Executable": "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/unfsd/bin/unfsd",
	"Arguments": ["-e", "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/unfsd/etc/exports", "-i", "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/unfsd/unfsd.pid", "-s", "-l", "{{.Host}}", "-d"],
	"Env": [
		"PROJECTS={{.GetUserHomeDir}}/Sites",
		"HOME={{.GetUserHomeDir}}",
		"HOMEDIR={{.GetUserHomeDir}}",
		"ADMINROOTDIR={{.GetAdminRootDir}}",
		"CACHEDIR={{.GetCacheDir}}",
		"BASEDIR={{.GetUserHomeDir}}",
		"USERCONFIGDIR={{.GetUserConfigDir}}"
	],

	"Option": {
		"KeepAlive": true,
		"RunAtLoad": false,
		"SessionCreate": false,
		"UserService": true
	},

	"Stderr": "{{.GetStderr}}",
	"Stdout": "{{.GetStdout}}"
}`
