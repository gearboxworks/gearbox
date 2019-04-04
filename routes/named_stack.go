package routes

import (
	"gearbox/gears"
	"gearbox/gearspecid"
	"gearbox/service"
	"gearbox/types"
)

type NamedStackMap map[types.Stackname]*NamedStack
type NamedStacks []*NamedStack

type NamedStack struct {
	Authority  types.AuthorityDomain `json:"authority"`
	StackName  types.Stackname       `json:"name"`
	RoleMap    interface{}           `json:"roles"`
	ServiceMap interface{}           `json:"services"`
}

func NewNamedStack(ns *gears.NamedStack) *NamedStack {
	return &NamedStack{
		Authority:  ns.Authority,
		StackName:  ns.Stackname,
		RoleMap:    newRoleMap(ns.RoleMap),
		ServiceMap: newServiceMap(ns.ServiceMap),
	}
}

type RoleType string

type roleMap map[RoleType]*role
type role struct {
	Role     gsid.Identifier `json:"role"`
	Type     RoleType        `json:"type"`
	Name     string          `json:"name"`
	Label    string          `json:"label"`
	Max      int             `json:"max"`
	Min      int             `json:"min"`
	Optional bool            `json:"optional,omitempty"`
	Examples []string        `json:"examples"`
}

func newRoleMap(srm gears.StackRoleMap) interface{} {
	rmr := make(map[RoleType]interface{}, len(srm))
	for gs, r := range srm {
		t := RoleType(r.GetRole())
		rmr[t] = &role{
			Role:     gs,
			Type:     t,
			Name:     r.Program,
			Label:    r.Label,
			Max:      r.Maximum,
			Min:      r.Minimum,
			Optional: r.Optional,
			Examples: r.Examples,
		}
	}
	return rmr
}

type serviceMap map[gsid.Identifier]*_service
type _service struct {
	OrgName   types.OrgName          `json:"org,omitempty"`
	Default   service.Identifier     `json:"default,omitempty"`
	Shareable gears.ShareableChoices `json:"shareable,omitempty"`
	Services  service.Identifiers    `json:"options,omitempty"`
}

func newServiceMap(sm gears.ServiceOptionsMap) interface{} {
	smr := make(map[gsid.Identifier]interface{}, len(sm))
	for gs, s := range sm {
		smr[gs] = &_service{
			OrgName:   s.OrgName,
			Default:   s.DefaultService.ServiceId,
			Shareable: s.Shareable,
			Services:  s.Services.ServiceIds(),
		}
	}
	return smr
}

func ConvertNamedStack(gbp *gears.NamedStack) *NamedStack {
	return &NamedStack{
		Authority:  gbp.Authority,
		StackName:  gbp.Stackname,
		RoleMap:    gbp.RoleMap,
		ServiceMap: gbp.ServiceMap,
	}
}
