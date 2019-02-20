package gearbox

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Ssh struct {
	Instance *ssh.Client

	// Status polling delays.
	NoWait        bool
	WaitDelay     time.Duration
	WaitRetries   int

	// SSH related.
	SshUsername  string
	SshPassword  string
	SshHost        string
	SshPort        string
	SshPublicKey   string

	SshOkString    string
	SshWait        time.Duration
}
type SshArgs Ssh


const SshDefaultUsername = "boxuser"
const SshDefaultPassword = "box"
const SshDefaultKeyFile = "./keyfile.pub"
const SshDefaultWaitDelay = time.Second
const SshDefaultSshHost = "127.0.0.1"
const SshDefaultSshPort = "2222"
const SshDefaultSshWait = time.Second


const SshUnknown = "unknown"
const SshHalted = "halted"
const SshRunning = "running"
const SshStarted = "started"
const SshGearBoxOK = "ok"


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
	if len(args)>0 {
		_args = args[0]
	}

	if _args.SshUsername == "" {
		_args.SshUsername = SshDefaultUsername
	}

	if _args.SshPassword == "" {
		_args.SshPassword = SshDefaultPassword
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = SshDefaultWaitDelay
	}

	if _args.SshHost == "" {
		_args.SshHost = SshDefaultSshHost
	}

	if _args.SshPort == "" {
		_args.SshPort = SshDefaultSshPort
	}

	if _args.SshWait == 0 {
		_args.SshWait = SshDefaultSshWait
	}

	if _args.SshPublicKey == "" {
		_args.SshPublicKey = SshDefaultKeyFile
	}


	sshConfig := &ssh.ClientConfig{}

	// Try SSH key file first.
	keyfile, err := readPublicKeyFile(_args.SshPublicKey)
	if err == nil && keyfile != nil {
		// Authenticate using SSH key.
		authenticate := []ssh.AuthMethod{keyfile}
		sshConfig = &ssh.ClientConfig{
			User: _args.SshUsername,
			Auth: authenticate,
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:time.Second * 10,
		}
	} else {
		sshConfig = &ssh.ClientConfig{
			User: _args.SshUsername,
			Auth: []ssh.AuthMethod{ssh.Password(_args.SshPassword)},
			// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:time.Second * 10,
		}
	}

	_args.Instance, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", _args.SshHost, _args.SshPort), sshConfig)
	if err != nil {
		fmt.Printf("Couldn't establish a connection to the remote server:'%s'", err)
		return nil, err
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
	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			return nil
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)
		if err != nil {
			return nil
		}

		// xterm-256color
		err = session.RequestPty("vt100", termHeight, termWidth, modes)
		if err != nil {
			return nil
		}
	}

	// Start remote shell
	err = session.Shell()
	if err != nil {
		fmt.Printf("Can't start shell: %s", err)
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

