package ja

import (
	"gearbox/apimodeler"
)

var _ ResourceIdentifier = (ResourceIdObjects)(nil)

func (ResourceIdObjects) IdentifiesResource() {}

type ResourceIdObjects []*ResourceIdObject

var _ ResourceIdentifier = (*ResourceIdObject)(nil)

func (*ResourceIdObject) IdentifiesResource() {}

type ResourceIdObject struct {
	ResourceId         `json:"id"`
	ResourceType       `json:"type"`
	MetaMap            `json:"meta,omitempty"`
	apimodeler.LinkMap `json:"links,omitempty"`
}

func NewResourceIdObject() *ResourceIdObject {
	rido := ResourceIdObject{
		MetaMap: make(MetaMap, 0),
		LinkMap: make(apimodeler.LinkMap, 0),
	}
	return &rido
}
