// +build darwin
package daemon

import (
	"bytes"
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/eventbroker/only"
	"gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"regexp"

	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"
)

// @TODO Consider using https://github.com/kardianos/service
// 	Daemon "github.com/kardianos/service"

type PlistData struct {
	Label       string
	Program     string
	ProgramArgs []string
	Path        string
	KeepAlive   bool
	RunAtLoad   bool
	PidFile		string
}

var PlistTemplate = `
<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<plist version='1.0'>
	<dict>
		<key>Label</key><string>{{.Label}}</string>
		<key>Program</key><string>{{.Program}}</string>
		<key>ProgramArguments</key>
		<array>
			<string>{{.Program}}</string>
			{{range .ProgramArgs}}<string>{{index .}}</string>
			{{end}}
		</array>
		<key>StandardOutPath</key><string>{{.Path}}/out.log</string>
		<key>StandardErrorPath</key><string>{{.Path}}/err.log</string>
		<key>PIDFile</key><string>{{.PidFile}}</string>
		<key>KeepAlive</key><{{.KeepAlive}}/>
		<key>RunAtLoad</key><{{.RunAtLoad}}/>
	</dict>
</plist>
`


func NewDaemon(OsSupport oss.OsSupporter, args ...Args) *Daemon {
	var _args Args
	if len(args) > 0 {
		_args = args[0]
	}

	OsSupport = OsSupport

	if _args.Boxname == "" {
		_args.Boxname = global.Brandname
	}

	if _args.ServiceData.Label == "" {
		_args.ServiceData.Label = "com.gearbox.heartbeat" // Reverse-DNS naming convention
	}

	if _args.ServiceData.Program == "" {
		execPath, _ := os.Executable()
		_args.ServiceData.Program = execPath
	}

	if _args.ServiceData.Path == "" {
		execCwd, _ := os.Getwd()
		if execCwd == "/" {
			execCwd = string(OsSupport.GetAdminRootDir())
		}
		_args.ServiceData.Path = execCwd
	}

	if len(_args.ServiceData.ProgramArgs) < 1 {
		_args.ServiceData.ProgramArgs = []string{"heartbeat", "daemon"}
	}

	if _args.ServiceFile == "" {
		_args.ServiceFile = fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), _args.ServiceData.Label)
	}

	daemon := &Daemon{}
	*daemon = Daemon(_args)

/*
	svcConfig := &SysSvc.Config{
		Name:        daemon.Boxname,
		DisplayName: "Gearbox",
		Description: "This is an example Go service.",
		Executable:	 daemon.ServiceData.Program,
		Arguments: daemon.ServiceData.ProgramArgs,
		Option: SysSvc.KeyValue{
			"KeepAlive": true,
			"RunAtLoad": true,
			"PIDFile": daemon.ServiceData.PidFile,
		},
	}

	prg := &program{}
	s, err := SysSvc.New(prg, svcConfig)
	if err != nil {
		//		log.Fatal(err)
	}
	//	logger, err = s.Logger(nil)
	if err != nil {
		//		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		//		logger.Error(err)
	}
*/
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

		t := template.Must(template.New("launchdConfig").Parse(PlistTemplate))

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

		if me.IsRunning() {
			sts = status.Success("%s Heartbeat - Service already started %s", global.Brandname, me.ServiceData.Label)
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

	if IsParentInit() {
		displayString += fmt.Sprintf("%s Heartbeat: Running from init\n", global.Brandname)
	}

	if me.IsLoaded() {
		displayString += fmt.Sprintf("%s Heartbeat - Service installed [%s]\n", global.Brandname, me.ServiceData.Label)
	} else {
		displayString += fmt.Sprintf("%s Heartbeat - Service NOT installed [%s]\n", global.Brandname, me.ServiceData.Label)
	}

	if me.IsRunning() {
		displayString += fmt.Sprintf("%s Heartbeat - Service running [%s]\n", global.Brandname, me.ServiceData.Label)
	} else {
		displayString += fmt.Sprintf("%s Heartbeat - Service NOT running [%s]\n", global.Brandname, me.ServiceData.Label)
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
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Service file doesn't exist - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		// cmd := exec.Command("launchctl", "list", me.ServiceData.Label)
		// err = cmd.Run()
		sts, _, _, exitCode := me.RunCommand("launchctl", "list", me.ServiceData.Label)
		if exitCode != 0 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("Heartbeat - Error checking launchd: %v", err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		yesNo = true
		sts = status.Success("Heartbeat - Service loaded.")
	}

	return yesNo
}


func (me *Daemon) IsRunning() (yesNo bool) {

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
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s Heartbeat - Service file doesn't exist - %s", global.Brandname, me.ServiceFile),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		// cmd := exec.Command("launchctl", "list", me.ServiceData.Label)
		// err = cmd.Run()
		sts, stdout, _, exitCode := me.RunCommand("launchctl", "list", me.ServiceData.Label)
		if exitCode != 0 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("Heartbeat - Error checking launchd: %v", err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		re := regexp.MustCompile(`"PID" = ([0-9]+);`)
		matches := re.FindStringSubmatch(stdout)
		if len(matches) == 2 {
			yesNo = true

			sts = status.Success("Heartbeat - Process running.")
			break
		}
	}

	return yesNo

/*
	for range only.Once {
		sts := EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("PPID:%v:\n", IsParentInit())

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

	}
*/
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