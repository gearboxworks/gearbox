package gearbox

/*
type Ssh struct {
	SshName  string
	Instance ssh.Client

	// Status polling delays.
	NoWait        bool
	WaitDelay     time.Duration
	WaitRetries   int

	// SSH related.
	SshHost        string
	SshPort        string
	SshPublicKey   string
	SshOkString    string
	SshWait        time.Duration
	SshWaitRetries int
	ShowSsh        bool
}
type SshArgs Ssh


const SshDefaultName = "Gearbox"
const SshDefaultWaitDelay = time.Second
const SshDefaultWaitRetries = 30
const SshDefaultSshHost = "127.0.0.1"
const SshDefaultSshPort = "2023"
const SshDefaultSshOkString = "Welcome to GearBox"
const SshDefaultShowSsh = false
const SshDefaultSshWait = time.Second
const SshDefaultSshWaitRetries = 5


const SshUnknown = "unknown"
const SshHalted = "halted"
const SshRunning = "running"
const SshStarted = "started"
const SshGearBoxOK = "ok"


// //////////////////////////////////////////////////////////////////////////////
// Gearbox related
func (me *Gearbox) StartSSH(sshArgs SshArgs) error {

	sshClient := NewSsh(*me, sshArgs)

	err := sshClient.StartSsh()

	return err
}


func (me *Gearbox) StopSSH(sshArgs SshArgs) error {

	sshClient := NewSsh(*me, sshArgs)

	err := sshClient.StopSsh()

	return err
}


func (me *Gearbox) RestartSSH(sshArgs SshArgs) error {

	sshClient := NewSsh(*me, sshArgs)

	err := sshClient.RestartSsh()

	return err
}



// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewSsh(gb Gearbox, args ...SshArgs) *Ssh {
	var _args SshArgs
	if len(args)>0 {
		_args = args[0]
	}

	if _args.SshName == "" {
		_args.SshName = SshDefaultName
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = SshDefaultWaitDelay
	}

	if _args.WaitRetries == 0 {
		_args.WaitRetries = SshDefaultWaitRetries
	}

	if _args.SshHost == "" {
		_args.SshHost = SshDefaultSshHost
	}

	if _args.SshPort == "" {
		_args.SshPort = SshDefaultSshPort
	}

	if _args.SshOkString == "" {
		_args.SshOkString = SshDefaultSshOkString
	}

	if _args.SshWait == 0 {
		_args.SshWait = SshDefaultSshWait
	}

	if _args.SshWaitRetries == 0 {
		_args.SshWaitRetries = SshDefaultSshWaitRetries
	}

	publicKey, err := PublicKeyFile(`privatekey.pem`)
	if err != nil {
		log.Println(err)
		return
	}

	_args.Instance = ssh.Client{
		Name: _args.SshName,
	}

	ssh.Dial()
	sshClient := &Ssh{}
	*sshClient = Ssh(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return sshClient
}


func (me *Ssh) WaitForState(waitForState string, displayString string) error {

	var waitCount int

	spinner := util.NewSpinner(util.SpinnerArgs{
		Text: displayString,
		ExitOK: displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
	spinner.Start()

	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
		state, err := me.StatusSsh()
		if err != nil {
			spinner.Stop(false)
			return err
		}
		if state == waitForState {
			spinner.Stop(true)
			break
		}

		time.Sleep(me.WaitDelay)
		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
	}

	return nil
}


func (me *Ssh) WaitForSsh(displayString string) (string, error) {

	if me == nil {
		// Throw software error.
		return "", nil
	}

	state, err := me.Instance.GetState()
	if err != nil {
		return state, err
	}

	if state == SshRunning {
		spinner := util.NewSpinner(util.SpinnerArgs{
			Text: displayString,
			ExitOK: displayString + " - OK",
			ExitNOK: displayString + " - FAILED",
		})

		if me.ShowSsh == false {
			// We want to display just a spinner instead of console output.
			spinner.Start()
		}

		// connect to this socket
		conn, err := net.Dial("tcp", me.SshHost + ":" + me.SshPort)
		defer conn.Close()
		if err != nil {
			return SshUnknown, err
		}

		err = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(me.WaitRetries)))
		if err != nil {
			return SshUnknown, err
		}

		var waitCount int
		var readBuffer []byte
		for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {

			bytesRead, err := conn.Read(readBuffer)
			if err != nil {
				state = SshUnknown
				spinner.Stop(false)
				break
			}

			if bytesRead > 0 {
				if me.Show == true {
					fmt.Printf("%s\n", readBuffer)
				}
			}

			match, _ := regexp.MatchString(me.OkString, string(readBuffer))
			if match == true {
				state = SshGearBoxOK
				spinner.Stop(true)
				break
			}

			time.Sleep(me.WaitDelay)
			spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
		}
	}

	return state, nil
}


func (me *Ssh) StartSsh() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	state, err := me.StatusSsh()
	if err != nil {
		return err
	}
	if state == SshRunning || state == SshStarted {
		return nil
	}

	err = me.Instance.Start()
	if err != nil {
		return err
	}
	if me.NoWait == false {
		err := me.WaitForState(SshRunning, me.SshName + " SSH: Starting")
		if err != nil {
			return err
		}

		state, err = me.WaitFor(me.SshName + " : Starting")
		if err != nil {
			return err
		}
	}

	return nil
}


func (me *Ssh) StopSsh() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	state, err := me.StatusSsh()
	if err != nil {
		return err
	}
	if state == SshHalted {
		return nil
	}

	err = me.Instance.Halt()
	if err != nil {
		return err
	}
	if me.NoWait == false {
		err := me.WaitForState(SshHalted, me.SshName + " SSH: Stopping")
		if err != nil {
			return err
		}

		state, err = me.WaitFor(me.SshName + " : Stopping")
		if err != nil {
			return err
		}
	}

	return nil
}


func (me *Ssh) RestartSsh() error {

	var stateBefore string
	var stateAfter string
	var err error

	if me == nil {
		// Throw software error.
		return nil
	}

	stateBefore, err = me.StatusSsh()
	if err != nil {
		return err
	}

	switch {
		case stateBefore == SshRunning:
			fallthrough
		case stateBefore == SshStarted:
			err := me.StopSsh()
			if err != nil {
				return err
			}
			err = me.StartSsh()
			if err != nil {
				return err
			}

		case stateBefore == SshHalted:
			err := me.StartSsh()
			if err != nil {
				return err
			}

		case stateBefore == SshUnknown:

	}

	stateAfter, err = me.StatusSsh()
	if err != nil {
		return err
	}

	if stateAfter != SshRunning {
		// Throw an error.
	}

	return err
}


func (me *Ssh) StatusSsh() (string, error) {

	if me == nil {
		// Throw software error.
		return "", nil
	}

	state, err := me.Instance.GetState()
	if err != nil {
		return state, err
	}
	if state == SshRunning {
		state, err = me.WaitFor(me.SshName + " : Status")
		if err != nil {
			return state, err
		}
	}

	return state, nil
}


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

// Wait Indicators.
const WINew = 0
const WIStart = 1
const WISpin = 2
const WIStopOK = 3
const WIStopNOK = 4

func DefaultWaitIndicator(position int) {

	// Do nothing.
	switch position {
		case WINew:

		case WIStart:

		case WISpin:

		case WIStopOK:

		case WIStopNOK:
	}

}


func UserWaitIndicator(position int) {

	switch position {
		case WINew:

		case WIStart:
			fmt.Printf("\nWaiting")

		case WISpin:
			fmt.Printf(".")

		case WIStopOK:
			fmt.Printf("OK\n")

		case WIStopNOK:
			fmt.Printf("Not OK!\n")
	}

}

*/