package gears

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"strings"
)

type NamedStackMap map[types.StackId]*NamedStack

type NamedStacks []*NamedStack

type NamedStack struct {
	Authority      types.AuthorityDomain `json:"authority"`
	Stackname      types.Stackname       `json:"name"`
	StackRoles     StackRoles            `json:"roles,omitempty"`
	ServiceOptions ServiceOptions        `json:"services,omitempty"`
	Gears          *Gears                `json:"-"`
	refreshed      bool
}

//func NewNamedStack(gears *Gears, stackid types.StackId) *NamedStack {
func NewNamedStack(stackid types.StackId) *NamedStack {
	stack := NamedStack{
		StackRoles:     make(StackRoles, 0),
		ServiceOptions: make(ServiceOptions, 0),
		//		Gears:           gears,
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
		nsrs, sts := me.Gears.GetNamedStackRoles(me.GetIdentifier())
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
func (me *NamedStack) GetStackRoles() (srs StackRoles, sts status.Status) {
	for range only.Once {
		srs = make(StackRoles, 0)
		for gs, rso := range me.StackRoles {
			if !strings.HasPrefix(string(gs), string(me.GetIdentifier())) {
				continue
			}
			srs = append(srs, rso)
		}
	}
	return srs, sts
}

func (me *NamedStack) String() string {
	return string(me.GetIdentifier())
}

func (me *NamedStack) LightweightClone() *NamedStack {
	stack := NamedStack{}
	stack = *me
	return &stack
}

//func (me *NamedStack) GetDefaultServices() (sm DefaultServiceMap, sts status.Status) {
//	for range only.Once {
//		sm = make(DefaultServiceMap, 0)
//		sts = me.Refresh()
//		if status.IsError(sts) {
//			break
//		}
//		for gs, s := range me.ServiceOptions {
//			if s.DefaultService == nil {
//				continue
//			}
//			sm[gs] = s.DefaultService.Clone()
//		}
//	}
//	return sm, sts
//}

func (me *NamedStack) AddStackRole(sr *StackRole) (sts status.Status) {
	me.StackRoles = append(me.StackRoles, sr)
	return sts
}

func (me *NamedStack) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *NamedStack) Refresh(gears *Gears) (sts status.Status) {
	for range only.Once {
		if !me.NeedsRefresh() {
			break
		}

		var nsrs StackRoles
		nsrs, sts = gears.GetNamedStackRoles(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.StackRoles = nsrs

		var sos ServiceOptions
		sos, sts = gears.GetNamedServiceOptions(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.ServiceOptions = sos

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
