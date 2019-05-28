package logger

import (
	"github.com/sirupsen/logrus"
	"strings"
)


func (me *Logger) printLog(level logrus.Level, fileName string, lineNumber int, textString string, opt ...interface{}) (returnCode bool) {

	returnCode = true

	textString = strings.TrimSuffix(textString, "\n")
	fields := logrus.Fields{
//		"tenant_name": notifyConfig.TenantName,
//		"tenant_guid": notifyConfig.TenantGUID,
//		"process_guid": notifyConfig.ProcessGUID,
		"filename": fileName,
		"line": lineNumber}

	switch {
		case level == DebugLevel:
			me.logrusInstance.WithFields(fields).Debugf(textString, opt...)

		case level == InfoLevel:
			me.logrusInstance.WithFields(fields).Infof(textString, opt...)

		case level == WarnLevel:
			me.logrusInstance.WithFields(fields).Warnf(textString, opt...)

		case level == ErrorLevel:
			me.logrusInstance.WithFields(fields).Errorf(textString, opt...)

		case level == FatalLevel:
			me.logrusInstance.WithFields(fields).Fatalf(textString, opt...)

		case level == PanicLevel:
			me.logrusInstance.WithFields(fields).Panicf(textString, opt...)
	}

	returnCode = false

	return
}

