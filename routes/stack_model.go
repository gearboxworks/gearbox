package routes

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/modeler"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

const AuthorityIdParam modeler.IdParam = "authority"
const StacknameIdParam modeler.IdParam = "stackname"

var NilStackModel = (*StackModel)(nil)
var _ modeler.Modeler = NilStackModel

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
	return "/stacks"
}

func (me *StackModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *StackModel) GetIdParams() modeler.IdParams {
	return modeler.IdParams{
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

func (me *StackModel) GetCollection(filterPath modeler.FilterPath) (collection modeler.Collection, sts status.Status) {
	for range only.Once {
		collection = make(modeler.Collection, 0)
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
			collection = append(collection, ns)
			if is.Error(sts) {
				break
			}
		}
	}
	return collection, sts
}

func (me *StackModel) GetCollectionIds() (itemIds modeler.ItemIds, sts status.Status) {
	for range only.Once {
		gbsm, sts := me.getGearboxStackMap()
		if is.Error(sts) {
			break
		}
		itemIds = make(modeler.ItemIds, len(gbsm))
		i := 0
		for _, gbs := range gbsm {
			itemIds[i] = modeler.ItemId(gbs.GetIdentifier())
			i++
		}
	}
	return itemIds, sts
}

func (me *StackModel) AddItem(item modeler.Item) (sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, _, sts = me.extractGearboxStack(item)
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

func (me *StackModel) UpdateItem(item modeler.Item) (sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, _, sts = me.extractGearboxStack(item)
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

func (me *StackModel) DeleteItem(stackid modeler.ItemId) (sts status.Status) {
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

func (me *StackModel) GetItem(stackid modeler.ItemId, ctx echo.Context) (collection modeler.Item, sts status.Status) {
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

func (me *StackModel) FilterItem(in modeler.Item, filterPath modeler.FilterPath) (out modeler.Item, sts status.Status) {
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

func (me *StackModel) GetCollectionFilterMap() modeler.FilterMap {
	return GetStackFilterMap()
}

func (me *StackModel) getGearboxStackMap() (gears.NamedStackMap, status.Status) {
	return me.Gearbox.GetNamedStackMap()
}

func (me *StackModel) extractGearboxStack(item modeler.Item) (gbs *gears.NamedStack, collection modeler.Collection, sts status.Status) {
	var ns *NamedStack
	for range only.Once {
		collection, sts = me.GetCollection(modeler.NoFilterPath)
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

func GetStackFilterMap() modeler.FilterMap {
	return modeler.FilterMap{}
}

func FilterStack(in *NamedStack, filterPath modeler.FilterPath) (out *NamedStack, sts status.Status) {
	return in, nil
}

func AssertStack(item modeler.Item) (s *NamedStack, sts status.Status) {
	s, ok := item.(*NamedStack)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Stack: %v", item),
		})
	}
	return s, sts
}
