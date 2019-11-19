package box

import (
	"bytes"
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Disk struct {
	Name   string
	Format string
	Size   string
}
type Disks []Disk

func (me *Box) CreateBox() (sts status.Status) {

	// fmt.Printf("CreateBox() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}
		steps := 0

		/*
			// 0. Destroy existing VM.
			fmt.Printf("#### %d. Destroy VM - ", steps); steps++
			destroyVm, sts := me.cmdDestroyVm()
			if !is.Success(sts) {
				fmt.Printf("VM still exists:%v\n", destroyVm["name"])
				// Convert to warning message.
				sts = sts.SetWarning(true)
				fmt.Printf("sts:%v\n", sts)
				break
			} else {
				fmt.Printf("\n")
			}
		*/

		// 1. Check if a VM already exists by that name.
		fmt.Printf("#### %d. Check for VM - ", steps)
		steps++
		listVm, sts := me.cmdListVm()
		if is.Success(sts) {
			fmt.Printf("VM exists:%v\n", listVm["name"])
			// Convert to warning message.
			sts = sts.SetWarn(true)
			fmt.Printf("sts:%v\n", sts)
			break
		} else {
			fmt.Printf("\n")
		}

		// 2. Create VM OVA file.
		fmt.Printf("#### %d. Create VM - ", steps)
		steps++
		createVm, sts := me.cmdCreateVm()
		if !is.Success(sts) {
			// Convert to warning message.
			sts = sts.SetWarn(true)
			fmt.Printf("sts:%v\n", sts)
			break
		} else {
			fmt.Printf("Created:%v\n", createVm["name"])
		}

		// 3. Modify VM OVA file.
		fmt.Printf("#### %d. Modify VM.\n", steps)
		steps++
		modifyVm, sts := me.cmdModifyVm()
		if !is.Success(sts) {
			// Convert to warning message.
			sts = sts.SetWarn(true)
			fmt.Printf("sts:%v\n", sts)
			break
		} else {
			fmt.Printf("Modified:%v\n", modifyVm["name"])
		}

		listVm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}
	}

	return sts
}

func (me *Box) StartBox() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		me.State.VM.WantState = VmStateRunning
		me.State.API.WantState = VmStateRunning

		// stdout, stderr, err := me.Run("showvminfo", vm, "--machinereadable")
		_, _, sts = me.Run("startvm", me.Boxname, "--type", "headless")
		if is.Error(sts) {
			break
		}

		_, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}
	}

	return kvm, sts
}

func (me *Box) StopBox() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		me.State.VM.WantState = VmStatePowerOff
		me.State.API.WantState = VmStatePowerOff

		// stdout, stderr, err := me.Run("showvminfo", vm, "--machinereadable")
		_, _, sts = me.Run("controlvm", me.Boxname, "acpipowerbutton")
		if is.Error(sts) {
			break
		}

		_, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}
	}

	return kvm, sts
}

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

func (me *Box) cmdListVm() (KeyValueMap, status.Status) {

	var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		// stdout, stderr, err := me.Run("showvminfo", vm, "--machinereadable")
		stdout, _, sts = me.Run("showvminfo", me.Boxname, "--machinereadable")
		if is.Error(sts) {
			break
		}

		kvs, ok := decodeResponse(stdout, '=')
		if ok == false {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM - Can't decode showvminfo output for '%s'.\n", global.Brandname, me.Boxname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    kvs,
			})
			break
		}

		kvm, ok = kvs.decodeShowVmInfo()
		if ok == false {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM - Can't decode showvminfo output for '%s'.\n", global.Brandname, me.Boxname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    kvm,
			})
			break
		}

		me.State.VM.CurrentState = kvm["VMState"]

		sts = status.Success("%s VM - Box exists and is in state %s.\n", global.Brandname, kvm["VMState"])
	}

	return kvm, sts
}

func (me *Box) cmdCreateVm() (KeyValueMap, status.Status) {

	// var stdout string
	var stderr string
	var sts status.Status
	var kvm KeyValueMap

	// fmt.Printf("cmdCreateVm() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		// stdout, stderr, sts = me.Run("createvm", "--name", me.Boxname, "--ostype", "Linux26_64", "--register", "--basefolder", me.VmBaseDir)
		_, stderr, sts = me.Run("createvm", "--name", me.Boxname, "--ostype", "Linux26_64", "--register", "--basefolder", me.VmBaseDir)
		if is.Error(sts) {
			if sts.Data() == "1" {
				// Logic inverted, because me.Run will normally return '1' if a Box does NOT exist.
				// In the 'createvm' case, it'll return a '1' if it does exist.
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("%s VM - Box already exists '%s'.\n", global.Brandname, me.Boxname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    stderr,
				})
			}
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s VM - Box created %s.\n", global.Brandname, me.Boxname)
	}

	return kvm, sts
}

