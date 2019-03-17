package gearbox

import "time"

const (
	BoxDefaultWaitDelay       = time.Second
	BoxDefaultWaitRetries     = 90
	BoxDefaultConsoleHost     = "127.0.0.1"
	BoxDefaultConsolePort     = "2023"
	BoxDefaultConsoleOkString = "Gearbox Heartbeat"
	BoxDefaultShowConsole     = false
	BoxDefaultConsoleReadWait = time.Second * 5
)

const (
	Boxname                = "Gearbox"
	BoxBasedir             = "/home/gearbox/projects"
	BoxOvaFileName         = "/box/vm/Gearbox.ova"
	PrimaryBasedirNickname = "primary"
)

const (
	BoxError      = "error"
	BoxUnknown    = "unknown"
	BoxAbsent     = "absent"
	BoxHalted     = "halted"
	BoxRunning    = "running"
	BoxStarted    = "started"
	BoxGearBoxOK  = "ok"
	BoxGearBoxNOK = "nok"
)

const (
	ConfigSchemaVersion = "1.0"
	ConfigHelpUrl       = "https://docs.gearbox.works/config"
)

const (
	ProjectFilename    = "gearbox.json"
	ProjectFileHelpUrl = "https://docs.gearbox.works/projects/" + ProjectFilename
)

const (
	HostApiPort        = "9999"
	HostAPIServiceName = "GearBox API"
	HostAPIVersion     = "0.1"
	HostAPIDocsUrl     = "https://docs.gearbox.works/api"
)
