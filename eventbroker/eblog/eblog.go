package eblog

import (
	"fmt"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"github.com/gearboxworks/go-status/only"
	"github.com/rifflock/lfshook"
	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"reflect"
	"runtime"
)

// To be used with MyCaller()
const (
	callerStack0           = iota
	callerStack1           = iota
	CallerCurrent          = iota
	CallerParent           = iota
	CallerGrandParent      = iota
	CallerGreatGrandParent = iota
)

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCallerExtended(whichCaller int) (fileName string, lineNumber int) {

	fileName = "unknown"

	if whichCaller == 0 {
		whichCaller = CallerParent
	}

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(whichCaller, fpcs)
	if n == 0 {
		return // Proper error handling would be better here.
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return
	}

	_, fileName, lineNumber, _ = runtime.Caller(whichCaller)

	fileName = fun.Name()

	return
}

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCaller(whichCaller int) (string, int) {

	pc, _, _, _ := runtime.Caller(whichCaller)
	e := runtime.FuncForPC(pc)
	fn := e.Name()
	_, ln := e.FileLine(pc)

	return fn, ln
}

type Caller struct {
	File       string
	LineNumber int
	Function   string
}
type Callers []Caller

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCallers(whichCaller int, howMany int) *Callers {

	if whichCaller == 0 {
		whichCaller = CallerParent
	}

	if howMany == 0 {
		howMany = 2
	}

	pc := make([]uintptr, howMany)
	count := runtime.Callers(whichCaller, pc)

	callers := make(Callers, count)
	for i, d := range pc {
		e := runtime.FuncForPC(d)
		callers[i].Function = e.Name()
		callers[i].File, callers[i].LineNumber = e.FileLine(d)
	}

	return &callers
}

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func (me *Callers) Print() string {

	var ret string

	if me == nil {
		return ""
	}

	for k, v := range *me {
		ret += fmt.Sprintf("[%d] %s %s:%d\n", k, v.File, v.Function, v.LineNumber)
	}

	return ret
}

// Check for a nil type.
func IsNil(i interface{}) bool {
	defer func() { recover() }()
	if i == nil || reflect.ValueOf(i).IsNil() {
		// It's a nil type.
		return true
	} else {
		// It's not a nil type.
		return false
	}
}

// Check for a nil type.
func IsNotNil(i interface{}) bool {
	defer func() { recover() }()
	if i == nil || reflect.ValueOf(i).IsNil() {
		// It's a nil type.
		return false
	} else {
		// It's not a nil type.
		return true
	}
}

// Check for a nil type.
func LogIfNil(i interface{}, format ...interface{}) bool {

	var ret bool

	switch {
	case reflect.ValueOf(i).String() == "":

	case i == nil:
		fallthrough
	case reflect.ValueOf(i).IsNil():
		ret = true

		//callers := " NIL:["
		//// Fetch last two callers.
		//for _, d := range *MyCallers(CallerParent, howMany) {
		//	callers += " <- " + d.Function + ":" + strconv.Itoa(d.LineNumber)
		//}
		//callers += "] "

		localLogger.printLog(logrus.ErrorLevel, "nil interface")
		//status.Success("nil interface" + callers).Log()
	}

	return ret
}

func Debug(client messages.MessageAddress, format string, a ...interface{}) {

	if localLogger == nil {
		return
	}

	fn, ln := MyCaller(CallerParent)

	localLogger.printLogOld(logrus.DebugLevel, fn, ln, string(client)+": "+format, a...)
}

const SkipNilCheck = ""
const howMany = 2

// Check for a nil type or err and log.
func LogIfError(address messages.MessageAddress, err error, format ...interface{}) bool {

	var ret bool

	if err != nil {
		ret = true
		//callers := "%s ERROR:["
		// callers := address.String() + " ERROR:%s ["
		// Fetch last two callers.
		//for _, d := range *MyCallers(CallerParent, howMany) {
		//	callers += " <- " + d.Function + ":" + strconv.Itoa(d.LineNumber)
		//}
		//callers += "] "

		if len(format) == 0 {
			localLogger.printLog(logrus.ErrorLevel, "%v", err)
			//status.Success(callers, err).Log()
			//fmt.Printf(callers + "\n", err)
		} else {
			localLogger.printLog(logrus.ErrorLevel, format[0].(string), format[1:]...)
			//status.Success(format[0].(string) + callers, format[1:]...).Log()
			//fmt.Printf(format[0].(string) + callers + "\n", format[1:]...)
		}
	}

	return ret
}

var localLogger *Logger

func NewLogger(OsPaths *ospaths.BasePaths, args ...Logger) (*Logger, error) {

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

		if _args.EntityId == nil {
			_args.EntityId = &_args.EntityName // messages.GenerateAddress()
		}

		_args.OsPaths = ospaths.New("")
		//_args.baseDir = _args.osPaths.UserConfigDir.AddToPath(DefaultBaseDir)
		//_args.pidFile = _args.baseDir.AddFileToPath(DefaultPidFile).String()

		_args.DebugMode = true

		_args.logrusInstance = logrus.New()
		_args.logrusInstance.SetFormatter(&logrus.JSONFormatter{})
		_args.currentLevel = DebugLevel
		_args.logrusInstance.SetLevel(_args.currentLevel)

		_args.LogFile.Enabled = true

		// Set sane values for File based logging.
		if _args.LogFile.Enabled == true {
			// fmt.Printf("Setting up file logging.")

			// Set sane defaults for permissions
			switch _args.LogFile.Permissions {
			case "":
				// fmt.Printf("Setting default permissions to '644'. Was '%s'", _args.LogFile.Permissions)
				_args.LogFile.Permissions = "644"
			}

			if _args.LogFile.Name == "" {
				_args.LogFile.Name = _args.OsPaths.EventBrokerLogDir.AddFileToPath(defaultLogFile).String()
			}

			// fmt.Printf("Logging to files.")
			pathMap := lfshook.PathMap{
				logrus.DebugLevel: _args.LogFile.Name,
				logrus.InfoLevel:  _args.LogFile.Name,
				logrus.WarnLevel:  _args.LogFile.Name,
				logrus.ErrorLevel: _args.LogFile.Name,
				logrus.FatalLevel: _args.LogFile.Name,
				logrus.PanicLevel: _args.LogFile.Name,
			}
			_args.logrusInstance.Hooks.Add(lfshook.NewHook(pathMap, &logrus.TextFormatter{})) // &logrus.JSONFormatter{},))
			// _args.logrusInstance.SetOutput(os.Stderr)
			// _args.logrusInstance.SetOutput(os.Stdout)
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

	//status.Logger = se

	return se, err
}

func (me *Logger) SetLevel(sl string) {
	me.currentLevel = me.GetLevel(sl)
	me.logrusInstance.SetLevel(me.currentLevel)
}

func (me *Logger) GetLevel(getLevel string) (returnLevel logrus.Level) {
	returnLevel, _ = logrus.ParseLevel(getLevel)

	return
}

/*
	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}
*/
