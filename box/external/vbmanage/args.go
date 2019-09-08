package vbmanage

import "strings"

type Args struct {
	Name string
	Type string
	Create bool
	Remove bool
	Disk bool
	Long bool
	Sorted bool
	MachineReadable bool
	HostOnlyIfs bool
}

func (me *Args) String() string {
	return strings.Join(me.ToParams()," ")
}

func (me *Args) ToParams() (params []string) {
	params = make([]string, 0)
	if me.Name != "" {
		params = append(params,me.Name)
	}
	if me.Type != "" {
		params = append(params,"--type")
		params = append(params,me.Type)
	}
	if me.PowerOff {
		params = append(params,"poweroff")
	}
	if me.AcpiPowerButton {
		params = append(params,"acpipowerbutton")
	}
	if me.Delete {
		params = append(params,"--delete")
	}
	if me.MachineReadable {
		params = append(params,"--machinereadable")
	}
	if me.Create {
		params = append(params,"create")
	}
	if me.Remove {
		params = append(params,"remove")
	}
	if me.Disk {
		params = append(params,"disk")
	}
	if me.Long {
		params = append(params,"long")
	}
	if me.Sorted {
		params = append(params,"sorted")
	}

	return []string{}
}
