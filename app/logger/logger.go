package logger

import (
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/daemon"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/jinzhu/copier"
	"github.com/rifflock/lfshook"
	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"io/ioutil"
	"log/syslog"
	"path/filepath"
)


func NewLogger(OsSupport oss.OsSupporter, args ...Logger) (Logger, status.Status) {

	var _args Logger
	var sts status.Status
	se := &Logger{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy Logger config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

//		if _args.DebugMode == nil {
//			*_args.DebugMode = false
//		}


		_args.logrusInstance = logrus.New()
		_args.logrusInstance.SetFormatter(&logrus.JSONFormatter{})
		_args.currentLevel = InfoLevel
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
				_args.LogFile.Name = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetUserConfigDir(), defaultLogFile))
			}

			// fmt.Printf("Logging to files.")
			pathMap := lfshook.PathMap{
				logrus.DebugLevel: _args.LogFile.Name,
				logrus.InfoLevel: _args.LogFile.Name,
				logrus.WarnLevel: _args.LogFile.Name,
				logrus.ErrorLevel: _args.LogFile.Name,
				logrus.FatalLevel: _args.LogFile.Name,
				logrus.PanicLevel: _args.LogFile.Name,
			}
			_args.logrusInstance.Hooks.Add(lfshook.NewHook(pathMap, &logrus.TextFormatter{})) // &logrus.JSONFormatter{},))
			// _args.logrusInstance.SetOutput(os.Stderr)
			// _args.logrusInstance.SetOutput(os.Stdout)
			_args.logrusInstance.SetOutput(ioutil.Discard)
		}


		// Set sane values for Syslog based logging.
		if _args.Syslog.Enabled == true { // Set sane defaults for Protocol
			switch _args.Syslog.Protocol {
				case "udp":
				case "tcp":

				default:
					//LogInfo("Setting default syslog protocol to 'udp'. Was '%s'", _args.Syslog.Protocol)
					_args.Syslog.Protocol = "udp"
			}

			// Set sane defaults for Port
			switch _args.Syslog.Port {
				case "":
					//LogInfo("Setting default syslog port to '514'. Was '%s'", _args.Syslog.Port)
					_args.Syslog.Port = "514"
			}

			// Setup syslog based logging.
			switch _args.Syslog.Hostname {
				case "":

				default:
					var err error

					_args.Syslog.hook, err = lSyslog.NewSyslogHook(_args.Syslog.Protocol,
						_args.Syslog.Hostname+":"+_args.Syslog.Port,
						syslog.LOG_INFO,
						_args.Syslog.Tag)

					if err != nil {
						fmt.Printf("Error establishing connection to syslog server '%s' - %s",
							_args.Syslog.Hostname+":"+_args.Syslog.Port,
							err)

					} else {
						fmt.Printf("Established connection to syslog server '%s'",
							_args.Syslog.Hostname+":"+_args.Syslog.Port)
						_args.logrusInstance.Hooks.Add(_args.Syslog.hook)
					}
			}
		}


		// Set defaults for Loggly based logging.
		if _args.Loggly.Enabled == true {
			if _args.Loggly.Token != "" {
				_args.Loggly.hook = logrusly.NewLogglyHook(_args.Loggly.Token, _args.Loggly.Server, logrus.InfoLevel, _args.Loggly.Tag)
				_args.logrusInstance.Hooks.Add(_args.Loggly.hook)

			} else {
				fmt.Printf("# Error: Loggly requires a customer token.\n")
			}
		}

		//execPath, _ := os.Executable()
		//execCwd := string(_args.OsSupport.GetAdminRootDir()) + "/heartbeat" // os.Getwd()

		//_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), defaultPidFile))

		*se = Logger(_args)
	}

	status.Logger = se

	return *se, sts
}


func (me *Logger) SetLevel(sl string) {
	me.currentLevel = me.GetLevel(sl)
	me.logrusInstance.SetLevel(me.currentLevel)
}


func (me *Logger) GetLevel(getLevel string) (returnLevel logrus.Level) {
	returnLevel, _ = logrus.ParseLevel(getLevel)

	return
}


func (me *Logger) Log(err error) {

	if err == nil {
		return
	}

	// Find the calling function.
	filename, linenumber := daemon.MyCaller(daemon.CallerGrandParent)

	me.printLog(InfoLevel, filename, linenumber, "%v", err)
}


func (me *Logger) Debug(msg status.Msg) {

	if msg == "" {
		return
	}

	// Find the calling function.
	filename, linenumber := daemon.MyCaller(daemon.CallerGrandParent)

	me.printLog(DebugLevel, filename, linenumber, msg)
}


func (me *Logger) Warn(msg status.Msg) {

	if msg == "" {
		return
	}

	// Find the calling function.
	filename, linenumber := daemon.MyCaller(daemon.CallerGrandParent)

	me.printLog(DebugLevel, filename, linenumber, msg)
}

func (me *Logger) Error(msg status.Msg) {

	if msg == "" {
		return
	}

	// Find the calling function.
	filename, linenumber := daemon.MyCaller(daemon.CallerGrandParent)

	me.printLog(DebugLevel, filename, linenumber, msg)
}

func (me *Logger) Fatal(msg status.Msg) {

	if msg == "" {
		return
	}

	// Find the calling function.
	filename, linenumber := daemon.MyCaller(daemon.CallerGrandParent)

	me.printLog(DebugLevel, filename, linenumber, msg)
}


func (me *Logger) Cause() error {
	panic("implement me")
}

func (me *Logger) Data() interface{} {
	panic("implement me")
}

func (me *Logger) Additional() string {
	panic("implement me")
}

func (me *Logger) ErrorCode() int {
	panic("implement me")
}

func (me *Logger) FullError() error {
	panic("implement me")
}

func (me *Logger) GetFullDetails() string {
	panic("implement me")
}

func (me *Logger) GetFullMessage() string {
	panic("implement me")
}

func (me *Logger) GetHelp(status.HelpType) string {
	panic("implement me")
}

func (me *Logger) Help() string {
	panic("implement me")
}

func (me *Logger) HttpStatus() int {
	panic("implement me")
}

func (me *Logger) IsError() bool {
	panic("implement me")
}

func (me *Logger) IsSuccess() bool {
	panic("implement me")
}

func (me *Logger) IsWarn() bool {
	panic("implement me")
}

func (me *Logger) LogTo() status.LogType {
	panic("implement me")
}

func (me *Logger) Message() string {
	panic("implement me")
}

func (me *Logger) SetCause(error) status.Status {
	panic("implement me")
}

func (me *Logger) SetData(interface{}) status.Status {
	panic("implement me")
}

func (me *Logger) SetAdditional(string, ...interface{}) status.Status {
	panic("implement me")
}

func (me *Logger) SetErrorCode(int) status.Status {
	panic("implement me")
}

func (me *Logger) SetHelp(status.HelpType, string, ...interface{}) status.Status {
	panic("implement me")
}

func (me *Logger) SetHttpStatus(int) status.Status {
	panic("implement me")
}

func (me *Logger) SetLogTo(status.LogType) status.Status {
	panic("implement me")
}

func (me *Logger) SetMessage(string, ...interface{}) status.Status {
	panic("implement me")
}

func (me *Logger) SetOtherHelp(status.HelpTypeMap) status.Status {
	panic("implement me")
}

func (me *Logger) SetSuccess(bool) status.Status {
	panic("implement me")
}

func (me *Logger) SetWarn(bool) status.Status {
	panic("implement me")
}
