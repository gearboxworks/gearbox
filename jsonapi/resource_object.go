package jsonapi

import (
	"encoding/json"
	"fmt"
	"gearbox/apiworks"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/util"
)

var NilResourceObject = (*ResourceObject)(nil)
var _ ResourceContainer = NilResourceObject
var _ RelationshipsLinkMapGetter = NilResourceObject
var _ apiworks.ItemModeler = NilResourceObject

func (*ResourceObject) ContainsResource() {}

type ResourceObject struct {
	ResourceIdObject
	apiworks.LinkMap `json:"links,omitempty"`
	AttributeMap     `json:"attributes"`
	RelationshipMap  `json:"relationships,omitempty"`
}

func AssertResourceObject(item apiworks.ItemModeler) (ro *ResourceObject, sts status.Status) {
	for range only.Once {
		var ok bool
		ro, ok = item.(*ResourceObject)
		if !ok {
			sts = status.OurBad("item '%s' is not a %T", item.GetId(), ro)
			break
		}
	}
	return ro, sts
}

func (me *ResourceObject) MarshalAttributeMap() (b []byte, sts status.Status) {
	for range only.Once {
		var err error
		b, err = json.Marshal(me.AttributeMap)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to marshal AttributeMap for resource object '%s'",
					me.GetId(),
				),
			})
			break
		}
	}
	return b, sts
}

func (me *ResourceObject) GetAttributeMap() AttributeMap {
	return me.AttributeMap
}

func (me *ResourceObject) GetRelationshipsLinkMap() (lm apiworks.LinkMap, sts status.Status) {
	lm = make(apiworks.LinkMap, 0)
	for fn, f := range me.RelationshipMap {
		link, ok := f.LinkMap[apiworks.SelfRelType]
		if !ok {
			panic(fmt.Sprintf("relationship '%s' does not have a 'self' link.", fn))
		}
		lm[apiworks.RelType(fn)] = link
	}
	return lm, sts
}

func (me *ResourceObject) GetId() apiworks.ItemId {
	return apiworks.ItemId(me.ResourceId)
}

func (me *ResourceObject) SetId(itemid apiworks.ItemId) (sts status.Status) {
	me.ResourceId = ResourceId(itemid)
	return sts
}

func (me *ResourceObject) GetType() apiworks.ItemType {
	return apiworks.ItemType(me.ResourceType)
}

func (me *ResourceObject) SetType(typ apiworks.ItemType) (sts status.Status) {
	me.ResourceType = ResourceType(typ)
	return nil
}

func (me *ResourceObject) GetItem() (apiworks.ItemModeler, status.Status) {
	panic("implement me")
	return nil, nil
}

func (me *ResourceObject) GetItemLinkMap(*apiworks.Context) (apiworks.LinkMap, status.Status) {
	panic("implement me")
	return nil, nil
}

func (me *ResourceObject) GetRelatedItems(ctx *apiworks.Context) (list apiworks.List, sts status.Status) {
	panic("implement me")
	return nil, nil
}

func NewResourceObject() *ResourceObject {
	ro := ResourceObject{
		ResourceIdObject: *NewResourceIdObject(),
	}
	ro.Renew()
	return &ro
}

func (me *ResourceObject) Renew() {
	if me.LinkMap == nil {
		me.LinkMap = make(apiworks.LinkMap, 0)
	}
	if me.AttributeMap == nil {
		me.AttributeMap = make(AttributeMap, 0)
	}
	if me.RelationshipMap == nil {
		me.RelationshipMap = make(RelationshipMap, 0)
	}
}

func (me *ResourceObject) getRelationshipTypesData(ctx *apiworks.Context, item apiworks.ItemModeler, list apiworks.List) (rm RelationshipMap) {
	for range only.Once {
		fnitms := make(map[Fieldname]apiworks.List, 0)
		rfs := ctx.Controller.GetRelatedFields()
		for _, rf := range rfs {
			if rf.Include == nil && rf.IncludeType == "" {
				panic("related fields has neither Include callback or IncludeType")
			}
			fn := Fieldname(rf.Fieldname)
			fnitms[fn] = make(apiworks.List, 0)
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
			r.Data = me.getResourceIdentifier(ctx, item, fnlst)
			r.LinkMap.AddLink(apiworks.SelfRelType, me.getRelationshipSelfLink(ctx, fn))
			rm[fn] = r
		}
	}
	return rm
}

//
//
//
func (me *ResourceObject) getRelationshipSelfLink(ctx *apiworks.Context, fieldname Fieldname) (link apiworks.Link) {
	for range only.Once {
		baseurl, sts := ctx.GetRequestPath()
		if is.Error(sts) {
			break
		}
		path := util.Dashify(string(fieldname))
		link = apiworks.Link(fmt.Sprintf("%s/%s/", baseurl, path))
	}
	return link
}

//
// Recursive function to return either a ResourceIdObject or a slice of ResourceIdObjects
//
func (me *ResourceObject) getResourceIdentifier(ctx *apiworks.Context, item apiworks.ItemModeler, list apiworks.List) (ri ResourceIdentifier) {
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
			oil := apiworks.List{item}
			rios[i] = me.getResourceIdentifier(ctx, item, oil).(*ResourceIdObject)
		}
		ri = rios
	}
	return ri
}

func (me *ResourceObject) SetRelatedItems(ctx *apiworks.Context, item apiworks.ItemModeler, list apiworks.List) (sts status.Status) {
	for range only.Once {
		me.RelationshipMap = me.getRelationshipTypesData(ctx, item, list)
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

func (me *ResourceObject) SetLinks(links apiworks.LinkMap) {
	me.LinkMap = links

}
