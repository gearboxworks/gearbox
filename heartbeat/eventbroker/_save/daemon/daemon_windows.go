// +build windows
package daemon

import (
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"path/filepath"
)

// @TODO Consider using https://github.com/kardianos/service
// 	Daemon "github.com/kardianos/service"

// @TODO These should be stubs, but it's just a copy of daemon_darwin.go

type Daemon struct {
	Boxname      string
	ServiceFile  string
	ServiceData	 plistData

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

	_args.ServiceData = plistData{
		Label:     "com.gearbox.heartbeat", // Reverse-DNS naming convention
		Program:   execPath,
		Path:      execCwd,
		KeepAlive: true,
		RunAtLoad: true,
	}

	if _args.ServiceFile == "" {
		_args.ServiceFile = fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), _args.ServiceData.Label)
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

		f, err := os.Create(me.ServiceFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Cannot create service file - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		t := template.Must(template.New("launchdConfig").Parse(plistTemplate))

		err = t.Execute(f, me.ServiceData)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Service file creation failed - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
		}
	}

	return sts
}


func (me *Daemon) Load() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.CreatePlist()
		if is.Error(sts) {
			break
		}

		if me.IsLoaded() {
			sts = status.Success("%s Heartbeat - Service already loaded %s", global.Brandname, me.ServiceData.Label)
			break
		}

		_, err := os.Open(me.ServiceFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Service file doesn't exist - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		cmd := exec.Command("launchctl", "load", "-F", me.ServiceFile)
		err = cmd.Run()
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Cannot load service file - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})

		} else {
			sts = status.Success("%s Heartbeat - Loaded service %s", global.Brandname, me.ServiceData.Label)
		}
	}

	return sts
}


func (me *Daemon) Unload() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		if !me.IsLoaded() {
			sts = status.Success("%s Heartbeat - Service already unloaded %s", global.Brandname, me.ServiceData.Label)
			break
		}

		_, err := os.Open(me.ServiceFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Service file doesn't exist - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		cmd := exec.Command("launchctl", "unload", me.ServiceFile)
		err = cmd.Run()
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Cannot unload service file - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})

		} else {
			sts = status.Success("%s Heartbeat - Unloaded service %s", global.Brandname, me.ServiceData.Label)
		}
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


func (me *Daemon) GetState() (sts status.Status) {

	displayString := ""

	if me.IsParentInit() {
		displayString += fmt.Sprintf("%s Heartbeat: Running from init\n", global.Brandname)
	}

	if me.IsLoaded() {
		displayString += fmt.Sprintf("%s Heartbeat - Service installed [%s]\n", global.Brandname, me.ServiceData.Label)
	} else {
		displayString += fmt.Sprintf("%s Heartbeat - Service NOT installed [%s]\n", global.Brandname, me.ServiceData.Label)
	}

	/*
		foo1, _ := process.Pids()
		for i, p := range foo1 {
			fmt.Printf("process.Pids:%v:	%v:\n", i, p)
		}

		foo2, _ := process.Processes()
		for _, p := range foo2 {
			c, _ := p.Cmdline()
			fmt.Printf("process.Processes:%v:	'%s'\n", p.Pid, c)
		}

		infoStat, _ := host.Info()
		fmt.Printf("Total processes: %v\n", infoStat.Procs)

		miscStat, _ := load.Misc()
		fmt.Printf("Running processes: %v\n", miscStat.ProcsRunning)
	*/

	/*
		fp := filepath.FromSlash(fmt.Sprintf("%s/%s", me.OsSupport.GetAdminRootDir(), s))
		if fp == "" {
			return nil
		}

		b, err := ioutil.ReadFile(fp)
		if err != nil {
			fmt.Print(err)
		}

		return b
	*/
	sts = status.Success(displayString)

	return sts
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


func (me *Daemon) IsParentInit() (bool) {

	ppid := os.Getppid()
	if ppid == 1 {
		return true
	}

	return false
}


func (me *Daemon) IsRunning() (bool) {

	fmt.Printf("PPID:%v:\n", me.IsParentInit())

	foo1, _ := process.Pids()
	for i, p := range foo1 {
		fmt.Printf("process.Pids:%v:	%v:\n", i, p)
	}

	foo2, _ := process.Processes()
	for _, p := range foo2 {
		c, _ := p.Cmdline()
		fmt.Printf("process.Processes:%v:	'%s'\n", p.Pid, c)
	}

	infoStat, _ := host.Info()
	fmt.Printf("Total processes: %v\n", infoStat.Procs)

	miscStat, _ := load.Misc()
	fmt.Printf("Running processes: %v\n", miscStat.ProcsRunning)

	return false
}


func (me *Daemon) IsLoaded() (yesNo bool) {

	for range only.Once {
		sts := EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.CreatePlist()
		if is.Error(sts) {
			break
		}

		_, err := os.Open(me.ServiceFile)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Service file doesn't exist - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		// cmd := exec.Command("launchctl", "list", me.ServiceData.Label)
		// err = cmd.Run()
		sts, _, _, exitCode := me.RunCommand("launchctl", "list", me.ServiceData.Label)
		if exitCode == 0 {
			yesNo = true
		} else {
			yesNo = false
		}
	}

	return yesNo
}


const defaultFailedCode = 1
func (me *Daemon) RunCommand(name string, args ...string) (sts status.Status, stdout string, stderr string, exitCode int) {

	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()

			sts = status.Success("%s Heartbeat EXEC - Return(%d)", global.Brandname, exitCode)

		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			// log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}

			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s Heartbeat EXEC - Return(%d)\n%s", global.Brandname, exitCode, err.Error()),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    exitCode,
			})
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
		sts = status.Success("%s Heartbeat EXEC - Return(%d)", global.Brandname, exitCode)
	}
	// log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)

	return
}