package eblog

import (
	// @TODO Since Gearbox won't need Loggly, can you make loggly support
	//       an external package that Gearbox won't need to include?
	//
	"github.com/sebest/logrusly"
)

type Args Logger

type Loggly struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	Server  string `json:"server"`
	Tag     string `json:"tag"`
	hook    *logrusly.LogglyHook
}

type LogFile struct {
	Enabled     bool   `json:"enabled"`
	Permissions string `json:"permissions"`
	Name        string `json:"name"`
}

// Disabled to work on GOOS=windows
//type Syslog struct {
//	Enabled  bool            `json:"enabled"`
//	Hostname string          `json:"hostname"`
//	Port     string          `json:"port"`
//	Protocol string          `json:"protocol"`
//	Priority syslog.Priority `json:"priority"`
//	Tag      string          `json:"tag"`
//	hook     *lSyslog.SyslogHook
//}
