package models

type Projects []*Project

//var ProjectsInstance = (Projects)(nil)
//var _ ab.ItemCollection = ProjectsInstance
//
//
//func (me Projects) GetFilterMap() ab.FilterMap {
//	return GetProjectFilterMap()
//}
//
//func (me Projects) GetCollectionItemIds() (ab.ItemIds, status.Status) {
//	itemIds := make(ab.ItemIds, len(me))
//	for i, p := range me {
//		itemIds[i] = p.GetIdentifier()
//	}
//	return itemIds, nil
//}
//
//func (me Projects) AddItem(item ab.ItemInstance) (collection ab.ItemCollection, sts status.Status) {
//	found := false
//	collection = me
//	var project *Project
//	for range only.Once {
//		item, sts = me.GetItem(item.GetIdentifier())
//		if !is.Error(sts) {
//			sts = status.Fail(&status.Args{
//				Cause:      sts,
//				Message:    fmt.Sprintf("project '%s' already exists", item.GetIdentifier()),
//				HttpStatus: http.StatusConflict,
//			})
//		}
//		project, sts = AssertProject(item)
//		if is.Error(sts) {
//			break
//		}
//		found = true
//		break
//	}
//	if !found {
//		collection = append(me, project)
//	}
//	return collection, sts
//}
//
//func (me Projects) UpdateItem(item ab.ItemInstance) (collection ab.ItemCollection, sts status.Status) {
//	updated := false
//	collection = me
//	for range only.Once {
//		item, sts = me.GetItem(item.GetIdentifier())
//		if is.Error(sts) {
//			break
//		}
//		project, sts := AssertProject(item)
//		if is.Error(sts) {
//			break
//		}
//		collection, sts = collection.AddItem(project)
//		if is.Error(sts) {
//			break
//		}
//		updated = true
//	}
//	if !updated {
//		sts = status.Fail(&status.Args{
//			Message:    fmt.Sprintf("project '%s' not found", item.GetIdentifier()),
//			HttpStatus: http.StatusNotFound,
//		})
//	}
//	return collection, nil
//}
//
//func (me Projects) DeleteItem(id ab.ItemId) (collection ab.ItemCollection, sts status.Status) {
//	deleted := false
//	collection = me
//	for i, p := range me {
//		if id != p.GetIdentifier() {
//			continue
//		}
//		collection = append(me[:i], me[i+1:]...)
//		deleted = true
//		break
//	}
//	if !deleted {
//		sts = status.Fail(&status.Args{
//			Message:    fmt.Sprintf("project '%s' not found", id),
//			HttpStatus: http.StatusNotFound,
//		})
//	}
//	return collection, nil
//}
//
//func (me Projects) GetItem(id ab.ItemId) (item ab.ItemInstance, sts status.Status) {
//	found := false
//	for _, p := range me {
//		if id != p.GetIdentifier() {
//			continue
//		}
//		item = p
//		found = true
//		break
//	}
//	if !found {
//		sts = status.Fail(&status.Args{
//			Message:    fmt.Sprintf("project '%s' not found", id),
//			HttpStatus: http.StatusNotFound,
//		})
//	}
//	return item, sts
//}
//
//func (me Projects) GetItemCollection(filterPath ab.FilterPath) (collection ab.ItemCollection, sts status.Status) {
//	collection = make(Projects, len(me))
//	for _, p := range me {
//		p, sts = FilterProject(p, filterPath)
//		if p == nil {
//			continue
//		}
//		if is.Error(sts) {
//			break
//		}
//		collection, sts = collection.AddItem(p)
//		if is.Error(sts) {
//			break
//		}
//	}
//	return collection, sts
//}
//
//func (me Projects) GetCollectionSlice(filterPath ab.FilterPath) (slice ab.ItemInstances, sts status.Status) {
//	slice = make(ab.ItemInstances, len(me))
//	for i, p := range me {
//		p, sts = FilterProject(p, filterPath)
//		if p == nil {
//			continue
//		}
//		if is.Error(sts) {
//			break
//		}
//		slice[i] = p
//	}
//	return slice, sts
//}
//
