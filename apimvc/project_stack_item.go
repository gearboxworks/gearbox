package apimvc

import (
	"gearbox/gear"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/service"
)

type ProjectStackItemDetailMap map[gearspec.Identifier]*ProjectStackItemDetails
type ProjectStackItemDetails []*ProjectStackItemDetail
type ProjectStackItemDetail struct {
	*ProjectStackItem
	Shareable  global.ShareableChoices `json:"shareable,omitempty"`
	ServiceIds gear.Identifiers        `json:"services,omitempty"`
}

type ProjectStackItems []*ProjectStackItem
type ProjectStackItem struct {
	GearspecId gearspec.Identifier `json:"gearspec_id,omitempty"`
	ServiceId  service.Identifier  `json:"service_id,omitempty"`
}

func NewProjectStackItemFromServiceModel(sm *ServiceModel) (si *ProjectStackItem) {
	return &ProjectStackItem{
		GearspecId: sm.GearspecId,
		ServiceId:  sm.ServiceId,
	}
}
func NewProjectStackItemDetailFromServiceModel(sm *ServiceModel) (sid *ProjectStackItemDetail) {
	return &ProjectStackItemDetail{
		ProjectStackItem: &ProjectStackItem{
			GearspecId: sm.GearspecId,
			ServiceId:  sm.ServiceId,
		},
		Shareable: "no", //@TODO fix it so this is initialized
	}
}
