package gears

import (
	"fmt"
	"gearbox/gearspecid"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"strings"
)

//type NamedStackMap map[Stackname]StackMap

type ServiceBag map[gsid.Identifier]interface{}

type NamedStackMap map[gsid.Identifier]*NamedStack

type NamedStacks []*NamedStack

type NamedStack struct {
	Authority  types.AuthorityDomain `json:"authority"`
	Stackname  types.Stackname       `json:"name"`
	RoleMap    StackRoleMap          `json:"roles,omitempty"`
	ServiceMap ServiceOptionsMap     `json:"services,omitempty"`
	refreshed  bool
	Gears      *Gears `json:"-"`
}

//
// Get the available service options for a given named stack
//
func (me *NamedStack) GetServiceOptionMap() (rsm ServiceOptionsMap, sts status.Status) {
	for range only.Once {
		rsm = make(ServiceOptionsMap, 0)
		for gs, rso := range me.ServiceMap {
			if !strings.HasPrefix(string(gs), string(me.GetIdentifier())) {
				continue
			}
			rsm[gs] = rso
		}
	}
	return rsm, sts
}

func NewNamedStack(g *Gears, stackid types.StackId) *NamedStack {
	stack := NamedStack{
		RoleMap:    make(StackRoleMap, 0),
		ServiceMap: make(ServiceOptionsMap, 0),
		Gears:      g,
	}
	// This will split authority out, or do nothing
	_ = stack.SetIdentifier(stackid)
	return &stack
}

func (me *NamedStack) String() string {
	return string(me.Stackname)
}

func (me *NamedStack) LightweightClone() *NamedStack {
	stack := NamedStack{}
	stack = *me
	return &stack
}

func (me *NamedStack) GetDefaultServices() (sm StackMap, sts status.Status) {
	for range only.Once {
		sm = make(StackMap, 0)
		sts = me.Refresh()
		if status.IsError(sts) {
			break
		}
		for gs, s := range me.ServiceMap {
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

		var som ServiceOptionsMap
		som, sts = me.Gears.GetNamedStackServiceOptionMap(me.GetIdentifier())
		if is.Error(sts) {
			break
		}
		me.ServiceMap = som

		me.refreshed = true
	}
	if is.Success(sts) {
		sts = status.Success("named stack '%s' refreshed", me.GetIdentifier())
	}
	return sts
}

func (me *NamedStack) SetIdentifier(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		gsi := gsid.NewGearspecId()
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
