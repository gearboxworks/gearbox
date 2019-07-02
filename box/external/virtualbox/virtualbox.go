package virtualbox

import (
	"fmt"
	"gearbox/ensure"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type VirtualBox struct {
}

func NewVirtualBox() *VirtualBox {
	vb := VirtualBox{}
	vb.Reset()
	return &vb
}

var logger Logger

func (me *VirtualBox) Reset() {
	logger = NewCommandLineLog()
}

// RunVm runs a VBoxManage command.
func RunVm(vm VirtualMachiner, args ...string) (exitCode string, err error) {

	var vboxManagePath string

	for range only.Once {

		path, err := exec.LookPath("VBoxManage")
		if err == nil {
			vboxManagePath = path
		} else {
			vboxManagePath = VBOXMANAGE
		}

		logger.Debug("EXEC[%v]: %v '%v'", vm.GetId(), vboxManagePath, strings.Join(args, ` `))
		cmd := exec.Command(vboxManagePath, args...)
		cmd.Stdout = logger.GetStdout()
		cmd.Stderr = logger.GetStderr()
		err = cmd.Run()
		if err == nil {
			break
		}
		exitCode = strings.TrimPrefix(err.Error(), "exit status ")

		switch exitCode {
		case exitCodeMissingVm:
			err = fmt.Errorf("VirtualBox[%s] command error '%v'", vm.GetName(), cmd.Stderr)

		default:
			err = fmt.Errorf("VirtualBox[%s] failed to run command '%v'",
				vm.GetName(),
				err.Error(),
			)
			fmt.Printf("stdout:%v\n", cmd.Stdout)
			fmt.Printf("stderr:%v\n", cmd.Stderr)
			fmt.Printf("returnCode:'%v'\n", exitCode)
		}

	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return
}

func CreateVm(vm VirtualMachiner) (states.State, error) {

	var err error
	var state states.State

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		state, err = CmdVmInfo(vm)
		switch state {
		case states.StateError:
			logger.Debug("%s: %v", vm.GetId(), err)

		case states.StateStopped:
			fallthrough
		case states.StateStarted:
			err = nil
			logger.Debug("VM '%s' already created", vm.GetId())

		case states.StateUnregistered:
			logger.Debug("creating VM '%s'", vm.GetId())
			err = CmdCreateVm(vm)
			if err != nil {
				break
			}

			logger.Debug("modifying VM '%s'", vm.GetId())
			err = CmdModifyVm(vm)
			if err != nil {
				break
			}

			state, err = CmdVmInfo(vm)

			logger.Debug("VM '%s' created OK", vm.GetId())

		case states.StateUnknown:
			logger.Debug(vm.GetId(), "%v", err)
		}

		if err != nil {
			logger.Debug(vm.GetId(), "VM couldn't be created with '%v'", err)
			_ = CmdDestroyVm(vm)
			break
		}
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return state, err
}

func DestroyVm(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		var state states.State
		state, err = CmdVmInfo(vm)
		switch state {
		case states.StateError:
			logger.Debug("VM '%s': %v", vm.GetId(), err)

		case states.StateStopped:
			logger.Debug("destroying VM '%s'", vm.GetId())
			err = CmdDestroyVm(vm)
			if err != nil {
				break
			}

			err = vm.DestroyConfig()
			if err != nil {
				break
			}

			logger.Debug(vm.GetId(), "VM destroyed OK")

		case states.StateStarted:
			err = msgs.MakeError("VM '%s' not stopped", vm.GetName())
			logger.Debug(vm.GetId(), "%v", err)

		case states.StateUnregistered:
			err = msgs.MakeError("VM '%s' doesn't exist", vm.GetName())
			logger.Debug(vm.GetId(), "%v", err)

		case states.StateUnknown:
			logger.Debug(vm.GetId(), "%v", err)

		}

	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func StartVm(vm VirtualMachiner, wait bool) (bool, error) {

	var err error
	var ok bool

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		var state states.State
		state, err = CmdVmInfo(vm)
		switch state {
		case states.StateError:
			logger.Debug("VM '%s': %v", vm.GetId(), err)

		case states.StateStopped:
			// stdout, stderr, err := RunVm(vm,"showvminfo", vm, "--machinereadable")
			_, err = RunVm(vm, "startvm", vm.GetName(), "--type", "headless")
			if err != nil {
				break
			}

			ok, err = WaitForVmState(vm, states.StateStarted)
			if err != nil {
				break
			}
			if !ok {
				break
			}
			logger.Debug("VM '%s' started OK", vm.GetId())

		case states.StateStarted:
			err = fmt.Errorf("VM '%s' already started", vm.GetName())
			logger.Debug("VM '%s': %v", vm.GetId(), err)
			fallthrough

		case states.StateUnregistered, states.StateUnknown:
			logger.Debug("VM '%s': %v", vm.GetId(), err)
		}

	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return ok, err
}

func StopVm(vm VirtualMachiner, force bool, wait bool) (bool, error) {

	var err error
	var ok bool

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		var state states.State
		state, err = CmdVmInfo(vm)
		switch state {
		case states.StateStopped:
			logger.Debug("VM '%s' already stopped", vm.GetId())

		case states.StateStarted:
			// stdout, stderr, err := RunVm(vm,"showvminfo", vm, "--machinereadable")
			if force {
				_, err = RunVm(vm, "controlvm", vm.GetName(), "poweroff")
			} else {
				_, err = RunVm(vm, "controlvm", vm.GetName(), "acpipowerbutton")
			}
			if err != nil {
				break
			}

			ok, err = WaitForVmState(vm, states.StateStarted)
			if err != nil {
				break
			}
			if !ok {
				break
			}

			logger.Debug("VM '%s' stopped OK", vm.GetId())

		case states.StateError, states.StateUnregistered, states.StateUnknown:
			logger.Debug("VM '%s': %v", vm.GetId(), err)
		}

	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return ok, err
}

func WaitForVmState(vm VirtualMachiner, want states.State) (bool, error) {

	var err error
	var ok bool

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		var state states.State
		for loop := 0; (state != want) && (loop < vm.GetRetryMax()); loop++ {
			state, err = CmdVmInfo(vm)
			if state == want {
				ok = true
			}
			time.Sleep(vm.GetRetryDelay())
		}
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return ok, err
}

func VmStatus(vm VirtualMachiner) (states.State, error) {

	var err error
	var state states.State

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		state, err = CmdVmInfo(vm)
	}
	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return state, err
}

func CmdDestroyVm(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "unregistervm", vm.GetName(), "--delete")
		if err != nil {
			break
		}

		err = CmdDestroyHostOnlyNet(vm)
		if err != nil {
			break
		}

		logger.Debug("destroyed VM '%s' OK", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdDestroyHostOnlyNet(vm VirtualMachiner) error {

	var err error

	// fmt.Printf("CmdCreateHostOnlyNet() ENTRY.\n")
	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}
		nic := vm.GetNic()
		if nic.Name == "" {
			logger.Debug("hostonlyif not created for VM '%s'", vm.GetId())
			break
		}

		_, err = RunVm(vm, "hostonlyif",
			"remove", nic.Name,
		)

		eblog.LogIfError(err) // Don't exit just yet, but log error.

		_, err = RunVm(vm, "dhcpserver", "remove",
			"--netname", fmt.Sprintf("HostInterfaceNetworking-%s", nic.Name),
		)
		if err != nil {
			break
		}

		nic.Name = ""
		vm.SetNic(nic)

		logger.Debug("hostonlyif for VM '%s' destroyed ok", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdModifyVm(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		err = CmdModifyVmBasic(vm)
		if err != nil {
			break
		}

		err = CmdCreateHostOnlyNet(vm)
		if err != nil {
			break
		}

		err = CmdModifyVmNetwork(vm)
		if err != nil {
			break
		}

		err = CmdModifyVmStorage(vm)
		if err != nil {
			break
		}

		err = CmdModifyVmIso(vm)
		if err != nil {
			break
		}

		logger.Debug("modified VM '%s' ok", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdModifyVmIso(vm VirtualMachiner) error {

	var err error

	// fmt.Printf("CmdModifyVmIso() ENTRY.\n")
	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "storagectl", vm.GetName(),
			"--name", "IDE",
			"--add", "ide",
			"--hostiocache", "on",
			"--bootable", "on",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "storageattach", vm.GetName(),
			"--storagectl", "IDE",
			"--port", "0",
			"--device", "0",
			"--type", "dvddrive",
			"--tempeject", "on",
			"--medium", vm.GetReleaser().GetFilepath(),
		)
		if err != nil {
			break
		}

		logger.Debug("modified VM '%s' ISO OK", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdModifyVmStorage(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "storagectl", vm.GetName(),
			"--name", "SATA",
			"--add", "sata",
			"--controller", "IntelAHCI",
			"--portcount", "4",
			"--hostiocache", "off",
			"--bootable", "on",
		)
		if err != nil {
			break
		}

		// SIGH - Needs to be not hard-coded.
		disks := Disks{}
		disks = append(disks, Disk{
			Name:   "Gearbox-Opt.vdi",
			Format: "VDI",
			Size:   "1024",
		})
		disks = append(disks, Disk{
			Name:   "Gearbox-Docker.vdi",
			Format: "VDI",
			Size:   "16384",
		})
		disks = append(disks, Disk{
			Name:   "Gearbox-Projects.vdi",
			Format: "VDI",
			Size:   "16384",
		})
		disks = append(disks, Disk{
			Name:   "Gearbox-Config.vdi",
			Format: "VDI",
			Size:   "1024",
		})

		for index, disk := range disks {

			fileName := osdirs.AddPaths(vm.GetVmDir(), vm.GetName())
			fileName = osdirs.AddFilef(fileName, disk.Name)
			order := strconv.Itoa(index)

			_, err = RunVm(vm, "createmedium", "disk",
				"--filename", fileName,
				"--size", disk.Size,
				"--format", disk.Format,
				"--variant", "Standard",
			)
			if err != nil {
				break
			}

			_, err = RunVm(vm, "storageattach", vm.GetName(),
				"--storagectl", "SATA",
				"--port", order,
				"--device", "0",
				"--type", "hdd",
				"--medium", fileName,
				"--hotpluggable", "off",
			)
			if err != nil {
				break
			}
		}

		logger.Debug("modified VM '%s' storage OK", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdModifyVmNetwork(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		if vm.GetNic().Name == "" {
			logger.Debug("hostonlyif not created for VM '%s'", vm.GetId())
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--nic2", "nat",
			"--nictype2", "82540EM",
			"--cableconnected2", "on",
			"--macaddress2", "auto",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--natnet2", "default",
			"--natpf2", "API,tcp,,9970,,9970",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--natnet2", "default",
			"--natpf2", fmt.Sprintf("SSH,tcp,,%s,,22", vm.GetSsh().GetPort()),
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--natnet2", "default",
			"--natpf2", "VMcontrol,tcp,,9971,,9971",
		)
		if err != nil {
			break
		}

		// _, err = RunVm(vm,"modifyvm", vm.Boxname, "--nic2", "bridged", "--bridgeadapter2", nic["Name"], "--nictype2", "82540EM", "--cableconnected2", "on", "--macaddress2", "auto", "--nicpromisc2", "deny")
		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--nic1", "hostonly",
			"--hostonlyadapter1", vm.GetNic().Name,
			"--nictype1", "82540EM",
			"--cableconnected1", "on",
			"--macaddress1", "auto",
			"--nicpromisc1", "deny",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--uart1", "0x3f8", "4",
			"--uartmode1", "tcpserver", vm.GetConsole().GetPort(),
		)
		if err != nil {
			break
		}

		logger.Debug("modified VM '%s' ok", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdFindHostOnlyNet(vm VirtualMachiner) (string, error) {

	var err error
	var ret string

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "list", "hostonlyifs", "-ls")
		if err != nil {
			break
		}

		var nic KeyValueMap
		dr, ok := decodeResponse(logger.GetStdout(), ':')
		if ok == true {
			nics, ok := dr.decodeNics()
			if ok == false {
				err = fmt.Errorf("no NICs found for VM '%s", vm.GetName())
				break
			}
			vm.SetNics(nics)
			var id string
			for id, nic = range vm.GetNics() {
				a, ok := nic["IPAddress"]
				if !ok {
					logger.Error("NIC '%s' for VM '%s' does not have an 'IPAddress' entry", nic, vm.GetId())
					continue
				}
				if a != vm.GetNic().Ip {
					continue
				}
				ret = id
				break
			}
		}

		logger.Debug("using NIC '%s' for VM '%s'", nic, vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return ret, err
}

func CmdCreateHostOnlyNet(vm VirtualMachiner) error {

	var err error

	// fmt.Printf("CmdCreateHostOnlyNet() ENTRY.\n")
	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}
		nic := NewHostOnlyNic()
		err = nic.initialize(vm)
		if err != nil {
			break
		}
		vm.SetNic(nic)

		if nic.Name != "" {
			logger.Debug("hostonlyif already created for VM '%s'", vm.GetId())
			break
		}

		// Such a pain in the neck this process...
		_, err = RunVm(vm, "hostonlyif", "create")
		if err != nil {
			break
		}

		// Now we need to parse the output and look for this:
		// Interface 'vboxnet2' was successfully created
		re := regexp.MustCompile(`Interface '(\w+)' was successfully created`)
		f := re.FindAllStringSubmatch(logger.GetStdout().String(), -1)
		if len(f) == 0 {
			err = fmt.Errorf("hostonlyif not created for VM '%s'", vm.GetName())
			break
		}
		if len(f[0]) == 0 {
			err = fmt.Errorf("hostonlyif not created for VM '%s'", vm.GetName())
			break
		}
		nic.Name = f[0][1]
		vm.SetNic(nic)

		_, err = RunVm(vm, "hostonlyif",
			"ipconfig", nic.Name,
			"--ip", nic.Ip,
			"--netmask", nic.Netmask,
		)
		if err != nil {
			break
		}

		args := []string{
			"--enable",
			"--netname", "HostInterfaceNetworking-" + nic.Name,
			"--ip", nic.Ip,
			"--netmask", nic.Netmask,
			"--lowerip", nic.DhcpLowerIp,
			"--upperip", nic.DhcpUpperIp,
		}
		create := append([]string{"dhcpserver", "add"}, args...)
		_, err = RunVm(vm, create...)
		if err != nil {
			// Now we need to parse the output and look for this:
			// DHCP server already exists
			re := regexp.MustCompile(`DHCP server already exists`)
			f := re.FindAllStringSubmatch(logger.GetStderr().String(), -1)
			if len(f) == 0 {
				// Some other error.
				break
			}

			// Seems we have one already. So we want to modify it.
			modify := append([]string{"dhcpserver", "modify"}, args...)
			_, err = RunVm(vm, modify...)
			if err != nil {
				break
			}
		}

		logger.Debug("hostonlyif created for VM '%s' ok", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdModifyVmBasic(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		err = osdirs.CheckFileExists(vm.GetIconFile())
		if err == nil {
			// Just don't add the icon if the file doesn't exist.
			// exitcode, Err = RunVm(vm, "modifyvm", vm.Boxname, "--description", vm.Boxname + " OS VM", "--iconfile", vm.Entry.IconFile)
			name := vm.GetName()
			_, err = RunVm(vm, "modifyvm", name,
				"--description", fmt.Sprintf("%s OS VM", name),
				"--iconfile", vm.GetIconFile(),
			)
			if err != nil {
				break
			}
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--ioapic", "on",
			"--acpi", "on",
			"--biosbootmenu", "disabled",
			"--biosapic", "apic",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--boot1", "dvd",
			"--boot2", "none",
			"--boot3", "none",
			"--boot4", "none",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--vrde", "off",
			"--autostart-enabled", "off",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--accelerate3d", "off",
			"--accelerate2dvideo", "off",
			"--mouse", "usbtablet",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--defaultfrontend", "headless",
			"--snapshotfolder", "default",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--memory", "2048",
			"--vram", "128",
		)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "modifyvm", vm.GetName(),
			"--audio", "none",
		)
		if err != nil {
			break
		}

		logger.Debug("modification of VM '%s''s basic config ok", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdCreateVm(vm VirtualMachiner) error {

	var err error

	// fmt.Printf("CmdCreateVm() ENTRY.\n")
	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		// exitcode, err = RunVm(vm,"createvm", "--name", vm.Boxname, "--ostype", "Linux26_64", "--register", "--basefolder", vm.VmBaseDir)
		_, err = RunVm(vm, "createvm",
			"--name", vm.GetName(),
			"--ostype", "Linux26_64",
			"--register",
			"--basefolder", vm.GetVmDir(),
		)
		if err != nil {
			break
		}

		logger.Debug("created VM '%s' ok", vm.GetId())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}

func CmdVmInfo(vm VirtualMachiner) (states.State, error) {

	var err error
	var exitCode string
	var state states.State
	// states.StateUnregistered
	// states.StateError
	// states.StateStopped
	// states.StateStarted
	// states.StateUnknown

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		for range only.Once {
			exitCode, err = RunVm(vm, "showvminfo", vm.GetName(), "--machinereadable")
			if err != nil {
				if exitCode == exitCodeMissingVm {
					state = states.StateUnregistered
					err = fmt.Errorf("VM '%s' is not registered", vm.GetName())
				} else {
					state = states.StateError
				}
				break
			}

			kvs, ok := decodeResponse(logger.GetStdout(), '=')
			if ok == false {
				state = states.StateError
				err = fmt.Errorf("VM '%s' can't decode showvminfo response", vm.GetName())
				break
			}

			var kvm KeyValueMap
			kvm, ok = kvs.decodeShowVmInfo()
			if ok == false {
				state = states.StateError
				err = fmt.Errorf("VM '%s' can't decode showvminfo output", vm.GetName())
				break
			}
			vm.SetInfo(kvm)

			switch vm.GetInfoValue("VMState") {
			case "poweroff":
				fallthrough
			case "paused":
				fallthrough
			case "saved":
				state = states.StateStopped

			case "running":
				state = states.StateStarted

			default:
				state = states.StateUnknown
				err = fmt.Errorf("VM '%s' is in an unknown state", vm.GetName())
			}

		}

		logger.Debug("VM '%s' is in state %s", vm.GetId(), state)
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return state, err
}
