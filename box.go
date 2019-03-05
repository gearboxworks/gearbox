package gearbox

import (
	"bufio"
	"fmt"
	"gearbox/util"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"net"
	"regexp"
	"strings"
	"time"
)

type Box struct {
	Name     string
	Instance virtualbox.VM
	Status   string

	// Status polling delays.
	NoWait      bool
	WaitDelay   time.Duration
	WaitRetries int

	// Console related.
	ConsoleHost     string
	ConsolePort     string
	ConsoleOkString string
	ConsoleReadWait time.Duration
	ShowConsole     bool
}
type BoxArgs Box

const BoxDefaultName = "Gearbox"
const BoxDefaultWaitDelay = time.Second
const BoxDefaultWaitRetries = 30
const BoxDefaultConsoleHost = "127.0.0.1"
const BoxDefaultConsolePort = "2023"
const BoxDefaultConsoleOkString = "GearBox API"
const BoxDefaultShowConsole = false
const BoxDefaultConsoleReadWait = time.Second * 5

const BoxError = "error"
const BoxUnknown = "unknown"
const BoxHalted = "halted"
const BoxRunning = "running"
const BoxStarted = "started"
const BoxOK = "ok"
const BoxNOK = "nok"

// //////////////////////////////////////////////////////////////////////////////
// Gearbox related
func (me *Gearbox) StartBox(boxArgs BoxArgs) error {

	box := NewBox(*me, boxArgs)

	err := box.Start()

	return err
}

func (me *Gearbox) StopBox(boxArgs BoxArgs) error {

	box := NewBox(*me, boxArgs)

	err := box.Stop()

	return err
}

func (me *Gearbox) RestartBox(boxArgs BoxArgs) error {

	box := NewBox(*me, boxArgs)

	err := box.Restart()

	return err
}

func (me *Gearbox) GetBoxStatus(boxArgs BoxArgs) (string, error) {

	box := NewBox(*me, boxArgs)

	err := box.GetStatus()
	if err != nil {
		return box.Status, err
	}

	err = box.GetApiStatus()
	if err != nil {
		return box.Status, err
	}

	switch box.Status {
	case BoxUnknown:
		fmt.Printf("\r👎 %s: BOX & API in an unknown state.\n", me.Config.BoxName)

	case BoxHalted:
		fmt.Printf("\r👎 %s: BOX halted. API halted.\n", me.Config.BoxName)

	case BoxRunning:
		fmt.Printf("\r👎 %s: BOX running. API halted.\n", me.Config.BoxName)

	case BoxStarted:
		fmt.Printf("\r👎 %s: BOX running. API halted.\n", me.Config.BoxName)

	case BoxOK:
		fmt.Printf("\r👍 %s: BOX running. API running.\n", me.Config.BoxName)

	case BoxNOK:
		fmt.Printf("\r👎 %s: BOX running. API halted.\n", me.Config.BoxName)
	}

	return box.Status, err
}

// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewBox(gb Gearbox, args ...BoxArgs) *Box {
	var _args BoxArgs
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.Name == "" {
		_args.Name = BoxDefaultName
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = BoxDefaultWaitDelay
	}

	if _args.WaitRetries == 0 {
		_args.WaitRetries = BoxDefaultWaitRetries
	}

	if _args.ConsoleHost == "" {
		_args.ConsoleHost = BoxDefaultConsoleHost
	}

	if _args.ConsolePort == "" {
		_args.ConsolePort = BoxDefaultConsolePort
	}

	if _args.ConsoleOkString == "" {
		_args.ConsoleOkString = BoxDefaultConsoleOkString
	}

	if _args.ConsoleReadWait == 0 {
		_args.ConsoleReadWait = BoxDefaultConsoleReadWait
	}

	_args.Instance = virtualbox.VM{
		Name: _args.Name,
	}

	box := &Box{}
	*box = Box(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return box
}

