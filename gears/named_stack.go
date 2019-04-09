package gears

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"strings"
)

type NamedStackMap map[types.StackId]*NamedStack

type NamedStacks []*NamedStack

type NamedStack struct {
	Authority       types.AuthorityDomain `json:"authority"`
	Stackname       types.Stackname       `json:"name"`
	RoleMap         StackRoleMap          `json:"roles,omitempty"`
	RoleServicesMap RoleServicesMap       `json:"services,omitempty"`
	refreshed       bool
	Gears           *Gears `json:"-"`
}

func NewNamedStack(g *Gears, stackid types.StackId) *NamedStack {
	stack := NamedStack{
		RoleMap:         make(StackRoleMap, 0),
		RoleServicesMap: make(RoleServicesMap, 0),
		Gears:           g,
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
		gsids = make(gearspec.Identifiers, 0)
		nsrm, sts := me.Gears.GetNamedStackRoleMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		for _, r := range nsrm {
			gsids = append(gsids, r.GetGearspecId())
		}
	}
	return gsids, sts
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

func (me *NamedStack) GetDefaultServices() (sm DefaultServiceMap, sts status.Status) {
	for range only.Once {
		sm = make(DefaultServiceMap, 0)
		sts = me.Refresh()
		if status.IsError(sts) {
			break
		}
		for gs, s := range me.RoleServicesMap {
			if s.DefaultService == nil {
				continue
			}
			sm[gs] = MakeService(s.DefaultService)
		}
	}
	return sm, sts
}

func MakeService(s *Service) *Service {
	var x *Service
	_, _ = x.Parse() // Cause a crash so we can implement what is needed here.
	return s
}

func (me *NamedStack) AddStackRole(sr *StackRole) (sts status.Status) {
	me.RoleMap[sr.GetGearspecId()] = sr
	return sts
}

func (me *NamedStack) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *NamedStack) Refresh() (sts status.Status) {
	for range only.Once {
		if !me.NeedsRefresh() {
			break
		}

		var nsrm StackRoleMap
		nsrm, sts = me.Gears.GetNamedStackRoleMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.RoleMap = nsrm

		var som RoleServicesMap
		som, sts = me.Gears.GetNamedStackServiceOptionMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.RoleServicesMap = som

		me.refreshed = true
	}
	if is.Success(sts) {
		sts = status.Success("named stack '%s' refreshed", me.GetIdentifier())
	}
	return sts
}

func (me *NamedStack) SetIdentifier(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		gsi := gearspec.NewGearspecId()
		sts := gsi.SetStackId(stackid)
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

func (me *Gears) ValidateNamedStackId(stackid types.StackId) (sts status.Status) {
	_, sts = me.FindNamedStackId(stackid)
	return sts
}

func (me *Gears) FindNamedStackId(stackid types.StackId) (sid types.StackId, sts status.Status) {
	for range only.Once {
		var stackid types.StackId
		var ok bool
		for _, nsid := range me.NamedStackIds {
			if nsid == stackid {
				sid = nsid
				ok = true
				break
			}
		}
		if !ok {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("named stack ID '%s' not found", stackid),
				HttpStatus: http.StatusNotFound,
				Help:       fmt.Sprintf("see valid named stack IDs at %s", JsonUrl),
			})
		} else {
			sts = status.Success("named stack ID '%s' found", stackid)
		}
	}
	return sid, sts
}
