package jsonapi

var NilResourceIdObject = (*ResourceIdObject)(nil)
var _ ResourceContainer = NilResourceIdObject

var _ ResourceIdentifier = (ResourceIdObjects)(nil)
var _ ResourceIdentifier = (*ResourceIdObject)(nil)

type ResourceIdObjects []*ResourceIdObject

func (ResourceIdObjects) IdentifiesResource() {}
func (*ResourceIdObject) IdentifiesResource() {}

type ResourceIdObject struct {
	ResourceId   `json:"id"`
	ResourceType `json:"type"`
	MetaMap      `json:"meta,omitempty"`
}

func (ResourceIdObject) ContainsResource() {}

func NewResourceIdObject() *ResourceIdObject {
	rido := ResourceIdObject{
		MetaMap: make(MetaMap, 0),
	}
	return &rido
}
