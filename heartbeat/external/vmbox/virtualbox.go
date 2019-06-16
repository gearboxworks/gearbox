package vmbox

import (
	"bytes"
	"fmt"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/states"
	"gearbox/only"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Disk struct {
	Name   string
	Format string
	Size   string
}
type Disks []Disk


func (me *Vm) vbCreate() (states.State, error) {

	var err error
	var state states.State

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		state, err = me.cmdVmInfo()
		switch state {
			case states.StateError:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateStopped:
				fallthrough
			case states.StateStarted:
				err = me.EntityName.ProduceError("VM already created")
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateUnregistered:
				eblog.Debug(me.EntityId, "creating VM")
				err = me.cmdCreateVm()
				if err != nil {
					break
				}

				eblog.Debug(me.EntityId, "modifying VM")
				err = me.cmdModifyVm()
				if err != nil {
					break
				}

				state, err = me.cmdVmInfo()
				if err != nil {
					err = me.EntityName.ProduceError("VM couldn't be created")
					eblog.Debug(me.EntityId, "%v", err)
					break
				}

				eblog.Debug(me.EntityId, "VM created OK")


			case states.StateUnknown:
				eblog.Debug(me.EntityId, "%v", err)
		}

	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return state, err
}


func (me *Vm) vbDestroy() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var state states.State
		state, err = me.cmdVmInfo()
		switch state {
			case states.StateError:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateStopped:
				eblog.Debug(me.EntityId, "destroying VM")
				err = me.cmdDestroyVm()
				if err != nil {
					break
				}

				eblog.Debug(me.EntityId, "VM destroyed OK")


			case states.StateStarted:
				err = me.EntityName.ProduceError("VM not stopped")
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateUnregistered:
				err = me.EntityName.ProduceError("VM doesn't exist")
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateUnknown:
				eblog.Debug(me.EntityId, "%v", err)
		}

	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) vbStart() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var state states.State
		state, err = me.cmdVmInfo()
		switch state {
			case states.StateError:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateStopped:
				// stdout, stderr, err := me.Run("showvminfo", vm, "--machinereadable")
				_, err = me.Run("startvm", me.EntityName.String(), "--type", "headless")
				if err != nil {
					break
				}

				state, err = me.cmdVmInfo()
				if err != nil {
					break
				}

				eblog.Debug(me.EntityId, "VM started OK")


			case states.StateStarted:
				err = me.EntityName.ProduceError("VM already started")
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateUnregistered:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateUnknown:
				eblog.Debug(me.EntityId, "%v", err)
		}

	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) vbStop() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var state states.State
		state, err = me.cmdVmInfo()
		switch state {
			case states.StateError:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateStopped:
				eblog.Debug(me.EntityId, "VM already stopped")


			case states.StateStarted:
				// stdout, stderr, err := me.Run("showvminfo", vm, "--machinereadable")
				_, err = me.Run("controlvm", me.EntityName.String(), "poweroff")		// @stupid vbox "acpipowerbutton")
				if err != nil {
					break
				}

				state, err = me.cmdVmInfo()
				if err != nil {
					break
				}

				eblog.Debug(me.EntityId, "VM stopped OK")


			case states.StateUnregistered:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateUnknown:
				eblog.Debug(me.EntityId, "%v", err)
		}

	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) vbStatus() (states.State, error) {

	var err error
	var state states.State

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		state, err = me.cmdVmInfo()
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return state, err
}


