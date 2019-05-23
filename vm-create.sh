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

