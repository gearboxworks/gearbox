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
	Authority       types.AuthorityDomain `json:"authority"`
	Stackname       types.Stackname       `json:"name"`
	RoleMap         StackRoleMap          `json:"roles,omitempty"`
	RoleServicesMap RoleServicesMap       `json:"services,omitempty"`
	Gears           *Gears                `json:"-"`
	refreshed       bool
}

//func NewNamedStack(gears *Gears, stackid types.StackId) *NamedStack {
func NewNamedStack(stackid types.StackId) *NamedStack {
	stack := NamedStack{
		RoleMap:         make(StackRoleMap, 0),
		RoleServicesMap: make(RoleServicesMap, 0),
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
		nsrm, sts := me.Gears.GetNamedStackRoleMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		for _, r := range nsrm {
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
// Get the available service options for a given named stack
//
func (me *NamedStack) GetRoleMap() (srm StackRoleMap, sts status.Status) {
	for range only.Once {
		srm = make(StackRoleMap, 0)
		for gs, rso := range me.RoleMap {
			if !strings.HasPrefix(string(gs), string(me.GetIdentifier())) {
				continue
			}
			srm[gs] = rso
		}
	}
	return srm, sts
}

//
// Get the available service options for a given named stack
//
func (me *NamedStack) GetServiceOptionMap() (rsm RoleServicesMap, sts status.Status) {
	for range only.Once {
		rsm = make(RoleServicesMap, 0)
		for gs, rso := range me.RoleServicesMap {
			if !strings.HasPrefix(string(gs), string(me.GetIdentifier())) {
				continue
			}
			rsm[gs] = rso
		}
	}
	return rsm, sts
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
//		for gs, s := range me.RoleServicesMap {
//			if s.DefaultService == nil {
//				continue
//			}
//			sm[gs] = s.DefaultService.Clone()
//		}
//	}
//	return sm, sts
//}

func (me *NamedStack) AddStackRole(sr *StackRole) (sts status.Status) {
	me.RoleMap[sr.GetGearspecId()] = sr
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

		var nsrm StackRoleMap
		nsrm, sts = gears.GetNamedStackRoleMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.RoleMap = nsrm

		var rsm RoleServicesMap
		rsm, sts = gears.GetNamedStackRoleServicesMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.RoleServicesMap = rsm

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
		sts := gsi.SetId(stackid)
		if status.IsError(sts) {
			break
		}
		me.Authority = gsi.Authority
		me.Stackname = gsi.Stackname
	}
	return sts
}

func (me *NamedStack) GetIdentifier() types.StackId {
	return types.StackId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
}
