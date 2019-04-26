package jsonapi

import (
	"gearbox/apimodeler"
	"gearbox/status"
	"gearbox/status/is"
)

var _ ResourceContainer = (ResourceObjects)(nil)

type ResourceObjects []*ResourceObject

func (me ResourceObjects) GetRelationshipsLinkMap() (lm apimodeler.LinkMap, sts status.Status) {
	lm = make(apimodeler.LinkMap, 0)
	for fn, _lm := range me {
		for rel, link := range _lm.LinkMap {
			if rel != apimodeler.SelfRelType {
				continue
			}
			lm[apimodeler.RelType(fn)] = link
		}
	}
	return lm, sts
}

func (ResourceObjects) ContainsResource() {}

func (me ResourceObjects) SetAttributes(attrs interface{}) (sts status.Status) {
	panic("Not yet implemented")
	return nil
}

func (me ResourceObjects) AppendResourceObject(ro *ResourceObject) (ResourceObjects, status.Status) {
	return append(me, ro), nil
}

func (me ResourceObjects) SetIds(ids ResourceIds) (sts status.Status) {
	for i, ro := range me {
		sts = ro.SetStackId(apimodeler.ItemId(ids[i]))
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func (me ResourceObjects) SetTypes(types ResourceTypes) (sts status.Status) {
	for i, ro := range me {
		sts = ro.SetType(apimodeler.ItemType(types[i]))
		if is.Error(sts) {
			break
		}
	}
	return sts
}
