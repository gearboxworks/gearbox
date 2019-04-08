package ja

import (
	"encoding/json"
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status"
)

var _ ResourceContainer = (*ResourceObject)(nil)

func (*ResourceObject) ContainsResource() {}

type ResourceObject struct {
	ResourceIdObject
	LinkMap         apimodeler.LinkMap `json:"links,omitempty"`
	AttributeMap    `json:"attributes"`
	RelationshipMap `json:"relationships,omitempty"`
}

func NewResourceObject() *ResourceObject {
	ro := ResourceObject{
		AttributeMap:     make(AttributeMap, 0),
		RelationshipMap:  make(RelationshipMap, 0),
		ResourceIdObject: *NewResourceIdObject(),
	}
	return &ro
}
func (me *ResourceObject) SetId(id ResourceId) status.Status {
	me.ResourceId = id
	return nil
}

func (me *ResourceObject) GetId() ResourceId {
	return me.ResourceId
}

func (me *ResourceObject) GetType() ResourceType {
	return me.ResourceType
}

func (me *ResourceObject) SetType(_typ ResourceType) (sts status.Status) {
	me.ResourceType = _typ
	return nil
}

func (me *ResourceObject) SetAttributes(attrs interface{}) (sts status.Status) {
	var attrMap AttributeMap
	for range only.Once {
		b, err := json.Marshal(attrs)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: "unable to marshal attributes",
			})
			break
		}
		attrMap = make(AttributeMap, 0)
		err = json.Unmarshal(b, &attrMap)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: "unable to unmarshal attributes",
			})
			break
		}
		me.AttributeMap = attrMap
	}
	return sts
}

func (me *ResourceObject) SetLinks(links apimodeler.LinkMap) {
	me.LinkMap = links

}