func (me *Vm) vbWaitForState(want states.State) (bool, error) {

	var err error
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var state states.State
		for loop := 0; (state != want) && (loop < me.Entry.retryMax); loop++ {
			state, err = me.cmdVmInfo()
			//switch state {
			//	case states.StateError:
			//		eblog.Debug(me.EntityId, "%v", err)
			//
			//	case states.StateUnregistered:
			//		eblog.Debug(me.EntityId, "%v", err)
			//
			//	case states.StateUnknown:
			//		eblog.Debug(me.EntityId, "%v", err)
			//}

			time.Sleep(me.Entry.retryDelay)
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ok, err
}



////////////////////////////////////////////////////////////////////////
func (me *Vm) cmdVmInfo() (states.State, error) {

	var err error
	var exitCode string
	var state states.State
	// states.StateUnregistered
	// states.StateError
	// states.StateStopped
	// states.StateStarted
	// states.StateUnknown

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			// stdout, stderr, err := me.Run("showvminfo", vm, "--machinereadable")
			exitCode, err = me.Run("showvminfo", me.EntityName.String(), "--machinereadable")
			if err != nil {
				if exitCode == exitCodeMissingVm {
					state = states.StateUnregistered
					err = me.EntityName.ProduceError("VM not registered")
				} else {
					state = states.StateError
				}
				break
			}

			kvs, ok := decodeResponse(me.Entry.cmdStdout, '=')
			if ok == false {
				state = states.StateError
				err = me.EntityName.ProduceError("can't decode showvminfo response")
				break
			}

			var kvm KeyValueMap
			kvm, ok = kvs.decodeShowVmInfo()
			if ok == false {
				state = states.StateError
				err = me.EntityName.ProduceError("can't decode showvminfo output")
				break
			}
			me.Entry.vmInfo = kvm

			switch me.Entry.vmInfo["VMState"] {
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
					err = me.EntityName.ProduceError("VM is in an unknown state")
			}

			// 	VmStatePowerOff 	= "poweroff"	// Valid VM state return from listvm
			//	VmStatePaused 		= "paused"		// Valid VM state return from listvm
			//	VmStateSaved 		= "saved"		// Valid VM state return from listvm
			//	VmStateRunning  	= "running"		// Valid VM state return from listvm
			//
			//if is.Error(sts) {
			//	me.State.VM.CurrentState = VmStateNotPresent
			//}
			//
			//if me.State.VM.WantState == VmStateInit {
			//	me.State.VM.WantState = me.State.VM.CurrentState
			//}
			//
			//	// First check on the VM.
			//	// state, err := me.VmInstance.GetState()
			//	switch kvm["VMState"] {
			//		case VmStateRunning:
			//			me.State.VM.CurrentState = VmStateRunning
			//
			//		case VmStatePowerOff:
			//			fallthrough
			//		case VmStateSaved:
			//			fallthrough
			//		case VmStatePaused:
			//			me.State.VM.CurrentState = VmStatePowerOff
			//			me.State.API.CurrentState = VmStatePowerOff
			//	}
			//me.State.VM.CurrentState = kvm["VMState"]
		}

		eblog.Debug(me.EntityId, "VM is in state %s", state)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return state, err
}


