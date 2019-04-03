package routes

import (
	"fmt"
	"gearbox"
	"gearbox/apibuilder"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

const AuthorityIdParam ab.IdParam = "authority"
const StacknameIdParam ab.IdParam = "stackname"

var NilStackConnector = (*StackConnector)(nil)
var _ ab.Connector = NilStackConnector

type StackConnector struct {
	Gearbox gearbox.Gearboxer
}

func NewStackConnector(gb gearbox.Gearboxer) *StackConnector {
	return &StackConnector{
		Gearbox: gb,
	}
}

func (me *StackConnector) Related() {
	return
}

func (me *StackConnector) GetBasepath() ab.Basepath {
	return "/stacks"
}

func (me *StackConnector) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *StackConnector) GetIdParams() ab.IdParams {
	return ab.IdParams{
		AuthorityIdParam,
		StacknameIdParam,
	}
}

func (me *StackConnector) GetNamedStackMap() (gears.NamedStackMap, status.Status) {
	return me.Gearbox.GetNamedStackMap()
}

func (me *StackConnector) getGearboxStackRoleMap() (gears.StackRoleMap, status.Status) {
	return me.Gearbox.GetStackRoleMap()
}

func (me *StackConnector) GetCollection(filterPath ab.FilterPath) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		collection = make(ab.Collection, 0)
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

func (me *StackConnector) GetCollectionIds() (itemIds ab.ItemIds, sts status.Status) {
	for range only.Once {
		gbsm, sts := me.getGearboxStackMap()
		if is.Error(sts) {
			break
		}
		itemIds = make(ab.ItemIds, len(gbsm))
		i := 0
		for _, gbs := range gbsm {
			itemIds[i] = ab.ItemId(gbs.GetIdentifier())
			i++
		}
	}
	return itemIds, sts
}

func (me *StackConnector) AddItem(item ab.Item) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, collection, sts = me.extractGearboxStack(item)
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
	return collection, sts
}

func (me *StackConnector) UpdateItem(item ab.Item) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, collection, sts = me.extractGearboxStack(item)
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
	return collection, sts

}

func (me *StackConnector) DeleteItem(stackid ab.ItemId) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		sts := me.Gearbox.DeleteNamedStack(types.StackId(stackid))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return collection, sts
}

func (me *StackConnector) GetItem(stackid ab.ItemId, ctx echo.Context) (collection ab.Item, sts status.Status) {
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

func (me *StackConnector) FilterItem(in ab.Item, filterPath ab.FilterPath) (out ab.Item, sts status.Status) {
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

func (me *StackConnector) GetCollectionFilterMap() ab.FilterMap {
	return GetStackFilterMap()
}

func (me *StackConnector) getGearboxStackMap() (gears.NamedStackMap, status.Status) {
	return me.Gearbox.GetNamedStackMap()
}

func (me *StackConnector) extractGearboxStack(item ab.Item) (gbs *gears.NamedStack, collection ab.Collection, sts status.Status) {
	var ns *NamedStack
	for range only.Once {
		collection, sts = me.GetCollection(ab.NoFilterPath)
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

func GetStackFilterMap() ab.FilterMap {
	return ab.FilterMap{}
}

func FilterStack(in *NamedStack, filterPath ab.FilterPath) (out *NamedStack, sts status.Status) {
	return in, nil
}

func AssertStack(item ab.Item) (s *NamedStack, sts status.Status) {
	s, ok := item.(*NamedStack)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Stack: %v", item),
		})
	}
	return s, sts
}
