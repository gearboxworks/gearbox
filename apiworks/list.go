package apiworks

import "gearbox/status"

type List []ItemModeler

//func (List) ContainsResource() {}
//
//func (me List) GetIds() (ItemIds, status.Status) {
//	itemIds := make(ItemIds, len(me))
//	for i, item := range me {
//		itemIds[i] = item.GetId()
//	}
//	return itemIds, nil
//}
//
func (me List) GetList(*Context, ...FilterPath) (List, status.Status) {
	return me, nil
}

//func GetListSlice(List List, sts status.Status) (List, status.Status) {
//	var slice = make(List, len(List))
//	if is.Success(sts) {
//		for _, item := range List {
//			slice = append(slice, item)
//		}
//	}
//	return slice, sts
//}
