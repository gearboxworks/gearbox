package vmbox

import (
	"gearbox/box/external/hypervisor"
)

var _ hypervisor.Disker = (*Disk)(nil)

type Disks []Disk
type Disk struct {
	Name   string
	Format string
	Size   string
}

func (me *Disk) GetName() string {
	return me.Name
}

func (me *Disk) GetFormat() string {
	return me.Format
}

func (me *Disk) GetSize() string {
	return me.Size
}