func (me *Box) cmdDestroyVm() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		// stdout, stderr, sts = me.Run("unregistervm", me.Boxname, "--delete")
		_, _, sts = me.Run("unregistervm", me.Boxname, "--delete")
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if !is.Success(sts) {
			// Basically, no VM means we deleted it.
			sts = status.Success("%s VM - Box destroyed %s.\n", global.Brandname, me.Boxname)
		}
	}

	return kvm, sts
}

func (me *Box) cmdModifyVmBasic() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	// fmt.Printf("cmdModifyVmBasic() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		// stdout, stderr, sts = me.Run("modifyvm", me.Boxname, "--description", me.Boxname + " OS VM", "--iconfile", string(me.OsBridge.GetAdminRootDir()) + "/" + IconLogo)
		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--description", me.Boxname+" OS VM", "--iconfile", string(me.OsBridge.GetAdminRootDir())+"/"+IconLogoPng)
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--ioapic", "on", "--acpi", "on", "--biosbootmenu", "disabled", "--biosapic", "apic")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--boot1", "dvd", "--boot2", "none", "--boot3", "none", "--boot4", "none")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--vrde", "off", "--autostart-enabled", "off")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--cpuhotplug", "on", "--cpus", "4", "--pae", "off", "--longmode", "on", "--largepages", "on", "--paravirtprovider", "default")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--accelerate3d", "off", "--accelerate2dvideo", "off", "--mouse", "usbtablet")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--defaultfrontend", "headless", "--snapshotfolder", "default")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--memory", "2048", "--vram", "128")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--audio", "none")
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s - VM modified %s.\n", global.Brandname, me.Boxname)
	}

	return kvm, sts
}

func (me *Box) cmdModifyVmNetwork() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	// fmt.Printf("cmdModifyVmNetwork() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		/*
			nic, sts := me.findFirstNic()
			if is.Error(sts) {
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("%s VM - No NIC found '%s'.\n", global.Brandname, me.Boxname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    "",
				})
				break
			}
		*/

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--nic1", "nat", "--nictype1", "82540EM", "--cableconnected1", "on", "--macaddress1", "auto")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--natnet1", "default", "--natpf1", "API,tcp,,9970,,9970")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--natnet1", "default", "--natpf1", "SSH,tcp,,2222,,22")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--natnet1", "default", "--natpf1", "VMcontrol,tcp,,9971,,9971")
		if is.Error(sts) {
			break
		}

		// _, _, sts = me.Run("modifyvm", me.Boxname, "--nic2", "bridged", "--bridgeadapter2", nic["Name"], "--nictype2", "82540EM", "--cableconnected2", "on", "--macaddress2", "auto", "--nicpromisc2", "deny")
		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--nic2", "hostonly", "--hostonlyadapter2", "vboxnet0", "--nictype2", "82540EM", "--cableconnected2", "on", "--macaddress2", "auto", "--nicpromisc2", "deny")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("modifyvm", me.Boxname,
			"--uart1", "0x3f8", "4", "--uartmode1", "tcpserver", me.ConsolePort)
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s - VM modified %s.\n", global.Brandname, me.Boxname)
	}

	return kvm, sts
}

func (me *Box) cmdModifyVmStorage() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	// fmt.Printf("cmdModifyVmStorage() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("storagectl", me.Boxname,
			"--name", "SATA", "--add", "sata", "--controller", "IntelAHCI", "--portcount", "4", "--hostiocache", "off", "--bootable", "on")
		if is.Error(sts) {
			break
		}

		// SIGH - Needs to be not hard-coded.
		disks := Disks{}
		disks = append(disks, Disk{
			Name:   "Gearbox-Opt.vmdk",
			Format: "VMDK",
			Size:   "1024",
		})
		disks = append(disks, Disk{
			Name:   "Gearbox-Docker.vmdk",
			Format: "VMDK",
			Size:   "16384",
		})
		disks = append(disks, Disk{
			Name:   "Gearbox-Projects.vmdk",
			Format: "VMDK",
			Size:   "16384",
		})
		disks = append(disks, Disk{
			Name:   "Gearbox-Config.vmdk",
			Format: "VMDK",
			Size:   "1024",
		})

		for index, disk := range disks {
			fileName := me.VmBaseDir + "/" + me.Boxname + "/" + disk.Name
			order := strconv.Itoa(index)
			_, _, sts = me.Run("createmedium", "disk", "--filename", fileName, "--size", disk.Size, "--format", disk.Format, "--variant", "Stream")
			if is.Error(sts) {
				break
			}
			_, _, sts = me.Run("storageattach", me.Boxname,
				"--storagectl", "SATA", "--port", order, "--device", "0", "--type", "hdd", "--medium", fileName, "--hotpluggable", "off")
			if is.Error(sts) {
				break
			}
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s - VM modified %s.\n", global.Brandname, me.Boxname)
	}

	return kvm, sts
}

