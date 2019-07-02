package eblog

import (
	"github.com/sirupsen/logrus"
)

const (
	SkipNilCheck      = ""
	howMany           = 2
	defaultLogFile    = "eventbroker.log"
	DefaultEntityName = "eblog"
)

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)
