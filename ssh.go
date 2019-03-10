package gearbox

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type SSH struct {
	Instance *ssh.Client

	//	// Status polling delays.
	//	NoWait        bool
	//	WaitDelay     time.Duration
	//	WaitRetries   int

	// SSH related.
	Username   string
	Password   string
	Host       string
	Port       string
	PublicKey  string
	StatusLine StatusLine
	//	OkString    string
	//	Wait        time.Duration
}
type StatusLine struct {
	Text        string
	Disable     bool
	UpdateDelay time.Duration
}

type SshArgs SSH

const SshDefaultUsername = "boxuser"
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
func NewSsh(gb Gearbox, args ...SshArgs) (*SSH, error) {

	var err error
	var _args SshArgs
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.Username == "" {
		_args.Username = SshDefaultUsername
	}

	if _args.Password == "" {
		_args.Password = SshDefaultPassword
	}

	if _args.StatusLine.UpdateDelay == 0 {
		_args.StatusLine.UpdateDelay = SshDefaultStatusLineUpdateDelay
	}

	if _args.Host == "" {
		_args.Host = SshDefaultSshHost
	}

	if _args.Port == "" {
		_args.Port = SshDefaultSshPort
	}

	if _args.PublicKey == "" {
		_args.PublicKey = SshDefaultKeyFile
	}

	sshConfig := &ssh.ClientConfig{}

	// Try SSH key file first.
	keyfile, err := readPublicKeyFile(_args.PublicKey)
	if err == nil && keyfile != nil {
		// Authenticate using SSH key.
		authenticate := []ssh.AuthMethod{keyfile}
		sshConfig = &ssh.ClientConfig{
			User: _args.Username,
			Auth: authenticate,
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 10,
		}
	} else {
		sshConfig = &ssh.ClientConfig{
			User: _args.Username,
			Auth: []ssh.AuthMethod{ssh.Password(_args.Password)},
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 10,
		}
	}

	_args.Instance, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", _args.Host, _args.Port), sshConfig)
	if err != nil {
		fmt.Printf("Gearbox SSH error: %s\n", err)
		return nil, err
	}

	sshClient := &SSH{}
	*sshClient = SSH(_args)

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

func (me *SSH) StartSsh() error {

	if me == nil || me.Instance == nil {
		// Throw software error.
		return nil
	}

	session, err := me.Instance.NewSession()
	defer session.Close()
	defer me.Instance.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	// Set IO
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	//	in, _ := session.StdinPipe()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 19200,
		ssh.TTY_OP_OSPEED: 19200,
	}

	// Request pseudo terminal
	termWidth := 0
	termHeight := 0
	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			return nil
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err = terminal.GetSize(fileDescriptor)
		if err != nil {
			return nil
		}

		// xterm-256color
		err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			return nil
		}
	}

	go me.StatusLineWorker(termHeight, termWidth)

	go me.exampleHostWorker()

	// Start remote shell
	err = session.Shell()
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
	session.Wait()

	return nil
}

func (me *SSH) StatusLineWorker(termHeight int, termWidth int) {
	const savePos = "[s"
	const restorePos = "[u"
	bottomPos := fmt.Sprintf("[%d;0H", termHeight)
	scrollFix := fmt.Sprintf("[1;%dr", termHeight-1)
	if me.StatusLine.Disable == false {
		fmt.Printf(scrollFix)
	}

	for {
		if me.StatusLine.Disable == false {
			fmt.Printf("%s%s%s%s", savePos, bottomPos, me.StatusLine.Text, restorePos)
		}

		time.Sleep(me.StatusLine.UpdateDelay)
	}
}

func (me *SSH) SetStatusLine(text string) {

	me.StatusLine.Text = text
	// fmt.Printf("%s%s%d%s", savePos, bottomPos, time.Now().Unix(), restorePos)
}

// Example host worker. This periodically changes the me.Text from the host side.
// The StatusLineWorker() will update the bottom line using the me.Text.
func (me *SSH) exampleHostWorker() {

	yellow := color.New(color.BgBlack, color.FgHiYellow).SprintFunc()
	magenta := color.New(color.BgBlack, color.FgHiMagenta).SprintFunc()
	green := color.New(color.BgBlack, color.FgHiGreen).SprintFunc()
	normal := color.New(color.BgWhite, color.FgHiBlack).SprintFunc()

	for {
		now := time.Now()
		dateStr := normal("Date:") + " " + yellow(fmt.Sprintf("%.4d/%.2d/%.2d", now.Year(), now.Month(), now.Day()))
		timeStr := normal("Time:") + " " + magenta(fmt.Sprintf("%.2d:%.2d:%.2d", now.Hour(), now.Minute(), now.Second()))
		statusStr := normal("Status:") + " " + green("OK")

		line := fmt.Sprintf("%s	%s %s", statusStr, dateStr, timeStr)

		me.SetStatusLine(line)

		time.Sleep(time.Second * 5)
	}
}
