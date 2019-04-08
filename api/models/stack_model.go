package models

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
)

const StacksBasepath types.Basepath = "/stacks"
const AuthorityIdParam apimodeler.IdParam = "authority"
const StacknameIdParam apimodeler.IdParam = "stackname"

var NilStackModel = (*StackModel)(nil)
var _ apimodeler.Modeler = NilStackModel

type StackModel struct {
	Gearbox gearbox.Gearboxer
}

func NewStackConnector(gb gearbox.Gearboxer) *StackModel {
	return &StackModel{
		Gearbox: gb,
	}
}

func (me *StackModel) Related() {
	return
}

func (me *StackModel) GetBasepath() types.Basepath {
	return StacksBasepath
}

func (me *StackModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *StackModel) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{
		AuthorityIdParam,
		StacknameIdParam,
	}
}

func (me *StackModel) GetNamedStackMap() (gears.NamedStackMap, status.Status) {
	return me.Gearbox.GetNamedStackMap()
}

func (me *StackModel) getGearboxStackRoleMap() (gears.StackRoleMap, status.Status) {
	return me.Gearbox.GetStackRoleMap()
}

func (me *StackModel) GetCollection(ctx apimodeler.Contexter, filterPath ...apimodeler.FilterPath) (coll apimodeler.Collection, sts status.Status) {
	var fp apimodeler.FilterPath
	if len(filterPath) > 0 {
		fp = filterPath[0]
	} else {
		fp = apimodeler.NoFilterPath
	}
	for range only.Once {
		gbsm, sts := me.GetNamedStackMap()
		if is.Error(sts) {
			break
		}
		for _, gbs := range gbsm {
			var ns *NamedStack
			ns, sts = FilterStack(ConvertNamedStack(gbs), fp)
			if is.Error(sts) {
				break
			}
			if ns == nil {
				continue
			}
			coll = append(coll, ns)
			if is.Error(sts) {
				break
			}
		}
	}
	return coll, sts
}

func (me *StackModel) FilterCollection(ctx apimodeler.Contexter, filterPath apimodeler.FilterPath) (coll apimodeler.Collection, sts status.Status) {
	for range only.Once {
		coll = make(apimodeler.Collection, 0)
		gbsm, sts := me.GetNamedStackMap()
		if is.Error(sts) {
			break
		}
		for _, gbs := range gbsm {
			var ns *NamedStack
			ns, sts = FilterStack(ConvertNamedStack(gbs), filterPath)
			if is.Error(sts) {
				break
			}
			if ns == nil {
				continue
			}
			coll = append(coll, ns)
			if is.Error(sts) {
				break
			}
		}
	}
	return coll, sts
}

func (me *StackModel) GetCollectionIds(ctx apimodeler.Contexter) (itemIds apimodeler.ItemIds, sts status.Status) {
	for range only.Once {
		gbsm, sts := me.getGearboxStackMap()
		if is.Error(sts) {
			break
		}
		itemIds = make(apimodeler.ItemIds, len(gbsm))
		i := 0
		for _, gbs := range gbsm {
			itemIds[i] = apimodeler.ItemId(gbs.GetIdentifier())
			i++
		}
	}
	return itemIds, sts
}

func (me *StackModel) AddItem(ctx apimodeler.Contexter, item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, _, sts = me.extractGearboxStack(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.AddNamedStack(gbs)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' added", gbs.GetIdentifier())
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func (me *StackModel) UpdateItem(ctx apimodeler.Contexter, item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, _, sts = me.extractGearboxStack(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.UpdateNamedStack(gbs)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' updated", item.GetId())
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts

}

func (me *StackModel) DeleteItem(ctx apimodeler.Contexter, stackid apimodeler.ItemId) (sts status.Status) {
	for range only.Once {
		sts := me.Gearbox.DeleteNamedStack(types.StackId(stackid))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts
}

func (me *StackModel) GetItem(ctx apimodeler.Contexter, stackid apimodeler.ItemId) (collection apimodeler.Itemer, sts status.Status) {
	var ns *NamedStack
	for range only.Once {
		gbns, sts := me.Gearbox.FindNamedStack(types.StackId(stackid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Stack '%s' not found", stackid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns = ConvertNamedStack(gbns)
		sts = status.Success("Stack '%s' found", stackid)
	}
	return ns, sts

}

func (me *StackModel) FilterItem(in apimodeler.Itemer, filterPath apimodeler.FilterPath) (out apimodeler.Itemer, sts status.Status) {
	for range only.Once {
		var ns *NamedStack
		ns, sts = AssertStack(in)
		if is.Error(sts) {
			break
		}
		out, sts = FilterStack(ns, filterPath)
	}
	return out, sts
}

func (me *StackModel) GetFilterMap() apimodeler.FilterMap {
	return GetStackFilterMap()
}

func (me *StackModel) getGearboxStackMap() (gears.NamedStackMap, status.Status) {
	return me.Gearbox.GetNamedStackMap()
}

func (me *StackModel) extractGearboxStack(ctx apimodeler.Contexter, item apimodeler.Itemer) (gbs *gears.NamedStack, collection apimodeler.Collection, sts status.Status) {
	var ns *NamedStack
	for range only.Once {
		collection, sts = me.GetCollection(ctx)
		if is.Error(sts) {
			break
		}
		ns, sts = AssertStack(item)
		if is.Error(sts) {
			break
		}
		gbs, sts = MakeGearboxStack(me.Gearbox, ns)
	}
	return gbs, collection, sts
}

func GetStackFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func FilterStack(in *NamedStack, filterPath apimodeler.FilterPath) (out *NamedStack, sts status.Status) {
	return in, nil
}

func AssertStack(item apimodeler.Itemer) (s *NamedStack, sts status.Status) {
	s, ok := item.(*NamedStack)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Stack: %v", item),
		})
	}
	return s, sts
}
