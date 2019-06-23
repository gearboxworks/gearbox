package apimvc

import (
	"gearbox/apiworks"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/service"
	"gearbox/types"
)

func NewStackMemberFromGearsServiceOptions(ctx *apiworks.Context, gsr *gears.StackRole, sids service.Identifiers) (rss *StackMember) {
	var dsi service.Identifier
	ds := gsr.GetDefaultService()
	if ds != nil {
		dsi = ds.GetIdentifier()
	}
	return &StackMember{
		DefaultServiceId: dsi,
		GearspecId:       gsr.GetGearspecId(),
		AuthorityDomain:  gsr.AuthorityDomain,
		Stackname:        gsr.Stackname,
		Role:             gsr.Role,
		Revision:         gsr.Revision,
		Shareable:        gsr.Shareable,
		ServiceIds:       sids,
	}
}

type StackMembers []*StackMember
type StackMember struct {
	GearspecId       gearspec.Identifier    `json:"gearspec_id"`
	AuthorityDomain  types.AuthorityDomain  `json:"authority"`
	Stackname        types.Stackname        `json:"stackname"`
	Role             types.StackRole        `json:"role,omitempty"`
	Revision         types.Revision         `json:"revision,omitempty"`
	DefaultServiceId service.Identifier     `json:"default_service,omitempty"`
	Shareable        global.ShareableChoice `json:"shareable,omitempty"`
	ServiceIds       service.Identifiers    `json:"services,omitempty"`
}
