package apimodeler

import "gearbox/status"

type Collection []Itemer

//func (Collection) ContainsResource() {}
//
//func (me Collection) GetIds() (ItemIds, status.Status) {
//	itemIds := make(ItemIds, len(me))
//	for i, item := range me {
//		itemIds[i] = item.GetId()
//	}
//	return itemIds, nil
//}
//
func (me Collection) GetCollection(Contexter, ...FilterPath) (Collection, status.Status) {
	return me, nil
}

//func GetCollectionSlice(collection Collection, sts status.Status) (Collection, status.Status) {
//	var slice = make(Collection, len(collection))
//	if is.Success(sts) {
//		for _, item := range collection {
//			slice = append(slice, item)
//		}
//	}
//	return slice, sts
//}
