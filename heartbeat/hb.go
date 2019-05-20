package heartbeat

/*

import (
	"fmt"
	"github.com/reiver/go-telnet"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
)






// //////////////////////////////////////////////////////////////////////////////

type fn func()


// Function: DaemonizeMe()
// Description: Daemonizes a Go executable.
// Return:
// TRUE - Instance is running as a background process.
// FALSE - Instance is running as a foreground process.
func DaemonizeMe(function func()) (bool) {

	// Pull in all processes - we want to check if the daemon is running already.
// Doesn't seem to work with go-daemon
	pidss, err := process.Processes()
	if err != nil {
		log.Fatal("# Unable to check PIDs: ", err)
	}
	fmt.Printf("# PID count: %d\n", len(pidss))

	// Scan for our daemon process.
	for _, myPID := range pidss {
		cmd, _ := myPID.Cmdline()
		// fmt.Printf("# %d\n", myPID.Pid)
		if cmd == pidName {
			fmt.Printf("# Daemon process already running, [%d].\n", myPID.Pid)
			foundPID = true
			//break
		}
	}
	pidss = nil

	// go-daemon doesn't actually create any pid files, so we're going to hack it in here.
	// detect if file exists
	var _, err = os.Stat("/tmp/gearbox.pid")

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Printf("DUCK1\n")
		var file, err = os.Create("/tmp/gearbox.pid")
		if err != nil {
			// Daemon process already running.
			fmt.Printf("DUCK2\n")
			return false
		}
		defer file.Close()
		fmt.Printf("DUCK3\n")

	} else {
		// Daemon process already running.
		fmt.Printf("DUCK4\n")
		return false
	}
	fmt.Printf("DUCK5\n")


	// We haven't found it - we can start one.
	fmt.Printf("Run daemon process.\n")
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{pidName},
	}


	if len(daemon.ActiveFlags()) > 0 {
		fmt.Printf("SHIT2\n")
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return false
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	defer cntxt.Release()
	if d != nil {
		// Daemon running.
		fmt.Printf("Daemon running.\n")
		return false
	}

	fmt.Printf("- - - - - - - - - - - - - - -")
	fmt.Printf("Daemon started")


	// Daemon process needs to run.
	if function != nil {
		// .
	}

	return true
}



// //////////////////////////////////////////////////////////////////////////////
func (vm *GearboxVM) startGearbox() (bool) {
	fmt.Printf("Menu: Start\n")
	err := vm.Start()
	if err != nil {
		fmt.Printf("Shit! an error: %s\n", err)
	}

	return true
}


func (vm *GearboxVM) stopGearbox() (bool) {
	fmt.Printf("Menu: Stop\n")
	vm.Halt()

	return true
}


func (vm *GearboxVM) consoleGearbox() (bool) {
	fmt.Printf("Menu: Console\n")

	return true
}


func (vm *GearboxVM) guiGearbox() (bool) {
	fmt.Printf("Menu: GUI\n")

	return true
}


func (vm *GearboxVM) sshGearbox() (bool) {
	fmt.Printf("Menu: SSH\n")
	// foo, _ := vm.GetSSH(ssh.Options{})

	return true
}


func (vm GearboxVM) getGearboxState() (state string) {

	var err error

	state, err = vm.GetState()
	// fmt.Printf("STATE:[%s]\t", state)
	if err != nil {
		state = "unknown"
	}

	return
}



// //////////////////////////////////////////////////////////////////////////////
func telnetClient() {
	var caller telnet.Caller = telnet.StandardCaller

	telnet.DialToAndCall("localhost:2023", caller)
}

*/