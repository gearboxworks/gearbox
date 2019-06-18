package eblog

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)


func (me *Logger) printLogOld(level logrus.Level, fileName string, lineNumber int, textString string, opt ...interface{}) (returnCode bool) {

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


func (me *Logger) printLog(level logrus.Level, textString string, opt ...interface{}) (returnCode bool) {

	returnCode = true

	textString = strings.TrimSuffix(textString, "\n")
	fields := logrus.Fields{}

	switch {
		case level == DebugLevel:
			me.logrusInstance.WithFields(fields).Debugf(textString, opt...)

		case level == InfoLevel:
			me.logrusInstance.WithFields(fields).Infof(textString, opt...)

		case level == WarnLevel:
			me.logrusInstance.WithFields(fields).Warnf(textString, opt...)

		case level == ErrorLevel:
			for i, d := range *MyCallers(CallerGrandParent, howMany) {
				fields["caller" + strconv.Itoa(i)] = d.Function + ":" + strconv.Itoa(d.LineNumber)
			}
			me.logrusInstance.WithFields(fields).Errorf(textString, opt...)

		case level == FatalLevel:
			for i, d := range *MyCallers(CallerGrandParent, howMany) {
				fields["caller" + strconv.Itoa(i)] = d.Function + ":" + strconv.Itoa(d.LineNumber)
			}
			me.logrusInstance.WithFields(fields).Fatalf(textString, opt...)

		case level == PanicLevel:
			for i, d := range *MyCallers(CallerGrandParent, howMany) {
				fields["caller" + strconv.Itoa(i)] = d.Function + ":" + strconv.Itoa(d.LineNumber)
			}
			me.logrusInstance.WithFields(fields).Panicf(textString, opt...)
	}

	returnCode = false

	return
}

