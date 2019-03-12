package gearbox

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

type Ssh struct {
	Instance *ssh.Client
	Session *ssh.Session

	//	// Status polling delays.
	//	NoWait        bool
	//	WaitDelay     time.Duration
	//	WaitRetries   int

	// SSH related.
	SshUsername  string
	SshPassword  string
	SshPublicKey string
	SshHost      string
	SshPort      string

	//	SshOkString    string
	//	SshWait        time.Duration
	SshStatusLine         string
	DisableStatusLine     bool
	StatusLineUpdateDelay time.Duration
	SshTermWidth int
	SshTermHeight int
	SshTerminateFlag	bool
}
type SshArgs Ssh

const SshDefaultUsername = "gearbox"
const SshDefaultPassword = "box"
const SshDefaultKeyFile = "./keyfile.pub"
const SshDefaultSshHost = "127.0.0.1"
const SshDefaultSshPort = "2222"
const SshDefaultStatusLineUpdateDelay = time.Second * 2

// //////////////////////////////////////////////////////////////////////////////
// Gearbox related
func (me *Gearbox) ConnectSSH(sshArgs SshArgs) error {

	sshClient, err := NewSsh(*me, sshArgs)
	if err != nil {
		return err
	}

	err = sshClient.StartSsh()

	return err
}

// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewSsh(gb Gearbox, args ...SshArgs) (*Ssh, error) {

	var err error
	var _args SshArgs
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.SshUsername == "" {
		_args.SshUsername = SshDefaultUsername
	}

	if _args.SshPassword == "" {
		_args.SshPassword = SshDefaultPassword
	}

	if _args.SshPublicKey == "" {
		_args.SshPublicKey = SshDefaultKeyFile
	}

	if _args.StatusLineUpdateDelay == 0 {
		_args.StatusLineUpdateDelay = SshDefaultStatusLineUpdateDelay
	}

	if _args.SshHost == "" {
		_args.SshHost = SshDefaultSshHost
	}

	if _args.SshPort == "" {
		_args.SshPort = SshDefaultSshPort
	}

	sshClient := &Ssh{}
	*sshClient = Ssh(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return sshClient, err
}

func readPublicKeyFile(file string) (ssh.AuthMethod, error) {

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		// fmt.Printf("# Error reading file '%s': %s\n", file, err)
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		// fmt.Printf("# Error parsing key '%s': %s\n", signer, err)
		return nil, err
	}

	return ssh.PublicKeys(signer), err
}

func (me *Ssh) StartSsh() error {

	var err error

	if me == nil {
		// Throw software error.
		return nil
	}

	sshConfig := &ssh.ClientConfig{}

	// Try SSH key file first.
	keyfile, err := readPublicKeyFile(me.SshPublicKey)
	if err == nil && keyfile != nil {
		// Authenticate using SSH key.
		authenticate := []ssh.AuthMethod{keyfile}
		sshConfig = &ssh.ClientConfig{
			User: me.SshUsername,
			Auth: authenticate,
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 10,
		}
	} else {
		sshConfig = &ssh.ClientConfig{
			User: me.SshUsername,
			Auth: []ssh.AuthMethod{ssh.Password(me.SshPassword)},
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 10,
		}
	}

	me.Instance, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", me.SshHost, me.SshPort), sshConfig)
	if err != nil {
		fmt.Printf("Gearbox SSH error: %s\n", err)
		return err
	}

	me.Session, err = me.Instance.NewSession()
	defer me.Session.Close()
	defer me.Instance.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	// Set IO
	me.Session.Stdout = os.Stdout
	me.Session.Stderr = os.Stderr
	me.Session.Stdin = os.Stdin
	//	in, _ := me.Session.StdinPipe()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 19200,
		ssh.TTY_OP_OSPEED: 19200,
	}

	// Request pseudo terminal
	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			return nil
		}
		defer terminal.Restore(fileDescriptor, originalState)

		me.SshTermWidth, me.SshTermHeight, err = terminal.GetSize(fileDescriptor)
		if err != nil {
			return nil
		}

		// xterm-256color
		err = me.Session.RequestPty("xterm-256color", me.SshTermHeight, me.SshTermWidth, modes)
		if err != nil {
			return nil
		}
	}

	go me.StatusLineWorker()

	go me.exampleHostWorker()

	// Start remote shell
	err = me.Session.Shell()
	if err != nil {
		fmt.Printf("Gearbox SSH error: %s\n", err)
	}

	/*
		// Loop around input <-> output.
		for {
			reader := bufio.NewReader(os.Stdin)
			str, _ := reader.ReadString('\n')
			fmt.Fprint(in, str)
		}
	*/
	me.Session.Wait()
	me.resetView()

	return nil
}

