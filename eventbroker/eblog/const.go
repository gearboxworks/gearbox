package eblog

import (
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
)


const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)

const defaultLogFile = "eventbroker.log"

const DefaultEntityName = "eblog"


type Logger struct {
	EntityId   *messages.MessageAddress
	EntityName messages.MessageAddress
	OsPaths    *ospaths.BasePaths
	//Sts        status.Status
	DebugMode  bool    `json:"debug"`
	Loggly     Loggly  `json:"loggly"`
	LogFile    LogFile `json:"logfile"`
	// Disabled to work on GOOS=windows
	//Syslog    Syslog  `json:"syslog"`

	//status.L
	logrusInstance *logrus.Logger
	currentLevel   logrus.Level
}
type Args Logger

type Loggly struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	Server  string `json:"server"`
	Tag     string `json:"tag"`
	hook *logrusly.LogglyHook
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

type LogFile struct {
	Enabled     bool   `json:"enabled"`
	Permissions string `json:"permissions"`
	Name        string `json:"name"`
}
