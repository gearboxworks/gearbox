package eblog

import (
	"fmt"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"github.com/gearboxworks/go-status/only"
	"github.com/rifflock/lfshook"
	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"strings"
)

var localLogger *Logger

type Logger struct {
	EntityId   msgs.Address
	EntityName msgs.Address
	OsPaths    *osdirs.BaseDirs
	//Sts        status.Status
	DebugMode bool    `json:"debug"`
	Loggly    Loggly  `json:"loggly"`
	LogFile   LogFile `json:"logfile"`
	// Disabled to work on GOOS=windows
	//Syslog    Syslog  `json:"syslog"`

	//status.L
	logrusInstance *logrus.Logger
	currentLevel   logrus.Level
}

// @TODO I have been following the convention that a New*() func never
//       returns an error, but always returns a valid object. IF an error
//       is possible it is a code smell if it is found in a constructor
//       meaning that there should be an .initialize() method instead.
//		 .
// 		 Can we restructure this?
//
func NewLogger(OsPaths *osdirs.BaseDirs, args ...Logger) (*Logger, error) {

	var _args Logger
	var err error
	se := &Logger{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.EntityName == "" {
			_args.EntityName = DefaultEntityName
		}

		if _args.EntityId == "" {
			_args.EntityId = _args.EntityName
		}

		_args.OsPaths = osdirs.New()

		_args.DebugMode = true

		_args.logrusInstance = logrus.New()
		_args.logrusInstance.SetFormatter(&logrus.JSONFormatter{})
		_args.currentLevel = DebugLevel
		_args.logrusInstance.SetLevel(_args.currentLevel)

		_args.LogFile.Enabled = true

		// Set sane values for File based logging.
		if _args.LogFile.Enabled == true {

			// Set sane defaults for permissions
			switch _args.LogFile.Permissions {
			case "":
				// fmt.Printf("Setting default permissions to '644'. Was '%s'", _args.LogFile.Permissions)
				// @TODO Can we reference consts in fileperms package vs. hardcoding 644?
				//       And if we might ever want to change, can we create a constant here instead?
				_args.LogFile.Permissions = "644"
			}

			if _args.LogFile.Name == "" {
				_args.LogFile.Name = osdirs.AddFilef(
					_args.OsPaths.EventBrokerLogDir,
					defaultLogFile,
				)
			}

			pathMap := lfshook.PathMap{
				logrus.DebugLevel: _args.LogFile.Name,
				logrus.InfoLevel:  _args.LogFile.Name,
				logrus.WarnLevel:  _args.LogFile.Name,
				logrus.ErrorLevel: _args.LogFile.Name,
				logrus.FatalLevel: _args.LogFile.Name,
				logrus.PanicLevel: _args.LogFile.Name,
			}
			_args.logrusInstance.Hooks.Add(lfshook.NewHook(pathMap, &logrus.TextFormatter{})) // &logrus.JSONFormatter{},))
			_args.logrusInstance.SetOutput(ioutil.Discard)
		}

		// Disabled to work on GOOS=windows
		// Set sane values for Syslog based logging.
		//if _args.Syslog.Enabled == true { // Set sane defaults for Protocol
		//	switch _args.Syslog.Protocol {
		//		case "udp":
		//		case "tcp":
		//
		//		default:
		//			//LogInfo("Setting default syslog protocol to 'udp'. Was '%s'", _args.Syslog.Protocol)
		//			_args.Syslog.Protocol = "udp"
		//	}
		//
		//	// Set sane defaults for Port
		//	switch _args.Syslog.Port {
		//		case "":
		//			//LogInfo("Setting default syslog port to '514'. Was '%s'", _args.Syslog.Port)
		//			_args.Syslog.Port = "514"
		//	}
		//
		//	// Setup syslog based logging.
		//	switch _args.Syslog.Hostname {
		//		case "":
		//
		//		default:
		//			var err error
		//
		//			_args.Syslog.hook, err = lSyslog.NewSyslogHook(_args.Syslog.Protocol,
		//				_args.Syslog.Hostname+":"+_args.Syslog.Port,
		//				syslog.LOG_INFO,
		//				_args.Syslog.Tag)
		//
		//			if err != nil {
		//				fmt.Printf("Error establishing connection to syslog server '%s' - %s",
		//					_args.Syslog.Hostname+":"+_args.Syslog.Port,
		//					err)
		//
		//			} else {
		//				fmt.Printf("Established connection to syslog server '%s'",
		//					_args.Syslog.Hostname+":"+_args.Syslog.Port)
		//				_args.logrusInstance.Hooks.Add(_args.Syslog.hook)
		//			}
		//	}
		//}

		// @TODO For Gearbox, we have no need nor no desire for Loggly.
		//       Do not default to Loggly.
		//
		// Set defaults for Loggly based logging.
		if _args.Loggly.Enabled == true {
			if _args.Loggly.Token != "" {
				_args.Loggly.hook = logrusly.NewLogglyHook(_args.Loggly.Token, _args.Loggly.Server, logrus.InfoLevel, _args.Loggly.Tag)
				_args.logrusInstance.Hooks.Add(_args.Loggly.hook)

			} else {
				fmt.Printf("# Error: Loggly requires a customer token.\n")
			}
		}

		*se = Logger(_args)

		if localLogger == nil {
			localLogger = se
		}
	}

	return se, err
}

func (me *Logger) GetLevel(getLevel string) (returnLevel logrus.Level) {
	returnLevel, _ = logrus.ParseLevel(getLevel)
	return
}

func (me *Logger) SetLevel(sl string) {
	me.currentLevel = me.GetLevel(sl)
	me.logrusInstance.SetLevel(me.currentLevel)
}

func (me *Logger) printLog(level logrus.Level, s string, opt ...interface{}) (returnCode bool) {

	returnCode = true

	s = strings.TrimSuffix(s, "\n")
	fields := logrus.Fields{}

	switch {
	case level == DebugLevel:
		me.logrusInstance.WithFields(fields).Debugf(s, opt...)

	case level == InfoLevel:
		me.logrusInstance.WithFields(fields).Infof(s, opt...)

	case level == WarnLevel:
		me.logrusInstance.WithFields(fields).Warnf(s, opt...)

	case level == ErrorLevel:
		for i, d := range *MyCallers(GrandParentCaller, howMany) {
			fields["caller"+strconv.Itoa(i)] = d.Function + ":" + strconv.Itoa(d.LineNumber)
		}
		me.logrusInstance.WithFields(fields).Errorf(s, opt...)

	case level == FatalLevel:
		for i, d := range *MyCallers(GrandParentCaller, howMany) {
			fields["caller"+strconv.Itoa(i)] = d.Function + ":" + strconv.Itoa(d.LineNumber)
		}
		me.logrusInstance.WithFields(fields).Fatalf(s, opt...)

	case level == PanicLevel:
		for i, d := range *MyCallers(GrandParentCaller, howMany) {
			fields["caller"+strconv.Itoa(i)] = d.Function + ":" + strconv.Itoa(d.LineNumber)
		}
		me.logrusInstance.WithFields(fields).Panicf(s, opt...)
	}

	returnCode = false

	return
}
