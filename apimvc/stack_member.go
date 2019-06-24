package apimvc

import (
	"gearbox/apiworks"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/service"
	"gearbox/types"
)

type StackMembers []*StackMember
type StackMember struct {
	GearspecId       gearspec.Identifier    `json:"gearspec_id"`
	AuthorityDomain  types.AuthorityDomain  `json:"authority"`
	Stackname        types.Stackname        `json:"stackname"`
	Specname         types.Specname         `json:"specname,omitempty"`
	Revision         types.Revision         `json:"revision,omitempty"`
	DefaultServiceId service.Identifier     `json:"default_service"`
	Shareable        global.ShareableChoice `json:"shareable,omitempty"`
	//	GearIds          service.Identifiers    `json:"available_gears,omitempty"`
}

func NewStackMemberFromGearspec(ctx *apiworks.Context, gsr *gears.Gearspec) (rss *StackMember) {
	var dsi service.Identifier
	ds := gsr.GetDefaultGear()
	if ds != nil {
		dsi = ds.GetIdentifier()
	}
	return &StackMember{
		DefaultServiceId: dsi,
		GearspecId:       gsr.GetGearspecId(),
		AuthorityDomain:  gsr.AuthorityDomain,
		Stackname:        gsr.Stackname,
		Specname:         gsr.Specname,
		Revision:         gsr.Revision,
		Shareable:        gsr.Shareable,
		//		GearIds:          gsr.Gears.GetGearIds(),
	}
}
