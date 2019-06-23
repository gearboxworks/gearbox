package gears

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"strings"
)

type NamedStackMap map[types.StackId]*NamedStack

type NamedStacks []*NamedStack

type NamedStack struct {
	Authority    types.AuthorityDomain `json:"authority"`
	Stackname    types.Stackname       `json:"stackname"`
	Gearspecs    Gearspecs             `json:"gearspecs,omitempty"`
	GearOptions  GearOptions           `json:"gears,omitempty"`
	GearRegistry *GearRegistry         `json:"-"`
	refreshed    bool
}

//func NewNamedStack(gears *GearRegistry, stackid types.StackId) *NamedStack {
func NewNamedStack(stackid types.StackId) *NamedStack {
	stack := NamedStack{
		Gearspecs:   make(Gearspecs, 0),
		GearOptions: make(GearOptions, 0),
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
		var gss gearspec.Gearspecs
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
func (me *NamedStack) GetGearspecs() (gss gearspec.Gearspecs, sts status.Status) {
	for range only.Once {
		gss = make(gearspec.Gearspecs, 0)
		nsrs, sts := me.GearRegistry.GetNamedRegistries(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		for _, r := range nsrs {
			gs := gearspec.NewGearspec()
			sts = gs.Parse(r.GetGearspecId())
			if is.Error(sts) {
				break
			}
			gss = append(gss, gs)
		}
	}
	return gss, sts
}

//
// Get the available stack roles for a given named stack
//
func (me *NamedStack) GetRegistries() (gss Gearspecs, sts status.Status) {
	for range only.Once {
		gss = make(Gearspecs, 0)
		for gs, rso := range me.Gearspecs {
			if !strings.HasPrefix(string(gs), string(me.GetIdentifier())) {
				continue
			}
			gss = append(gss, rso)
		}
	}
	return gss, sts
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

func (me *NamedStack) Refresh(gears *GearRegistry) (sts status.Status) {
	for range only.Once {
		if !me.NeedsRefresh() {
			break
		}

		var nsrs Gearspecs
		nsrs, sts = gears.GetNamedRegistries(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.Gearspecs = nsrs

		var sos GearOptions
		sos, sts = gears.GetNamedGearOptions(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.GearOptions = sos

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
