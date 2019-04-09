package models

import (
	"gearbox/gearid"
	"gearbox/gears"
	gsid "gearbox/gearspecid"
	"gearbox/types"
)

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
			Name:     r.Name,
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
	Default   gearid.GearIdentifier  `json:"default,omitempty"`
	Shareable gears.ShareableChoices `json:"shareable,omitempty"`
	Services  gearid.GearIdentifiers `json:"choices,omitempty"`
}

func newServiceMap(sm gears.RoleServicesMap) interface{} {
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
		Stackname:  gbp.Stackname,
		ServiceMap: gbp.RoleServicesMap,
	}
}
