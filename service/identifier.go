package service

import (
	"gearbox/gear"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"sort"
)

type Identifiers []Identifier

func (me Identifiers) Sort() {
	sns := make([]string, len(me))
	for i, sn := range me {
		sns[i] = string(sn)
	}
	sort.Strings(sns)
	for i, sn := range sns {
		me[i] = Identifier(sn)
	}
}

type Identifier string

func (me Identifier) GetPersistableServiceValue() (s Servicer, sts status.Status) {
	for range only.Once {
		gid := gear.NewGear()
		sts = gid.Parse(gear.Identifier(me))
		if is.Error(sts) {
			break
		}
		s = Identifier(gid.GetIdentifier())
	}
	return s, sts
}
func (me Identifier) GetServiceValue() (Servicer, status.Status) {
	return me, nil
}
func (me Identifier) GetServiceId() (Identifier, status.Status) {
	return me, nil
}
