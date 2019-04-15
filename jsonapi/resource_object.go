package ja

import (
	"encoding/json"
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

var _ ResourceContainer = (*ResourceObject)(nil)
var _ apimodeler.ItemModeler = (*ResourceObject)(nil)

func (*ResourceObject) ContainsResource() {}

type ResourceObject struct {
	ResourceIdObject
	apimodeler.LinkMap `json:"links,omitempty"`
	AttributeMap       `json:"attributes"`
	RelationshipMap    `json:"relationships,omitempty"`
}

func (me *ResourceObject) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.ResourceId)
}

func (me *ResourceObject) SetId(itemid apimodeler.ItemId) (sts status.Status) {
	me.ResourceId = ResourceId(itemid)
	return sts
}

func (me *ResourceObject) GetType() apimodeler.ItemType {
	return apimodeler.ItemType(me.ResourceType)
}

func (me *ResourceObject) SetType(typ apimodeler.ItemType) (sts status.Status) {
	me.ResourceType = ResourceType(typ)
	return nil
}

func (me *ResourceObject) GetItem() (apimodeler.ItemModeler, status.Status) {
	panic("implement me")
	return nil, nil
}

func (me *ResourceObject) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	panic("implement me")
	return nil, nil
}

func (me *ResourceObject) GetRelatedItems(ctx *apimodeler.Context) (list apimodeler.List, sts status.Status) {
	panic("implement me")
	return nil, nil
}

func NewResourceObject() *ResourceObject {
	ro := ResourceObject{
		AttributeMap:     make(AttributeMap, 0),
		RelationshipMap:  make(RelationshipMap, 0),
		ResourceIdObject: *NewResourceIdObject(),
	}
	return &ro
}

func (me *ResourceObject) getRelationshipTypesData(ctx *apimodeler.Context, list apimodeler.List) (rm RelationshipMap) {

	for range only.Once {
		fnitms := make(map[Fieldname]apimodeler.List, 0)
		rfs := ctx.Controller.GetRelatedFields()
		for _, rf := range rfs {
			if rf.Include == nil && rf.IncludeType == "" {
				panic("related fields has neither Include callback or IncludeType")
			}
			fn := Fieldname(rf.Fieldname)
			fnitms[fn] = make(apimodeler.List, 0)
			for _, item := range list {
				fn = Fieldname(rf.Fieldname)
				if item.GetType() == rf.IncludeType {
					fnitms[fn] = append(fnitms[fn], item)
				} else if rf.Include != nil && rf.Include(item) {
					fnitms[fn] = append(fnitms[fn], item)
				}
			}
		}
		rm = make(RelationshipMap, 0)
		for fn, fnlst := range fnitms {
			r := NewRelationship()
			r.Data = me.getResourceIdentifier(fnlst)
			rm[fn] = r
		}
	}
	return rm
}

//
// Recursive function to return either a ResourceIdObject or a slice of ResourceIdObjects
//
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
		me.RelationshipMap = me.getRelationshipTypesData(ctx, list)
		for i, item := range list {
			ro := NewResourceObject()
			sts = ro.SetId(item.GetId())
			if is.Error(sts) {
				break
			}
			sts = ro.SetType(item.GetType())
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
		sts = ctx.RootDocumentor.SetRelated(ctx, list)
	}
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
