package apimodels

import (
	"gearbox/api/global"
	"gearbox/gear"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/types"
)

type RoleType string

type roleMap map[RoleType]*role
type role struct {
	Role     gearspec.Identifier `json:"role"`
	Type     RoleType            `json:"type"`
	Name     string              `json:"name"`
	Label    string              `json:"label"`
	Max      int                 `json:"max"`
	Min      int                 `json:"min"`
	Optional bool                `json:"optional,omitempty"`
	Examples []string            `json:"examples"`
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

type serviceMap map[gearspec.Identifier]*_service
type _service struct {
	OrgName   types.OrgName           `json:"org,omitempty"`
	Default   gear.Identifier         `json:"default,omitempty"`
	Shareable global.ShareableChoices `json:"shareable,omitempty"`
	Services  gear.Identifiers        `json:"choices,omitempty"`
}

func newServiceMap(sm gears.RoleServicesMap) interface{} {
	smr := make(map[gearspec.Identifier]interface{}, len(sm))
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
