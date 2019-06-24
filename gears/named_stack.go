package gears

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

type NamedStackMap map[types.StackId]*NamedStack

type NamedStacks []*NamedStack

func (me NamedStacks) GetMap() (nsm NamedStackMap) {
	nsm = make(NamedStackMap, len(me))
	for _, ns := range me {
		nsm[ns.GetIdentifier()] = ns
	}
	return nsm
}

type NamedStack struct {
	Authority types.AuthorityDomain `json:"authority"`
	Stackname types.Stackname       `json:"stackname"`
	Gearspecs Gearspecs             `json:"gearspecs,omitempty"`
	refreshed bool
}

//func NewNamedStack(gears *GearRegistry, stackid types.StackId) *NamedStack {
func NewNamedStack(stackid types.StackId) *NamedStack {
	stack := NamedStack{
		Gearspecs: make(Gearspecs, 0),
	}
	// This will split authority out, or do nothing
	_ = stack.SetIdentifier(stackid)
	return &stack
}

//
// Get the list of gearspec identifiers
//
func (me *NamedStack) GetGearspecIds() (gsids gearspec.Identifiers, sts status.Status) {
	for range only.Once {
		var gss Gearspecs
		gss, sts = me.GetGearspecs()
		if is.Error(sts) {
			break
		}
		for _, r := range gss {
			gsids = append(gsids, r.GetIdentifier())
		}
	}
	return gsids, sts
}

//
// Get the list of gearspecs
//
func (me *NamedStack) GetGearspecs() (gss Gearspecs, sts status.Status) {
	return me.Gearspecs, sts
}

func (me *NamedStack) String() string {
	return string(me.GetIdentifier())
}

func (me *NamedStack) LightweightClone() *NamedStack {
	stack := NamedStack{}
	stack = *me
	return &stack
}

//func (me *NamedStack) GetDefaultGears() (sm DefaultGearMap, sts status.Status) {
//	for range only.Once {
//		sm = make(DefaultGearMap, 0)
//		sts = me.Refresh()
//		if status.IsError(sts) {
//			break
//		}
//		for gs, s := range me.GearOptions {
//			if s.DefaultGear == nil {
//				continue
//			}
//			sm[gs] = s.DefaultGear.Clone()
//		}
//	}
//	return sm, sts
//}

func (me *NamedStack) AddGearspec(gs *Gearspec) (sts status.Status) {
	me.Gearspecs = append(me.Gearspecs, gs)
	return sts
}

func (me *NamedStack) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *NamedStack) Refresh(grs *GearRegistry) (sts status.Status) {
	for range only.Once {
		if !me.NeedsRefresh() {
			break
		}

		var nsgss Gearspecs
		nsgss, sts = grs.FilterByNamedStack(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.Gearspecs = nsgss

		me.refreshed = true
	}
	if is.Success(sts) {
		sts = status.Success("named stack '%s' refreshed", me.GetIdentifier())
	}
	return sts
}

func (me *NamedStack) SetIdentifier(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		gsi := gearspec.NewGearspec()
		sts := gsi.SetStackId(stackid)
		if status.IsError(sts) {
			break
		}
		me.Authority = gsi.AuthorityDomain
		me.Stackname = gsi.Stackname
	}
	return sts
}

func (me *NamedStack) GetIdentifier() types.StackId {
	return types.StackId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
}

func (me *NamedStack) Fixup(gr *GearRegistry) (sts status.Status) {
	stackid := me.GetIdentifier()
	me.Gearspecs = make(Gearspecs, 0)
	for _, gs := range gr.Gearspecs {
		if gs.GetStackId() != stackid {
			continue
		}
		me.Gearspecs = append(me.Gearspecs, gs)
	}
	return sts
}