func (me *Box) cmdModifyVmIso() (KeyValueMap, status.Status) {

	//var stdout string
	//var stderr string
	var sts status.Status
	var kvm KeyValueMap

	// fmt.Printf("cmdModifyVmIso() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("storagectl", me.Boxname,
			"--name", "IDE", "--add", "ide", "--hostiocache", "on", "--bootable", "on")
		if is.Error(sts) {
			break
		}

		_, _, sts = me.Run("storageattach", me.Boxname,
			"--storagectl", "IDE", "--port", "0", "--device", "0", "--type", "dvddrive", "--tempeject", "on", "--medium", me.VmIsoFile)
		if is.Error(sts) {
			fmt.Printf("Error: %v\n", sts.Data())
			break
		}

		kvm, sts = me.cmdListVm()
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s - VM modified %s.\n", global.Brandname, me.Boxname)
	}

	return kvm, sts
}

func (me *Box) cmdModifyVm() (KeyValueMap, status.Status) {

	var sts status.Status
	var kvm KeyValueMap

	// fmt.Printf("cmdModifyVm() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdModifyVmBasic()
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdModifyVmNetwork()
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdModifyVmStorage()
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdModifyVmIso()
		if is.Error(sts) {
			break
		}

		kvm, sts = me.cmdListVm()
		sts = status.Success("%s - VM modifed %s.\n", global.Brandname, me.Boxname)
	}

	return kvm, sts
}

func (me *Box) findFirstNic() (KeyValueMap, status.Status) {

	var stdout string
	//var stderr string
	var sts status.Status
	var nic KeyValueMap

	// fmt.Printf("findFirstNic() ENTRY.\n")
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		stdout, _, sts = me.Run("list", "bridgedifs", "-s")
		if is.Error(sts) {
			break
		}

		dr, ok := decodeResponse(stdout, ':')
		if ok == true {
			nics, ok := dr.decodeBridgeIfs()
			if ok == false {
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("%s VM - No NIC found.\n", global.Brandname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    "",
				})
				break
			}

			for _, nic = range nics {
				if nic["FirstNic"] == "Yes" {
					break
				}
			}
		}

		if nic == nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM - No NIC found.\n", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    "",
			})
			break
		}
	}

	return nic, sts
}

// Run runs a VBoxManage command.
func (me *Box) Run(args ...string) (string, string, status.Status) {

	var vboxManagePath string
	var sts status.Status

	// If vBoxManage is not found in the system path, fall back to the
	// hard coded path.
	if path, err := exec.LookPath("VBoxManage"); err == nil {
		vboxManagePath = path
	} else {
		vboxManagePath = VBOXMANAGE
	}

	fmt.Printf("EXEC:%v \"%v\"\n", vboxManagePath, strings.Join(args, `" "`))
	cmd := exec.Command(vboxManagePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if err != nil {
		returnCode := strings.TrimPrefix(err.Error(), "exit status ")
		switch returnCode {
		case "1":
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM - No such Box called '%s'.\n", global.Brandname, me.Boxname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    stderr.String(),
			})
			//fmt.Printf("stdout:%v\n", stdout.String())
			//fmt.Printf("stderr:%v\n", stderr.String())
			//fmt.Printf("returnCode:'%v'\n", returnCode)

		default:
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM - Failed to run command '%v'.\n", global.Brandname, err.Error()),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    returnCode,
			})
			fmt.Printf("stdout:%v\n", stdout.String())
			fmt.Printf("stderr:%v\n", stderr.String())
			fmt.Printf("returnCode:'%v'\n", returnCode)
		}
	}

	return stdout.String(), stderr.String(), sts
}

// RunCombinedError runs a VBoxManage command.  The output is stdout and the the
// combined err/stderr from the command.
func (me *Box) RunCombinedError(args ...string) (string, error) {

	wout, werr, err := me.Run(args...)
	if err != nil {
		if werr != "" {
			return wout, fmt.Errorf("%s: %s", err, werr)
		}
		return wout, err
	}

	return wout, nil
}