func (me *Box) WaitForBoxState(waitForState string, displayString string) error {

	var waitCount int

	spinner := util.NewSpinner(util.SpinnerArgs{
		Text:    displayString,
		ExitOK:  displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
	spinner.Start()

	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
		err := me.GetStatus()
		if err != nil {
			spinner.Stop(false)
			return err
		}
		if me.Status == waitForState {
			spinner.Stop(true)
			break
		}

		time.Sleep(me.WaitDelay)
		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
	}

	return nil
}

func (me *Box) WaitForConsole(displayString string, waitFor time.Duration) error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.GetStatus()
	if err != nil {
		return err
	}

	// TRUE - show the spinner on console.
	displaySpinner := me.ShowConsole == false && displayString != ""

	if me.Status == BoxRunning {
		spinner := util.NewSpinner(util.SpinnerArgs{
			Text:    displayString,
			ExitOK:  displayString + " - OK",
			ExitNOK: displayString + " - FAILED",
		})

		if displaySpinner == true {
			// We want to display just a spinner instead of console output.
			spinner.Start()
		}

		// Connect to this console
		conn, err := net.Dial("tcp", me.ConsoleHost+":"+me.ConsolePort)
		defer conn.Close()
		if err != nil {
			return err
		}

		exitWhen := time.Now().Add(time.Second * waitFor)

		readBuffer := make([]byte, 512)
		for waitCount := 0; time.Now().Unix() < exitWhen.Unix(); waitCount++ {
			err = conn.SetDeadline(time.Now().Add(me.ConsoleReadWait))
			if err != nil {
				return err
			}

			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
			// bytesRead, err := conn.Read(readBuffer)
			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
			// bytesRead := len(readBuffer)
			if err != nil {
				me.Status = BoxNOK
				if displaySpinner == true {
					spinner.Stop(false)
				}
				break
			}

			if bytesRead > 0 {
				if me.ShowConsole == true {
					fmt.Printf("%s", string(readBuffer[:bytesRead]))
				}

				apiSplit := strings.Split(string(readBuffer[:bytesRead]), ";")
				if len(apiSplit) > 1 {
					match, _ := regexp.MatchString(me.ConsoleOkString, apiSplit[1])
					if match == true {
						if apiSplit[2] == "OK" {
							me.Status = BoxOK
						} else {
							me.Status = BoxNOK
						}
						if displaySpinner == true {
							spinner.Stop(true)
						}
						break
					}
				}
			}

			time.Sleep(me.WaitDelay)
			if displaySpinner == true {
				spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
			}
		}

		if me.ShowConsole == true {
			fmt.Printf("\n\n# Exiting Console.\n")
		}
	}

	return nil
}

func (me *Box) Start() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.GetStatus()
	if err != nil {
		return err
	}
	if me.Status == BoxRunning || me.Status == BoxStarted {
		return nil
	}

	err = me.Instance.Start()
	if err != nil {
		return err
	}
	if me.NoWait == false {
		err := me.WaitForBoxState(BoxRunning, me.Name+" BOX: Starting")
		if err != nil {
			return err
		}

		// Check final state of the system from the top down.
		err = me.WaitForConsole(me.Name+" : Starting", 30)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Box) Stop() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.GetStatus()
	if err != nil {
		return err
	}
	if me.Status == BoxHalted {
		return nil
	}

	err = me.Instance.Halt()
	if err != nil {
		return err
	}

	if me.NoWait == false {
		err := me.WaitForBoxState(BoxHalted, me.Name+" BOX: Stopping")
		if err != nil {
			return err
		}

		// Check final state of the system from the top down.
		err = me.WaitForConsole(me.Name+" : Stopping", 30)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Box) Restart() error {

	var err error

	if me == nil {
		// Throw software error.
		return nil
	}

	err = me.GetStatus()
	if err != nil {
		return err
	}

	switch me.Status {
	case BoxOK:
		fallthrough
	case BoxNOK:
		fallthrough
	case BoxRunning:
		fallthrough
	case BoxStarted:
		err := me.Stop()
		if err != nil {
			return err
		}
		err = me.Start()
		if err != nil {
			return err
		}

	case BoxHalted:
		err := me.Start()
		if err != nil {
			return err
		}

	case BoxUnknown:

	}

	if me.Status != BoxRunning {
		// Throw an error.
	}

	return err
}

func (me *Box) GetStatus() error {

	if me == nil {
		// Throw software error.
		me.Status = BoxUnknown
		return nil
	}

	state, err := me.Instance.GetState()
	if err != nil {
		me.Status = BoxError
	} else {
		me.Status = state
	}

	return err
}

func (me *Box) GetApiStatus() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	if me.Status == BoxRunning {
		err := me.WaitForConsole("", 10)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func scanForAPI(text string) bool {

	r, err := regexp.Compile("^.*%$")
	// r, err := regexp.Compile("Welcome to GearBox.*")
	if err == nil {
		if r.MatchString(text) {

			switch {
				case text == "Welcome to GearBox":
					return true
			}
		}
	}

	return false
}
*/
