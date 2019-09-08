package vbmanage

type Strings = []string

type ModifyVmArgs struct {
	Name string
}

type CreateVmArgs struct {
	Name string
}

type ControlVmArgs struct {
	Name string
	PowerOff bool
	AcpiPowerButton bool
}

type UnregisterVmArgs struct {
	Name string
	Delete bool
}

type DhcpServerArgs struct {
	Name string
	Remove bool
	Netname string
}

type StorageCtlArgs struct {
	Name string
	Add string
	Controller string
	PortCount int
	HostIoCache bool  // true=="on"
	Bootable bool  // true=="on"
}

type StorageAttachArgs struct {
	Name string
}

type CreateMediumArgs struct {
	Name string
	Disk bool
}

type HostOnlyIfArgs struct {
	Name string
	Disk bool
	Remove bool
}

type ShowVmInfoArgs struct {
	Name string
	MachineReadable bool
}

