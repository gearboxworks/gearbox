package routes

import (
	"fmt"
	"gearbox"
	"gearbox/apibuilder"
	"gearbox/gears"
	"gearbox/status"
	"gearbox/types"
)

const StackTypeName = "stack"

var StackInstance = (*NamedStack)(nil)
var _ ab.Item = StackInstance

func (me *NamedStack) GetType() ab.ItemType {
	return StackTypeName
}

func (me *NamedStack) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStack) GetId() ab.ItemId {
	return ab.ItemId(fmt.Sprintf("%s/%s", me.Authority, me.StackName))
}

func (me *NamedStack) GetItem() (ab.Item, status.Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStack) (gbns *gears.NamedStack, sts status.Status) {
	gbns = gears.NewNamedStack(gb.GetGears(), types.StackId(ns.GetId()))
	sts = gbns.Refresh()
	return gbns, sts
}