/*
#!/bin/bash

VM_NAME="GearboxNew"

SATA_DISK0_NAME="Gearbox-Opt.vmdk"
SATA_DISK0_FORMAT=VMDK
SATA_DISK0_SIZE=1024
SATA_DISK0_ORDER=0

SATA_DISK1_NAME="Gearbox-Docker.vmdk"
SATA_DISK1_FORMAT=VMDK
SATA_DISK1_SIZE=16384
SATA_DISK1_ORDER=1

SATA_DISK2_NAME="Gearbox-Projects.vmdk"
SATA_DISK2_FORMAT=VMDK
SATA_DISK2_SIZE=16384
SATA_DISK2_ORDER=2

SATA_DISK3_NAME="Gearbox-Config.vmdk"
SATA_DISK3_FORMAT=VMDK
SATA_DISK3_SIZE=1024
SATA_DISK3_ORDER=3

# This can potentially create several interface names.
# On my Mac it returns:
# "en0: Ethernet"
BRIDGED_HOST="$(VBoxManage list bridgedifs -s | awk '/^Name:/{gsub(/^Name: +/, ""); print}')"

# /dev/sda	/opt/gearbox		ext4	noauto,defaults	0 0
# /dev/sdb	/var/lib/docker		ext4	noauto,defaults	0 0
# /dev/sdc	/var/lib/gearbox	ext4	noauto,defaults	0 0
# /dev/sdd	/etc/gearbox		ext4	noauto,defaults	0 0


# Create the base VM and register.
VBoxManage createvm --name ${VM_NAME} --ostype Linux26_64 --register --basefolder "${VM_DIR}"

VM_CFG_FILE="$(VBoxManage showvminfo ${VM_NAME} --machinereadable | awk -F= '/^CfgFile/{gsub(/"/, ""); print$2}')"
VM_BASE_DIR="$(dirname "${VM_CFG_FILE}")"
cd "${VM_BASE_DIR}"

# Misc options.
VBoxManage modifyvm ${VM_NAME} --description "Gearbox OS VM" # --iconfile "${ICON_FILE_NAME}"
VBoxManage modifyvm ${VM_NAME} --ioapic on --acpi on --biosbootmenu disabled --biosapic apic
VBoxManage modifyvm ${VM_NAME} --boot1 dvd --boot2 none --boot3 none --boot4 none
VBoxManage modifyvm ${VM_NAME} --vrde off --autostart-enabled off
VBoxManage modifyvm ${VM_NAME} --cpuhotplug on --cpus 4 --pae off --longmode on --largepages on --paravirtprovider default
VBoxManage modifyvm ${VM_NAME} --accelerate3d off --accelerate2dvideo off --mouse usbtablet
VBoxManage modifyvm ${VM_NAME} --defaultfrontend headless --snapshotfolder default
VBoxManage modifyvm ${VM_NAME} --memory 2048 --vram 128
VBoxManage modifyvm ${VM_NAME} --audio none
VBoxManage modifyvm ${VM_NAME} --nic1 nat --nictype1 82540EM --cableconnected1 on --macaddress1 auto
VBoxManage modifyvm ${VM_NAME} --natnet1 default --natpf1 API,tcp,,9970,,9970
VBoxManage modifyvm ${VM_NAME} --natnet1 default --natpf1 SSH,tcp,,2222,,22
VBoxManage modifyvm ${VM_NAME} --natnet1 default --natpf1 VMcontrol,tcp,,9971,,9971
VBoxManage modifyvm ${VM_NAME} --nic2 bridged --bridgeadapter2 "${BRIDGED_HOST}" --nictype2 82540EM --cableconnected2 on --macaddress2 auto --nicpromisc2 deny

# Setup console UART
VBoxManage modifyvm ${VM_NAME} --uart1 0x3f8 4 --uartmode1 tcpserver 2023

# Create a SATA controller instance.
VBoxManage storagectl ${VM_NAME} --name "SATA" --add sata --controller IntelAHCI --portcount 4 --hostiocache off --bootable on

# Create virtual disks.
VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK0_NAME}" --size ${SATA_DISK0_SIZE} --format ${SATA_DISK0_FORMAT} --variant Stream
VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK1_NAME}" --size ${SATA_DISK1_SIZE} --format ${SATA_DISK1_FORMAT} --variant Stream
VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK2_NAME}" --size ${SATA_DISK2_SIZE} --format ${SATA_DISK2_FORMAT} --variant Stream
VBoxManage createmedium disk --filename "${VM_BASE_DIR}/${SATA_DISK3_NAME}" --size ${SATA_DISK3_SIZE} --format ${SATA_DISK3_FORMAT} --variant Stream

# Attach virtual disks to SATA controller.
VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK0_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK0_NAME}" --hotpluggable off
VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK1_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK1_NAME}" --hotpluggable off
VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK2_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK2_NAME}" --hotpluggable off
VBoxManage storageattach ${VM_NAME} --storagectl "SATA" --port ${SATA_DISK3_ORDER} --device 0 --type hdd --medium "${VM_BASE_DIR}/${SATA_DISK3_NAME}" --hotpluggable off

# Create IDE bus for ISO.
VBoxManage storagectl ${VM_NAME} --name "IDE" --add ide --hostiocache on --bootable on
VBoxManage storageattach ${VM_NAME} --storagectl "IDE" --port 0 --device 0 --type dvddrive --tempeject on  --medium $HOME/.gearbox/box/iso/gearbox-0.5.0.iso
*/
