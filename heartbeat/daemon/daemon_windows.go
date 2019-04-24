// +build windows
package daemon

import (
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/status"
	"gearbox/status/is"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"path/filepath"
)

// @TODO Consider using https://github.com/kardianos/service
// 	Daemon "github.com/kardianos/service"


type Daemon struct {
	Boxname      string
	// BoxInstance  *box.Box
	PlistFile    string
	PlistData	 plistData

	OsSupport    oss.OsSupporter
}
type Args Daemon

type plistData struct {
	Label     string
	Program   string
	Path      string
	KeepAlive bool
	RunAtLoad bool
}


var plistTemplate = `
`


func NewDaemon(OsSupport oss.OsSupporter, args ...Args) *Daemon {
	var _args Args
	if len(args) > 0 {
		_args = args[0]
	}

	_args.OsSupport = OsSupport

	if _args.Boxname == "" {
		_args.Boxname = global.Brandname
	}

	execPath, _ := os.Executable()
	execCwd, _ := os.Getwd()

	_args.PlistData = plistData{
		Label:     "com.gearbox.heartbeat", // Reverse-DNS naming convention
		Program:   execPath,
		Path:      execCwd,
		KeepAlive: true,
		RunAtLoad: true,
	}

	if _args.PlistFile == "" {
		_args.PlistFile = fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), _args.PlistData.Label)
	}

	daemon := &Daemon{}
	*daemon = Daemon(_args)

	return daemon
}


func (me *Daemon) CreatePlist() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		// @TODO Windows STUB

/*
		f, err := os.Create(me.PlistFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Cannot create Plist file - %s", global.Brandname, me.PlistFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		t := template.Must(template.New("launchdConfig").Parse(plistTemplate))

		err = t.Execute(f, me.PlistData)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Plist file creation failed - %s", global.Brandname, me.PlistFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
		}
*/
	}

	return sts
}


func (me *Daemon) Load() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		// @TODO Windows STUB

/*
		sts = me.CreatePlist()
		if is.Error(sts) {
			break
		}

		_, err := os.Open(me.PlistFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Plist file doesn't exist - %s", global.Brandname, me.PlistFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		cmd := exec.Command("launchctl", "load", me.PlistFile)
		err = cmd.Run()
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Cannot load Plist file - %s", global.Brandname, me.PlistFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
		}
*/
	}

	return sts
}


func (me *Daemon) Unload() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		// @TODO Windows STUB

/*
		_, err := os.Open(me.PlistFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Plist file doesn't exist - %s", global.Brandname, me.PlistFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		cmd := exec.Command("launchctl", "unload", me.PlistFile)
		err = cmd.Run()
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Cannot unload Plist file - %s", global.Brandname, me.PlistFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
		}
*/
	}

	return sts
}


func (me *Daemon) getFile(s string) []byte {

	fp := filepath.FromSlash(fmt.Sprintf("%s/%s", me.OsSupport.GetAdminRootDir(), s))
	if fp == "" {
		return nil
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Print(err)
	}

	return b
}


func (me *Daemon) GetState(s string) []byte {

	foo, err := process.Processes()

	fp := filepath.FromSlash(fmt.Sprintf("%s/%s", me.OsSupport.GetAdminRootDir(), s))
	if fp == "" {
		return nil
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Print(err)
	}

	return b
}


func EnsureNotNil(bx *Daemon) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    nil,
		})
	}
	return sts
}
