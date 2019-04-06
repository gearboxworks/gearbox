package routes

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/status"
	"gearbox/types"
)

const StackTypeName = "stack"

var StackInstance = (*NamedStack)(nil)
var _ apimodeler.Itemer = StackInstance

func (me *NamedStack) GetType() apimodeler.ItemType {
	return StackTypeName
}

func (me *NamedStack) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStack) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(fmt.Sprintf("%s/%s", me.Authority, me.StackName))
}

func (me *NamedStack) GetItem() (apimodeler.Itemer, status.Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStack) (gbns *gears.NamedStack, sts status.Status) {
	gbns = gears.NewNamedStack(gb.GetGears(), types.StackId(ns.GetId()))
	sts = gbns.Refresh()
	return gbns, sts
}