// StatusLineWorker() - handles the actual updates to the status line
func (me *Ssh) StatusLineWorker() {

	me.setView()
	// w := gob.NewEncoder(me.Session)
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGWINCH)

	for ; me.SshTerminateFlag == false; {
		// Handle terminal windows size changes properly.
		fileDescriptor := int(os.Stdin.Fd())
		width, height, _ := terminal.GetSize(fileDescriptor)
		if (me.SshTermWidth != width) || (me.SshTermHeight != height) {
			me.SshTermWidth = width
			me.SshTermHeight = height
			// me.Session.Signal(syscall.SIGWINCH)
			me.Session.WindowChange(height, width)
		} else {
			// Only update if we haven't seen a SIGWINCH - just to wait for things to settle.
			me.displayStatusLine()
		}

		time.Sleep(me.StatusLineUpdateDelay)
	}

}


func (me *Ssh) SetStatusLine(text string) {

	me.SshStatusLine = text
}


func (me *Ssh) displayStatusLine() {
	const savePos = "\033[s"
	const restorePos = "\033[u"
	bottomPos := fmt.Sprintf("\033[%d;0H", me.SshTermHeight)
	// topPos := fmt.Sprintf("\033[0;0H")

	if me.DisableStatusLine == false {
		fmt.Printf("%s%s%s%s", savePos, bottomPos, me.SshStatusLine, restorePos)
	}
}


func (me *Ssh) setView() {
	const clearScreen = "\033[H\033[2J"
	scrollFixBottom := fmt.Sprintf("\033[1;%dr", me.SshTermHeight-1)
	// scrollFixTop := fmt.Sprintf("\033[2;%dr", termHeight)

	if me.DisableStatusLine == false {
		fmt.Printf(scrollFixBottom)
		fmt.Printf(clearScreen)
	}
}


func (me *Ssh) resetView() {
	const savePos = "\033[s"
	const restorePos = "\033[u"
	scrollFixBottom := fmt.Sprintf("\033[1;%dr", me.SshTermHeight)
	// scrollFixTop := fmt.Sprintf("\033[2;%dr", termHeight)

	if me.DisableStatusLine == false {
		fmt.Printf(savePos)
		fmt.Printf(scrollFixBottom)
		fmt.Printf(restorePos)

		me.SshStatusLine = ""
		for i := 0; i <= me.SshTermWidth; i++ {
			me.SshStatusLine += " "
		}
		me.displayStatusLine()
	}

}


func stripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)

	return re.ReplaceAllString(str, "")
}


// Example host worker. This periodically changes the me.SshStatusLine from the host side.
// The StatusLineWorker() will update the bottom line using the me.SshStatusLine.
func (me *Ssh) exampleHostWorker() {

	yellow := color.New(color.BgBlack, color.FgHiYellow).SprintFunc()
	magenta := color.New(color.BgBlack, color.FgHiMagenta).SprintFunc()
	green := color.New(color.BgBlack, color.FgHiGreen).SprintFunc()
	normal := color.New(color.BgWhite, color.FgHiBlack).SprintFunc()

	for ; me.SshTerminateFlag == false; {
		//now := time.Now()
		//dateStr := normal("Date:") + " " + yellow(fmt.Sprintf("%.4d/%.2d/%.2d", now.Year(), now.Month(), now.Day()))
		//timeStr := normal("Time:") + " " + magenta(fmt.Sprintf("%.2d:%.2d:%.2d", now.Hour(), now.Minute(), now.Second()))
		statusStr := normal("Status:") + " " + green("OK")
		infoStr := yellow("You are connected to") + " " + magenta("Gearbox OS")

		//line := fmt.Sprintf("%s	%s %s", statusStr, dateStr, timeStr)
		line := fmt.Sprintf("%s - %s", infoStr, statusStr)

		// Add spaces to ensure it's right justified.
		spaces := ""
		lineLen := len(stripAnsi(line))
		for i:=0; i < me.SshTermWidth - lineLen; i++ {
			spaces += " "
		}

		me.SetStatusLine(spaces + line) // + fmt.Sprintf("W:%d L:%d S:%d C:%d", me.SshTermWidth, len(line), len(spaces), lineLen))

		time.Sleep(time.Second * 5)
	}
}