func (me *Vm) cmdCreateVm() error {

	var err error

	// fmt.Printf("cmdCreateVm() ENTRY.\n")
	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// stdout, stderr, sts = me.Run("createvm", "--name", me.Boxname, "--ostype", "Linux26_64", "--register", "--basefolder", me.VmBaseDir)
		_, err = me.Run("createvm", "--name", me.EntityName.String(), "--ostype", "Linux26_64", "--register", "--basefolder", me.baseDir.String())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "created VM OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) cmdModifyVm() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.cmdModifyVmBasic()
		if err != nil {
			break
		}

		err = me.cmdModifyVmNetwork()
		if err != nil {
			break
		}

		err = me.cmdModifyVmStorage()
		if err != nil {
			break
		}

		err = me.cmdModifyVmIso()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "modified VM OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) cmdDestroyVm() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// stdout, stderr, sts = me.Run("unregistervm", me.Boxname, "--delete")
		_, err = me.Run("unregistervm", me.EntityName.String(), "--delete")
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) cmdModifyVmBasic() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// stdout, stderr, sts = me.Run("modifyvm", me.Boxname, "--description", me.Boxname + " OS VM", "--iconfile", string(me.OsSupport.GetAdminRootDir()) + "/" + IconLogo)
		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--description", me.EntityName.String() + " OS VM", "--iconfile", me.osPaths.UserConfigDir.AddFileToPath(IconLogoPng).String())
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--ioapic", "on", "--acpi", "on", "--biosbootmenu", "disabled", "--biosapic", "apic")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--boot1", "dvd", "--boot2", "none", "--boot3", "none", "--boot4", "none")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--vrde", "off", "--autostart-enabled", "off")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--cpuhotplug", "on", "--cpus", "4", "--pae", "off", "--longmode", "on", "--largepages", "on", "--paravirtprovider", "default")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--accelerate3d", "off", "--accelerate2dvideo", "off", "--mouse", "usbtablet")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--defaultfrontend", "headless", "--snapshotfolder", "default")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--memory", "2048", "--vram", "128")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--audio", "none")
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "modified basic VM config OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) cmdModifyVmNetwork() error {

	var err error

	// fmt.Printf("cmdModifyVmNetwork() ENTRY.\n")
	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--nic1", "nat", "--nictype1", "82540EM", "--cableconnected1", "on", "--macaddress1", "auto")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--natnet1", "default", "--natpf1", "API,tcp,,9970,,9970")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--natnet1", "default", "--natpf1", "SSH,tcp,," + me.Entry.SshPort + ",,22")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--natnet1", "default", "--natpf1", "VMcontrol,tcp,,9971,,9971")
		if err != nil {
			break
		}

		// _, _, sts = me.Run("modifyvm", me.Boxname, "--nic2", "bridged", "--bridgeadapter2", nic["Name"], "--nictype2", "82540EM", "--cableconnected2", "on", "--macaddress2", "auto", "--nicpromisc2", "deny")
		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--nic2", "hostonly", "--hostonlyadapter2", "vboxnet0", "--nictype2", "82540EM", "--cableconnected2", "on", "--macaddress2", "auto", "--nicpromisc2", "deny")
		if err != nil {
			break
		}

		_, err = me.Run("modifyvm", me.EntityName.String(),
			"--uart1", "0x3f8", "4", "--uartmode1", "tcpserver", me.Entry.ConsolePort)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "modified VM OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) cmdModifyVmStorage() error {

	var err error

	// fmt.Printf("cmdModifyVmStorage() ENTRY.\n")
	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		_, err = me.Run("storagectl", me.EntityName.String(),
			"--name", "SATA", "--add", "sata", "--controller", "IntelAHCI", "--portcount", "4", "--hostiocache", "off", "--bootable", "on")
		if err != nil {
			break
		}

		// SIGH - Needs to be not hard-coded.
		disks := Disks{}
		disks = append(disks, Disk{
			Name: "Gearbox-Opt.vmdk",
			Format: "VMDK",
			Size: "1024",
		})
		disks = append(disks, Disk{
			Name: "Gearbox-Docker.vmdk",
			Format: "VMDK",
			Size: "16384",
		})
		disks = append(disks, Disk{
			Name: "Gearbox-Projects.vmdk",
			Format: "VMDK",
			Size: "16384",
		})
		disks = append(disks, Disk{
			Name: "Gearbox-Config.vmdk",
			Format: "VMDK",
			Size: "1024",
		})

		for index, disk := range disks {
			fileName := me.baseDir.String() + "/" + me.EntityName.String() + "/" + disk.Name
			order := strconv.Itoa(index)

			_, err = me.Run("createmedium", "disk", "--filename", fileName, "--size", disk.Size, "--format", disk.Format, "--variant", "Stream")
			if err != nil {
				break
			}

			_, err = me.Run("storageattach", me.EntityName.String(),
				"--storagectl", "SATA", "--port", order, "--device", "0", "--type", "hdd", "--medium", fileName, "--hotpluggable", "off")
			if err != nil {
				break
			}
		}

		eblog.Debug(me.EntityId, "modified VM storage OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) cmdModifyVmIso() error {

	var err error

	// fmt.Printf("cmdModifyVmIso() ENTRY.\n")
	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		_, err = me.Run("storagectl", me.EntityName.String(),
			"--name", "IDE", "--add", "ide", "--hostiocache", "on", "--bootable", "on")
		if err != nil {
			break
		}

		_, err = me.Run("storageattach", me.EntityName.String(),
			"--storagectl", "IDE", "--port", "0", "--device", "0", "--type", "dvddrive", "--tempeject", "on", "--medium", me.osRelease.File.String())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "modified VM ISO OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) findFirstNic() error {

	var err error

	// fmt.Printf("findFirstNic() ENTRY.\n")
	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		_, err = me.Run("list", "bridgedifs", "-s")
		if err != nil {
			break
		}

		var nic KeyValueMap
		dr, ok := decodeResponse(me.Entry.cmdStdout, ':')
		if ok == true {
			me.Entry.vmNics, ok = dr.decodeBridgeIfs()
			if ok == false {
				err = me.EntityName.ProduceError("no NICs found")
				break
			}

			for _, nic = range me.Entry.vmNics {
				if nic["FirstNic"] == "Yes" {
					break
				}
			}
		}

		if nic == nil {
			err = me.EntityName.ProduceError("no NICs found")
			break
		}

		eblog.Debug(me.EntityId, "using NIC '%s' for VM", nic)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Run runs a VBoxManage command.
func (me *Vm) Run(args ...string) (exitCode string, err error) {

	var vboxManagePath string

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// If vBoxManage is not found in the system path, fall back to the
		// hard coded path.
		if path, err := exec.LookPath("VBoxManage"); err == nil {
			vboxManagePath = path
		} else {
			vboxManagePath = VBOXMANAGE
		}

		eblog.Debug(me.EntityId, "EXEC:%v '%v'", vboxManagePath, strings.Join(args, `" "`))
		cmd := exec.Command(vboxManagePath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Stdout, cmd.Stderr = &me.Entry.cmdStdout, &me.Entry.cmdStderr
		err = cmd.Run()
		if err != nil {
			exitCode = strings.TrimPrefix(err.Error(), "exit status ")

			switch exitCode {
				case exitCodeMissingVm:
					err = me.EntityName.ProduceError("VM unregistered")
					//fmt.Printf("stdout:%v\n", stdout.String())
					//fmt.Printf("stderr:%v\n", stderr.String())
					//fmt.Printf("returnCode:'%v'\n", returnCode)

				default:
					err = me.EntityName.ProduceError("failed to run command '%v'", err.Error())
					fmt.Printf("stdout:%v\n", me.Entry.cmdStdout.String())
					fmt.Printf("stderr:%v\n", me.Entry.cmdStderr.String())
					fmt.Printf("returnCode:'%v'\n", exitCode)
			}
		}

		//eblog.Debug(me.EntityId, "started task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return
}

const (
	exitCodeOK = "0"
	exitCodeMissingVm = "1"
	exitCodeCmdError = "2"
)



type KeyValue struct {
	Key	string
	Value string
}
type KeyValues []KeyValue

type KeyValueMap map[string]string
type KeyValuesMap map[string]KeyValueMap


func (kvs *KeyValues) decodeBridgeIfs() (KeyValuesMap, bool) {

	ok := false
	kvm := make(KeyValuesMap)
	firstNic := false

	currentKv := KeyValueMap{}
	for _, kv := range *kvs {
		// Output always ends with a 'VBoxNetworkName' key.
		if kv.Key == "VBoxNetworkName" && kv.Value != "" {
			ok = true
			if (firstNic == true) && (currentKv["Status"] == "Yes") {
				currentKv["FirstNic"] = "Yes"
				firstNic = true
			}

			kvm[kv.Value] = currentKv
			currentKv = KeyValueMap{}

		} else {
			currentKv[kv.Key] = kv.Value
		}
	}

	return kvm, ok
}


func (kvs *KeyValues) decodeShowVmInfo() (KeyValueMap, bool) {

	ok := false
	kvm := make(KeyValueMap)

	for _, kv := range *kvs {
		// Output always start with a 'name' key.
		if kv.Key == "name" && kv.Value != "" {
			ok = true
		}
		if kv.Key != "" && kv.Value != "" {
			kvm[kv.Key] = kv.Value
		}
	}

	return kvm, ok
}


func lineToKey(s string, splitOn rune) (ld KeyValue, ok bool) {

	ok = false
	foundKey := false
	stripSpace := false
	f := func(c rune) bool {
		switch {
		case c == splitOn:
			if foundKey {
				return false
			} else {
				foundKey = true
				stripSpace = true
				return true
			}

		case unicode.IsSpace(c) && stripSpace:
			return true

		default:
			stripSpace = false
			return false
			//return unicode.IsSpace(c)
		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(s, f)
	if len(items) == 1 {
		ld.Key = items[0]
		ld.Value = ""
		ok = true
	} else if len(items) == 2 {
		ld.Key = items[0]
		ld.Value = items[1]
		ok = true
	}
	ld.Key = strings.TrimSuffix(strings.TrimPrefix(ld.Key, `"`), `"`)
	ld.Key = strings.TrimSuffix(strings.TrimPrefix(ld.Key, `{`), `}`)
	ld.Value = strings.TrimSuffix(strings.TrimPrefix(ld.Value, `"`), `"`)
	ld.Value = strings.TrimSuffix(strings.TrimPrefix(ld.Value, `{`), `}`)

	return
}


func decodeResponse(s bytes.Buffer, splitOn rune) (dr KeyValues, ok bool) {

	ok = false

	lines := strings.Split(s.String(), "\n")
	for _, l := range lines {
		kv, lineOk := lineToKey(l, splitOn)
		if lineOk == false {
			continue

		} else {
			// Return true if we have at least one key/value pair found.
			ok = true
			dr = append(dr, kv)
		}
		// fmt.Printf("items[%d]: '%s' = '%s'\n", i, foo.Key, foo.Value)
	}

	return
}




//// RunCombinedError runs a VBoxManage command.  The output is stdout and the the
//// combined err/stderr from the command.
//func (me *Vm) RunCombinedError(args ...string) (string, error) {
//
//	wout, werr, err := me.Run(args...)
//	if err != nil {
//		if werr != "" {
//			return wout, fmt.Errorf("%s: %s", err, werr)
//		}
//		return wout, err
//	}
//
//	return wout, nil
//}


//#!/bin/bash
//
//VM_NAME="GearboxNew"
//
//SATA_DISK0_NAME="Gearbox-Opt.vmdk"
//SATA_DISK0_FORMAT=VMDK
//SATA_DISK0_SIZE=1024
//SATA_DISK0_ORDER=0
//
//SATA_DISK1_NAME="Gearbox-Docker.vmdk"
//SATA_DISK1_FORMAT=VMDK
//SATA_DISK1_SIZE=16384
//SATA_DISK1_ORDER=1
//
//SATA_DISK2_NAME="Gearbox-Projects.vmdk"
//SATA_DISK2_FORMAT=VMDK
//SATA_DISK2_SIZE=16384
//SATA_DISK2_ORDER=2
//
//SATA_DISK3_NAME="Gearbox-Config.vmdk"
//SATA_DISK3_FORMAT=VMDK
//SATA_DISK3_SIZE=1024
//SATA_DISK3_ORDER=3
//
//# This can potentially create several interface names.
//# On my Mac it returns:
//# "en0: Ethernet"
//BRIDGED_HOST="$(VBoxManage list bridgedifs -s | awk '/^Name:/{gsub(/^Name: +/, ""); print}')"
//
//# /dev/sda	/opt/gearbox		ext4	noauto,defaults	0 0
//# /dev/sdb	/var/lib/docker		ext4	noauto,defaults	0 0
//# /dev/sdc	/var/lib/gearbox	ext4	noauto,defaults	0 0
//# /dev/sdd	/etc/gearbox		ext4	noauto,defaults	0 0
//
//
//# Create the base VM and register.
//VBoxManage createvm --name ${VM_NAME} --ostype Linux26_64 --register --basefolder "${VM_DIR}"
//
//VM_CFG_FILE="$(VBoxManage showvminfo ${VM_NAME} --machinereadable | awk -F= '/^CfgFile/{gsub(/"/, ""); print$2}')"
//VM_BASE_DIR="$(dirname "${VM_CFG_FILE}")"
//cd "${VM_BASE_DIR}"
//
//# Misc options.
//VBoxManage modifyvm ${VM_NAME} --description "Gearbox OS VM" # --iconfile "${ICON_FILE_NAME}"
//VBoxManage modifyvm ${VM_NAME} --ioapic on --acpi on --biosbootmenu disabled --biosapic apic
//VBoxManage modifyvm ${VM_NAME} --boot1 dvd --boot2 none --boot3 none --boot4 none
//VBoxManage modifyvm ${VM_NAME} --vrde off --autostart-enabled off
//VBoxManage modifyvm ${VM_NAME} --cpuhotplug on --cpus 4 --pae off --longmode on --largepages on --paravirtprovider default
//VBoxManage modifyvm ${VM_NAME} --accelerate3d off --accelerate2dvideo off --mouse usbtablet
//VBoxManage modifyvm ${VM_NAME} --defaultfrontend headless --snapshotfolder default
//VBoxManage modifyvm ${VM_NAME} --memory 2048 --vram 128
//VBoxManage modifyvm ${VM_NAME} --audio none
//VBoxManage modifyvm ${VM_NAME} --nic1 nat --nictype1 82540EM --cableconnected1 on --macaddress1 auto
//VBoxManage modifyvm ${VM_NAME} --natnet1 default --natpf1 API,tcp,,9970,,9970
//VBoxManage modifyvm ${VM_NAME} --natnet1 default --natpf1 SSH,tcp,,2222,,22
//VBoxManage modifyvm ${VM_NAME} --natnet1 default --natpf1 VMcontrol,tcp,,9971,,9971
//VBoxManage modifyvm ${VM_NAME} --nic2 bridged --bridgeadapter2 "${BRIDGED_HOST}" --nictype2 82540EM --cableconnected2 on --macaddress2 auto --nicpromisc2 deny
//
//# Setup console UART
//VBoxManage modifyvm ${VM_NAME} --uart1 0x3f8 4 --uartmode1 tcpserver 2023
//
//# Create a SATA controller instance.
//VBoxManage storagectl ${VM_NAME} --name "SATA" --add sata --controller IntelAHCI --portcount 4 --hostiocache off --bootable on
//
//# Create virtual disks.
//VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK0_NAME}" --size ${SATA_DISK0_SIZE} --format ${SATA_DISK0_FORMAT} --variant Stream
//VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK1_NAME}" --size ${SATA_DISK1_SIZE} --format ${SATA_DISK1_FORMAT} --variant Stream
//VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK2_NAME}" --size ${SATA_DISK2_SIZE} --format ${SATA_DISK2_FORMAT} --variant Stream
//VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK3_NAME}" --size ${SATA_DISK3_SIZE} --format ${SATA_DISK3_FORMAT} --variant Stream
//
//# Attach virtual disks to SATA controller.
//VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK0_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK0_NAME}" --hotpluggable off
//VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK1_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK1_NAME}" --hotpluggable off
//VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK2_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK2_NAME}" --hotpluggable off
//VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK3_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK3_NAME}" --hotpluggable off
//
//# Create IDE bus for ISO.
//VBoxManage storagectl ${VM_NAME} --name "IDE" --add ide --hostiocache on --bootable on
//VBoxManage storageattach ${VM_NAME} --storagectl "IDE" --port 0 --device 0 --type dvddrive --tempeject on  --medium $HOME/.gearbox/box/iso/gearbox-0.5.0.iso

