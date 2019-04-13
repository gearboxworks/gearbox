package ja

import (
	"encoding/json"
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

var _ ResourceContainer = (*ResourceObject)(nil)

func (*ResourceObject) ContainsResource() {}

type ResourceObject struct {
	ResourceIdObject
	apimodeler.LinkMap `json:"links,omitempty"`
	AttributeMap       `json:"attributes"`
	RelationshipMap    `json:"relationships,omitempty"`
}

func NewResourceObject() *ResourceObject {
	ro := ResourceObject{
		AttributeMap:     make(AttributeMap, 0),
		RelationshipMap:  make(RelationshipMap, 0),
		ResourceIdObject: *NewResourceIdObject(),
	}
	return &ro
}

func (me *ResourceObject) getRelationshipTypesData(list apimodeler.List) RelationshipMap {
	fnitms := make(map[Fieldname]apimodeler.List, 0)
	for _, item := range list {
		fn := Fieldname(item.GetType())
		if _, ok := fnitms[fn]; !ok {
			fnitms[fn] = make(apimodeler.List, 0)
		}
		fnitms[fn] = append(fnitms[fn], item)
	}
	rm := make(RelationshipMap, 0)
	for fn, fnlst := range fnitms {
		r := NewRelationship()
		r.Data = me.getResourceIdentifier(fnlst)
		rm[fn] = r
	}
	return rm
}

func (me *ResourceObject) getResourceIdentifier(list apimodeler.List) (ri ResourceIdentifier) {
	switch len(list) {
	case 0:
		break
	case 1:
		rio := &ResourceIdObject{}
		item := list[0]
		rio.ResourceId = ResourceId(item.GetId())
		rio.ResourceType = ResourceType(item.GetType())
		ri = rio
	default:
		rios := make(ResourceIdObjects, len(list))
		for i, item := range list {
			oil := apimodeler.List{item}
			rios[i] = me.getResourceIdentifier(oil).(*ResourceIdObject)
		}
		ri = rios
	}
	return ri
}

func (me *ResourceObject) SetRelatedItems(ctx *apimodeler.Context, list apimodeler.List) (sts status.Status) {
	for range only.Once {
		me.RelationshipMap = me.getRelationshipTypesData(list)
		for i, item := range list {
			ro := NewResourceObject()
			sts = ro.SetId(ResourceId(item.GetId()))
			if is.Error(sts) {
				break
			}
			sts = ro.SetType(ResourceType(item.GetType()))
			if is.Error(sts) {
				break
			}
			b, err := json.Marshal(item)
			if err != nil {
				sts = status.OurBad("cannot marshal related item '%s", item.GetId())
				break
			}
			am := make(AttributeMap, 0)
			err = json.Unmarshal(b, &am)
			if err != nil {
				sts = status.OurBad("cannot unmarshal related item '%s'", item.GetId())
				break
			}
			ro.AttributeMap = am
			ii := IncludedItem(*ro)
			list[i] = &ii
		}
		sts = ctx.RootDocumentor.SetIncluded(ctx, list)
	}
	return nil
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
