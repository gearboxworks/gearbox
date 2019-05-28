package logger

import (
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

var _ status.MsgLogger = (*Logger)(nil)

type Logger struct {
	Boxname   string
	OsSupport oss.OsSupporter
	Sts       status.Status
	DebugMode bool    `json:"debug"`
	Loggly    Loggly  `json:"loggly"`
	Syslog    Syslog  `json:"syslog"`
	LogFile   LogFile `json:"logfile"`

	status.L
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

type Syslog struct {
	Enabled  bool            `json:"enabled"`
	Hostname string          `json:"hostname"`
	Port     string          `json:"port"`
	Protocol string          `json:"protocol"`
	Priority syslog.Priority `json:"priority"`
	Tag      string          `json:"tag"`
	hook     *lSyslog.SyslogHook
}

type LogFile struct {
	Enabled     bool   `json:"enabled"`
	Permissions string `json:"permissions"`
	Name        string `json:"name"`
}

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)

const defaultLogFile = "logs/gearbox.log"
